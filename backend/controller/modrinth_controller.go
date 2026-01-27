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
// Params: q, limit, offset, serverId
func (ctrl *ModrinthController) Search(c echo.Context) error {
	query := c.QueryParam("q")
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")
	serverId := c.QueryParam("serverId")

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
			// Loader filter
			if server.Type == core.TypeFabric {
				facets = append(facets, "categories:fabric")
			} else if server.Type == core.TypeForge {
				facets = append(facets, "categories:forge")
			}
			// Add "project_type:mod" to ensure we get mods
			facets = append(facets, "project_type:mod")
		}
	}

	resp, err := ctrl.service.SearchMods(query, limit, offset, facets)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// POST /api/modrinth/install
type InstallModRequest struct {
	ServerID  string `json:"serverId"`
	ProjectID string `json:"projectId"`
}

func (ctrl *ModrinthController) Install(c echo.Context) error {
	var req InstallModRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := ctrl.service.InstallMod(req.ServerID, req.ProjectID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "installed"})
}
