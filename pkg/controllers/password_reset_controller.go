package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/middleware/security"
	customValidator "github.com/iamsuteerth/skyfox-backend/pkg/middleware/validator"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type PasswordResetController struct {
	resetPasswordService services.PasswordResetService
	skyCustomerService   services.SkyCustomerService
}

func NewPasswordResetController(passwordResetService services.PasswordResetService, skyCustomerService services.SkyCustomerService) *PasswordResetController {
	return &PasswordResetController{
		resetPasswordService: passwordResetService,
		skyCustomerService:   skyCustomerService,
	}
}

func (c *PasswordResetController) ForgotPassword(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	var resetRequest request.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&resetRequest); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			customValidator.HandleValidationErrors(ctx, err)
			return
		}

		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request body", err), requestID)
		return
	}

	err := c.resetPasswordService.ForgotPassword(
		ctx.Request.Context(),
		resetRequest.Email,
		resetRequest.ResetToken,
		resetRequest.NewPassword,
	)

	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	utils.SendOKResponse(ctx, "Password has been reset successfully", requestID, nil)
}

func (c *PasswordResetController) ChangePassword(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify credentials", err), requestID)
		return
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("INVALID_TOKEN", "Invalid token claims", nil), requestID)
		return
	}

	var changePasswordRequest request.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&changePasswordRequest); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			customValidator.HandleValidationErrors(ctx, err)
			return
		}

		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request body", err), requestID)
		return
	}

	err = c.resetPasswordService.ChangePassword(
		ctx.Request.Context(),
		username,
		changePasswordRequest.CurrentPassword,
		changePasswordRequest.NewPassword,
	)

	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	utils.SendOKResponse(ctx, "Password updated successfully", requestID, nil)
}
