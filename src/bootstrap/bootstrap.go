package bootstrap

import (
	"backend/config"
	"backend/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Bootstrap() {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		fmt.Printf("Error load config: %s", err)
		return
	}

	app := gin.Default()
	routes.InitRoutes(app)
	app.Run(":" + config.App.Port)
}
