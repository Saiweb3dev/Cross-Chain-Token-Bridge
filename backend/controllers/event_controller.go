package controllers

import (
    "net/http"
    "time"
    "log"
    "github.com/gin-gonic/gin"
)

type EventData struct {
    CallerAddress   string    `json:"callerAddress"`
    Event           string    `json:"event"`
    BlockNumber     uint64    `json:"blockNumber"`
    TransactionHash string    `json:"transactionHash"`
    Timestamp       time.Time `json:"timestamp"`
    // Add other fields as needed
}

var (
    totalRequests      int64
    totalProcessingTime time.Duration
)

func HandleContractEvent(c *gin.Context) {
    start := time.Now()
    var eventData EventData
    if err := c.ShouldBindJSON(&eventData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    log.Printf("Received event data: %+v", eventData)

    time.Sleep(time.Millisecond * 10)

    // For now, just log the received data
    c.JSON(http.StatusOK, gin.H{
        "message": "Event data received successfully",
        "data":    eventData,
    })

    duration := time.Since(start)

    // Update metrics
    totalRequests++
    totalProcessingTime += duration

    log.Printf("Request processed in %v", duration)
    log.Printf("Average processing time: %v", totalProcessingTime/time.Duration(totalRequests))
}

// Add a new endpoint to get performance metrics
func GetPerformanceMetrics(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "totalRequests":        totalRequests,
        "averageProcessingTime": totalProcessingTime / time.Duration(totalRequests),
    })
}