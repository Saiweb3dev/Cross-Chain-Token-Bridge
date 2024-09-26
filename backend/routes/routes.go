package routes

import (
    "backend/controllers"
    "github.com/gin-gonic/gin"
)

// SetupRouter sets up the main router for the API
func SetupRouter() *gin.Engine {
    // Create a new default Gin engine
    router := gin.Default()

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
        apiRoutes.GET("/contracts/:index", controllers.GetContractData)
    }

    // Return the configured router
    return router
}