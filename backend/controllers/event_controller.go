package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	totalRequests       int64
	totalProcessingTime time.Duration
)

func HandleMintEvent(c *gin.Context) {
	handleEvent(c, "Mint")
}

func HandleBurnEvent(c *gin.Context) {
	handleEvent(c, "Burn")
}

func HandleTokensReleasedEvent(c *gin.Context) {
	handleEvent(c, "TokensReleased")
}

func HandleTokensLockedEvent(c *gin.Context) {
	handleEvent(c, "TokensLocked")
}

func HandleMessageSentEvent(c *gin.Context) {
	handleEvent(c, "MessageSent")
}

func HandleMessageReceivedEvent(c *gin.Context) {
	handleEvent(c, "MessageReceived")
}

func handleEvent(c *gin.Context, eventName string) {
	start := time.Now()

	var eventData models.EventData
	if err := c.ShouldBindJSON(&eventData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	eventData.EventName = eventName
	log.Printf("Received %s event data: %+v", eventName, eventData)

	db := database.GetDatabase()
	collection := db.Collection("events")

	ensureIndexes(collection)

	now := time.Now()
	eventData.CreatedAt = now.Format(time.RFC3339Nano)
eventData.UpdatedAt = now.Format(time.RFC3339Nano)

	_, err := collection.InsertOne(context.Background(), eventData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store event data", "details": err.Error()})
		return
	}

	storedData, _ := json.Marshal(eventData)
	log.Printf("Stored %s event data: %s", eventName, storedData)

	c.JSON(http.StatusOK, gin.H{
		"message": eventName + " event data received and stored successfully",
		"data":    eventData,
	})

	duration := time.Since(start)
	totalRequests++
	totalProcessingTime += duration

	log.Printf("Request processed in %v", duration)
	log.Printf("Average processing time: %v", totalProcessingTime/time.Duration(totalRequests))
}

func ensureIndexes(collection *mongo.Collection) {
	ctx := context.Background()
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{"event_name", 1}, {"timestamp", -1}}},
		{Keys: bson.D{{"caller_address", 1}, {"event_name", 1}, {"timestamp", -1}}},
		{Keys: bson.D{{"chain_id", 1}, {"event_name", 1}, {"timestamp", -1}}},
		{Keys: bson.D{{"contract_address", 1}, {"event_name", 1}, {"timestamp", -1}}},
		{Keys: bson.D{{"message_id", 1}}, Options: options.Index().SetSparse(true)},
		{Keys: bson.D{{"to_from_user", 1}, {"event_name", 1}, {"timestamp", -1}}, Options: options.Index().SetSparse(true)},
	}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	_, err := collection.Indexes().CreateMany(ctx, indexes, opts)
	if err != nil {
		log.Printf("Error creating indexes: %v", err)
	} else {
		log.Println("Indexes created successfully")
	}
}

func GetLastEventData(c *gin.Context) {
	callerAddress := c.Param("callerAddress")
	eventName := c.Query("eventName")

	if callerAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Caller address is required"})
		return
	}

	log.Printf("Retrieving last event data for caller address: %s, event name: %s", callerAddress, eventName)

	db := database.GetDatabase()
	collection := db.Collection("events")

	var lastEventData models.EventData

	filter := bson.M{"caller_address": callerAddress}
	if eventName != "" {
		filter["event_name"] = eventName
	}

	err := collection.FindOne(
		context.Background(),
		filter,
		options.FindOne().SetSort(bson.M{"timestamp": -1}),
	).Decode(&lastEventData)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "No events found for this caller address and event type"})
			return
		}
		log.Printf("Error retrieving last event data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last event data"})
		return
	}

	lastEventDataJSON, _ := json.Marshal(lastEventData)
	log.Printf("Retrieved last event data: %s", string(lastEventDataJSON))

	c.JSON(http.StatusOK, gin.H{
		"data": lastEventData,
	})
}

func GetPerformanceMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"totalRequests":         totalRequests,
		"averageProcessingTime": totalProcessingTime / time.Duration(totalRequests),
	})
}
