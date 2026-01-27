package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ZiplEix/crafteur/core"
)

type VersionService struct {
	cache       []core.MojangVersion
	lastFetch   time.Time
	mu          sync.RWMutex
	ManifestURL string
}

func NewVersionService() *VersionService {
	return &VersionService{
		ManifestURL: "https://launchermeta.mojang.com/mc/game/version_manifest.json",
	}
}

func (s *VersionService) GetVersions() ([]core.MojangVersion, error) {
	s.mu.RLock()
	if time.Since(s.lastFetch) < 1*time.Hour && len(s.cache) > 0 {
		defer s.mu.RUnlock()
		return s.cache, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	resp, err := http.Get(s.ManifestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var manifest core.MojangManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}

	// Filter releases only for now, can be changed later
	var releases []core.MojangVersion
	for _, v := range manifest.Versions {
		if v.Type == "release" {
			releases = append(releases, v)
		}
	}

	s.cache = releases
	s.lastFetch = time.Now()

	return releases, nil
}

func (s *VersionService) GetDownloadURL(versionID string) (string, error) {
	versions, err := s.GetVersions()
	if err != nil {
		return "", err
	}

	var versionURL string
	for _, v := range versions {
		if v.ID == versionID {
			versionURL = v.URL
			break
		}
	}

	if versionURL == "" {
		return "", fmt.Errorf("version %s not found", versionID)
	}

	resp, err := http.Get(versionURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var pkg core.VersionPackage
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return "", err
	}

	if pkg.Downloads.Server.URL == "" {
		return "", fmt.Errorf("server download not available for version %s", versionID)
	}

	return pkg.Downloads.Server.URL, nil
}
