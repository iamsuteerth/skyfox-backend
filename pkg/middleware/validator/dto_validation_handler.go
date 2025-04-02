package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

func HandleStructValidationError(err error, requestID string) *utils.AppError {
	var structErrors []string
	for _, fieldErr := range err.(validator.ValidationErrors) {
		structErrors = append(structErrors, fieldError{fieldErr}.String())
	}

	if len(structErrors) > 0 {
		return utils.NewBadRequestError("VALIDATION_FAILED", structErrors[0], nil)
	}

	return utils.NewBadRequestError("VALIDATION_FAILED", "Validation failed", err)
}

type fieldError struct {
	err validator.FieldError
}

func (e fieldError) String() string {
	var sb strings.Builder
	sb.WriteString("field '" + e.err.Field() + "'")
	sb.WriteString(", condition: " + validationErrorToText(e.err))
	if e.err.Value() != nil && e.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", provided: %v", e.err.Value()))
	}
	return sb.String()
}

func validationErrorToText(fieldErr validator.FieldError) string {
	switch fieldErr.ActualTag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldErr.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", fieldErr.Field(), fieldErr.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", fieldErr.Field(), fieldErr.Param())
	case "email":
		return "Invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", fieldErr.Field(), fieldErr.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than %s", fieldErr.Field(), fieldErr.Param())
	case "phoneNumber":
		return "Invalid phone number"
	case "maxSeats":
		return "Number of seats exceeds maximum allowed"
	}
	return fmt.Sprintf("%s is not valid", fieldErr.Field())
}
