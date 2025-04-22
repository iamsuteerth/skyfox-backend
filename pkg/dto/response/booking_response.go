package response

import (
	"time"
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
