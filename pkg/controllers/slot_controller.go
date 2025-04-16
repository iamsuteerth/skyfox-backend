package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type SlotController struct {
	slotService services.SlotService
}

func NewSlotController(slotService services.SlotService) *SlotController {
	return &SlotController{
		slotService: slotService,
	}
}

func (sc *SlotController) GetAvailableSlots(c *gin.Context) {
	requestID := utils.GetRequestID(c)

	dateStr := c.Query("date")
	date, err := utils.GetDateFromDateStringDefaultToday(dateStr)

	if err != nil {
		utils.HandleErrorResponse(c, utils.NewBadRequestError("INVALID_DATE", "Invalid date format. Use YYYY-MM-DD", err), requestID)
		return
	}

	slots, err := sc.slotService.GetAvailableSlots(c.Request.Context(), date)
	if err != nil {
		log.Error().Err(err).Str("date", dateStr).Msg("Failed to get available slots")
		utils.HandleErrorResponse(c, err, requestID)
		return
	}

	var slotResponses []response.SlotResponse
	for _, slot := range slots {
		slotResponses = append(slotResponses, response.NewSlotResponse(slot))
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message:   "Available slots retrieved successfully",
		RequestID: requestID,
		Status:    "SUCCESS",
		Data:      slotResponses,
	})
}
