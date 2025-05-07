package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
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

	// Add performance middleware
	router.Use(gzip.Gzip(gzip.BestCompression)) // Enable gzip compression
	router.Use(cors.New(cors.Config{            // Configure CORS
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Configure server settings
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Env.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
