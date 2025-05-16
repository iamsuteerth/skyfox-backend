package paymentservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/govalues/decimal"
	"github.com/iamsuteerth/skyfox-backend/pkg/config"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type validationErrorResponse struct {
	Errors    []validationError `json:"errors"`
	RequestID string            `json:"request_id"`
	Status    string            `json:"status"`
}

type validationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type PaymentService interface {
	ProcessPayment(ctx context.Context, cardNumber, cvv, expiry, name string, amount decimal.Decimal) (string, error)
}

type paymentService struct {
	config config.PaymentServiceConfig
}

func NewPaymentService(cfg config.PaymentServiceConfig) PaymentService {
	return &paymentService{
		config: cfg,
	}
}

func (s *paymentService) ProcessPayment(ctx context.Context, cardNumber, cvv, expiry, name string, amount decimal.Decimal) (string, error) {
	url := fmt.Sprintf("%s/payment", s.config.BaseURL)

	paymentReq := request.PaymentRequest{
		CardNumber: cardNumber,
		CVV:        cvv,
		Expiry:     expiry,
		Name:       name,
		Amount:     amount,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	payloadBytes, err := json.Marshal(paymentReq)
	if err != nil {
		return "", utils.NewInternalServerError("REQUEST_PREPARATION_FAILED", "Failed to prepare payment request", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", utils.NewInternalServerError("REQUEST_CREATION_FAILED", "Failed to create payment request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.config.APIKey != "" {
		req.Header.Set("x-api-key", s.config.APIKey)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", utils.NewInternalServerError("PAYMENT_SERVICE_ERROR", "Failed to connect to payment service", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", utils.NewInternalServerError("RESPONSE_READ_ERROR", "Failed to read payment response", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var paymentResp response.PaymentResponse
		if err := json.Unmarshal(body, &paymentResp); err != nil {
			return "", utils.NewInternalServerError("JSON_PARSE_ERROR", "Failed to parse payment response", err)
		}

		if paymentResp.Status != "SUCCESS" {
			return "", utils.NewInternalServerError("PAYMENT_FAILED", "Payment was not successful", nil)
		}

		return paymentResp.TransactionID, nil

	case http.StatusUnprocessableEntity:
		var validationResp validationErrorResponse
		if err := json.Unmarshal(body, &validationResp); err != nil {
			return "", utils.NewInternalServerError("JSON_PARSE_ERROR", "Failed to parse validation errors", err)
		}

		errorMsg := "Payment validation failed: "
		for i, valErr := range validationResp.Errors {
			if i > 0 {
				errorMsg += ", "
			}
			errorMsg += fmt.Sprintf("%s (%s)", valErr.Message, valErr.Field)
		}

		return "", utils.NewBadRequestError("PAYMENT_VALIDATION_FAILED", errorMsg, nil)

	case http.StatusForbidden:
		return "", utils.NewInternalServerError("PAYMENT_AUTH_FAILED", "Payment service authentication failed", nil)

	default:
		return "", utils.NewInternalServerError(
			"PAYMENT_SERVICE_ERROR",
			fmt.Sprintf("Payment service returned unexpected status code %d", resp.StatusCode),
			fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body)),
		)
	}
}
