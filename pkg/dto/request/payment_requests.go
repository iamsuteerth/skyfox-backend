package request

import "github.com/govalues/decimal"

type PaymentRequest struct {
	CardNumber string          `json:"card_number"`
	CVV        string          `json:"cvv"`
	Expiry     string          `json:"expiry"`
	Name       string          `json:"name"`
	Amount     decimal.Decimal `json:"amount"`
	Timestamp  string          `json:"timestamp"`
}
