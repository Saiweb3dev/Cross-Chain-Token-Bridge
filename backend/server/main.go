// main.go

package main

import (
	"backend/config"
	"backend/routes"
	"backend/services"
	"log"
)

func main() {
	// Initialize config
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Start the event listener in a separate goroutine
	go func() {
		if err := services.ListenForEvents(); err != nil {
			log.Printf("Error in event listener: %v", err)
		}
	}()

	// Setup and run the HTTP server
	r := routes.SetupRouter()
	r.Run(config.ServerAddress())
}