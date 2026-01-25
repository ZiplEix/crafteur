package services

import (
	"archive/zip"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ZiplEix/crafteur/minecraft"
)

type FileService struct {
	manager *minecraft.Manager
	dataDir string
}

type FileInfo struct {
	Name    string    `json:"name"`
	IsDir   bool      `json:"is_dir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

func NewFileService(manager *minecraft.Manager, dataDir string) *FileService {
	return &FileService{
		manager: manager,
		dataDir: dataDir,
	}
}

// resolvePath securely resolves a path within a server's directory
func (s *FileService) resolvePath(serverID, requestPath string) (string, error) {
	// Base server directory
	serverDir := filepath.Join(s.dataDir, serverID)

	// Join and clean the requested path
	fullPath := filepath.Join(serverDir, requestPath)
	cleanPath := filepath.Clean(fullPath)

	// Security check: Ensure the resolved path starts with the server directory
	// This prevents directory traversal (../) attacks
	if !strings.HasPrefix(cleanPath, serverDir) {
		return "", fmt.Errorf("accès interdit : chemin invalide")
	}

	return cleanPath, nil
}

func (s *FileService) ListFiles(serverID, path string) ([]FileInfo, error) {
	fullPath, err := s.resolvePath(serverID, path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	files := make([]FileInfo, 0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}
	return files, nil
}

func (s *FileService) CreateDirectory(serverID, path, name string) error {
	fullPath, err := s.resolvePath(serverID, filepath.Join(path, name))
	if err != nil {
		return err
	}
	return os.MkdirAll(fullPath, 0755)
}

func (s *FileService) DeletePath(serverID, path string) error {
	fullPath, err := s.resolvePath(serverID, path)
	if err != nil {
		return err
	}
	return os.RemoveAll(fullPath)
}

func (s *FileService) UploadFile(serverID, targetPath string, fileHeader *multipart.FileHeader) error {
	fullPath, err := s.resolvePath(serverID, filepath.Join(targetPath, fileHeader.Filename))
	if err != nil {
		return err
	}

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

func (s *FileService) Unzip(serverID, path, filename string) error {
	zipPath, err := s.resolvePath(serverID, filepath.Join(path, filename))
	if err != nil {
		return err
	}

	targetDir, err := s.resolvePath(serverID, path)
	if err != nil {
		return err
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// ZipSlip protection
		fpath := filepath.Join(targetDir, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			// Skip potentially dangerous file path
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
