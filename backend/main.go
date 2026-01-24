package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed public/*
var embedFrontend embed.FS

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embedFrontend, "public")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())

	api := e.Group("/api")
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok", "message": "Go Backend is running"})
	})

	assetHandler := http.FileServer(getFileSystem())
	e.GET("/*", echo.WrapHandler(assetHandler))

	e.Logger.Fatal(e.Start(":8080"))
}
