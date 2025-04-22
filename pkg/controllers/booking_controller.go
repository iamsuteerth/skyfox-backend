package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type BookingController struct {
	bookingService      services.BookingService
	adminBookingService services.AdminBookingService
}

func NewBookingController(bookingService services.BookingService, adminBookingService services.AdminBookingService) *BookingController {
	return &BookingController{
		bookingService:      bookingService,
		adminBookingService: adminBookingService,
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

func (bc *BookingController) CreateAdminBooking(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	var bookingRequest request.AdminBookingRequest
	if err := ctx.ShouldBindJSON(&bookingRequest); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			utils.HandleErrorResponse(ctx, utils.NewValidationError(validationErrs), requestID)
			return
		}
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err), requestID)
		return
	}

	booking, err := bc.adminBookingService.CreateAdminBooking(ctx.Request.Context(), bookingRequest)
	if err != nil {
		log.Error().Err(err).Interface("request", bookingRequest).Msg("Failed to create admin booking")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	utils.SendCreatedResponse(ctx, "Booking created successfully", requestID, booking)
}
