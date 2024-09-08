package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"backend/routes"  // Update this import path
)

func main() {
	// Load .env file
	godotenv.Load()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a new Gin router with default middleware
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	r.Run(":" + port)
}