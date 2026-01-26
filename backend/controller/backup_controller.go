package controller

import (
	"fmt"
	"net/http"

	"github.com/ZiplEix/crafteur/services"
	"github.com/labstack/echo/v4"
)

type BackupController struct {
	backupService *services.BackupService
}

func NewBackupController(backupService *services.BackupService) *BackupController {
	return &BackupController{
		backupService: backupService,
	}
}

func (c *BackupController) ListBackups(ctx echo.Context) error {
	serverID := ctx.Param("id")
	if serverID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID is required"})
	}

	backups, err := c.backupService.ListBackups(serverID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, backups)
}

func (c *BackupController) CreateBackup(ctx echo.Context) error {
	serverID := ctx.Param("id")
	if serverID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID is required"})
	}

	// NOTE: In production, this should be a background job. Synchronous for now per requirements.
	err := c.backupService.CreateBackup(serverID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create backup: %v", err)})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Backup created successfully"})
}

func (c *BackupController) DownloadBackup(ctx echo.Context) error {
	serverID := ctx.Param("id")
	filename := ctx.Param("filename")

	if serverID == "" || filename == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID and filename are required"})
	}

	path, err := c.backupService.GetBackupPath(serverID, filename)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return ctx.Attachment(path, filename)
}

func (c *BackupController) DeleteBackup(ctx echo.Context) error {
	serverID := ctx.Param("id")
	filename := ctx.Param("filename")

	if serverID == "" || filename == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID and filename are required"})
	}

	err := c.backupService.DeleteBackup(serverID, filename)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Backup deleted successfully"})
}
