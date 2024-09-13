// config/config.go

package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethClient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type ContractDetails struct {
	ABI       abi.ABI
	Addresses map[string]common.Address
}

var (
	contractDetails *ContractDetails
	infuraURL       string
	infuraWSURL     string
)


func Init() error {

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get Infura URLs from environment variables
	infuraURL = os.Getenv("INFURA_URL")
	infuraWSURL = os.Getenv("INFURA_WEBSOCKET_URL")

	if infuraURL == "" || infuraWSURL == "" {
		return fmt.Errorf("INFURA_URL or INFURA_WEBSOCKET_URL not set in .env file")
	}

	var err error
	contractDetails, err = loadContractDetails()
	return err
}

func GetContractDetails() *ContractDetails {
	return contractDetails
}

func loadContractDetails() (*ContractDetails, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}

	// Construct paths to the JSON files
	abiPath := filepath.Join(wd, "..", "contractDetails", "tokenContractABI.json")
	addressPath := filepath.Join(wd, "..", "contractDetails", "tokenContractAddress.json")

	// Read and parse ABI
	abiFile, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %v", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(abiFile)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Read and parse addresses
	addressFile, err := os.ReadFile(addressPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read address file: %v", err)
	}

	var addresses map[string]string
	if err := json.Unmarshal(addressFile, &addresses); err != nil {
		return nil, fmt.Errorf("failed to parse address JSON: %v", err)
	}

	// Convert string addresses to common.Address
	addressMap := make(map[string]common.Address)
	for network, addr := range addresses {
		addressMap[network] = common.HexToAddress(addr)
	}

	return &ContractDetails{
		ABI:       contractAbi,
		Addresses: addressMap,
	}, nil
}

func GetEthereumConnection() (*ethClient.Client, error) {
	return ethClient.Dial(infuraURL)
}

func GetEthereumWebSocketConnection() (*ethClient.Client, error) {
	return ethClient.Dial(infuraWSURL)
}


func ServerAddress() string {
	return ":8080"
}