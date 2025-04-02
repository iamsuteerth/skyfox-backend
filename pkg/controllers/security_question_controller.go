package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
