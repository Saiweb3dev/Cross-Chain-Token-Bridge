package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variables to track performance metrics
var (
	totalRequests      int64
	totalProcessingTime time.Duration
)


// Handles contract event data insertion
func HandleContractEvent(c *gin.Context) {
	// Start timing the operation
	start := time.Now()

	// Bind incoming JSON data to Event struct
	var eventData models.Event
	if err := c.ShouldBindJSON(&eventData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Log received event data
	log.Printf("Received event data: %+v", eventData)

	// Get database connection
	db := database.GetDatabase()

	// Create sanitized collection name
	collectionName := sanitizeCollectionName(eventData.CallerAddress)
	collection := db.Collection(collectionName)

	// Define filter and update operations
	filter := bson.M{"id": eventData.ID}
	update := bson.M{
		"$set": bson.M{
			"id":                 eventData.ID,
			"ChainId":            eventData.ChainID,
			"contract_address": eventData.ContractAddress,
			"event_name":       eventData.EventName,
			"caller_address":   eventData.CallerAddress,
			"block_number":     eventData.BlockNumber,
			"transaction_hash": eventData.TransactionHash,
			"timestamp":        eventData.Timestamp,
			"amount_from_event": eventData.AmountFromEvent,
			"to_from_event":     eventData.ToFromEvent,
			"created_at":       eventData.CreatedAt,
			"updated_at":       eventData.UpdatedAt,
		},
	}

	// Perform update operation
	result, err := collection.UpdateOne(
		context.Background(),
		filter,
		update,
		options.Update().SetUpsert(true),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store event data", "details": err.Error()})
		return
	}

	// Log update result
	log.Printf("Updated %d document(s) in collection %s", result.ModifiedCount, collectionName)

	// Simulate delay (for demonstration purposes only)
	time.Sleep(time.Millisecond * 10)

	// Log stored data
	storedData, _ := json.Marshal(eventData)
	log.Printf("Stored event data in collection %s: %s", collectionName, storedData)

	// Send success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Event data received and stored successfully",
		"data":    eventData,
	})

	// Calculate and log performance metrics
	duration := time.Since(start)
	totalRequests++
	totalProcessingTime += duration

	log.Printf("Request processed in %v", duration)
	log.Printf("Average processing time: %v", totalProcessingTime/time.Duration(totalRequests))
}

// Retrieves the last event data for a given caller address
func GetLastEventData(c *gin.Context) {
	// Extract caller address from URL parameters
	callerAddress := c.Param("callerAddress")
	if callerAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Caller address is required"})
		return
	}

	// Log the retrieval attempt
	log.Printf("Retrieving last event data for caller address: %s", callerAddress)

	// Get database connection
	db := database.GetDatabase()

	// Create sanitized collection name
	collectionName := sanitizeCollectionName(callerAddress)
	collection := db.Collection(collectionName)

	// Initialize variable to hold the last event data
	var lastEventData models.Event

	// Find the most recent event
	err := collection.FindOne(
		context.Background(),
		bson.M{"caller_address": callerAddress},
		options.FindOne().SetSort(bson.M{"timestamp": -1}),
	).Decode(&lastEventData)

	// Handle potential errors
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, gin.H{"error": "No events found for this caller address"})
			return
		}
		log.Printf("Error retrieving last event data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last event data"})
		return
	}

	// Marshal and log the retrieved data
	lastEventDataJSON, _ := json.Marshal(lastEventData)
	log.Printf("Retrieved last event data: %s", string(lastEventDataJSON))

	// Send the retrieved data
	c.JSON(http.StatusOK, gin.H{
		"data": lastEventData,
	})
}

// Sanitizes a collection name by replacing invalid characters and converting to lowercase
func sanitizeCollectionName(name string) string {
	// Remove dots and hyphens, replace with underscores
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")

	// Convert to lowercase
	return strings.ToLower(name)
}

// Returns performance metrics
func GetPerformanceMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"totalRequests":        totalRequests,
		"averageProcessingTime": totalProcessingTime / time.Duration(totalRequests),
	})
}
