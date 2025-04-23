package response

import (
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
)

type AdminBookingResponse struct {
	BookingID    int       `json:"booking_id"`
	ShowID       int       `json:"show_id"`
	CustomerName string    `json:"customer_name"`
	PhoneNumber  string    `json:"phone_number"`
	SeatNumbers  []string  `json:"seat_numbers"`
	AmountPaid   float64   `json:"amount_paid"`
	PaymentType  string    `json:"payment_type"`
	BookingTime  time.Time `json:"booking_time"`
	Status       string    `json:"status"`
}

type InitializeBookingResponse struct {
	BookingID       int       `json:"booking_id"`
	ShowID          int       `json:"show_id"`
	SeatNumbers     []string  `json:"seat_numbers"`
	AmountDue       float64   `json:"amount_due"`
	ExpirationTime  time.Time `json:"expiration_time"`
	TimeRemainingMs int64     `json:"time_remaining_ms"`
}

type CustomerBookingResponse struct {
	BookingID     int          `json:"booking_id"`
	ShowID        int          `json:"show_id"`
	ShowDate      string       `json:"show_date"`
	ShowTime      string       `json:"show_time"`
	Movie         models.Movie `json:"movie"`
	SeatNumbers   []string     `json:"seat_numbers"`
	AmountPaid    float64      `json:"amount_paid"`
	PaymentType   string       `json:"payment_type"`
	BookingTime   time.Time    `json:"booking_time"`
	Status        string       `json:"status"`
	TransactionID string       `json:"transaction_id,omitempty"`
}
