package response

type PaymentResponse struct {
	Message       string `json:"message"`
	RequestID     string `json:"request_id"`
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
}
