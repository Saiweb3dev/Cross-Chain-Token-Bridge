package controllers

import (
    "backend/services"
    "github.com/gin-gonic/gin"
)

func HandleContractEvent(c *gin.Context) {
    var eventData interface{}
    if err := c.BindJSON(&eventData); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    err := services.ProcessContractEvent(eventData)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Event processed successfully"})
}