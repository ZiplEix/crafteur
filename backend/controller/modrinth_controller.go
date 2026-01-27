package controller

import (
	"net/http"
	"strconv"

	"github.com/ZiplEix/crafteur/core"
	"github.com/ZiplEix/crafteur/services"
	"github.com/labstack/echo/v4"
)

type ModrinthController struct {
	service       *services.ModrinthService
	serverService *services.ServerService
}

func NewModrinthController(s *services.ModrinthService, ss *services.ServerService) *ModrinthController {
	return &ModrinthController{
		service:       s,
		serverService: ss,
	}
}

// GET /api/modrinth/search
// Params: q, limit, offset, serverId, type (mod|plugin)
func (ctrl *ModrinthController) Search(c echo.Context) error {
	query := c.QueryParam("q")
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")
	serverId := c.QueryParam("serverId")
	projectType := c.QueryParam("type")

	if projectType == "" {
		projectType = "mod"
	}

	limit := 20
	if l, err := strconv.Atoi(limitStr); err == nil {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(offsetStr); err == nil {
		offset = o
	}

	var facets []string

	// Auto-filters based on server context
	if serverId != "" {
		server, err := ctrl.serverService.GetServer(serverId)
		if err == nil {
			// Version filter
			if server.Version != "" {
				facets = append(facets, "versions:"+server.Version)
			}

			// Loader filter logic moved partially to Service, but we can stick to adding specific categories here if needed.
			// However, for simplified logic, let's let Service handle project_type specific facets.
			// But wait, "categories:fabric" is specific to mods.
			if projectType == "mod" {
				if server.Type == core.TypeFabric {
					facets = append(facets, "categories:fabric")
				} else if server.Type == core.TypeForge {
					facets = append(facets, "categories:forge")
				}
			}
			// Service handles adding "project_type:..." and plugin specific categories.
		}
	}

	resp, err := ctrl.service.SearchProjects(query, limit, offset, facets, projectType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// POST /api/modrinth/install
type InstallModRequest struct {
	ServerID  string `json:"serverId"`
	ProjectID string `json:"projectId"`
	Type      string `json:"type"` // "mod" or "plugin"
}

func (ctrl *ModrinthController) Install(c echo.Context) error {
	var req InstallModRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Type == "" {
		req.Type = "mod"
	}

	if err := ctrl.service.InstallProject(req.ServerID, req.ProjectID, req.Type); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "installed"})
}
