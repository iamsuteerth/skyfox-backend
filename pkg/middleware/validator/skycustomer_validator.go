package validator

import (
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("customName", ValidateName)
	v.RegisterValidation("customUsername", ValidateUsername)
	v.RegisterValidation("customPhone", ValidatePhoneNumber)
	v.RegisterValidation("customPassword", ValidatePassword)
	v.RegisterValidation("securityAnswer", validateSecurityAnswer)
}

func ValidateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	name = strings.TrimSpace(name)

	if len(name) < 3 || len(name) > 70 {
		return false
	}

	words := strings.Fields(name)
	if len(words) > 4 {
		return false
	}

	for _, c := range name {
		if c != ' ' && !unicode.IsLetter(c) {
			return false
		}
	}

	return !strings.Contains(name, "  ")
}

func ValidateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	if len(username) < 3 || len(username) > 30 {
		return false
	}

	if len(username) > 0 && unicode.IsDigit(rune(username[0])) {
		return false
	}

	if strings.Contains(username, " ") {
		return false
	}

	if strings.ToLower(username) != username {
		return false
	}

	specialChars := "!@#$%^&*()-_+={}[]|\\:;\"'<>,.?/"
	for i := 0; i < len(username)-1; i++ {
		if strings.ContainsRune(specialChars, rune(username[i])) &&
			strings.ContainsRune(specialChars, rune(username[i+1])) {
			return false
		}
	}

	return true
}

func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	number := fl.Field().String()

	if len(number) != 10 {
		return false
	}

	for _, c := range number {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper := false
	for _, c := range password {
		if unicode.IsUpper(c) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return false
	}

	hasLower := false
	for _, c := range password {
		if unicode.IsLower(c) {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return false
	}

	specialChars := "!@#$%^&*()-_+={}[]|\\:;\"'<>,.?/"
	hasSpecial := false
	for _, c := range password {
		if strings.ContainsRune(specialChars, c) {
			hasSpecial = true
			break
		}
	}

	return hasSpecial
}

func validateSecurityAnswer(fl validator.FieldLevel) bool {
	answer := fl.Field().String()
	return len(strings.TrimSpace(answer)) >= 3
}

var ValidationErrorMessages = map[string]string{
	"required":       "This field is required",
	"email":          "Invalid email format",
	"customName":     "Name must be 3-70 characters, max 4 words, letters only, no consecutive spaces",
	"customUsername": "Username must be 3-30 characters, lowercase, no spaces, cannot start with a number, no consecutive special characters",
	"customPhone":    "Phone number must be exactly 10 digits",
	"customPassword": "Password must be at least 8 characters with at least one uppercase letter and one special character",
	"securityAnswer": "Security answer must be at least 3 characters long",
	"min":            "Value must be at least %s characters long",
	"max":            "Value must be at most %s characters long",
}

func GetValidationErrorMessage(tag string, param string) string {
	if msg, exists := ValidationErrorMessages[tag]; exists {
		if strings.Contains(msg, "%s") && param != "" {
			return strings.Replace(msg, "%s", param, 1)
		}
		return msg
	}
	return "Validation failed for field"
}

func HandleValidationErrors(ctx *gin.Context, err error) {
	requestID := utils.GetRequestID(ctx)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make([]utils.ValidationError, 0)

		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			errors = append(errors, utils.ValidationError{
				Field:   field,
				Message: GetValidationErrorMessage(tag, param),
			})
		}

		ctx.JSON(400, utils.StandardizedErrorResponse{
			Errors:    errors,
			RequestID: requestID,
			Status:    "ERROR",
			Code:      "VALIDATION_ERROR",
			Message:   "Validation failed",
		})
		return
	}

	utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err), requestID)
}
