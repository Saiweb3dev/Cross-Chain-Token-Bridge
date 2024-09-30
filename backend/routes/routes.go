package routes

import (
    "backend/controllers"
    "github.com/gin-gonic/gin"
     "github.com/gin-contrib/cors"
)

// SetupRouter sets up the main router for the API
func SetupRouter() *gin.Engine {
    // Create a new default Gin engine
    router := gin.Default()

    // Configure CORS
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"} // Add your frontend URL here
    config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

    // Use CORS middleware
    router.Use(cors.New(config))

    // Group all API routes under "/api"
    apiRoutes := router.Group("/api")
    {
        // Event routes
        apiRoutes.POST("/events/mint", controllers.HandleMintEvent)
        apiRoutes.POST("/events/burn", controllers.HandleBurnEvent)
        apiRoutes.POST("/events/tokens-released", controllers.HandleTokensReleasedEvent)
        apiRoutes.POST("/events/tokens-locked", controllers.HandleTokensLockedEvent)
        apiRoutes.POST("/events/message-sent", controllers.HandleMessageSentEvent)
        apiRoutes.POST("/events/message-received", controllers.HandleMessageReceivedEvent)
        apiRoutes.GET("/events/:callerAddress/last", controllers.GetLastEventData)
        apiRoutes.GET("/metrics", controllers.GetPerformanceMetrics)

        // New contract routes

         apiRoutes.GET("/contract/:chainID/:index", controllers.GetContractData)
    }

    // Return the configured router
    return router
}