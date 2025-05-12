package models

import (
	"time"
)

type CustomerWallet struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
