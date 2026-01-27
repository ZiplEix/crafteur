package controller

import (
	"net/http"

	"github.com/ZiplEix/crafteur/services"
	"github.com/labstack/echo/v4"
)

type WorldController struct {
	worldService *services.WorldService
}

func NewWorldController(worldService *services.WorldService) *WorldController {
	return &WorldController{
		worldService: worldService,
	}
}

func (c *WorldController) ListWorlds(ctx echo.Context) error {
	serverID := ctx.Param("id")
	if serverID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID is required"})
	}

	worlds, err := c.worldService.ListWorlds(serverID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, worlds)
}

type CreateWorldRequest struct {
	Name string `json:"name"`
}

func (c *WorldController) CreateWorld(ctx echo.Context) error {
	serverID := ctx.Param("id")
	if serverID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID is required"})
	}

	var req CreateWorldRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Name == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "World name is required"})
	}

	if err := c.worldService.CreateWorld(serverID, req.Name); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "World created"})
}

func (c *WorldController) ActivateWorld(ctx echo.Context) error {
	serverID := ctx.Param("id")
	worldName := ctx.Param("name")

	if serverID == "" || worldName == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID and World Name are required"})
	}

	if err := c.worldService.ActivateWorld(serverID, worldName); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "World activated"})
}

func (c *WorldController) DeleteWorld(ctx echo.Context) error {
	serverID := ctx.Param("id")
	worldName := ctx.Param("name")

	if serverID == "" || worldName == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID and World Name are required"})
	}

	if err := c.worldService.DeleteWorld(serverID, worldName); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "World deleted"})
}
