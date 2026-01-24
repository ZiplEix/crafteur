package services

import (
	"fmt"
	"path/filepath"

	"github.com/ZiplEix/crafteur/core"
	"github.com/ZiplEix/crafteur/database"
	"github.com/ZiplEix/crafteur/minecraft"
)

type ServerService struct {
	manager *minecraft.Manager
}

func NewServerService(m *minecraft.Manager) *ServerService {
	return &ServerService{
		manager: m,
	}
}

func (s *ServerService) LoadServersAtStartup() error {
	configs, err := database.GetAllServers()
	if err != nil {
		return err
	}

	fmt.Printf("Chargement de %d serveurs depuis la BDD...\n", len(configs))

	for _, cfg := range configs {
		runDir := filepath.Join("./data/servers", fmt.Sprintf("server_%d", cfg.ID))

		inst := s.manager.AddInstance(cfg.ID, runDir, "server.jar")

		fmt.Printf(" -> Serveur chargé : %s (ID: %d)\n", cfg.Name, cfg.ID)
		_ = inst
	}
	return nil
}

func (s *ServerService) CreateNewServer(name string, sType core.ServerType, port int, ram int) (*core.ServerConfig, error) {
	cfg := &core.ServerConfig{
		Name:        name,
		Type:        sType,
		Port:        port,
		RAM:         ram,
		JavaVersion: 21,
	}

	id, err := database.CreateServer(cfg)
	if err != nil {
		return nil, err
	}
	cfg.ID = id

	runDir := filepath.Join("./data/servers", fmt.Sprintf("server_%d", id))

	s.manager.AddInstance(id, runDir, "server.jar")

	return cfg, nil
}
