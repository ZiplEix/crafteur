package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/ZiplEix/crafteur/controller"
	"github.com/ZiplEix/crafteur/database"
	"github.com/ZiplEix/crafteur/minecraft"
	"github.com/ZiplEix/crafteur/routes"
	"github.com/ZiplEix/crafteur/services"
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
	database.InitDB()

	mcManager := minecraft.NewManager()
	serverService := services.NewServerService(mcManager)

	if err := serverService.LoadServersAtStartup(); err != nil {
		log.Fatal("Can't load servers at startup:", err)
	}

	serverCtrl := controller.NewServerController(serverService)

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	routes.Register(e, serverCtrl)

	assetHandler := http.FileServer(getFileSystem())
	e.GET("/*", echo.WrapHandler(assetHandler))

	e.Logger.Fatal(e.Start(":8080"))
}
