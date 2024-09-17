package models

import (

)

type Event struct {
	ID               string    `json:"id" bson:"id"`
	ChainID               string    `json:"ChainId" bson:"ChainId"`
	ContractAddress  string    `json:"contract_address" bson:"contract_address"`
	EventName        string    `json:"event_name" bson:"event_name"`
	CallerAddress    string    `json:"caller_address" bson:"caller_address"`
	BlockNumber      uint64    `json:"block_number" bson:"block_number"`
	TransactionHash  string    `json:"transaction_hash" bson:"transaction_hash"`
	Timestamp        string `json:"timestamp" bson:"timestamp"`
	AmountFromEvent  string    `json:"amount_from_event" bson:"amount_from_event"` 
	ToFromEvent      string    `json:"to_from_event" bson:"to_from_event"`         
	CreatedAt        string `json:"created_at" bson:"created_at"`
	UpdatedAt        string `json:"updated_at" bson:"updated_at"`
}
