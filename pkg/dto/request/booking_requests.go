package request

type AdminBookingRequest struct {
	ShowID       int      `json:"show_id" binding:"required,numeric"`
	CustomerName string   `json:"customer_name" binding:"required,customName"`
	PhoneNumber  string   `json:"phone_number" binding:"required,customPhone"`
	SeatNumbers  []string `json:"seat_numbers" binding:"required,min=1,dive,min=2,max=3"`
	AmountPaid   float64  `json:"amount_paid" binding:"gt=0"`
}
