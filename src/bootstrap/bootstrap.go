package bootstrap

import (
	"backend/config"
	"backend/routes"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Bootstrap() {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		fmt.Printf("Error load config: %s", err)
		return
	}

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},                                                 // Allow all origins
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},           // Allow specific methods
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"}, // Allow specific headers
	}))

	routes.InitRoutes(app)
	app.Run(":" + config.App.Port)
}
