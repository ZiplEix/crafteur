package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZiplEix/crafteur/internal/filesystem"
	"github.com/ZiplEix/crafteur/internal/process"
)

func main() {
	// Basic configuration
	dataDir := "./data"
	jarName := "server.jar"
	// Official Minecraft 1.21.11 (Vanilla) URL
	jarUrl := "https://piston-data.mojang.com/v1/objects/64bb6d763bed0a9f1d632ec347938594144943ed/server.jar"

	// 1. Prepare directories
	if err := filesystem.EnsureDir(dataDir); err != nil {
		log.Fatal(err)
	}

	// 2. Download server
	jarPath := filepath.Join(dataDir, jarName)
	if err := filesystem.DownloadFile(jarUrl, jarPath); err != nil {
		log.Fatal("Download error:", err)
	}

	// 3. Accept EULA
	if err := filesystem.CreateEula(dataDir); err != nil {
		log.Fatal("EULA error:", err)
	}

	// 4. Initialize process
	mcServer := process.NewServer(dataDir, jarName)
	
	if err := mcServer.Start(); err != nil {
		log.Fatal("Impossible to start the server:", err)
	}

	// 5. Loop to read my keyboard and send it to the server
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type commands (ex: 'stop', 'list', 'say hello') :")
	
	for scanner.Scan() {
		input := scanner.Text()
		mcServer.WriteCommand(input)
		
		if strings.TrimSpace(input) == "stop" {
			break // On sort de la boucle pour attendre la fin du programme
		}
	}

	// Wait for the server to shut down properly
	mcServer.Wait()
	fmt.Println("Server shut down. Bye !")
}
