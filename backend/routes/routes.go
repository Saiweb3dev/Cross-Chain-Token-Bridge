package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	eventRoutes := router.Group("/events")
	{
		eventRoutes.POST("/",controllers.HandleContractEvent)
	}
	return router
}