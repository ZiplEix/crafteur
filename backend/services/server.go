package services

import (
	"fmt"
	"path/filepath"

	"github.com/ZiplEix/crafteur/core"
	"github.com/ZiplEix/crafteur/database"
	"github.com/ZiplEix/crafteur/minecraft"
	"github.com/google/uuid"
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
	newID := uuid.New().String()

	cfg := &core.ServerConfig{
		ID:          newID,
		Name:        name,
		Type:        sType,
		Port:        port,
		RAM:         ram,
		JavaVersion: 21,
	}

	if err := database.CreateServer(cfg); err != nil {
		return nil, err
	}

	runDir := filepath.Join("./data/servers", fmt.Sprintf("server_%d", newID))

	s.manager.AddInstance(newID, runDir, "server.jar")

	return cfg, nil
}

func (s *ServerService) StartServer(id string) error {
	inst, exists := s.manager.GetInstance(id)
	if !exists {
		return fmt.Errorf("serveur introuvable (id: %d)", id)
	}
	return inst.Start()
}

func (s *ServerService) StopServer(id string) error {
	inst, exists := s.manager.GetInstance(id)
	if !exists {
		return fmt.Errorf("serveur introuvable")
	}
	return inst.Stop()
}

func (s *ServerService) SendCommand(id string, cmd string) error {
	inst, exists := s.manager.GetInstance(id)
	if !exists {
		return fmt.Errorf("serveur introuvable")
	}
	return inst.SendCommand(cmd)
}

func (s *ServerService) GetConsoleStream(id string) (chan string, error) {
	inst, exists := s.manager.GetInstance(id)
	if !exists {
		return nil, fmt.Errorf("serveur introuvable")
	}
	return inst.Output, nil
}

func (s *ServerService) GetAllServers() ([]core.ServerConfig, error) {
	return database.GetAllServers()
}
