package controllers

import (
	"strings"

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
    
    var req request.RevenueDashboardRequest
    if err := ctx.ShouldBindQuery(&req); err != nil {
        utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_PARAMS", "Invalid query parameters", err), requestID)
        return
    }
    
    if req.Timeframe != "" && req.Timeframe != "all" && (req.Month != nil || req.Year != nil) {
        utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_PARAMS", "Timeframe cannot be combined with month or year filters", nil), requestID)
        return
    }
    
    queryString := ctx.Request.URL.RawQuery
    paramOrder := extractParamOrder(queryString)
    req.ParamOrder = paramOrder
    
    revenueData, err := c.revenueService.GetRevenue(ctx, req)
    if err != nil {
        utils.HandleErrorResponse(ctx, err, requestID)
        return
    }
    
    utils.SendOKResponse(ctx, "Revenue data fetched successfully", requestID, revenueData)
}


func extractParamOrder(queryString string) []string {
	if queryString == "" {
		return []string{}
	}

	var paramOrder []string
	parts := strings.Split(queryString, "&")

	for _, part := range parts {
		if paramPair := strings.Split(part, "="); len(paramPair) > 0 {
			paramName := paramPair[0]
			paramOrder = append(paramOrder, paramName)
		}
	}

	return paramOrder
}
