package models

import (
	"time"

	"github.com/govalues/decimal"
)

type WalletTransaction struct {
	ID              int64           `json:"id"`
	WalletID        int64           `json:"wallet_id"`
	Username        string          `json:"username"`
	BookingID       *int64          `json:"booking_id,omitempty"`
	TransactionID   string          `json:"transaction_id"`
	Amount          decimal.Decimal `json:"amount"`
	Timestamp       time.Time       `json:"timestamp"`
	TransactionType string          `json:"transaction_type"` // "ADD" or "DEDUCT"
}
