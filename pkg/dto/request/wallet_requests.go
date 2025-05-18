package request

import "github.com/govalues/decimal"

type AddWalletFundsRequest struct {
	Amount         decimal.Decimal `json:"amount" binding:"required"`
	CardNumber     string          `json:"card_number" binding:"required"`
	CVV            string          `json:"cvv" binding:"required"`
	ExpiryMonth    string          `json:"expiry_month" binding:"required"`
	ExpiryYear     string          `json:"expiry_year" binding:"required"`
	CardholderName string          `json:"cardholder_name" binding:"required"`
}
