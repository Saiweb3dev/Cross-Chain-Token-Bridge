package routes

import (
    "backend/controllers"
    "github.com/gin-gonic/gin"
)

// SetupRouter sets up the main router for the API
func SetupRouter() *gin.Engine {
    // Create a new default Gin engine
    router := gin.Default()

    // Group all event-related routes under "/api"
    eventRoutes := router.Group("/api")
    {
        eventRoutes.POST("/events/mint", controllers.HandleMintEvent)
        eventRoutes.POST("/events/burn", controllers.HandleBurnEvent)
        eventRoutes.POST("/events/tokens-released", controllers.HandleTokensReleasedEvent)
        eventRoutes.POST("/events/tokens-locked", controllers.HandleTokensLockedEvent)
        eventRoutes.POST("/events/message-sent", controllers.HandleMessageSentEvent)
        eventRoutes.POST("/events/message-received", controllers.HandleMessageReceivedEvent)
        eventRoutes.GET("/events/:callerAddress/last", controllers.GetLastEventData)
        eventRoutes.GET("/metrics", controllers.GetPerformanceMetrics)
    }

    // Return the configured router
    return router
}
