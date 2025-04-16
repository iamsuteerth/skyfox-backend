package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamsuteerth/skyfox-backend/pkg/middleware/security"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type AdminStaffController struct {
	adminStaffProfileService services.AdminStaffProfileService
}

func NewAdminStaffController(adminStaffProfileService services.AdminStaffProfileService) *AdminStaffController {
	return &AdminStaffController{
		adminStaffProfileService: adminStaffProfileService,
	}
}

func (c *AdminStaffController) GetAdminProfile(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify credentials", err), requestID)
		return
	}

	tokenUsername, ok := claims["username"].(string)
	if !ok || tokenUsername == "" {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("INVALID_TOKEN", "Invalid token claims", nil), requestID)
		return
	}

	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		utils.HandleErrorResponse(ctx, utils.NewForbiddenError("FORBIDDEN", "Access denied. Admin role required", nil), requestID)
		return
	}

	profile, err := c.adminStaffProfileService.GetProfile(ctx.Request.Context(), tokenUsername)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse{
		Message:   "Admin profile retrieved successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data:      profile,
	})
}

func (c *AdminStaffController) GetStaffProfile(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify credentials", err), requestID)
		return
	}

	tokenUsername, ok := claims["username"].(string)
	if !ok || tokenUsername == "" {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("INVALID_TOKEN", "Invalid token claims", nil), requestID)
		return
	}

	role, ok := claims["role"].(string)
	if !ok || role != "staff" {
		utils.HandleErrorResponse(ctx, utils.NewForbiddenError("FORBIDDEN", "Access denied. Staff role required", nil), requestID)
		return
	}

	profile, err := c.adminStaffProfileService.GetProfile(ctx.Request.Context(), tokenUsername)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse{
		Message:   "Staff profile retrieved successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data:      profile,
	})
}
