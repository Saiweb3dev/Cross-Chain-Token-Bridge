package mainserver


import (
    "log"
    "backend/config"
    "backend/routes"
    "backend/services"
    "backend/database"
    "github.com/zsais/go-gin-prometheus"
)

func RunMainServer() {
    // Initialize config
    if err := config.Init(); err != nil {
        log.Fatalf("Failed to initialize config: %v", err)
    }

    go services.StartContractEventMonitor("80002", "Token")

    database.ConnectToMongoDB()

    // Setup and run the HTTP server
    r := routes.SetupRouter()

    // Add prometheus middleware
    p := ginprometheus.NewPrometheus("gin")
    p.Use(r)
    
    log.Println("Main server is running on", config.ServerAddress())
    if err := r.Run(config.ServerAddress()); err != nil {
        log.Fatalf("Failed to run main server: %v", err)
    }
}
