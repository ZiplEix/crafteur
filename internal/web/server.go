package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/ZiplEix/crafteur/internal/database"
	"github.com/ZiplEix/crafteur/internal/filesystem"
	"github.com/ZiplEix/crafteur/internal/process"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed templates/*.html
var templateFS embed.FS

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartWebServer() {
	e := echo.New()

	// Middleware (Logs, Recover)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configuration des templates
	templates := template.Must(template.ParseFS(templateFS, "templates/*.html"))
	e.Renderer = &TemplateRegistry{
		Templates: templates,
	}

	// --- ROUTES ---

	// 1. Dashboard : Liste de tous les serveurs
	e.GET("/", func(c echo.Context) error {
		servers, err := database.GetAllServers()
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "dashboard.html", servers)
	})

	// 2. Vue Détail d'un serveur (Console)
	e.GET("/server/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		// On vérifie si le serveur est chargé en mémoire
		srv := process.GetServer(id)
		if srv == nil {
			return c.String(http.StatusNotFound, "Serveur introuvable ou non chargé")
		}

		// On passe l'objet serveur à la vue
		return c.Render(http.StatusOK, "manage.html", srv)
	})

	// 3. Actions (HTMX)
	e.POST("/server/:id/start", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		srv := process.GetServer(id)
		if srv != nil {
			go srv.Start()
		}
		return c.NoContent(http.StatusOK)
	})

	e.POST("/server/:id/stop", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		srv := process.GetServer(id)
		if srv != nil {
			srv.WriteCommand("stop")
		}
		return c.NoContent(http.StatusOK)
	})

	e.POST("/server/:id/command", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		cmd := c.FormValue("command")
		srv := process.GetServer(id)
		if srv != nil {
			srv.WriteCommand(cmd)
		}
		return c.NoContent(http.StatusOK)
	})

	// 4. WebSocket (Console Live)
	e.GET("/server/:id/ws", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		srv := process.GetServer(id)
		if srv == nil {
			return echo.ErrNotFound
		}

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		// On stream le channel Output du serveur vers le websocket
		for msg := range srv.Output {
			err := ws.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				break
			}
		}
		return nil
	})

	// 5. Settings (Vue partielle pour HTMX)
	e.GET("/server/:id/settings", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		srv := process.GetServer(id)

		propsPath := fmt.Sprintf("%s/server.properties", srv.ServerDir)
		props, err := filesystem.LoadProperties(propsPath)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// On passe ID et Props à la vue partielle
		data := map[string]any{
			"ID":    id,
			"Props": props,
		}
		return c.Render(http.StatusOK, "settings.html", data)
	})

	e.POST("/server/:id/save-settings", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		srv := process.GetServer(id)

		// Parsing du formulaire
		values, _ := c.FormParams()
		newProps := make(map[string]string)
		for k, v := range values {
			if len(v) > 0 {
				newProps[k] = v[0]
			}
		}

		err := filesystem.SaveProperties(srv.ServerDir+"/server.properties", newProps)
		if err != nil {
			return c.HTML(http.StatusOK, "<span class='text-red-500'>Erreur !</span>")
		}
		return c.HTML(http.StatusOK, "<span class='text-green-400'>Sauvegardé !</span>")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
