// main.go

package main

import (
	
	"log"
	"math/big"
	"time"

	"backend/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"backend/routes"
	"backend/services"
)

func main() {
	// Initialize config
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	StartContractMonitor()
	// Start the event listener in a separate goroutine
	go func() {
		if err := services.ListenForEvents(); err != nil {
			log.Printf("Error in event listener: %v", err)
		}
	}()

	// Setup and run the HTTP server
	r := routes.SetupRouter()
	r.Run(config.ServerAddress())
}

func StartContractMonitor() {
	go func() {
		for {
			// Connect to the Ethereum network
			client, err := config.GetEthereumConnection()
			if err != nil {
				log.Printf("Failed to connect to the network: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// Get contract details
			contractDetails := config.GetContractDetails()
			if contractDetails == nil {
				log.Println("Contract details not available")
				time.Sleep(5 * time.Second)
				continue
			}

			// Use the Polygon Mumbai testnet address (chain ID 80002)
			contractAddress, ok := contractDetails.Addresses["80002"]
			if !ok {
				log.Println("Polygon Mumbai testnet address not found")
				time.Sleep(5 * time.Second)
				continue
			}

			// Create a new instance of the contract
			contract := bind.NewBoundContract(contractAddress, contractDetails.ABI, client, client, client)

			// Get token name
			var name string
			err = contract.Call(&bind.CallOpts{}, &[]interface{}{&name}, "name")
			if err != nil {
				log.Printf("Failed to get token name: %v", err)
			} else {
				log.Printf("Token Name: %s", name)
			}

			// Get token symbol
			var symbol string
			err = contract.Call(&bind.CallOpts{}, &[]interface{}{&symbol}, "symbol")
			if err != nil {
				log.Printf("Failed to get token symbol: %v", err)
			} else {
				log.Printf("Token Symbol: %s", symbol)
			}

			// Get total supply
			var totalSupply *big.Int
			err = contract.Call(&bind.CallOpts{}, &[]interface{}{&totalSupply}, "totalSupply")
			if err != nil {
				log.Printf("Failed to get total supply: %v", err)
			} else {
				log.Printf("Total Supply: %s", totalSupply.String())
			}

			// Get contract owner
			var owner common.Address
			err = contract.Call(&bind.CallOpts{}, &[]interface{}{&owner}, "owner")
			if err != nil {
				log.Printf("Failed to get contract owner: %v", err)
			} else {
				log.Printf("Contract Owner: %s", owner.Hex())
			}

			// Get available supply
			var availableSupply *big.Int
			err = contract.Call(&bind.CallOpts{}, &[]interface{}{&availableSupply}, "availableSupply")
			if err != nil {
				log.Printf("Failed to get available supply: %v", err)
			} else {
				log.Printf("Available Supply: %s", availableSupply.String())
			}

			// Get paused status
			var paused bool
			err = contract.Call(&bind.CallOpts{}, &[]interface{}{&paused}, "paused")
			if err != nil {
				log.Printf("Failed to get paused status: %v", err)
			} else {
				log.Printf("Contract Paused: %v", paused)
			}

			// Sleep for a while before the next check
			time.Sleep(1 * time.Minute)
		}
	}()
}