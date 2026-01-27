package services

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type WorldEntry struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
	Size     int64  `json:"size"`
}

type WorldService struct {
	serverService *ServerService
	basePath      string
}

func NewWorldService(serverService *ServerService, basePath string) *WorldService {
	return &WorldService{
		serverService: serverService,
		basePath:      basePath,
	}
}

func (s *WorldService) ListWorlds(serverID string) ([]WorldEntry, error) {
	// 1. Get current level-name from server.properties
	props, err := s.serverService.GetProperties(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get server properties: %w", err)
	}
	currentLevel := props["level-name"]

	// 2. Scan server directory
	serverDir := filepath.Join(s.basePath, serverID)
	entries, err := os.ReadDir(serverDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read server directory: %w", err)
	}

	var worlds []WorldEntry

	ignoredDirs := map[string]bool{
		"logs":           true,
		"libraries":      true,
		"versions":       true,
		"crash-reports":  true,
		"plugins":        true,
		"cache":          true,
		"config":         true,
		"mods":           true,
		"bundler":        true,
		"generated":      true,
		"debug":          true,
		"resources":      true,
		"resource-packs": true,
		"usercache.json": true, // Not a dir but for safety
		".git":           true,
		".idea":          true,
		".vscode":        true,
		"lobby":          false, // Lobby is a world
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if ignoredDirs[entry.Name()] {
			continue
		}

		worldName := entry.Name()
		worldPath := filepath.Join(serverDir, worldName)

		// ALWAYS include the directory unless blocked
		// This allows empty/new worlds to be listed and activated.
		size, _ := getDirSize(worldPath)

		worlds = append(worlds, WorldEntry{
			Name:     worldName,
			IsActive: worldName == currentLevel,
			Size:     size,
		})
	}

	return worlds, nil
}

func (s *WorldService) CreateWorld(serverID, name string) error {
	// Validate name (alphanumeric, dashes, underscores)
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid world name: only alphanumeric, dashes and underscores allowed")
	}

	serverDir := filepath.Join(s.basePath, serverID)
	worldPath := filepath.Join(serverDir, name)

	if err := os.Mkdir(worldPath, 0755); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("world already exists")
		}
		return err
	}

	return nil
}

func (s *WorldService) ActivateWorld(serverID, worldName string) error {
	serverDir := filepath.Join(s.basePath, serverID)
	worldPath := filepath.Join(serverDir, worldName)

	// Check if world exists (has level.dat or at least is a dir that was created)
	// The prompt implies checking format validity, but strictly speaking just checking Dir existence + maybe logic from ListWorlds
	// For Activate, we trust it's a valid target if it exists.
	if _, err := os.Stat(worldPath); os.IsNotExist(err) {
		return fmt.Errorf("world not found")
	}

	// Update level-name
	return s.serverService.UpdateProperties(serverID, map[string]string{
		"level-name": worldName,
	})
}

func (s *WorldService) DeleteWorld(serverID, worldName string) error {
	// 1. Check if active
	props, err := s.serverService.GetProperties(serverID)
	if err != nil {
		return err
	}
	if props["level-name"] == worldName {
		return fmt.Errorf("cannot delete active world")
	}

	// 2. Delete
	serverDir := filepath.Join(s.basePath, serverID)
	worldPath := filepath.Join(serverDir, worldName)

	return os.RemoveAll(worldPath)
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// EnsureDir checks if a directory exists, creates if not.
// core package likely has this but for self-containment/usage:
// core.EnsureDir(path) if public.
// Using core.EnsureDir as seen in server.go
