package filesystem

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// DownloadFile downloads a file from a URL to a local path
func DownloadFile(url string, destPath string) error {
	// Check if the file already exists to avoid re-downloading
	if _, err := os.Stat(destPath); err == nil {
		fmt.Println("The file already exists, no download necessary.")
		return nil
	}

	fmt.Printf("Downloading %s...\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// CreateEula writes eula=true to the server directory
func CreateEula(serverDir string) error {
	eulaPath := filepath.Join(serverDir, "eula.txt")
	content := []byte("eula=true\n")
	return os.WriteFile(eulaPath, content, 0644)
}
