package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/middleware/security"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type BookingController struct {
	bookingService         services.BookingService
	adminBookingService    services.AdminBookingService
	customerBookingService services.CustomerBookingService
}

func NewBookingController(bookingService services.BookingService, adminBookingService services.AdminBookingService, customerBookingService services.CustomerBookingService) *BookingController {
	return &BookingController{
		bookingService:         bookingService,
		adminBookingService:    adminBookingService,
		customerBookingService: customerBookingService,
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

func (bc *BookingController) InitializeCustomerBooking(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	var bookingRequest request.InitializeBookingRequest
	if err := ctx.ShouldBindJSON(&bookingRequest); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			utils.HandleErrorResponse(ctx, utils.NewValidationError(validationErrs), requestID)
			return
		}
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err), requestID)
		return
	}

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify user credentials", err), requestID)
		return
	}

	username, _ := claims["username"].(string)

	booking, err := bc.customerBookingService.InitializeBooking(ctx.Request.Context(), username, bookingRequest)
	if err != nil {
		log.Error().Err(err).Interface("request", bookingRequest).Msg("Failed to initialize booking")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	utils.SendCreatedResponse(ctx, "Booking initialized successfully", requestID, booking)
}

func (bc *BookingController) ProcessPayment(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	var paymentRequest request.ProcessPaymentRequest
	if err := ctx.ShouldBindJSON(&paymentRequest); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			utils.HandleErrorResponse(ctx, utils.NewValidationError(validationErrs), requestID)
			return
		}
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid payment data", err), requestID)
		return
	}

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify user credentials", err), requestID)
		return
	}

	username, _ := claims["username"].(string)

	confirmedBooking, err := bc.customerBookingService.ProcessPayment(ctx.Request.Context(), username, paymentRequest)
	if err != nil {
		log.Error().Err(err).Interface("request", paymentRequest).Msg("Failed to process payment")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	utils.SendOKResponse(ctx, "Payment processed successfully", requestID, confirmedBooking)
}

func (bc *BookingController) GetQRCode(ctx *gin.Context) {
    requestID := utils.GetRequestID(ctx)
    
    bookingIDStr := ctx.Param("id")
    bookingID, err := strconv.Atoi(bookingIDStr)
    if err != nil {
        utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_BOOKING_ID", "Invalid booking ID format", err), requestID)
        return
    }
    
    claims, err := security.GetTokenClaims(ctx)
    if err != nil {
        utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify user credentials", err), requestID)
        return
    }
    
    role := claims["role"].(string)
    username := claims["username"].(string)
    
    booking, err := bc.bookingService.GetBookingById(ctx.Request.Context(), bookingID)
    if err != nil {
        utils.HandleErrorResponse(ctx, err, requestID)
        return
    }
    
    if role != "admin" && role != "staff" {
        if booking.CustomerUsername == nil || *booking.CustomerUsername != username {
            utils.HandleErrorResponse(ctx, utils.NewForbiddenError("FORBIDDEN", "Access denied to this booking", nil), requestID)
            return
        }
    }
    
    qrCode, err := bc.bookingService.GenerateQRCode(ctx.Request.Context(), bookingID)
    if err != nil {
        utils.HandleErrorResponse(ctx, err, requestID)
        return
    }
    
    response := map[string]interface{}{
        "qr_code": qrCode,
        "mime_type": "image/png",
        "encoding": "base64",
    }
    
    utils.SendOKResponse(ctx, "QR code generated successfully", requestID, response)
}

func (bc *BookingController) GetPDF(ctx *gin.Context) {
    requestID := utils.GetRequestID(ctx)
    
    bookingIDStr := ctx.Param("id")
    bookingID, err := strconv.Atoi(bookingIDStr)
    if err != nil {
        utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_BOOKING_ID", "Invalid booking ID format", err), requestID)
        return
    }
    
    claims, err := security.GetTokenClaims(ctx)
    if err != nil {
        utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify user credentials", err), requestID)
        return
    }
    
    role := claims["role"].(string)
    username := claims["username"].(string)
    
    booking, err := bc.bookingService.GetBookingById(ctx.Request.Context(), bookingID)
    if err != nil {
        utils.HandleErrorResponse(ctx, err, requestID)
        return
    }
    
    if role != "admin" && role != "staff" {
        if booking.CustomerUsername == nil || *booking.CustomerUsername != username {
            utils.HandleErrorResponse(ctx, utils.NewForbiddenError("FORBIDDEN", "Access denied to this booking", nil), requestID)
            return
        }
    }
    
    pdf, err := bc.bookingService.GeneratePDF(ctx.Request.Context(), bookingID)
    if err != nil {
        utils.HandleErrorResponse(ctx, err, requestID)
        return
    }
    
    response := map[string]interface{}{
        "pdf": pdf,
        "mime_type": "application/pdf",
        "encoding": "base64",
    }
    
    utils.SendOKResponse(ctx, "PDF ticket generated successfully", requestID, response)
}
