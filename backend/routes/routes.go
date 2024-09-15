package routes

import (
    "backend/controllers"
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()
    eventRoutes := router.Group("/api")  // Changed to /api for consistency
    {
        eventRoutes.POST("/events", controllers.HandleContractEvent)
    }
    return router
}