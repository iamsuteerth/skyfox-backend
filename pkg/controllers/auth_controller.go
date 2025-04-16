package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
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
	var loginRequest request.LoginRequest
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

	loginResponse := response.LoginResponse{
		User: response.UserInfo{
			Username: user.Username,
			Role:     user.Role,
		},
		Token: token,
	}

	utils.SendOKResponse(ctx, "Login successful", requestID, loginResponse)
}
