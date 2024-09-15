package services

import (
    "context"
    "fmt"
    "log"
    "time"
	"bytes"
    "encoding/json"
    "net/http"

    "backend/config"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

type EventData struct {
    CallerAddress   string    `json:"callerAddress"`
    Event           string    `json:"event"`
    BlockNumber     uint64    `json:"blockNumber"`
    TransactionHash string    `json:"transactionHash"`
    Timestamp       time.Time `json:"timestamp"`
}

// StartContractEventMonitor is the main function that sets up and runs the Ethereum event monitoring service.
func StartContractEventMonitor() {
    go func() {
        maxRetries := 5
        retryDelay := 5 * time.Second

        // Loop until we successfully connect to the Ethereum node
        for attempt := 0; ; attempt++ {
            // Attempt to create a new Ethereum client connection
            client, err := config.GetEthereumWebSocketConnection()
            if err != nil {
                log.Printf("Attempt %d failed to connect: %v", attempt+1, err)
                if attempt >= maxRetries {
                    log.Fatalf("Max retries reached. Unable to establish WebSocket connection.")
                }
                time.Sleep(retryDelay)
                continue
            }

            // Verify the connection by checking the network ID
            networkID, err := client.NetworkID(context.Background())
            if err != nil {
                log.Printf("Attempt %d failed to verify connection: %v", attempt+1, err)
                client.Close()
                time.Sleep(retryDelay)
                continue
            }

            log.Printf("Successfully connected to Ethereum node (Network ID: %s)", networkID.String())

            // Get contract details from configuration
            contractDetails := config.GetContractDetails()
            if contractDetails == nil {
                log.Println("Contract details not available")
                time.Sleep(5 * time.Second)
                continue
            }

            // Check for the specific contract address we're interested in
            contractAddress, ok := contractDetails.Addresses["80002"]
            if !ok {
                log.Println("Polygon Mumbai testnet address not found")
                time.Sleep(5 * time.Second)
                continue
            }

            // Start listening for events on the specified contract
            err = listenForEvents(client, contractAddress, contractDetails.ABI)
            if err != nil {
                log.Printf("Error listening for events: %v", err)
                client.Close()
                time.Sleep(retryDelay)
                continue
            }

            break // Exit the loop if successful
        }
    }()
}

// listenForEvents sets up a subscription to filter logs for a specific contract address
func listenForEvents(client *ethclient.Client, contractAddress common.Address, contractABI abi.ABI) error {
    // Create a filter query to listen only to events from the specified contract
    query := ethereum.FilterQuery{
        Addresses: []common.Address{contractAddress},
    }

    // Create a channel to receive logged events
    logs := make(chan types.Log)
    // Subscribe to filtered logs
    sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
    if err != nil {
        return fmt.Errorf("failed to subscribe to logs: %v", err)
    }
    defer sub.Unsubscribe() // Ensure we unsubscribe when done

    // Continuously process incoming logs
    for {
        select {
        case err := <-sub.Err():
            return fmt.Errorf("subscription error: %v", err)
        case vLog := <-logs:
					// Debug: Print the entire log object
            go processLog(vLog, contractABI,client)
        }
    }
}

// processLog takes a single log entry and processes it according to the contract ABI
func processLog(vLog types.Log, contractABI abi.ABI,client *ethclient.Client) {
    // Get the event type from the log topics
    event, err := contractABI.EventByID(vLog.Topics[0])
    if err != nil {
        log.Printf("Failed to get event: %v", err)
        return
    }

		// Fetch the full transaction details
    tx, isPending, err := client.TransactionByHash(context.Background(), vLog.TxHash)
    if err != nil {
        log.Printf("Failed to fetch transaction details: %v", err)
        return
    }

		log.Printf("______________________________")

		log.Printf("Transaction Hash: %s", vLog.TxHash.Hex())
    log.Printf("Is Pending: %t", isPending)
    log.Printf("Nonce: %d", tx.Nonce())
    log.Printf("Gas Price: %s", tx.GasPrice())
    log.Printf("Value: %s", tx.Value())
    log.Printf("Input: %x", tx.Data())

		log.Printf("______________________________")

    // Log basic event information
    log.Printf("Event: %s", event.Name)
    log.Printf("Block Number: %d", vLog.BlockNumber)
    log.Printf("Transaction Hash: %s", vLog.TxHash.Hex())
    log.Printf("Log Index: %d", vLog.Index)
    
    // Log current timestamp
    timestamp := time.Now().UTC()
    log.Printf("Timestamp: %s", timestamp.Format(time.RFC3339))

    // Check both data and topics for the caller's address
    var callerAddress common.Address
    if len(event.Inputs) > 0 && event.Inputs[0].Name == "from" {
        // Try to convert the raw data to a string
        fromString := string(vLog.Data)
        if len(fromString) == 42 { // Length of a hex-encoded address
            callerAddress = common.HexToAddress(fromString)
        }
    }

    // Check topics for the caller's address
    if callerAddress == (common.Address{}) {
        topic := vLog.Topics[1] // Assuming the second topic contains the caller's address
        if bytes.Equal(topic[:], common.LeftPadBytes([]byte{0x12}, 32)[:]) {
            // This is a special topic indicating the transaction was sent from an EOA
            callerAddress = common.BytesToAddress(vLog.TxHash.Bytes()[:20])
        } else {
            // Try to extract the address from the topic
            callerAddress = common.BytesToAddress(topic[:])
        }
    }

    log.Printf("Caller Address: %s", callerAddress.String())



    // Unpack the event data based on the ABI definition
    data, err := event.Inputs.Unpack(vLog.Data)
    if err != nil {
        log.Printf("Failed to unpack event data: %v", err)
        return
    }

    // Log each field of the event
    for i, input := range event.Inputs {
        if i >= len(data) {
            log.Printf("%s: (no value)", input.Name)
        } else {
            log.Printf("%s: %v", input.Name, data[i])
        }
    }
    eventData := EventData{
        CallerAddress:   callerAddress.String(),
        Event:           event.Name,
        BlockNumber:     vLog.BlockNumber,
        TransactionHash: vLog.TxHash.Hex(),
        Timestamp:       time.Now().UTC(),
    }

    // Log the data being sent to the API
    log.Printf("Sending event data to API: %+v", eventData)

    err = sendEventDataToAPI(eventData)
    if err != nil {
        log.Printf("Failed to send event data to API: %v", err)
    } else {
        log.Println("Event data sent to API successfully")
    }

    log.Println("--------------------")
}

func sendEventDataToAPI(data EventData) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal event data: %v", err)
    }

    resp, err := http.Post("http://localhost:8080/api/events", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to send data to API: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
    }
    log.Printf("API response status: %s", resp.Status)

    return nil
}
