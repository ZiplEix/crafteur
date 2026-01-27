package main

import (
	"embed"

	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

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

	// Check for CLI commands
	if len(os.Args) > 1 && os.Args[1] == "create-user" {
		if len(os.Args) != 4 {
			fmt.Println("Usage: ./crafteur create-user <username> <password>")
			os.Exit(1)
		}

		username := os.Args[2]
		password := os.Args[3]

		if err := services.Register(username, password); err != nil {
			fmt.Printf("Error creating user: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("User '%s' created successfully.\n", username)
		os.Exit(0)
	}

	mcManager := minecraft.NewManager()
	versionService := services.NewVersionService()

	// File Service assuming data is in "backend/data/servers" relative to running dir, or just "data/servers"
	fileService := services.NewFileService(mcManager, "data/servers")

	// Fabric Service
	fabricService := services.NewFabricService()
	// Paper Service
	paperService := services.NewPaperService()

	// Server Service now needs FileService, FabricService and PaperService
	serverService := services.NewServerService(mcManager, versionService, fileService, fabricService, paperService)

	if err := serverService.LoadServersAtStartup(); err != nil {
		log.Fatal("Can't load servers at startup:", err)
	}
	playerService := services.NewPlayerService(mcManager, "data")
	logService := services.NewLogService("data/servers")
	backupService := services.NewBackupService("data/servers", "data/backups")
	schedulerService := services.NewSchedulerService(serverService)
	worldService := services.NewWorldService(serverService, "data/servers")
	addonService := services.NewAddonService(serverService, "data/servers")
	modrinthService := services.NewModrinthService(serverService)

	serverCtrl := controller.NewServerController(serverService)
	fileCtrl := controller.NewFileController(fileService)
	playerCtrl := controller.NewPlayerController(playerService, serverService)
	logCtrl := controller.NewLogController(logService)
	backupCtrl := controller.NewBackupController(backupService)
	schedulerCtrl := controller.NewSchedulerController(schedulerService)
	worldCtrl := controller.NewWorldController(worldService)
	addonCtrl := controller.NewAddonController(addonService)
	modrinthCtrl := controller.NewModrinthController(modrinthService, serverService)

	e := echo.New()

	// Load tasks on startup
	if err := schedulerService.LoadTasks(); err != nil {
		e.Logger.Error("Failed to load scheduled tasks:", err)
	}
	schedulerService.Start()
	defer schedulerService.Stop()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	routes.Register(e, serverCtrl, fileCtrl, playerCtrl, logCtrl, backupCtrl, schedulerCtrl, worldCtrl, addonCtrl, modrinthCtrl)

	assetHandler := http.FileServer(getFileSystem())
	e.GET("/*", echo.WrapHandler(assetHandler))

	e.Logger.Fatal(e.Start(":8080"))
}
