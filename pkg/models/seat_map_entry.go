package models

type SeatMapEntry struct {
	SeatNumber string  `json:"seat_number"`
	SeatRow    string  `json:"seat_row"`
	SeatColumn string  `json:"seat_column"`
	SeatType   string  `json:"seat_type"`
	Price      float64 `json:"price"`
	Occupied   bool    `json:"occupied"`
}
