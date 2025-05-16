package response

import (
	"time"

	"github.com/govalues/decimal"
)

type CustomerBookingInfo struct {
	BookingID   int             `json:"booking_id"`
	ShowID      int             `json:"show_id"`
	ShowDate    string          `json:"show_date"`
	ShowTime    string          `json:"show_time"`
	SeatNumbers []string        `json:"seat_numbers"`
	AmountPaid  decimal.Decimal `json:"amount_paid"`
	PaymentType string          `json:"payment_type"`
	BookingTime time.Time       `json:"booking_time"`
	Status      string          `json:"status"`
}
