package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type SecurityQuestionController struct {
	securityQuestionService services.SecurityQuestionService
}

func NewSecurityQuestionController(securityQuestionService services.SecurityQuestionService) *SecurityQuestionController {
	return &SecurityQuestionController{
		securityQuestionService: securityQuestionService,
	}
}

func (c *SecurityQuestionController) GetSecurityQuestions(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	questions, err := c.securityQuestionService.GetAllSecurityQuestions(ctx.Request.Context())
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse{
		Message:   "Security questions retrieved successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data:      questions,
	})
}

func (c *SecurityQuestionController) GetSecurityQuestionByEmail(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)
	var req request.GetSecurityQuestionByEmailRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			utils.HandleErrorResponse(ctx, utils.NewValidationError(validationErrs), requestID)
			return
		}
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err), requestID)
		return
	}

	question, err := c.securityQuestionService.GetSecurityQuestionByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse{
		Message:   "Security question retrieved successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data:      question,
	})
}
