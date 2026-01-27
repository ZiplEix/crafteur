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

// SearchMods: facets examples: `["versions:1.20.1"]`, `["categories:fabric"]`
func (s *ModrinthService) SearchMods(query string, limit int, offset int, facets []string) (*core.ModrinthSearchResponse, error) {
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

	// Facets must be a JSON array of arrays of strings: [["versions:1.20.1"], ["categories:fabric"]]
	// We receive a flat list of simple filters that we need to group or format correctly.
	// For simplicity, let's assume the caller passes properly formatted facet groups if needed,
	// OR we assume AND logic between passed facets.
	// Modrinth API expects: facets=[["versions:1.20.1"],["categories:fabric"]] for AND
	if len(facets) > 0 {
		// Convert flat slice to slice of slices for AND logic
		var facetGroups [][]string
		for _, f := range facets {
			facetGroups = append(facetGroups, []string{f})
		}
		facetJson, err := json.Marshal(facetGroups)
		if err == nil {
			q.Set("facets", string(facetJson))
		}
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

func (s *ModrinthService) InstallMod(serverID string, projectID string) error {
	// 1. Get Server Config
	server, err := s.serverService.GetServer(serverID)
	if err != nil {
		return err
	}

	// 2. Prepare Filters for Version Search
	// loaders: ["fabric"] or ["forge"] (vanilla servers usually don't install mods, but let's assume fabric if compatible)
	// game_versions: ["1.20.1"]

	var loaders []string
	if server.Type == core.TypeFabric {
		loaders = append(loaders, "fabric")
	} else if server.Type == core.TypeForge {
		loaders = append(loaders, "forge")
	} else {
		// If vanilla, maybe don't allow? Or default to fabric?
		// For now let's just use the type as loader if it matches, or error out.
		return fmt.Errorf("server type %s supports no mods or is not configured", server.Type)
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
		return fmt.Errorf("no compatible version found for %s on %s", loaders[0], server.Version)
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
	targetDir := filepath.Join(s.serverService.GetDataDir(), serverID, "mods")
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
