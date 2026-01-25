package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func DownloadFile(url string, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func CreateEula(serverDir string) error {
	eulaPath := filepath.Join(serverDir, "eula.txt")
	return os.WriteFile(eulaPath, []byte("eula=true\n"), 0644)
}
