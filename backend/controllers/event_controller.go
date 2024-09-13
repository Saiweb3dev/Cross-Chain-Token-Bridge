package controllers

import (
    "github.com/gin-gonic/gin"
)

func HandleContractEvent(c *gin.Context) {
    var eventData map[string]interface{}
    if err := c.BindJSON(&eventData); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Event processed successfully"})
}


