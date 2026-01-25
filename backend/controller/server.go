package controller

import (
	"net/http"

	"github.com/ZiplEix/crafteur/core"
	"github.com/ZiplEix/crafteur/minecraft"
	"github.com/ZiplEix/crafteur/services"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ServerController struct {
	service  *services.ServerService
	upgrader websocket.Upgrader
}

func NewServerController(s *services.ServerService) *ServerController {
	return &ServerController{
		service: s,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// GET /api/servers
func (ctrl *ServerController) Index(c echo.Context) error {
	servers, err := ctrl.service.GetAllServers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, servers)
}

// GET /api/servers/:id
func (ctrl *ServerController) GetOne(c echo.Context) error {
	id := c.Param("id")
	detail, err := ctrl.service.GetServerDetail(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Server not found"})
	}
	return c.JSON(http.StatusOK, detail)
}

// POST /api/servers
type CreateServerRequest struct {
	Name string          `json:"name"`
	Type core.ServerType `json:"type"`
	Port int             `json:"port"`
	RAM  int             `json:"ram"`
}

func (ctrl *ServerController) Create(c echo.Context) error {
	var req CreateServerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	if req.Port < 1024 || req.RAM < 512 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid port or insufficient RAM"})
	}

	newServer, err := ctrl.service.CreateNewServer(req.Name, req.Type, req.Port, req.RAM)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, newServer)
}

// POST /api/servers/:id/start
func (ctrl *ServerController) Start(c echo.Context) error {
	id := c.Param("id")
	if err := ctrl.service.StartServer(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "starting"})
}

// POST /api/servers/:id/stop
func (ctrl *ServerController) Stop(c echo.Context) error {
	id := c.Param("id")
	if err := ctrl.service.StopServer(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "stopping"})
}

// POST /api/servers/:id/command
func (ctrl *ServerController) Command(c echo.Context) error {
	id := c.Param("id")
	cmd := c.FormValue("command")

	if err := ctrl.service.SendCommand(id, cmd); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "sent"})
}

// WS /api/servers/:id/ws
func (ctrl *ServerController) Console(c echo.Context) error {
	id := c.Param("id")

	stream, cleanup, err := ctrl.service.SubscribeConsole(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Server not found"})
	}
	defer cleanup()

	ws, err := ctrl.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Send history
	history, _ := ctrl.service.GetServerLogHistory(id)
	for _, line := range history {
		msg := minecraft.WSMessage{Type: "log", Data: line}
		if err := ws.WriteJSON(msg); err != nil {
			break
		}
	}

	for msg := range stream {
		if err := ws.WriteJSON(msg); err != nil {
			break
		}
	}
	return nil
}
