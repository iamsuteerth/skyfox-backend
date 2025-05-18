package request

import "github.com/govalues/decimal"

type AdminBookingRequest struct {
	ShowID       int             `json:"show_id" binding:"required,numeric"`
	CustomerName string          `json:"customer_name" binding:"required,customName"`
	PhoneNumber  string          `json:"phone_number" binding:"required,customPhone"`
	SeatNumbers  []string        `json:"seat_numbers" binding:"required,min=1,dive,min=2,max=3"`
	AmountPaid   decimal.Decimal `json:"amount_paid" binding:"required"`
}

type InitializeBookingRequest struct {
	ShowID      int      `json:"show_id" binding:"required,numeric"`
	SeatNumbers []string `json:"seat_numbers" binding:"required,min=1,dive,min=2,max=3"`
}

type ProcessPaymentRequest struct {
	BookingID      int    `json:"booking_id" binding:"required,numeric"`
	PaymentMethod  string `json:"payment_method" binding:"required,oneof=Card Wallet"`
	CardNumber     string `json:"card_number" binding:"required,len=16,numeric"`
	CVV            string `json:"cvv" binding:"required,len=3,numeric"`
	ExpiryMonth    string `json:"expiry_month" binding:"required,min=1,max=2,numeric"`
	ExpiryYear     string `json:"expiry_year" binding:"required,len=2,numeric"`
	CardholderName string `json:"cardholder_name" binding:"required,customName"`
}
