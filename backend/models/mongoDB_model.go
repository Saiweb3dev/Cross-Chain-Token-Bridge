package models

import (
	"time"
)

type Event struct {
	ID              string    `json:"ID" bson:"id"`
	ContractAddress string    `json:"ContractAddress" bson:"contract_address"`
	EventName       string    `json:"EventName" bson:"event_name"`
	CallerAddress   string    `json:"CallerAddress" bson:"caller_address"`
	BlockNumber     int64     `json:"BlockNumber" bson:"block_number"`
	TransactionHash string    `json:"TransactionHash" bson:"transaction_hash"`
	Timestamp       time.Time `json:"Timestamp" bson:"timestamp"`
	CreatedAt       time.Time `json:"CreatedAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"UpdatedAt,omitempty" bson:"updated_at,omitempty"`
}
