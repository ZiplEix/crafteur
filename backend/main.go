package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/ZiplEix/crafteur/database"
	"github.com/ZiplEix/crafteur/minecraft"
	"github.com/ZiplEix/crafteur/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed public/*
var embedFrontend embed.FS

func init() {
	database.InitDB()
}

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embedFrontend, "public")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func main() {
	mcManager := minecraft.NewManager()
	serverService := services.NewServerService(mcManager)

	if err := serverService.LoadServersAtStartup(); err != nil {
		log.Fatal("Can't load servers at startup:", err)
	}

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
