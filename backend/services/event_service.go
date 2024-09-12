// services/event_service.go

package services

import (
	"backend/config"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func ListenForEvents() error {
	contractDetails := config.GetContractDetails()
	if contractDetails == nil {
		return fmt.Errorf("contract details not initialized")
	}

	// Specify the network you're using
	network := "80002"  // This is the chain ID for Polygon Amoy testnet
	contractAddress, ok := contractDetails.Addresses[network]
	if !ok {
		return fmt.Errorf("no contract address found for network: %s", network)
	}

	client, err := config.GetEthereumConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum: %v", err)
	}
	defer client.Close()

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %v", err)
	}
	defer sub.Unsubscribe()

	fmt.Printf("Listening for events on contract: %s\n", contractAddress.Hex())

	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %v", err)
		case vLog := <-logs:
			event, err := contractDetails.ABI.EventByID(vLog.Topics[0])
			if err != nil {
				log.Printf("Error parsing event: %v", err)
				continue
			}
			
			fmt.Printf("Event detected: %s\n", event.Name)
			
			eventData := make(map[string]interface{})
			err = contractDetails.ABI.UnpackIntoMap(eventData, event.Name, vLog.Data)
			if err != nil {
				log.Printf("Error unpacking event data: %v", err)
				continue
			}

			// Add indexed parameters
			for i, arg := range event.Inputs {
				if arg.Indexed {
					eventData[arg.Name] = vLog.Topics[i+1].Hex()
				}
			}

			jsonData, err := json.MarshalIndent(eventData, "", "  ")
			if err != nil {
				log.Printf("Error marshaling event data: %v", err)
				continue
			}

			fmt.Printf("Event data:\n%s\n\n", string(jsonData))
		}
	}
}



func ProcessContractEvent(eventData map[string]interface{}) {
	jsonData, err := json.MarshalIndent(eventData, "", "  ")
	if err != nil {
		log.Printf("Error marshaling event data: %v", err)
		return
	}
	fmt.Printf("Received event data:\n%s\n\n", string(jsonData))
}