package request

type PaymentRequest struct {
	CardNumber string  `json:"card_number"`
	CVV        string  `json:"cvv"`
	Expiry     string  `json:"expiry"`
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Timestamp  string  `json:"timestamp"`
}
