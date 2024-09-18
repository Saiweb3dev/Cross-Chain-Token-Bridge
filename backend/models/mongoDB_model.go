package models

import (

)

type EventData struct {
	ID               string    `json:"id" bson:"id"`
	ChainID               string    `json:"ChainId" bson:"ChainId"`
	ContractAddress  string    `json:"contract_address" bson:"contract_address"`
	EventName        string    `json:"event_name" bson:"event_name"`
	CallerAddress    string    `json:"caller_address" bson:"caller_address"`
	BlockNumber      uint64    `json:"block_number" bson:"block_number"`
	TransactionHash  string    `json:"transaction_hash" bson:"transaction_hash"`
	Timestamp        string `json:"timestamp" bson:"timestamp"`         
	CreatedAt        string `json:"created_at" bson:"created_at"`
	UpdatedAt        string `json:"updated_at" bson:"updated_at"`
	// Fields for Mint, Burn, TokensReleased, TokensLocked events
	ToFromUser               string    `bson:"to_from_user,omitempty" json:"to_from_user,omitempty"`
	Amount                   string    `bson:"amount,omitempty" json:"amount,omitempty"`
	// Additional fields for MessageSent event
	MessageID                string    `bson:"message_id,omitempty" json:"message_id,omitempty"`
	DestinationChainSelector uint64    `bson:"destination_chain_selector,omitempty" json:"destination_chain_selector,omitempty"`
	Receiver                 string    `bson:"receiver,omitempty" json:"receiver,omitempty"`
	Text                     string    `bson:"text,omitempty" json:"text,omitempty"`
	Client                   string    `bson:"client,omitempty" json:"client,omitempty"`
	FeeToken                 string    `bson:"fee_token,omitempty" json:"fee_token,omitempty"`
	Fees                     string    `bson:"fees,omitempty" json:"fees,omitempty"`
	 // Additional fields for MessageReceived event
	 SourceChainSelector      uint64    `bson:"source_chain_selector,omitempty" json:"source_chain_selector,omitempty"`
	 Sender                   string    `bson:"sender,omitempty" json:"sender,omitempty"`
}
