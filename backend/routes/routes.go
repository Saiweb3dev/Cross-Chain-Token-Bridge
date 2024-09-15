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
        // Route to handle POST requests for storing new event data
        eventRoutes.POST("/events", controllers.HandleContractEvent)
        
        // Route to retrieve the last event data for a specific caller address
        eventRoutes.GET("/events/:callerAddress/last", controllers.GetLastEventData)
    }

    // Return the configured router
    return router
}
