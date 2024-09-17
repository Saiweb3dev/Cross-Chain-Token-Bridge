package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"backend/config"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EventData represents the structure of event data to be sent to the API
type EventData struct {
	ID               string    `json:"id"`
	CallerAddress    string    `json:"caller_address"`
	EventName        string    `json:"event_name"`
	ContractAddress  string    `json:"contract_address"`
	BlockNumber      uint64    `json:"block_number"`
	TransactionHash  string    `json:"transaction_hash"`
	Timestamp        string     `json:"timestamp"`
	AmountFromEvent  string    `json:"amount_from_event"`
	ToFromEvent      string    `json:"to_from_event"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
}

// StartContractEventMonitor initializes and runs the Ethereum event monitoring service
func StartContractEventMonitor(chainID string, contractType string) {
	go monitorEvents(chainID, contractType)
}

// monitorEvents attempts to connect to the Ethereum node and listen for events
func monitorEvents(chainID string, contractType string) {
	maxRetries := 5
	retryDelay := 5 * time.Second

	for attempt := 0; ; attempt++ {
		// Attempt to connect to Ethereum node
		client, err := config.GetEthereumWebSocketConnection(chainID)
		if err != nil {
			handleConnectionError(err, attempt, maxRetries, retryDelay)
			continue
		}

		// Get contract details and start listening for events
		contractAddress, err := config.GetContractAddress(chainID)
		if err != nil {
			log.Println(err)
			time.Sleep(retryDelay)
			continue
		}

		contractABI, err := config.GetABI(contractType)
		if err != nil {
			log.Printf("Error loading ABI for contract type '%s': %v", contractType, err)
			time.Sleep(retryDelay)
			continue
		}

		err = listenForEvents(client, contractAddress, contractABI)
		if err != nil {
			log.Printf("Error listening for events: %v", err)
			client.Close()
			time.Sleep(retryDelay)
			continue
		}

		break // Exit the loop if successful
	}
}


// handleConnectionError logs the error and exits if max retries are reached
func handleConnectionError(err error, attempt, maxRetries int, retryDelay time.Duration) {
	log.Printf("Attempt %d failed: %v", attempt+1, err)
	if attempt >= maxRetries {
		log.Fatalf("Max retries reached. Unable to establish connection.")
	}
	time.Sleep(retryDelay)
}


// listenForEvents sets up a subscription to filter logs for the contract
func listenForEvents(client *ethclient.Client, contractAddress common.Address, contractABI abi.ABI) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %v", err)
		case vLog := <-logs:
			go processLog(vLog, contractABI)
		}
	}
}

// processLog handles a single log entry according to the contract ABI
func processLog(vLog types.Log, contractABI abi.ABI) {
	event, err := contractABI.EventByID(vLog.Topics[0])
	if err != nil {
		log.Printf("Failed to get event: %v", err)
		return
	}

	callerAddress := getCallerAddress(event, vLog)
	log.Printf("Caller Address: %s", callerAddress.String())

	processedInputs := processEventInputs(event, vLog)
	log.Printf("Processed Event Inputs: %+v", processedInputs)

	eventData := createEventData(vLog, event, callerAddress, processedInputs)
	logEventData(eventData)

	if eventData.EventName != "Transfer" {
		sendEventDataToAPI(eventData)
	} else {
		log.Println("Transfer event detected. Skipping API call.")
	}

	log.Println("--------------------")
}

// getCallerAddress extracts the caller's address from the log
func getCallerAddress(event *abi.Event, vLog types.Log) common.Address {
	if len(event.Inputs) > 0 && event.Inputs[0].Name == "from" {
		fromString := string(vLog.Data)
		if len(fromString) == 42 {
			return common.HexToAddress(fromString)
		}
	}

	topic := vLog.Topics[1]
	if bytes.Equal(topic[:], common.LeftPadBytes([]byte{0x12}, 32)[:]) {
		return common.BytesToAddress(vLog.TxHash.Bytes()[:20])
	}
	return common.BytesToAddress(topic[:])
}

// createEventData creates an EventData struct from log information
func createEventData(vLog types.Log, event *abi.Event, callerAddress common.Address, processedInputs map[string]interface{}) EventData {
	eventData := EventData{
		ID:               fmt.Sprintf("%x", vLog.TxHash),
		CallerAddress:    callerAddress.String(),
		EventName:        event.Name,
		ContractAddress:  vLog.Address.Hex(),
		BlockNumber:      vLog.BlockNumber,
		TransactionHash:  vLog.TxHash.Hex(),
		Timestamp:        formatTimestamp(time.Now().UTC()),
		CreatedAt:        formatTimestamp(time.Now().UTC()),
		UpdatedAt:        formatTimestamp(time.Now().UTC()),
	}

	if amount, ok := processedInputs["amount"].(string); ok {
		eventData.AmountFromEvent = amount
	}
	if to, ok := processedInputs["to"].(string); ok {
		eventData.ToFromEvent = to
	}
	

	return eventData
}

// Helper function to format timestamp
func formatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 MST")
}


// logEventData logs the specific event data
func logEventData(eventData EventData) {
	log.Printf("Event: %s", eventData.EventName)
	log.Printf("Amount: %s", eventData.AmountFromEvent)
	log.Printf("To: %s", eventData.ToFromEvent)
	log.Printf("Sending event data to API: %+v", eventData)
}

// sendEventDataToAPI sends the event data to the specified API endpoint
func sendEventDataToAPI(data EventData) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal event data: %v", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/api/events", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send data to API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned non-OK status: %d", resp.StatusCode)
		return
	}
	log.Printf("API response status: %s", resp.Status)
}

// processEventInputs extracts and processes both indexed and non-indexed event inputs
func processEventInputs(event *abi.Event, vLog types.Log) map[string]interface{} {
	result := make(map[string]interface{})

	// Process indexed inputs
	for i, input := range event.Inputs {
		if input.Indexed && i+1 < len(vLog.Topics) {
			result[input.Name] = processIndexedInput(input, vLog.Topics[i+1])
		}
	}

	// Process non-indexed inputs
	if len(vLog.Data) > 0 {
		nonIndexedData, err := event.Inputs.NonIndexed().Unpack(vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack non-indexed data: %v", err)
			return result
		}

		nonIndexedIndex := 0
		for _, input := range event.Inputs {
			if !input.Indexed && nonIndexedIndex < len(nonIndexedData) {
				result[input.Name] = processNonIndexedInput(input, nonIndexedData[nonIndexedIndex])
				nonIndexedIndex++
			}
		}
	}

	return result
}

// processIndexedInput handles different types of indexed inputs
func processIndexedInput(input abi.Argument, topic common.Hash) interface{} {
	switch input.Type.T {
	case abi.AddressTy:
		return common.HexToAddress(topic.Hex()).Hex()
	case abi.UintTy, abi.IntTy:
		return topic.Big().String()
	default:
		return topic.Hex()
	}
}

// processNonIndexedInput handles different types of non-indexed inputs
func processNonIndexedInput(input abi.Argument, value interface{}) interface{} {
	switch v := value.(type) {
	case common.Address:
		return v.Hex()
	case *big.Int:
		return v.String()
	case bool, string:
		return v
	case []byte:
		return fmt.Sprintf("0x%x", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}