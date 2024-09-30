package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethClient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type ChainConfig struct {
	ChainID               string `json:"chain_id"`
	RPCURLEnvVar          string `json:"rpc_url_env_var"`
	WebsocketURLEnv       string `json:"websocket_url_env_var"`
	TokenContractAddrEnv  string `json:"token_contract_addr_env"`
	VaultContractAddrEnv  string `json:"vault_contract_addr_env"`
	RouterContractAddrEnv string `json:"router_contract_addr_env"`
}

type Config struct {
	Chains           map[string]*ChainConfig `json:"chains"`
	GlobalABIFiles   map[string]string       `json:"global_abi_files"`
}

var globalConfig *Config

func Init() error {
	err := godotenv.Load(".env")
    if err != nil {
        // If not found, try to load from parent directory
        err = godotenv.Load("../.env")
        if err != nil {
            log.Printf("Warning: Error loading .env file: %v", err)
            // Continue execution even if .env file is not found
        }
    }

	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		return fmt.Errorf("CONFIG_FILE_PATH not set in .env file")
	}

	var loadErr error
	globalConfig, loadErr = loadConfig(configFilePath)
	if loadErr != nil {
		return loadErr
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
	var abiFileName string
	switch contractType {
	case "Token":
			abiFileName = "tokenContractABI.json"
	case "Vault":
			abiFileName = "vaultContractABI.json"
	case "Router":
			abiFileName = "messangerContractABI.json"
	default:
			return abi.ABI{}, fmt.Errorf("unknown contract type: %s", contractType)
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

func GetContractAddress(chainID string, contractType string) (string, error) {
	config, err := GetChainConfig(chainID)
	if err != nil {
			return "", err
	}
	
	var envVar string
	switch contractType {
	case "Token":
			envVar = config.TokenContractAddrEnv
	case "Vault":
			envVar = config.VaultContractAddrEnv
	case "Router":
			envVar = config.RouterContractAddrEnv
	default:
			return "", fmt.Errorf("unknown contract type: %s", contractType)
	}
	
	addrStr := os.Getenv(envVar)
	if addrStr == "" {
			return "", fmt.Errorf("Contract address environment variable '%s' not set", envVar)
	}
	return addrStr, nil
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
	// Try multiple possible locations for the ABI file
	possiblePaths := []string{
			filepath.Join(".", "contractDetails", fileName),
			filepath.Join("..", "contractDetails", fileName),
			filepath.Join("..", "..", "contractDetails", fileName),
	}

	var abiFile []byte
	var err error
	for _, path := range possiblePaths {
			abiFile, err = os.ReadFile(path)
			if err == nil {
					break
			}
	}

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
