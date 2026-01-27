package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ZiplEix/crafteur/core"
)

type FabricService struct{}

func NewFabricService() *FabricService {
	return &FabricService{}
}

func (s *FabricService) GetLatestInstaller() (string, error) {
	resp, err := http.Get("https://meta.fabricmc.net/v2/versions/installer")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var installers []core.FabricInstaller
	if err := json.NewDecoder(resp.Body).Decode(&installers); err != nil {
		return "", err
	}

	for _, inst := range installers {
		if inst.Stable {
			return inst.Url, nil
		}
	}

	if len(installers) > 0 {
		return installers[0].Url, nil
	}

	return "", fmt.Errorf("no fabric installer found")
}

func (s *FabricService) InstallFabric(serverDir string, mcVersion string, loaderVersion string) (string, error) {
	// 1. Get Installer URL
	installerUrl, err := s.GetLatestInstaller()
	if err != nil {
		return "", fmt.Errorf("failed to get installer url: %w", err)
	}

	// 2. Download Installer
	installerPath := filepath.Join(serverDir, "fabric-installer.jar")
	if err := core.DownloadFile(installerUrl, installerPath); err != nil {
		return "", fmt.Errorf("failed to download installer: %w", err)
	}
	defer os.Remove(installerPath)

	// 3. Prepare Command
	// java -jar fabric-installer.jar server -mcversion <mcVersion> -dir <serverDir> -downloadJar
	// loaderVersion is optional in CLI, if not set it uses stable.
	// If loaderVersion is provided we might need -loader <version>
	// But requirements said: "Si loaderVersion est vide, utilise 'latest'".
	// The CLI args: server -mcversion <mcVersion> [-loader <loaderVersion>] -dir <dir> -downloadJar

	args := []string{"-jar", installerPath, "server", "-mcversion", mcVersion, "-dir", serverDir}
	if loaderVersion != "" && loaderVersion != "latest" {
		args = append(args, "-loader", loaderVersion)
	}

	cmd := exec.Command("java", args...)
	// Capture output for debugging if needed
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("fabric installation failed: %s", string(output))
	}

	// 4. Return launch jar name
	// Fabric usually creates fabric-server-launch.jar
	return "fabric-server-launch.jar", nil
}
