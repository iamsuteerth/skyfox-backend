// pkg/controllers/auth_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type AuthController struct {
	userService services.UserService
}

func NewAuthController(userService services.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	requestID := utils.GetRequestID(ctx)

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request body", err), requestID)
		return
	}

	user, token, err := c.userService.Login(ctx, loginRequest.Username, loginRequest.Password)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse{
		Message:   "Login successful",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data: gin.H{
			"user": gin.H{
				"username": user.Username,
				"role":     user.Role,
			},
			"token": token,
		},
	})
}
