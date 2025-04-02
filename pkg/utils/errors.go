package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppError struct {
	HTTPCode int
	Code     string
	Message  string
	Err      error
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

	ctx.JSON(appErr.HTTPCode, ErrorResponse{
		Error:     appErr.Message,
		RequestID: requestID,
	})
}

func GetRequestID(ctx *gin.Context) string {
	requestID := ctx.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}
