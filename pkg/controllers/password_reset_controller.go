package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	customValidator "github.com/iamsuteerth/skyfox-backend/pkg/middleware/validator"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type PasswordResetController struct {
	resetPasswordService services.PasswordResetService
}

func NewPasswordResetController(passwordResetService services.PasswordResetService) *PasswordResetController {
	return &PasswordResetController{
		resetPasswordService: passwordResetService,
	}
}

func (c *PasswordResetController) ForgotPassword(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	var resetRequest request.ResetPasswordRequest
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

	ctx.JSON(200, utils.SuccessResponse{
		Status:    "SUCCESS",
		Message:   "Password has been reset successfully",
		RequestID: requestID,
	})
}
