package response

import "time"

type CustomerBookingInfo struct {
	BookingID   int       `json:"booking_id"`
	ShowID      int       `json:"show_id"`
	ShowDate    string    `json:"show_date"`
	ShowTime    string    `json:"show_time"`
	SeatNumbers []string  `json:"seat_numbers"`
	AmountPaid  float64   `json:"amount_paid"`
	PaymentType string    `json:"payment_type"`
	BookingTime time.Time `json:"booking_time"`
	Status      string    `json:"status"`
}
