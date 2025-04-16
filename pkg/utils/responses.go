package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type StandardizedErrorResponse struct {
	Status    string            `json:"status"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	RequestID string            `json:"request_id"`
	Errors    []ValidationError `json:"errors,omitempty"`
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

type AppError struct {
	HTTPCode         int
	Code             string
	Message          string
	Err              error
	ValidationErrors []ValidationError
}

func SendCreatedResponse(ctx *gin.Context, message string, requestID string, data interface{}) {
	sendSuccessResponse(ctx, http.StatusCreated, message, requestID, data)
}

func SendOKResponse(ctx *gin.Context, message string, requestID string, data interface{}) {
	sendSuccessResponse(ctx, http.StatusOK, message, requestID, data)
}

func sendSuccessResponse(ctx *gin.Context, statusCode int, message string, requestID string, data interface{}) {
	ctx.JSON(statusCode, SuccessResponse{
		Message:   message,
		RequestID: requestID,
		Status:    "SUCCESS",
		Data:      data,
	})
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

func NewForbiddenError(code string, message string, err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusForbidden,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

func NewValidationError(validationErrors validator.ValidationErrors) *AppError {
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

func HandleErrorResponse(ctx *gin.Context, err error, requestID string) {
	var response StandardizedErrorResponse
	response.RequestID = requestID
	response.Status = "ERROR"

	appErr, ok := err.(*AppError)
	if !ok {
		appErr = NewInternalServerError("INTERNAL_ERROR", "An unexpected error occurred", err)
	}

	response.Code = appErr.Code
	response.Message = appErr.Message

	if len(appErr.ValidationErrors) > 0 {
		response.Errors = appErr.ValidationErrors
	}

	ctx.JSON(appErr.HTTPCode, response)
}

func GetRequestID(ctx *gin.Context) string {
	requestID := ctx.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}

func getValidationErrorMessage(tag string, param string) string {
	messages := map[string]string{
		"required":       "This field is required",
		"email":          "Invalid email format",
		"min":            "Must be at least %s characters",
		"max":            "Must be at most %s characters",
		"customName":     "Name must be 3-70 characters, max 4 words, letters only, no consecutive spaces",
		"customUsername": "Username must be 3-30 characters, lowercase, no spaces, cannot start with a number, no consecutive special characters",
		"customPhone":    "Phone number must be exactly 10 digits",
		"customPassword": "Password must be at least 8 characters with at least one uppercase letter and one special character",
		"securityAnswer": "Security answer must be at least 3 characters long",
	}

	if msg, exists := messages[tag]; exists {
		if strings.Contains(msg, "%s") && param != "" {
			return fmt.Sprintf(msg, param)
		}
		return msg
	}
	return "Invalid value"
}
