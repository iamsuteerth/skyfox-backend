package models

type AdminBookedCustomer struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Number    string `json:"number"`
	BookingId int    `json:"booking_id"`
}
