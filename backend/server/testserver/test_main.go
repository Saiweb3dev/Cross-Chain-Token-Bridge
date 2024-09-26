package testserver

import (
    "log"

    "backend/controllers"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

//go run server/main.go -server=test

type ContractData struct {
    Name             string      `json:"name"`
    Description      string      `json:"description"`
    Details          string      `json:"details"`
    ContractAddress  string      `json:"contractAddress"`
    ABI              interface{} `json:"abi"`
}

// func getContractData(c *gin.Context) {
//     // chainID := c.Param("chainID")
//     index := c.Param("index")

//     // For now, we're ignoring the chainID in this test server
//     // In a real scenario, you might use it to fetch different data based on the chain

//     contractDataMap := map[string]ContractData{
//         "Token": {
//             Name:            "Custom Token",
//             Description:     "Custom Token.",
//             Details:         "Custom Token cryptocurrency.",
//             ContractAddress: "0x1234567890123456789012345678901234567890",
//             ABI:             []string{"function transfer(address to, uint256 amount) public"},
//         },
//         "Vault": {
//             Name:            "Vault",
//             Description:     "Vault",
//             Details:         "Vault is a decentralized Vault holdings.",
//             ContractAddress: "0x0987654321098765432109876543210987654321",
//             ABI:             []string{"function deposit(uint256 amount) public"},
//         },
//         "Router": {
//             Name:            "Router",
//             Description:     "Router",
//             Details:         "Router is a relationship with other cryptocurrencies.",
//             ContractAddress: "0x1122334455667788990011223344556677889900",
//             ABI:             []string{"function swap(address fromToken, address toToken, uint256 amount) public"},
//         },
//     }

//     if data, exists := contractDataMap[index]; exists {
//         c.JSON(http.StatusOK, data)
//     } else {
//         c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
//     }
// }

func RunTestServer() {
    router := gin.Default()
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
        AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
        AllowCredentials: true,
        MaxAge:           12 * 60 * 60,
    }))

    // Updated route to match the new structure
    router.GET("/api/contract/:chainID/:index", controllers.GetContractData)

    log.Println("Test server is running on http://localhost:8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to run test server: %v", err)
    }
}