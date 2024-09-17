package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethClient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type ChainConfig struct {
	ChainID         string `json:"chain_id"`
	RPCURLEnvVar    string `json:"rpc_url_env_var"`
	WebsocketURLEnv string `json:"websocket_url_env_var"`
	ContractAddrEnv string `json:"contract_addr_env"`
}

type Config struct {
	Chains           map[string]*ChainConfig `json:"chains"`
	GlobalABIFiles   map[string]string       `json:"global_abi_files"`
}

var globalConfig *Config

func Init() error {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		return fmt.Errorf("CONFIG_FILE_PATH not set in .env file")
	}

	var err error
	globalConfig, err = loadConfig(configFilePath)
	if err != nil {
		return err
	}

	return nil
}

func GetChainConfig(chainID string) (*ChainConfig, error) {
	if globalConfig.Chains == nil {
		return nil, fmt.Errorf("no chain configurations loaded")
	}
	chainConfig, exists := globalConfig.Chains[chainID]
	if !exists {
		return nil, fmt.Errorf("configuration for chain %s not found", chainID)
	}
	return chainConfig, nil
}

func GetABI(contractType string) (abi.ABI, error) {
	if globalConfig.GlobalABIFiles == nil {
		return abi.ABI{}, fmt.Errorf("global ABI files not found in configuration")
	}

	abiFileName, ok := globalConfig.GlobalABIFiles[contractType]
	if !ok {
		return abi.ABI{}, fmt.Errorf("ABI file for contract type '%s' not found in global configuration", contractType)
	}

	return loadABI(abiFileName)
}

func GetEthereumConnection(chainID string) (*ethClient.Client, error) {
	config, err := GetChainConfig(chainID)
	if err != nil {
		return nil, err
	}
	rpcURL := os.Getenv(config.RPCURLEnvVar)
	if rpcURL == "" {
		return nil, fmt.Errorf("RPC URL environment variable '%s' not set", config.RPCURLEnvVar)
	}
	return ethClient.Dial(rpcURL)
}

func GetEthereumWebSocketConnection(chainID string) (*ethClient.Client, error) {
	config, err := GetChainConfig(chainID)
	if err != nil {
		return nil, err
	}
	wsURL := os.Getenv(config.WebsocketURLEnv)
	if wsURL == "" {
		return nil, fmt.Errorf("Websocket URL environment variable '%s' not set", config.WebsocketURLEnv)
	}
	return ethClient.Dial(wsURL)
}

func GetContractAddress(chainID string) (common.Address, error) {
	config, err := GetChainConfig(chainID)
	if err != nil {
		return common.Address{}, err
	}
	addrStr := os.Getenv(config.ContractAddrEnv)
	if addrStr == "" {
		return common.Address{}, fmt.Errorf("Contract address environment variable '%s' not set", config.ContractAddrEnv)
	}
	return common.HexToAddress(addrStr), nil
}

func ServerAddress() string {
	return ":8080"
}

func loadConfig(filePath string) (*Config, error) {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %v", err)
	}

	return &config, nil
}

func loadABI(fileName string) (abi.ABI, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to get working directory: %v", err)
	}

	abiPath := filepath.Join(wd, "..", "contractDetails", fileName)

	// Read and parse ABI
	abiFile, err := os.ReadFile(abiPath)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read ABI file: %v", err)
	}

	var abiJSON abi.ABI
	err = json.Unmarshal(abiFile, &abiJSON)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return abiJSON, nil
}
