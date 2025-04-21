package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type BookingController struct {
	bookingService services.BookingService
}

func NewBookingController(bookingService services.BookingService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
	}
}

func (bc *BookingController) GetSeatMap(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	showIDStr := ctx.Param("show_id")
	showID, err := strconv.Atoi(showIDStr)
	if err != nil {
		log.Error().Err(err).Str("showID", showIDStr).Msg("Invalid show ID format")
		utils.HandleErrorResponse(ctx,
			utils.NewBadRequestError("INVALID_SHOW_ID", "Show ID must be a valid integer", err),
			requestID)
		return
	}

	seatMap, err := bc.bookingService.GetSeatMapForShow(ctx.Request.Context(), showID)
	if err != nil {
		log.Error().Err(err).Int("showID", showID).Msg("Failed to get seat map")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	seatsByRow := organizeSeatsByRow(seatMap)

	utils.SendOKResponse(ctx, "Seat map retrieved successfully", requestID, seatsByRow)
}

func organizeSeatsByRow(seatMap []models.SeatMapEntry) map[string]interface{} {
	seatsByRow := make(map[string][]map[string]interface{})

	for _, seat := range seatMap {
		seatsByRow[seat.SeatRow] = append(seatsByRow[seat.SeatRow], map[string]interface{}{
			"seat_number": seat.SeatNumber,
			"column":      seat.SeatColumn,
			"type":        seat.SeatType,
			"price":       seat.Price,
			"occupied":    seat.Occupied,
		})
	}

	return map[string]interface{}{
		"seat_map": seatsByRow,
	}
}
