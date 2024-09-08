package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// Ping handles GET request to /ping
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// Echo handles POST request to /echo
func Echo(c *gin.Context) {
	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requestBody)
}

// GetUser handles GET request to /user/:id
func GetUser(c *gin.Context) {
	id := c.Param("id")
	// In a real app, you'd fetch the user from a database
	c.JSON(http.StatusOK, gin.H{"id": id, "name": "John Doe"})
}

// CreateUser handles POST request to /user
func CreateUser(c *gin.Context) {
	var newUser struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// In a real app, you'd save the user to a database
	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": newUser})
}   