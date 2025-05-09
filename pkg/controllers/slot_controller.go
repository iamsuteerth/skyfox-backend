package controllers

import (
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

func (sc *SlotController) GetAvailableSlots(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	dateStr := ctx.Query("date")
	date, err := utils.GetDateFromDateStringDefaultToday(dateStr)

	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_DATE", "Invalid date format. Use YYYY-MM-DD", err), requestID)
		return
	}

	slots, err := sc.slotService.GetAvailableSlots(ctx.Request.Context(), date)
	if err != nil {
		log.Error().Err(err).Str("date", dateStr).Msg("Failed to get available slots")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	var slotResponses []response.SlotResponse
	for _, slot := range slots {
		slotResponses = append(slotResponses, response.NewSlotResponse(slot))
	}

	utils.SendOKResponse(ctx, "Available slots retrieved successfully", requestID, slotResponses)
}

func (sc *SlotController) GetAllSlots(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	slots, err := sc.slotService.GetAllSlots(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all slots")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	var slotResponses []response.SlotResponse
	for _, slot := range slots {
		slotResponses = append(slotResponses, response.NewSlotResponse(slot))
	}

	utils.SendOKResponse(ctx, "Slots retrieved successfully", requestID, slotResponses)
}
