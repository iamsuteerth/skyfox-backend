package models

type BookingSeatMapping struct {
	Id         int    `json:"id"`
	BookingId  int    `json:"booking_id"`
	SeatNumber string `json:"seat_number"`
}
