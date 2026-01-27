package main

import (
	"embed"
	"flag"
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

	// CLI Flags for user creation
	createUser := flag.String("create-user", "", "Create a new user")
	password := flag.String("password", "", "Password for the new user")
	flag.Parse()

	if *createUser != "" {
		if *password == "" {
			fmt.Println("Error: Password is required when creating a user")
			os.Exit(1)
		}
		err := services.Register(*createUser, *password)
		if err != nil {
			fmt.Printf("Error creating user: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("User '%s' created successfully.\n", *createUser)
		os.Exit(0)
	}

	mcManager := minecraft.NewManager()
	serverService := services.NewServerService(mcManager)

	if err := serverService.LoadServersAtStartup(); err != nil {
		log.Fatal("Can't load servers at startup:", err)
	}

	// File Service assuming data is in "backend/data/servers" relative to running dir, or just "data/servers"
	// Based on instance.go NewInstance, RunDir seems to be passed by manager.
	// The Manager likely knows the root. Let's assume data root is "./data/servers" for now or extraction from manager if possible.
	// Looking at previous ls output: backend/data/servers exists.
	fileService := services.NewFileService(mcManager, "data/servers")
	playerService := services.NewPlayerService(mcManager, "data")
	logService := services.NewLogService("data/servers")
	backupService := services.NewBackupService("data/servers", "data/backups")
	schedulerService := services.NewSchedulerService(serverService)

	serverCtrl := controller.NewServerController(serverService)
	fileCtrl := controller.NewFileController(fileService)
	playerCtrl := controller.NewPlayerController(playerService, serverService)
	logCtrl := controller.NewLogController(logService)
	backupCtrl := controller.NewBackupController(backupService)
	schedulerCtrl := controller.NewSchedulerController(schedulerService)

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

	routes.Register(e, serverCtrl, fileCtrl, playerCtrl, logCtrl, backupCtrl, schedulerCtrl)

	assetHandler := http.FileServer(getFileSystem())
	e.GET("/*", echo.WrapHandler(assetHandler))

	e.Logger.Fatal(e.Start(":8080"))
}
