package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AppError struct {
	HTTPCode         int
	Code             string
	Message          string
	Err              error
	ValidationErrors []ValidationError
}

func (ae AppError) Error() string {
	if ae.Message != "" {
		return ae.Message
	}
	if ae.Err != nil {
		return ae.Err.Error()
	}
	return ""
}

type ErrorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id"`
}

type ValidationErrorResponse struct {
	Errors    []ValidationError `json:"errors"`
	RequestID string            `json:"request_id"`
	Status    string            `json:"status"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
	Status    string      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
}

func NewNotFoundError(code string, message string, err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusNotFound,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

func NewBadRequestError(code string, message string, err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusBadRequest,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

func NewInternalServerError(code string, message string, err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusInternalServerError,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

func NewUnauthorizedError(code string, message string, err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusUnauthorized,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

func HandleErrorResponse(ctx *gin.Context, err error, requestID string) {
	appErr, ok := err.(*AppError)
	if !ok {
		appErr = NewInternalServerError("INTERNAL_ERROR", "An unexpected error occurred", err)
	}

	if appErr.Code == "VALIDATION_ERROR" && len(appErr.ValidationErrors) > 0 {
		ctx.JSON(appErr.HTTPCode, ValidationErrorResponse{
			Errors:    appErr.ValidationErrors,
			RequestID: requestID,
			Status:    "REJECT",
		})
	} else {
		ctx.JSON(appErr.HTTPCode, ErrorResponse{
			Error:     appErr.Message,
			RequestID: requestID,
		})
	}
}

func GetRequestID(ctx *gin.Context) string {
	requestID := ctx.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}

func NewValidationError(validationErrors validator.ValidationErrors) *AppError {
	// Convert validation errors to a slice of ValidationError
	errors := make([]ValidationError, 0)
	for _, err := range validationErrors {
		errors = append(errors, ValidationError{
			Field:   err.Field(),
			Message: getValidationErrorMessage(err.Tag(), err.Param()),
		})
	}

	return &AppError{
		HTTPCode:         http.StatusBadRequest,
		Code:             "VALIDATION_ERROR",
		Message:          "Validation failed",
		Err:              fmt.Errorf("validation errors: %v", validationErrors),
		ValidationErrors: errors,
	}
}

// Helper function to get a human-readable error message for a validation tag
func getValidationErrorMessage(tag string, param string) string {
	messages := map[string]string{
		"required":       "This field is required",
		"email":          "Invalid email format",
		"min":            "Must be at least %s characters",
		"max":            "Must be at most %s characters",
		"customName":     "Name must be 3-70 characters, max 4 words, letters only, no consecutive spaces",
		"customUsername": "Username must be 3-30 characters, lowercase, no spaces, cannot start with a number, no consecutive special characters",
		"customPhone":    "Phone number must be exactly 10 digits",
	}

	if msg, exists := messages[tag]; exists {
		if strings.Contains(msg, "%s") && param != "" {
			return fmt.Sprintf(msg, param)
		}
		return msg
	}

	return "Invalid value"
}
