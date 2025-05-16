package models

import (
	"time"

	"github.com/govalues/decimal"
)

type CustomerWallet struct {
	ID        int64           `json:"id"`
	Username  string          `json:"username"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
