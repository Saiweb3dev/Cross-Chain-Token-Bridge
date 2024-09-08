package routes

import (
	"github.com/gin-gonic/gin"
	"backend/controllers"  // Update this import path
)

func SetupRoutes(r *gin.Engine) {
	// Ping routes
	r.GET("/ping", controllers.Ping)
	r.POST("/echo", controllers.Echo)

	// User routes (example)
	userGroup := r.Group("/user")
	{
		userGroup.GET("/:id", controllers.GetUser)
		userGroup.POST("/", controllers.CreateUser)
	}
}