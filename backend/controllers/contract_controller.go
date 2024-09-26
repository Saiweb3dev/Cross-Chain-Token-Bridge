package controllers

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "backend/config"
)

type ContractData struct {
    Name             string      `json:"name"`
    Description      string      `json:"description"`
    Details          string      `json:"details"`
    ContractAddress  string      `json:"contractAddress"`
    ABI              interface{} `json:"abi"`
}

// GetContractData handles GET requests for contract data
func GetContractData(c *gin.Context) {
    index := c.Param("index")
    chainID := c.Param("chainID")

    log.Printf("Received request for contract data. Index: %s, ChainID: %s", index, chainID)

    // Fetch ABI and contract address from config
    abi, err := config.GetABI(index)
    if err != nil {
        log.Printf("Failed to fetch ABI: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch ABI: %v", err)})
        return
    }

    contractAddress, err := config.GetContractAddress(chainID)
    if err != nil {
        log.Printf("Failed to fetch contract address: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch contract address: %v", err)})
        return
    }

    // Prepare the contract data
    var contractData ContractData
    switch index {
    case "Token":
        contractData = ContractData{
            Name:        "Token",
            Description: "A customizable cryptocurrency token.",
            Details:     "Allows users to mint and burn tokens.",
        }
    case "Vault":
        contractData = ContractData{
            Name:        "Vault",
            Description: "A decentralized vault for storing cryptocurrencies.",
            Details:     "Provides secure storage and withdrawal services.",
        }
    case "Router":
        contractData = ContractData{
            Name:        "Router",
            Description: "An automated trading router for swapping cryptocurrencies.",
            Details:     "Facilitates trades between different token pairs.",
        }
    default:
        c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
        return
    }

    // Set the fetched ABI and contract address
    contractData.ContractAddress = contractAddress.Hex()
    contractData.ABI = abi
    log.Printf("Sending contract data response for %s", index)
    c.JSON(http.StatusOK, contractData)
}