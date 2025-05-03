// In pkg/controllers/dashboard_revenue_controller.go
package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
    "github.com/iamsuteerth/skyfox-backend/pkg/services"
    "github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type DashboardRevenueController struct {
    revenueService services.RevenueService
}

func NewDashboardRevenueController(revenueService services.RevenueService) *DashboardRevenueController {
    return &DashboardRevenueController{
        revenueService: revenueService,
    }
}

func (c *DashboardRevenueController) GetRevenue(ctx *gin.Context) {
    requestID := utils.GetRequestID(ctx)
    
    // Parse query parameters
    var req request.RevenueDashboardRequest
    if err := ctx.ShouldBindQuery(&req); err != nil {
        utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_PARAMS", "Invalid query parameters", err), requestID)
        return
    }
    
    // Get revenue data
    revenueData, err := c.revenueService.GetRevenue(ctx, req)
    if err != nil {
        utils.HandleErrorResponse(ctx, err, requestID)
        return
    }
    
    utils.SendOKResponse(ctx, "Revenue data fetched successfully", requestID, revenueData)
}
