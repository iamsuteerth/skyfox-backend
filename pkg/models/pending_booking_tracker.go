package models

import (
	"time"
)

type PendingBookingTracker struct {
	BookingId      int       `json:"booking_id"`
	ExpirationTime time.Time `json:"expiration_time"`
}
