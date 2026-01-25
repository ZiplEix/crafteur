package controller

import (
	"net/http"

	"github.com/ZiplEix/crafteur/services"
	"github.com/labstack/echo/v4"
)

type FileController struct {
	fileService *services.FileService
}

func NewFileController(fileService *services.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

func (ctrl *FileController) ListFiles(c echo.Context) error {
	serverID := c.Param("id")
	path := c.QueryParam("path")
	if path == "" {
		path = "/"
	}

	files, err := ctrl.fileService.ListFiles(serverID, path)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, files)
}

func (ctrl *FileController) CreateDirectory(c echo.Context) error {
	serverID := c.Param("id")
	var req struct {
		Path string `json:"path"`
		Name string `json:"name"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := ctrl.fileService.CreateDirectory(serverID, req.Path, req.Name); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "directory created"})
}

func (ctrl *FileController) DeletePath(c echo.Context) error {
	serverID := c.Param("id")
	path := c.QueryParam("path")

	if path == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "path is required"})
	}

	if err := ctrl.fileService.DeletePath(serverID, path); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "path deleted"})
}

func (ctrl *FileController) UploadFile(c echo.Context) error {
	serverID := c.Param("id")
	path := c.FormValue("path")
	if path == "" {
		path = "/"
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "file is required"})
	}

	if err := ctrl.fileService.UploadFile(serverID, path, file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "file uploaded"})
}

func (ctrl *FileController) Unzip(c echo.Context) error {
	serverID := c.Param("id")
	var req struct {
		Path     string `json:"path"`
		Filename string `json:"filename"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := ctrl.fileService.Unzip(serverID, req.Path, req.Filename); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "unzipped successfully"})
}
