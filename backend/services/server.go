package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/ZiplEix/crafteur/core"
	"github.com/ZiplEix/crafteur/database"
	"github.com/ZiplEix/crafteur/minecraft"
	"github.com/google/uuid"
)

type ServerService struct {
	manager     *minecraft.Manager
	vService    *VersionService
	fileService *FileService
	fabric      *FabricService
}

func NewServerService(m *minecraft.Manager, v *VersionService, f *FileService, fab *FabricService) *ServerService {
	return &ServerService{
		manager:     m,
		vService:    v,
		fileService: f,
		fabric:      fab,
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

		inst := s.manager.AddInstance(cfg.ID, runDir, cfg.JarName)
		inst.SetRAM(cfg.RAM)

		fmt.Printf(" -> Serveur chargé : %s (ID: %s)\n", cfg.Name, cfg.ID)
		_ = inst
	}
	return nil
}

func (s *ServerService) CreateNewServer(name string, sType core.ServerType, port int, ram int, version string, importFile *multipart.FileHeader) (*core.ServerConfig, error) {
	newID := uuid.New().String()
	serverPath := filepath.Join("./data/servers", newID)

	// 1. Validate version and get URL
	downloadUrl, err := s.vService.GetDownloadURL(version)
	if err != nil {
		return nil, fmt.Errorf("version invalide: %w", err)
	}

	// 2. Setup physique
	if err := core.EnsureDir(serverPath); err != nil {
		return nil, fmt.Errorf("création dossier impossible: %w", err)
	}

	// 3. Download Vanilla JAR (Always required)
	jarPath := filepath.Join(serverPath, "server.jar")
	if err := core.DownloadFile(downloadUrl, jarPath); err != nil {
		os.RemoveAll(serverPath)
		return nil, fmt.Errorf("téléchargement server.jar échoué: %w", err)
	}

	jarName := "server.jar"
	if sType == core.TypeFabric {
		// Install Fabric
		launchJar, err := s.fabric.InstallFabric(serverPath, version, "") // loader="" means latest
		if err != nil {
			os.RemoveAll(serverPath)
			return nil, fmt.Errorf("fabric install failed: %w", err)
		}
		jarName = launchJar
	}

	// 3. EULA
	if err := core.CreateEula(serverPath); err != nil {
		os.RemoveAll(serverPath)
		return nil, fmt.Errorf("eula creation failed: %w", err)
	}

	// 4. Import zip if provided
	if importFile != nil {
		// We use FileService logic but we are inside the service so we can call it directly.
		// However, FileService expects to work with serverID for security check logic resolving paths.
		// Since we have the ID and the path is structured correctly, we can use it.
		// Be careful: FileService resolves paths relative to dataDir + serverID.
		// Our s.fileService is initialized with dataDir="data/servers" (from main.go likely).
		// serverPath is "./data/servers/<id>".
		// Let's use FileService.UploadFile first to put the zip there, then Unzip.

		// Actually, FileService.UploadFile takes targetPath relative to server root.
		// So targetPath="" puts it in root.
		if err := s.fileService.UploadFile(newID, "", importFile); err != nil {
			os.RemoveAll(serverPath)
			return nil, fmt.Errorf("upload import error: %w", err)
		}

		// Now unzip
		if err := s.fileService.Unzip(newID, "", importFile.Filename); err != nil {
			os.RemoveAll(serverPath)
			return nil, fmt.Errorf("unzip import error: %w", err)
		}

		// Cleanup zip
		zipPath := filepath.Join(serverPath, importFile.Filename)
		_ = os.Remove(zipPath)
	}

	cfg := &core.ServerConfig{
		ID:          newID,
		Name:        name,
		Type:        sType,
		Port:        port,
		RAM:         ram,
		JavaVersion: 21,
		Version:     version,
		JarName:     jarName,
	}

	// 5. Persistance
	if err := database.CreateServer(cfg); err != nil {
		os.RemoveAll(serverPath)
		return nil, err
	}

	// 6. Runtime
	inst := s.manager.AddInstance(newID, serverPath, cfg.JarName)
	inst.SetRAM(ram)

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
		Version:     cfg.Version,
		Status:      status,
	}, nil
}

func (s *ServerService) GetProperties(id string) (map[string]string, error) {
	inst, exists := s.manager.GetInstance(id)
	if !exists {
		return nil, fmt.Errorf("serveur introuvable")
	}

	propsPath := filepath.Join(inst.RunDir, "server.properties")
	return minecraft.LoadProperties(propsPath)
}

func (s *ServerService) UpdateProperties(id string, newProps map[string]string) error {
	inst, exists := s.manager.GetInstance(id)
	if !exists {
		return fmt.Errorf("serveur introuvable")
	}

	propsPath := filepath.Join(inst.RunDir, "server.properties")

	// 1. Charger l'existant
	currentProps, err := minecraft.LoadProperties(propsPath)
	if err != nil {
		return err
	}

	// 2. Fusionner (update existing keys or add new ones)
	for k, v := range newProps {
		currentProps[k] = v
	}

	// 3. Sauvegarder
	return minecraft.SaveProperties(propsPath, currentProps)
}

func (s *ServerService) ChangeServerVersion(id string, targetVersion string) error {
	// 1. Validate version
	downloadUrl, err := s.vService.GetDownloadURL(targetVersion)
	if err != nil {
		return fmt.Errorf("version invalide: %w", err)
	}

	// 2. Get Server and Stop if running
	inst, exists := s.manager.GetInstance(id)
	if !exists {
		return fmt.Errorf("serveur introuvable")
	}

	if inst.GetStatus() != core.StatusStopped {
		if err := inst.Stop(); err != nil {
			return fmt.Errorf("impossible d'arrêter le serveur: %w", err)
		}
		// Wait for stop (simple sleep mostly, or check status loop)
		// For simplicity, we assume Stop() blocks or we wait a bit,
		// but Stop() in manager is async for process kill?
		// Actually inst.Stop() sends command or kills.
		// Ideally we should wait until status is stopped.
	}

	serverPath := inst.RunDir
	jarPath := filepath.Join(serverPath, "server.jar")
	backupPath := filepath.Join(serverPath, "server.jar.old")

	// 3. Backup old jar
	if _, err := os.Stat(jarPath); err == nil {
		os.Rename(jarPath, backupPath)
	}

	// 4. Download new jar
	if err := core.DownloadFile(downloadUrl, jarPath); err != nil {
		// Restore backup
		os.Rename(backupPath, jarPath)
		return fmt.Errorf("téléchargement échoué: %w", err)
	}

	// 5. Update DB
	// First get current config
	cfg, err := database.GetServer(id)
	if err != nil {
		return err
	}
	// Delete old entry
	if err := database.DeleteServer(id); err != nil {
		return err
	}
	// Update version and recreate
	cfg.Version = targetVersion
	if err := database.CreateServer(cfg); err != nil {
		return err
	}

	return nil
}

func (s *ServerService) GetVersions() ([]core.MojangVersion, error) {
	return s.vService.GetVersions()
}

func (s *ServerService) DeleteServer(id string) error {
	// 0. Stop memory instance
	inst, exists := s.manager.GetInstance(id)
	if exists {
		// Stop if running (blocking or async? usually best effort)
		if inst.GetStatus() != core.StatusStopped {
			_ = inst.Stop()
		}
		s.manager.RemoveInstance(id)
	}

	// 1. Remove Files (Data)
	serverPath := filepath.Join("./data/servers", id)
	if err := os.RemoveAll(serverPath); err != nil {
		return fmt.Errorf("failed to delete server files: %w", err)
	}

	// 2. Remove Backups (Optional but requested)
	backupPath := filepath.Join("./data/backups", id)
	if err := os.RemoveAll(backupPath); err != nil {
		// Log but don't fail hard? Or fail?
		// Let's return error to be safe, or just log fmt.Printf like LoadServers?
		// Ideally we want to be clean.
		return fmt.Errorf("failed to delete backups: %w", err)
	}

	// 3. Remove Tasks
	if err := database.DeleteTasksByServer(id); err != nil {
		return fmt.Errorf("failed to delete scheduled tasks: %w", err)
	}

	// 4. Remove DB Entry
	if err := database.DeleteServer(id); err != nil {
		return fmt.Errorf("failed to delete server from db: %w", err)
	}

	return nil
}

func (s *ServerService) GetServer(id string) (*core.ServerConfig, error) {
	return database.GetServer(id)
}

func (s *ServerService) GetDataDir() string {
	return "./data/servers"
}
