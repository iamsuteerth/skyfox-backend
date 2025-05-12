package models

import (
	"time"
)

type WalletTransaction struct {
	ID              int64     `json:"id"`
	WalletID        int64     `json:"wallet_id"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"` // "ADD", "DEDUCT", "REFUND"
	BookingID       *int64    `json:"booking_id,omitempty"`
	TransactionID   string    `json:"transaction_id"`
	CreatedAt       time.Time `json:"created_at"`
}
