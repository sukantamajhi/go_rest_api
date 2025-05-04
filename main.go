package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sukantamajhi/go_rest_api/config"
	"github.com/sukantamajhi/go_rest_api/database"
	"github.com/sukantamajhi/go_rest_api/routers"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Load configuration
	config.LoadConfig()

	if config.AppConfig.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	database.Connect_to_db()
	defer database.CloseDB()

	// Setup router
	router := routers.SetupRouter()

	// Start server
	log.Printf("Server starting on port %d", config.AppConfig.Port)
	router.Run(":" + strconv.Itoa(config.AppConfig.Port))
}
