package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/ZiplEix/crafteur/core"
)

type PaperService struct{}

func NewPaperService() *PaperService {
	return &PaperService{}
}

// Response structs for Paper MC API
type paperBuildsResponse struct {
	ProjectId   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Builds      []int  `json:"builds"`
}

func (s *PaperService) GetLatestBuild(version string) (int, error) {
	// GET https://api.papermc.io/v2/projects/paper/versions/{version}
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s", version)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("paper api returned status: %d", resp.StatusCode)
	}

	var data paperBuildsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	if len(data.Builds) == 0 {
		return 0, fmt.Errorf("no builds found for version %s", version)
	}

	// Last build is the latest
	return data.Builds[len(data.Builds)-1], nil
}

func (s *PaperService) GetDownloadURL(version string, build int) string {
	// https://api.papermc.io/v2/projects/paper/versions/{version}/builds/{build}/downloads/paper-{version}-{build}.jar
	return fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s/builds/%d/downloads/paper-%s-%d.jar", version, build, version, build)
}

func (s *PaperService) InstallPaper(serverDir string, version string) (string, error) {
	// 1. Get Latest Build
	build, err := s.GetLatestBuild(version)
	if err != nil {
		return "", fmt.Errorf("failed to get latest paper build: %w", err)
	}

	// 2. Download
	url := s.GetDownloadURL(version, build)
	filename := fmt.Sprintf("paper-%s-%d.jar", version, build)
	targetPath := filepath.Join(serverDir, filename)

	if err := core.DownloadFile(url, targetPath); err != nil {
		return "", fmt.Errorf("failed to download paper jar: %w", err)
	}

	// 3. Return filename
	return filename, nil
}
