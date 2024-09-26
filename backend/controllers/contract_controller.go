package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
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

    // This is where you'd typically fetch data from a database or other source
    // For this example, we're using a map to simulate stored data
    contractDataMap := map[string]ContractData{
        "Custom_Token": {
            Name:            "Custom Token",
            Description:     "Custom Token.",
            Details:         "Custom Token cryptocurrency.",
            ContractAddress: "0x1234567890123456789012345678901234567890",
            ABI:             []string{"function transfer(address to, uint256 amount) public"}, // Simplified ABI for example
        },
        "Vault": {
            Name:            "Vault",
            Description:     "Vault",
            Details:         "Vault is a decentralized Vault holdings.",
            ContractAddress: "0x0987654321098765432109876543210987654321",
            ABI:             []string{"function deposit(uint256 amount) public"}, // Simplified ABI for example
        },
        "Router": {
            Name:            "Router",
            Description:     "Router",
            Details:         "Router is a relationship with other cryptocurrencies.",
            ContractAddress: "0x1122334455667788990011223344556677889900",
            ABI:             []string{"function swap(address fromToken, address toToken, uint256 amount) public"}, // Simplified ABI for example
        },
    }

    if data, exists := contractDataMap[index]; exists {
        c.JSON(http.StatusOK, data)
    } else {
        c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
    }
}