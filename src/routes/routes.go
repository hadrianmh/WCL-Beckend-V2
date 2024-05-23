package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(app *gin.Engine) {
	route := app

	ApiV1Group := route.Group("/api/v1")

	ApiV1Group.GET("/", controllers.Home)
	ApiV1Group.GET("/ping", controllers.Ping)
	ApiV1Group.POST("/auth", controllers.Auth)
}
