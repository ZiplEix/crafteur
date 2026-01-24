package process

import (
	"fmt"
	"path/filepath"

	"github.com/ZiplEix/crafteur/internal/database"
	"github.com/ZiplEix/crafteur/internal/filesystem"
)

// Pool is a map of active server instances : Map[ID] -> *Server
var Pool = make(map[int]*Server)

// LoadServersFromDB reads the database and initializes the structures (without starting them)
func LoadServersFromDB() error {
	models, err := database.GetAllServers()
	if err != nil {
		return err
	}

	for _, model := range models {
		// Define the specific directory for this server : data/servers/server_1
		serverDir := filepath.Join("./data/servers", fmt.Sprintf("server_%d", model.ID))
		filesystem.EnsureDir(serverDir)

		// Create the Server object
		srv := NewServer(serverDir, "server.jar")

		srv.ID = model.ID
		srv.Port = model.Port

		// Add to pool
		Pool[model.ID] = srv
		fmt.Printf("Server loaded in memory : %s (Port %d)\n", model.Name, model.Port)
	}
	return nil
}

// GetServer retrieves a server by its ID
func GetServer(id int) *Server {
	return Pool[id]
}
