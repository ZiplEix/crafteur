package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type BackupEntry struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

type BackupService struct {
	sourcePath string
	backupPath string
}

func NewBackupService(sourcePath, backupPath string) *BackupService {
	return &BackupService{
		sourcePath: sourcePath,
		backupPath: backupPath,
	}
}

func (s *BackupService) CreateBackup(serverID string) error {
	if strings.Contains(serverID, "..") {
		return fmt.Errorf("invalid server ID")
	}

	sourceDir := filepath.Join(s.sourcePath, serverID)
	destDir := filepath.Join(s.backupPath, serverID)

	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	zipFilename := fmt.Sprintf("backup-%s.zip", timestamp)
	zipPath := filepath.Join(destDir, zipFilename)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Skip exclusions
		if info.Name() == "session.lock" {
			return nil
		}
		if relPath == "." {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Ensure we keep directory structure
		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Read file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func (s *BackupService) ListBackups(serverID string) ([]BackupEntry, error) {
	if strings.Contains(serverID, "..") {
		return nil, fmt.Errorf("invalid server ID")
	}

	backupDir := filepath.Join(s.backupPath, serverID)
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []BackupEntry{}, nil
		}
		return nil, err
	}

	var backups []BackupEntry
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".zip") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, BackupEntry{
			Name:      entry.Name(),
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	// Sort by CreatedAt descending
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	return backups, nil
}

func (s *BackupService) GetBackupPath(serverID, filename string) (string, error) {
	if strings.Contains(serverID, "..") || strings.Contains(filename, "..") {
		return "", fmt.Errorf("invalid path")
	}

	path := filepath.Join(s.backupPath, serverID, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("backup not found")
	}

	return path, nil
}

func (s *BackupService) DeleteBackup(serverID, filename string) error {
	if strings.Contains(serverID, "..") || strings.Contains(filename, "..") {
		return fmt.Errorf("invalid path")
	}

	path := filepath.Join(s.backupPath, serverID, filename)
	return os.Remove(path)
}
