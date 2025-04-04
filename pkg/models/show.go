package models

import "time"

type Show struct {
	Id      int       `json:"id"`
	MovieId string    `json:"movieId"`
	Date    time.Time `json:"date"` 
	Slot    Slot      `json:"slot"`
	SlotId  int       `json:"slotId"`
	Cost    float64   `json:"cost"`
}
