package models

type TicketData struct {
	BookingID     int      `json:"booking_id"`
	ShowName      string   `json:"show_name"`
	ShowDate      string   `json:"show_date"`
	ShowTime      string   `json:"show_time"`
	CustomerName  string   `json:"customer_name"`
	ContactNumber string   `json:"contact_number"`
	AmountPaid    float64  `json:"amount_paid"`
	NumberOfSeats int      `json:"number_of_seats"`
	SeatNumbers   []string `json:"seat_numbers"`
	Status        string   `json:"status"`
	BookingTime   string   `json:"booking_time"`
	PaymentType   string   `json:"payment_type"`
}
