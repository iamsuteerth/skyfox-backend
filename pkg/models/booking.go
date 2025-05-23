package models

import (
	"time"

	"github.com/govalues/decimal"
)

type Booking struct {
	Id               int             `json:"id"`
	Date             time.Time       `json:"date"`
	ShowId           int             `json:"show_id"`
	CustomerId       *int            `json:"customer_id"`
	CustomerUsername *string         `json:"customer_username"`
	NoOfSeats        int             `json:"no_of_seats"`
	AmountPaid       decimal.Decimal `json:"amount_paid"`
	Status           string          `json:"status"`
	BookingTime      time.Time       `json:"booking_time"`
	PaymentType      string          `json:"payment_type"`
}
