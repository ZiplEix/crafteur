package controller

import (
	"net/http"

	"github.com/ZiplEix/crafteur/services"
	"github.com/labstack/echo/v4"
)

type LogController struct {
	logService *services.LogService
}

func NewLogController(logService *services.LogService) *LogController {
	return &LogController{
		logService: logService,
	}
}

func (c *LogController) ListLogs(ctx echo.Context) error {
	serverID := ctx.Param("id")
	if serverID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID is required"})
	}

	logs, err := c.logService.ListLogFiles(serverID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, logs)
}

func (c *LogController) GetLogContent(ctx echo.Context) error {
	serverID := ctx.Param("id")
	filename := ctx.QueryParam("filename")

	if serverID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID is required"})
	}
	if filename == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Filename is required"})
	}

	content, err := c.logService.ReadLogFile(serverID, filename)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.String(http.StatusOK, content)
}
