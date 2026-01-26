package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ZiplEix/crafteur/minecraft"
)

type PlayerCacheItem struct {
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
	ExpiresOn string `json:"expiresOn"`
	Online    bool   `json:"online"`
}

type PlayerService struct {
	Manager *minecraft.Manager
	DataDir string
}

func NewPlayerService(manager *minecraft.Manager, dataDir string) *PlayerService {
	return &PlayerService{
		Manager: manager,
		DataDir: dataDir,
	}
}

func (s *PlayerService) readFile(serverID, filename string, v interface{}) error {
	path := filepath.Join(s.DataDir, "servers", serverID, filename)

	// If file does not exist, return empty (not error)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filename, err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	return nil
}

func (s *PlayerService) GetPlayerCache(serverID string) ([]PlayerCacheItem, error) {
	var originalCache []minecraft.PlayerCacheEntry
	// Initialize as empty slice so it marshals to [] instead of null
	originalCache = []minecraft.PlayerCacheEntry{}

	if err := s.readFile(serverID, "usercache.json", &originalCache); err != nil {
		return nil, err
	}

	// Enrich with Online status
	enrichedCache := make([]PlayerCacheItem, 0, len(originalCache))

	inst, exists := s.Manager.GetInstance(serverID)

	for _, entry := range originalCache {
		isOnline := false
		if exists {
			isOnline = inst.IsPlayerOnline(entry.Name)
		}

		enrichedCache = append(enrichedCache, PlayerCacheItem{
			Name:      entry.Name,
			UUID:      entry.UUID,
			ExpiresOn: entry.ExpiresOn,
			Online:    isOnline,
		})
	}

	return enrichedCache, nil
}

func (s *PlayerService) GetOps(serverID string) ([]minecraft.OpEntry, error) {
	var ops []minecraft.OpEntry
	ops = []minecraft.OpEntry{}

	if err := s.readFile(serverID, "ops.json", &ops); err != nil {
		return nil, err
	}
	return ops, nil
}

func (s *PlayerService) GetBanned(serverID string) ([]minecraft.BanEntry, error) {
	var banned []minecraft.BanEntry
	banned = []minecraft.BanEntry{}

	if err := s.readFile(serverID, "banned-players.json", &banned); err != nil {
		return nil, err
	}
	return banned, nil
}
