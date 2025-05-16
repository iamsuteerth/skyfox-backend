package models

import (
	"time"

	"github.com/govalues/decimal"
)

type PaymentTransaction struct {
	Id            int             `json:"id"`
	BookingId     int             `json:"booking_id"`
	TransactionId string          `json:"transaction_id"`
	PaymentMethod string          `json:"payment_method"`
	Amount        decimal.Decimal `json:"amount"`
	Status        string          `json:"status"`
	ProcessedAt   time.Time       `json:"processed_at"`
}
