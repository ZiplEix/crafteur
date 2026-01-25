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
