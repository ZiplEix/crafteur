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
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Type        ServerType `json:"type"`
	Port        int        `json:"port"`
	RAM         int        `json:"ram"`
	JavaVersion int        `json:"java_version"`
	Version     string     `json:"version"` // Minecraft version (e.g. 1.20.4)
}

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ServerDetailResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Type        ServerType   `json:"type"`
	Port        int          `json:"port"`
	RAM         int          `json:"ram"`
	JavaVersion int          `json:"java_version"`
	Version     string       `json:"version"`
	Status      ServerStatus `json:"status"`
}
