package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sukantamajhi/go_rest_api/config"
	"github.com/sukantamajhi/go_rest_api/database"
	"github.com/sukantamajhi/go_rest_api/routers"
)

func main() {
	args := os.Args

	config.LoadConfig()

	fmt.Println("Args: ", args)

	if config.Env.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	database.Connect_to_db()
	defer database.CloseDB()

	fmt.Println("Starting server on port:", config.Env.Port)

	router := routers.SetupRouter()
	router.Run(":" + strconv.Itoa(config.Env.Port))
}
