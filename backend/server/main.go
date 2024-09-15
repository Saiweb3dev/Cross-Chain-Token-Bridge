package main

import (
     "log"

    "backend/config"
    "backend/routes"
	"backend/services"
    "backend/database"

)

func main() {
    // Initialize config
    if err := config.Init(); err != nil {
        log.Fatalf("Failed to initialize config: %v", err)
    }

    go services.StartContractEventMonitor()

    database.ConnectToMongoDB()

    // Setup and run the HTTP server
    r := routes.SetupRouter()
    r.Run(config.ServerAddress())
}
