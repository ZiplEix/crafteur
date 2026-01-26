package routes

import (
	"github.com/ZiplEix/crafteur/controller"
	"github.com/ZiplEix/crafteur/services"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, serverCtrl *controller.ServerController, fileCtrl *controller.FileController, playerCtrl *controller.PlayerController, logCtrl *controller.LogController, backupCtrl *controller.BackupController) {
	api := e.Group("/api")

	// Public Routes
	api.POST("/login", controller.Login)

	// Protected Routes
	protected := api.Group("")
	config := echojwt.Config{
		SigningKey:  services.SecretKey,
		TokenLookup: "cookie:auth_token",
	}
	protected.Use(echojwt.WithConfig(config))

	protected.GET("/me", controller.Me)
	protected.POST("/logout", controller.Logout)

	protected.GET("/servers", serverCtrl.Index)
	protected.GET("/servers/:id", serverCtrl.GetOne)
	protected.POST("/servers", serverCtrl.Create)

	protected.POST("/servers/:id/start", serverCtrl.Start)
	protected.POST("/servers/:id/stop", serverCtrl.Stop)
	protected.POST("/servers/:id/command", serverCtrl.Command)

	protected.GET("/servers/:id/ws", serverCtrl.Console)
	protected.GET("/servers/:id/properties", serverCtrl.GetProperties)
	protected.POST("/servers/:id/properties", serverCtrl.UpdateProperties)

	// Player Routes
	protected.GET("/servers/:id/players/cache", playerCtrl.GetCache)
	protected.GET("/servers/:id/players/ops", playerCtrl.GetOps)
	protected.GET("/servers/:id/players/banned", playerCtrl.GetBanned)
	protected.POST("/servers/:id/players/action", playerCtrl.HandleAction)

	// File Routes
	protected.GET("/servers/:id/files", fileCtrl.ListFiles)
	protected.POST("/servers/:id/files/directory", fileCtrl.CreateDirectory)
	protected.DELETE("/servers/:id/files", fileCtrl.DeletePath)
	protected.POST("/servers/:id/files/upload", fileCtrl.UploadFile)
	protected.POST("/servers/:id/files/unzip", fileCtrl.Unzip)

	// Log Routes
	protected.GET("/servers/:id/logs", logCtrl.ListLogs)
	protected.GET("/servers/:id/logs/content", logCtrl.GetLogContent)

	// Backup Routes
	protected.GET("/servers/:id/backups", backupCtrl.ListBackups)
	protected.POST("/servers/:id/backups", backupCtrl.CreateBackup)
	protected.GET("/servers/:id/backups/:filename", backupCtrl.DownloadBackup)
	protected.DELETE("/servers/:id/backups/:filename", backupCtrl.DeleteBackup)
}
