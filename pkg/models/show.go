package models

import (
	"time"

	"github.com/govalues/decimal"
)

type Show struct {
	Id      int             `json:"id"`
	MovieId string          `json:"movieId"`
	Date    time.Time       `json:"date"`
	Slot    Slot            `json:"slot"`
	SlotId  int             `json:"slotId"`
	Cost    decimal.Decimal `json:"cost"`
}
