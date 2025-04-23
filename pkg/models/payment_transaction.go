package models

import (
	"time"
)

type PaymentTransaction struct {
	Id            int       `json:"id"`
	BookingId     int       `json:"booking_id"`
	TransactionId string    `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	ProcessedAt   time.Time `json:"processed_at"`
}
