package response

import "github.com/govalues/decimal"

type WalletResponse struct {
	Username  string          `json:"username"`
	Balance   decimal.Decimal `json:"balance"`
	UpdatedAt string          `json:"updated_at"`
}

type WalletTransactionResponse struct {
	ID              int64           `json:"id"`
	Amount          decimal.Decimal `json:"amount"`
	TransactionType string          `json:"transaction_type"` // "ADD" or "DEDUCT" to match enum
	BookingID       *int64          `json:"booking_id,omitempty"`
	TransactionID   string          `json:"transaction_id"`
	Timestamp       string          `json:"timestamp"`
}

type WalletTransactionsResponse struct {
	Transactions []WalletTransactionResponse `json:"transactions"`
}
