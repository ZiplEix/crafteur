package controller

import (
	"fmt"
	"net/http"

	"github.com/ZiplEix/crafteur/services"
	"github.com/labstack/echo/v4"
)

type PlayerController struct {
	playerService *services.PlayerService
	serverService *services.ServerService
}

func NewPlayerController(ps *services.PlayerService, ss *services.ServerService) *PlayerController {
	return &PlayerController{
		playerService: ps,
		serverService: ss,
	}
}

func (c *PlayerController) GetCache(ctx echo.Context) error {
	serverID := ctx.Param("id")
	cache, err := c.playerService.GetPlayerCache(serverID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, cache)
}

func (c *PlayerController) GetOps(ctx echo.Context) error {
	serverID := ctx.Param("id")
	ops, err := c.playerService.GetOps(serverID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, ops)
}

func (c *PlayerController) GetBanned(ctx echo.Context) error {
	serverID := ctx.Param("id")
	bans, err := c.playerService.GetBanned(serverID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, bans)
}

type PlayerActionRequest struct {
	PlayerName string `json:"player"`
	Action     string `json:"action"` // op, deop, ban, pardon, kick, whitelist_add, whitelist_remove
	Reason     string `json:"reason,omitempty"`
}

func (c *PlayerController) HandleAction(ctx echo.Context) error {
	serverID := ctx.Param("id")
	var req PlayerActionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.PlayerName == "" || req.Action == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Missing player or action"})
	}

	var cmd string
	switch req.Action {
	case "op":
		cmd = fmt.Sprintf("op %s", req.PlayerName)
	case "deop":
		cmd = fmt.Sprintf("deop %s", req.PlayerName)
	case "ban":
		if req.Reason != "" {
			cmd = fmt.Sprintf("ban %s %s", req.PlayerName, req.Reason)
		} else {
			cmd = fmt.Sprintf("ban %s", req.PlayerName)
		}
	case "pardon":
		cmd = fmt.Sprintf("pardon %s", req.PlayerName)
	case "kick":
		if req.Reason != "" {
			cmd = fmt.Sprintf("kick %s %s", req.PlayerName, req.Reason)
		} else {
			cmd = fmt.Sprintf("kick %s", req.PlayerName)
		}
	case "whitelist_add":
		cmd = fmt.Sprintf("whitelist add %s", req.PlayerName)
	case "whitelist_remove":
		cmd = fmt.Sprintf("whitelist remove %s", req.PlayerName)
	default:
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Unknown action"})
	}

	if err := c.serverService.SendCommand(serverID, cmd); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"status": "ok", "command": cmd})
}
