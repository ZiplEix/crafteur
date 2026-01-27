package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type AddonService struct {
	serverService *ServerService
	dataDir       string
}

func NewAddonService(serverService *ServerService, dataDir string) *AddonService {
	return &AddonService{
		serverService: serverService,
		dataDir:       dataDir,
	}
}

func (s *AddonService) GetAddonPath(serverID string, addonType string) (string, error) {
	serverRoot := filepath.Join(s.dataDir, serverID)

	// Ensure server exists (simple check via root dir)
	if _, err := os.Stat(serverRoot); os.IsNotExist(err) {
		return "", fmt.Errorf("server not found")
	}

	var targetPath string

	switch addonType {
	case "mods":
		targetPath = filepath.Join(serverRoot, "mods")
	case "plugins":
		targetPath = filepath.Join(serverRoot, "plugins")
	case "datapacks":
		// Need world name from properties
		props, err := s.serverService.GetProperties(serverID)
		if err != nil {
			return "", fmt.Errorf("failed to read server.properties: %w", err)
		}
		levelName := props["level-name"]
		if levelName == "" {
			levelName = "world" // Default fallback
		}
		targetPath = filepath.Join(serverRoot, levelName, "datapacks")
	default:
		return "", fmt.Errorf("invalid addon type: %s", addonType)
	}

	// Ensure directory exists
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create addon directory: %w", err)
	}

	return targetPath, nil
}

func (s *AddonService) ListAddons(serverID, addonType string) ([]FileInfo, error) {
	path, err := s.GetAddonPath(serverID, addonType)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make([]FileInfo, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Filter by extension based on type
		if addonType == "datapacks" {
			if !strings.HasSuffix(name, ".zip") {
				continue
			}
		} else {
			// mods and plugins are usually .jar
			if !strings.HasSuffix(name, ".jar") {
				continue
			}
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Name:    name,
			IsDir:   false,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}
	return files, nil
}

func (s *AddonService) UploadAddon(serverID, addonType string, fileHeader *multipart.FileHeader) error {
	targetDir, err := s.GetAddonPath(serverID, addonType)
	if err != nil {
		return err
	}

	// Validate extension
	ext := filepath.Ext(fileHeader.Filename)
	if addonType == "datapacks" {
		if ext != ".zip" {
			return fmt.Errorf("invalid file type for datapack, allowed: .zip")
		}
	} else {
		if ext != ".jar" {
			return fmt.Errorf("invalid file type, allowed: .jar")
		}
	}

	fullPath := filepath.Join(targetDir, fileHeader.Filename)

	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (s *AddonService) DeleteAddon(serverID, addonType, filename string) error {
	path, err := s.GetAddonPath(serverID, addonType)
	if err != nil {
		return err
	}

	// Security: Prevent directory traversal
	if filepath.Base(filename) != filename {
		return fmt.Errorf("invalid filename")
	}

	fullPath := filepath.Join(path, filename)
	return os.Remove(fullPath)
}
