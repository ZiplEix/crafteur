package services

import (
	"fmt"
	"os"
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
		runDir := filepath.Join("./data/servers", cfg.ID)

		// Verification dossier existe
		if _, err := os.Stat(runDir); os.IsNotExist(err) {
			fmt.Printf(" -> Serveur %s (ID: %s) introuvable sur le disque (path: %s), ignoré.\n", cfg.Name, cfg.ID, runDir)
			continue
		}

		inst := s.manager.AddInstance(cfg.ID, runDir, "server.jar")

		fmt.Printf(" -> Serveur chargé : %s (ID: %s)\n", cfg.Name, cfg.ID)
		_ = inst
	}
	return nil
}

func (s *ServerService) CreateNewServer(name string, sType core.ServerType, port int, ram int) (*core.ServerConfig, error) {
	newID := uuid.New().String()
	serverPath := filepath.Join("./data/servers", newID)

	// 1. Setup physique
	if err := core.EnsureDir(serverPath); err != nil {
		return nil, fmt.Errorf("création dossier impossible: %w", err)
	}

	// 2. Download (Hardcoded Vanilla 1.21.4)
	jarPath := filepath.Join(serverPath, "server.jar")
	downloadUrl := "https://piston-data.mojang.com/v1/objects/64bb6d763bed0a9f1d632ec347938594144943ed/server.jar"

	if err := core.DownloadFile(downloadUrl, jarPath); err != nil {
		os.RemoveAll(serverPath)
		return nil, fmt.Errorf("téléchargement échoué: %w", err)
	}

	// 3. EULA
	if err := core.CreateEula(serverPath); err != nil {
		os.RemoveAll(serverPath)
		return nil, fmt.Errorf("eula creation failed: %w", err)
	}

	cfg := &core.ServerConfig{
		ID:          newID,
		Name:        name,
		Type:        sType,
		Port:        port,
		RAM:         ram,
		JavaVersion: 21,
	}

	// 4. Persistance
	if err := database.CreateServer(cfg); err != nil {
		os.RemoveAll(serverPath)
		return nil, err
	}

	// 5. Runtime
	s.manager.AddInstance(newID, serverPath, "server.jar")

	return cfg, nil
}

func (s *ServerService) StartServer(id string) error {
	inst, exists := s.manager.GetInstance(id)
	if !exists {
		return fmt.Errorf("serveur introuvable (id: %s)", id)
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

func (s *ServerService) SubscribeConsole(id string) (chan minecraft.WSMessage, func(), error) {
	inst, exists := s.manager.GetInstance(id)
	if !exists {
		return nil, nil, fmt.Errorf("serveur introuvable")
	}
	ch := inst.Subscribe()
	cleanup := func() {
		inst.Unsubscribe(ch)
	}
	return ch, cleanup, nil
}

func (s *ServerService) GetServerLogHistory(id string) ([]string, error) {
	inst, exists := s.manager.GetInstance(id)
	if !exists {
		return nil, fmt.Errorf("serveur introuvable")
	}
	return inst.GetHistory(), nil
}

func (s *ServerService) GetAllServers() ([]core.ServerConfig, error) {
	return database.GetAllServers()
}

func (s *ServerService) GetServerDetail(id string) (*core.ServerDetailResponse, error) {
	// 1. Config from DB
	cfg, err := database.GetServer(id)
	if err != nil {
		return nil, err
	}

	// 2. Runtime status
	status := core.StatusStopped
	inst, exists := s.manager.GetInstance(id)
	if exists {
		status = inst.GetStatus()
	}

	return &core.ServerDetailResponse{
		ID:          cfg.ID,
		Name:        cfg.Name,
		Type:        cfg.Type,
		Port:        cfg.Port,
		RAM:         cfg.RAM,
		JavaVersion: cfg.JavaVersion,
		Status:      status,
	}, nil
}
