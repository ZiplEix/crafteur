package controller

import (
	"net/http"
	"time"

	"github.com/ZiplEix/crafteur/core"
	"github.com/ZiplEix/crafteur/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SchedulerController struct {
	schedulerService *services.SchedulerService
}

func NewSchedulerController(schedulerService *services.SchedulerService) *SchedulerController {
	return &SchedulerController{
		schedulerService: schedulerService,
	}
}

func (c *SchedulerController) ListTasks(ctx echo.Context) error {
	serverID := ctx.Param("id")
	if serverID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID is required"})
	}

	tasks, err := c.schedulerService.GetTasksByServer(serverID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, tasks)
}

func (c *SchedulerController) CreateTask(ctx echo.Context) error {
	serverID := ctx.Param("id")
	if serverID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID is required"})
	}

	var req struct {
		Name           string `json:"name"`
		Action         string `json:"action"`
		Payload        string `json:"payload"`
		CronExpression string `json:"cron_expression"`
		OneShot        bool   `json:"one_shot"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	task := &core.ScheduledTask{
		ID:             uuid.New().String(),
		ServerID:       serverID,
		Name:           req.Name,
		Action:         req.Action,
		Payload:        req.Payload,
		CronExpression: req.CronExpression,
		OneShot:        req.OneShot,
		LastRun:        time.Time{}, // Zero time
	}

	if err := c.schedulerService.CreateTask(task); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, task)
}

func (c *SchedulerController) DeleteTask(ctx echo.Context) error {
	serverID := ctx.Param("id")
	taskID := ctx.Param("taskId")

	if serverID == "" || taskID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID and Task ID are required"})
	}

	// Verify serverID match if needed, for strictness. For now delete by ID is sufficient but good to check ownership?
	// Simplified: just delete by ID.

	if err := c.schedulerService.DeleteTask(taskID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Task deleted"})
}
