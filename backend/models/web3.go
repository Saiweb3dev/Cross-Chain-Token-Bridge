package models

import "time"

type Web3Data struct {
    ID          uint   `json:"id" gorm:"primarykey"`
    Address     string `json:"address" gorm:"unique;not null"`
    Balance     int64  `json:"balance"`
    Nonce       int64  `json:"nonce"`
    TransactionCount int64 `json:"transaction_count"`
    LastUpdated time.Time `json:"last_updated"`
}
