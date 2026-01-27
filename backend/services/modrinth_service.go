package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ZiplEix/crafteur/core"
)

type ModrinthService struct {
	serverService *ServerService
}

func NewModrinthService(serverService *ServerService) *ModrinthService {
	return &ModrinthService{
		serverService: serverService,
	}
}

// SearchProjects: facets examples: `["versions:1.20.1"]`, `["categories:fabric"]`
// projectType: "mod" or "plugin"
func (s *ModrinthService) SearchProjects(query string, limit int, offset int, facets []string, projectType string) (*core.ModrinthSearchResponse, error) {
	baseUrl := "https://api.modrinth.com/v2/search"
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	if query != "" {
		q.Set("query", query)
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}

	// Dynamic Facuts
	// Base facets (coming from caller usually version): facets
	// We want to append our projectType specific facets.

	// Prepare simplified facet list.
	// We assume `facets` passed in contains the version filter e.g. "versions:1.20.1"

	// Filter by project type + loader/category
	var typeFacets []string
	if projectType == "plugin" {
		typeFacets = []string{"project_type:plugin"}
		// Plugin usually means bukkit/spigot/paper
		// categories:bukkit OR categories:spigot OR categories:paper
		// Modrinth API facets: [[A, B]] means A OR B. [[A], [B]] means A AND B.
		// We want (Version) AND (Plugin) AND (Bukkit OR Spigot OR Paper)
	} else {
		typeFacets = []string{"project_type:mod"}
		// Mod usually means Fabric or Forge (handled by caller passing "categories:fabric" in basic facets?
		// Or we should handle it here like the controller did?)
		// The controller was passing facets based on server type.
		// Let's rely on controller passing the right loader facet for mods.
	}

	// Construct the final facet structure
	// finalFacets = [ [passedFacets...], [project_type], [Categories...] ]
	// Actually caller (Controller) constructs "versions:..." and "categories:fabric"
	// We just need to add project_type filter.

	var finalFacets [][]string

	// Add passed facets (AND groups)
	// We assumed caller passes flat string list and we grouped them individually (AND).
	// But wait, the controller logic needs to be aligned.
	// Previous: SearchMods took flat `facets` and did `facetGroups = append(facetGroups, []string{f})` -> AND logic.

	for _, f := range facets {
		finalFacets = append(finalFacets, []string{f})
	}

	// Add Project Type
	if len(typeFacets) > 0 {
		finalFacets = append(finalFacets, typeFacets)
	}

	// For Plugins, we want to expand the search to Bukkit/Spigot/Paper categories if not already present?
	// The user request says:
	// If plugin: Filter on `["categories:bukkit", "categories:spigot", "categories:paper"]` (OR logic).
	if projectType == "plugin" {
		pluginCats := []string{"categories:bukkit", "categories:spigot", "categories:paper"}
		finalFacets = append(finalFacets, pluginCats)
	}

	facetJson, err := json.Marshal(finalFacets)
	if err == nil {
		q.Set("facets", string(facetJson))
	}

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("modrinth api error: %d", resp.StatusCode)
	}

	var searchResp core.ModrinthSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}

func (s *ModrinthService) InstallProject(serverID string, projectID string, projectType string) error {
	// 1. Get Server Config
	server, err := s.serverService.GetServer(serverID)
	if err != nil {
		return err
	}

	// 2. Prepare Filters for Version Search
	var loaders []string

	if projectType == "plugin" {
		// For plugins, loaders are "bukkit", "spigot", "paper", "purpur"?
		// Actually Modrinth uses "bukkit", "spigot", "paper" as loaders too for plugins.
		// Let's add them all to match broadly.
		loaders = append(loaders, "bukkit", "spigot", "paper", "purpur")
		// Ideally we match the server software but plugins are generally cross-compatible on these platforms.
	} else {
		// Mods
		if server.Type == core.TypeFabric {
			loaders = append(loaders, "fabric")
		} else if server.Type == core.TypeForge {
			loaders = append(loaders, "forge")
		} else {
			return fmt.Errorf("server type %s supports no mods or is not configured", server.Type)
		}
	}

	gameVersions := []string{server.Version}

	// 3. Query Modrinth Versions
	// GET /project/{id}/version
	baseUrl := fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", projectID)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}

	q := u.Query()
	if loadersJson, err := json.Marshal(loaders); err == nil {
		q.Set("loaders", string(loadersJson))
	}
	if gameVersionsJson, err := json.Marshal(gameVersions); err == nil {
		q.Set("game_versions", string(gameVersionsJson))
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("modrinth version api error: %d", resp.StatusCode)
	}

	var versions []core.ModrinthVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return err
	}

	if len(versions) == 0 {
		loaderStr := loaders[0]
		if len(loaders) > 1 {
			loaderStr = fmt.Sprintf("%v", loaders)
		}
		return fmt.Errorf("no compatible version found for %s on %s", loaderStr, server.Version)
	}

	// 4. Pick best version (first one is usually latest compatible)
	bestVersion := versions[0]

	// 5. Find primary file
	var fileToDownload *core.ModrinthFile
	for _, f := range bestVersion.Files {
		if f.Primary {
			fileToDownload = &f
			break
		}
	}
	// Fallback if no primary
	if fileToDownload == nil && len(bestVersion.Files) > 0 {
		fileToDownload = &bestVersion.Files[0]
	}

	if fileToDownload == nil {
		return fmt.Errorf("no files found in version")
	}

	// 6. Download
	// Determine folder based on type
	folderName := "mods"
	if projectType == "plugin" {
		folderName = "plugins"
	}

	targetDir := filepath.Join(s.serverService.GetDataDir(), serverID, folderName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	targetPath := filepath.Join(targetDir, fileToDownload.Filename)

	// Check if exists
	if _, err := os.Stat(targetPath); err == nil {
		// Already exists
		return nil
	}

	if err := core.DownloadFile(fileToDownload.Url, targetPath); err != nil {
		return err
	}

	return nil
}
