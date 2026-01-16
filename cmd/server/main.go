package main

import (
	"log"
	"path/filepath"

	"github.com/ZiplEix/crafteur/internal/filesystem"
	"github.com/ZiplEix/crafteur/internal/process"
	"github.com/ZiplEix/crafteur/internal/web"
)

func main() {
	dataDir := "./data"
	jarName := "server.jar"
	// URL officielle Minecraft 1.21.11
	jarUrl := "https://piston-data.mojang.com/v1/objects/64bb6d763bed0a9f1d632ec347938594144943ed/server.jar"

	// 1. Initialisation fichiers
	if err := filesystem.EnsureDir(dataDir); err != nil {
		log.Fatal(err)
	}

	jarPath := filepath.Join(dataDir, jarName)
	if err := filesystem.DownloadFile(jarUrl, jarPath); err != nil {
		log.Println("Vérification du fichier serveur...")
	}
	filesystem.CreateEula(dataDir)

	// 2. Préparation du serveur (mais on ne le lance pas tout de suite !)
	mcServer := process.NewServer(dataDir, jarName)

	// 3. Lancement du serveur Web (qui contrôlera mcServer)
	web.StartWebServer(mcServer)
}
