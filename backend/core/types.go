package core

type ServerStatus string

const (
	StatusStopped  ServerStatus = "STOPPED"
	StatusStarting ServerStatus = "STARTING"
	StatusRunning  ServerStatus = "RUNNING"
	StatusStopping ServerStatus = "STOPPING"
)

type ServerType string

const (
	TypeVanilla ServerType = "vanilla"
	TypeFabric  ServerType = "fabric"
)

type ServerConfig struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Type        ServerType `json:"type"`
	Port        int        `json:"port"`
	RAM         int        `json:"ram"`
	JavaVersion int        `json:"java_version"`
}
