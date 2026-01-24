package routes

import (
	"github.com/ZiplEix/crafteur/controller"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, serverCtrl *controller.ServerController) {
	api := e.Group("/api")

	api.GET("/servers", serverCtrl.Index)
	api.POST("/servers", serverCtrl.Create)

	api.POST("/servers/:id/start", serverCtrl.Start)
	api.POST("/servers/:id/stop", serverCtrl.Stop)
	api.POST("/servers/:id/command", serverCtrl.Command)

	api.GET("/servers/:id/ws", serverCtrl.Console)
}
