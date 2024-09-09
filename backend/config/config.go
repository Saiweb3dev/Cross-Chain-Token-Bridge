package config

import (
    ethClient "github.com/ethereum/go-ethereum/ethclient"
)

func GetEthereumConnection() (*ethClient.Client, error) {
    return ethClient.Dial("https://sepolia.infura.io/v3/your-project-id")
}

func ServerAddress() string {
    return ":8080"
}