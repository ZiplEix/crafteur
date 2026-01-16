package web

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/ZiplEix/crafteur/internal/process"
	"github.com/gorilla/websocket"
)

// On utilise embed pour inclure le HTML dans le binaire final (Single Binary !)
//
//go:embed templates/*
var templateFS embed.FS

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartWebServer(mcServer *process.Server) {
	// 1. Route pour la page d'accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFS(templateFS, "templates/index.html")
		tmpl.Execute(w, nil)
	})

	// 2. Actions HTMX (Start / Stop)
	http.HandleFunc("/action/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			go mcServer.Start()
		}
	})

	http.HandleFunc("/action/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			mcServer.WriteCommand("stop")
		}
	})

	// 3. Commandes envoyées depuis la console Web
	http.HandleFunc("/action/command", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			cmd := r.FormValue("command")
			mcServer.WriteCommand(cmd)
		}
	})

	// 4. WebSocket pour les logs
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		// Boucle infinie : on lit le channel du process et on l'envoie au navigateur
		for msg := range mcServer.Output {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				break
			}
		}
	})

	log.Printf("Web server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
