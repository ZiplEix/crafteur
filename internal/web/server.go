package web

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/ZiplEix/crafteur/internal/filesystem"
	"github.com/ZiplEix/crafteur/internal/process"
	"github.com/gorilla/websocket"
)

//go:embed templates/*
var templateFS embed.FS

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartWebServer(mcServer *process.Server) {
	// main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFS(templateFS, "templates/index.html")
		tmpl.Execute(w, nil)
	})

	// start and stop server
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

	// send command to server
	http.HandleFunc("/action/command", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			cmd := r.FormValue("command")
			mcServer.WriteCommand(cmd)
		}
	})

	// websocket for logs
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		for msg := range mcServer.Output {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				break
			}
		}
	})

	// settings html block
	http.HandleFunc("/view/settings", func(w http.ResponseWriter, r *http.Request) {
		propsPath := "./data/server.properties"

		props, err := filesystem.LoadProperties(propsPath)
		if err != nil {
			http.Error(w, "Impossible de lire properties: "+err.Error(), 500)
			return
		}

		tmpl, _ := template.ParseFS(templateFS, "templates/settings.html")
		tmpl.Execute(w, props)
	})

	// console html block
	http.HandleFunc("/view/console", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<div class="bg-black border border-slate-700 rounded-lg p-4 h-96 overflow-y-auto mb-4 flex flex-col-reverse" id="console-container">
            <div id="logs"></div>
        </div>
		<form hx-post="/action/command" hx-target="this" hx-swap="none" class="flex gap-2">
            <input type="text" name="command" class="flex-1 bg-slate-800 border border-slate-600 rounded px-4 py-2 text-white" placeholder="Commande..." required>
            <button type="submit" class="bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded">Envoyer</button>
        </form>
		<script>
			// Reconnexion rapide du socket si on revient sur la vue
			if(window.setupSocket) window.setupSocket(); 
		</script>
		`
		w.Write([]byte(html))
	})

	// save settings
	http.HandleFunc("/action/save-settings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
			newProps := make(map[string]string)
			for key, values := range r.Form {
				newProps[key] = values[0]
			}

			err := filesystem.SaveProperties("./data/server.properties", newProps)
			if err != nil {
				w.Write([]byte("<span class='text-red-500'>Erreur sauvegarde!</span>"))
				return
			}

			w.Write([]byte("<span class='text-green-400'>Sauvegardé ! Redémarrez le serveur.</span>"))
		}
	})

	log.Printf("Web server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
