package controller

import (
	"net/http"

	"github.com/ZiplEix/crafteur/services"
	"github.com/labstack/echo/v4"
)

type AddonController struct {
	service *services.AddonService
}

func NewAddonController(s *services.AddonService) *AddonController {
	return &AddonController{
		service: s,
	}
}

// GET /api/servers/:id/addons/:type
func (ctrl *AddonController) Index(c echo.Context) error {
	id := c.Param("id")
	addonType := c.Param("type")

	files, err := ctrl.service.ListAddons(id, addonType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, files)
}

// POST /api/servers/:id/addons/:type
// POST /api/servers/:id/addons/:type
func (ctrl *AddonController) Upload(c echo.Context) error {
	id := c.Param("id")
	addonType := c.Param("type")

	// Parse Multipart Form
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to parse form"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No files found"})
	}

	var errs []string
	for _, file := range files {
		if err := ctrl.service.UploadAddon(id, addonType, file); err != nil {
			errs = append(errs, file.Filename+": "+err.Error())
		}
	}

	if len(errs) > 0 {
		// Partial success or full failure
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "completed_with_errors",
			"errors": errs,
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "uploaded"})
}

// DELETE /api/servers/:id/addons/:type/:filename
func (ctrl *AddonController) Delete(c echo.Context) error {
	id := c.Param("id")
	addonType := c.Param("type")
	filename := c.Param("filename")

	if err := ctrl.service.DeleteAddon(id, addonType, filename); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
