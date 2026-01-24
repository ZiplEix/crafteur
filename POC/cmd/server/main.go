package main

import (
	"log"

	"github.com/ZiplEix/crafteur/internal/database"
	"github.com/ZiplEix/crafteur/internal/filesystem"
	"github.com/ZiplEix/crafteur/internal/process"
	"github.com/ZiplEix/crafteur/internal/web"
)

func main() {
	filesystem.EnsureDir("./data")
	filesystem.EnsureDir("./data/servers")

	database.Init()

	if err := process.LoadServersFromDB(); err != nil {
		log.Fatal("Cannot load servers from DB: ", err)
	}

	if len(process.Pool) == 0 {
		log.Println("Aucun serveur trouvé. Création du serveur par défaut...")
		id, _ := database.CreateServer("Survival 1", "vanilla", 25565)

		// On recharge le pool pour le prendre en compte
		process.LoadServersFromDB()

		// On setup les fichiers pour ce nouveau serveur (téléchargement)
		// Note : C'est temporaire, on fera ça proprement dans l'UI après
		srv := process.GetServer(int(id))
		jarUrl := "https://piston-data.mojang.com/v1/objects/64bb6d763bed0a9f1d632ec347938594144943ed/server.jar"
		filesystem.DownloadFile(jarUrl, srv.ServerDir+"/server.jar")
		filesystem.CreateEula(srv.ServerDir)
	}

	web.StartWebServer()
}
