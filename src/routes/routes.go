package routes

import (
	"backend/controllers"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func InitRoutes(app *gin.Engine) {
	route := app

	// General v1
	ApiV1 := route.Group("/api/v1")
	ApiV1.GET("/", controllers.Home)
	ApiV1.GET("/ping", controllers.Ping)
	ApiV1.POST("/auth", controllers.Auth)

	// Dashboard v1
	ApiV1Dashboard := route.Group("api/v1/dashboard")
	ApiV1Dashboard.Use(utils.AuthenticateJWT())
	{
		ApiV1Dashboard.GET("", controllers.Dashboard)

	}

}
