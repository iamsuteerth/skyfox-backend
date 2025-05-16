package models

import "github.com/govalues/decimal"

type SeatMapEntry struct {
	SeatNumber string          `json:"seat_number"`
	SeatRow    string          `json:"seat_row"`
	SeatColumn string          `json:"seat_column"`
	SeatType   string          `json:"seat_type"`
	Price      decimal.Decimal `json:"price"`
	Occupied   bool            `json:"occupied"`
}
