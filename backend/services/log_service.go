package services

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type LogFileEntry struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

type LogService struct {
	basePath string
}

func NewLogService(basePath string) *LogService {
	return &LogService{
		basePath: basePath,
	}
}

func (s *LogService) ListLogFiles(serverID string) ([]LogFileEntry, error) {
	// Prevent directory traversal
	if strings.Contains(serverID, "..") {
		return nil, fmt.Errorf("invalid server ID")
	}

	logsPath := filepath.Join(s.basePath, serverID, "logs")
	entries, err := os.ReadDir(logsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []LogFileEntry{}, nil
		}
		return nil, err
	}

	var logs []LogFileEntry
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		logs = append(logs, LogFileEntry{
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	// Sort by ModTime descending (newest first)
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].ModTime.After(logs[j].ModTime)
	})

	return logs, nil
}

func (s *LogService) ReadLogFile(serverID, filename string) (string, error) {
	if strings.Contains(serverID, "..") || strings.Contains(filename, "..") {
		return "", fmt.Errorf("invalid path")
	}

	filePath := filepath.Join(s.basePath, serverID, "logs", filename)

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var reader io.Reader = file

	// Check if gzip
	if strings.HasSuffix(filename, ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return "", err
		}
		defer gzReader.Close()
		reader = gzReader
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
