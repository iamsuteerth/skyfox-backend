package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
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
	checkInService         services.CheckInService
	bookingCSVService      services.BookingCSVService
}

func NewBookingController(bookingService services.BookingService, adminBookingService services.AdminBookingService, customerBookingService services.CustomerBookingService, checkInService services.CheckInService, bookingCSVService services.BookingCSVService) *BookingController {
	return &BookingController{
		bookingService:         bookingService,
		adminBookingService:    adminBookingService,
		customerBookingService: customerBookingService,
		checkInService:         checkInService,
		bookingCSVService:      bookingCSVService,
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
            var relevantErrors validator.ValidationErrors
            for _, fieldErr := range validationErrs {
                if fieldErr.Field() == "BookingID" || fieldErr.Field() == "PaymentMethod" {
                    relevantErrors = append(relevantErrors, fieldErr)
                }
            }
            
            if len(relevantErrors) > 0 {
                utils.HandleErrorResponse(ctx, utils.NewValidationError(relevantErrors), requestID)
                return
            }
        } else {
            utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid payment data", err), requestID)
            return
        }
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

func (bc *BookingController) CancelBooking(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)
	bookingIDStr := ctx.Param("id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		utils.HandleErrorResponse(ctx,
			utils.NewBadRequestError("INVALID_BOOKING_ID", "Booking ID must be a valid integer", err),
			requestID)
		return
	}

	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx,
			utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify user credentials", err),
			requestID)
		return
	}

	username, _ := claims["username"].(string)

	err = bc.customerBookingService.CancelPendingBooking(ctx.Request.Context(), username, bookingID)
	if err != nil {
		log.Error().Err(err).Int("bookingID", bookingID).Str("username", username).Msg("Failed to cancel booking")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}

	utils.SendOKResponse(ctx, "Booking cancelled successfully", requestID, nil)
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
		"qr_code":   qrCode,
		"mime_type": "image/png",
		"encoding":  "base64",
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
		"pdf":       pdf,
		"mime_type": "application/pdf",
		"encoding":  "base64",
	}

	utils.SendOKResponse(ctx, "PDF ticket generated successfully", requestID, response)
}

func (c *BookingController) GetCustomerBookings(ctx *gin.Context) {
	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("INVALID_TOKEN", "Unauthorized", err), utils.GetRequestID(ctx))
		return
	}
	username, _ := claims["username"].(string)

	bookings, err := c.customerBookingService.GetBookingsForCustomer(ctx, username)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, utils.GetRequestID(ctx))
		return
	}
	utils.SendOKResponse(ctx, "Bookings fetched successfully", utils.GetRequestID(ctx), bookings)
}

func (c *BookingController) GetCustomerLatestBooking(ctx *gin.Context) {
	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("INVALID_TOKEN", "Unauthorized", err), utils.GetRequestID(ctx))
		return
	}
	username, _ := claims["username"].(string)

	booking, err := c.customerBookingService.GetLatestBookingForCustomer(ctx, username)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, utils.GetRequestID(ctx))
		return
	}
	utils.SendOKResponse(ctx, "Latest booking fetched successfully", utils.GetRequestID(ctx), booking)
}

func (c *BookingController) GetCheckInBookings(ctx *gin.Context) {
	bookings, err := c.checkInService.FindConfirmedBookings(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, utils.GetRequestID(ctx))
		return
	}
	utils.SendOKResponse(ctx, "Confirmed bookings fetched successfully", utils.GetRequestID(ctx), bookings)
}

func (c *BookingController) BulkCheckInBookings(ctx *gin.Context) {
	var req request.BulkCheckInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_INPUT", "Invalid input", err), utils.GetRequestID(ctx))
		return
	}

	checkedIn, alreadyDone, invalid, err := c.checkInService.MarkBookingsCheckedIn(ctx, req.BookingIDs)
	if err != nil {
		utils.HandleErrorResponse(ctx, err, utils.GetRequestID(ctx))
		return
	}

	resp := response.BulkCheckInResponse{
		CheckedIn:   checkedIn,
		AlreadyDone: alreadyDone,
		Invalid:     invalid,
	}
	utils.SendOKResponse(ctx, "Bulk check-in attempted", utils.GetRequestID(ctx), resp)
}

func (c *BookingController) SingleCheckInBooking(ctx *gin.Context) {
	var req request.SingleCheckInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_INPUT", "Invalid input", err), utils.GetRequestID(ctx))
		return
	}

	checkedIn, alreadyDone, invalid, err := c.checkInService.MarkBookingsCheckedIn(ctx, []int{req.BookingID})
	if err != nil {
		utils.HandleErrorResponse(ctx, err, utils.GetRequestID(ctx))
		return
	}

	resp := response.BulkCheckInResponse{
		CheckedIn:   checkedIn,
		AlreadyDone: alreadyDone,
		Invalid:     invalid,
	}
	msg := "Booking checked in successfully"
	if len(alreadyDone) > 0 {
		msg = "Booking was already checked in"
	} else if len(invalid) > 0 {
		msg = "Check-in failed: invalid booking (already expired/invalid status/or show ended)"
	}
	utils.SendOKResponse(ctx, msg, utils.GetRequestID(ctx), resp)
}

func (bc *BookingController) DownloadBookingsCSV(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)

	var month, year *int

	monthStr := ctx.Query("month")
	if monthStr != "" {
		mVal, err := strconv.Atoi(monthStr)
		if err != nil || mVal < 1 || mVal > 12 {
			utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_MONTH", "Month must be a number between 1 and 12", err), requestID)
			return
		}
		month = &mVal
	}

	yearStr := ctx.Query("year")
	if yearStr != "" {
		yVal, err := strconv.Atoi(yearStr)
		if err != nil {
			utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_YEAR", "Year must be a valid number", err), requestID)
			return
		}
		year = &yVal
	}

	filename := "bookings.csv"
	if month != nil && year != nil {
		monthName := time.Month(*month).String()
		filename = fmt.Sprintf("bookings_%s_%d.csv", monthName, *year)
	} else if month != nil {
		monthName := time.Month(*month).String()
		filename = fmt.Sprintf("bookings_%s.csv", monthName)
	} else if year != nil {
		filename = fmt.Sprintf("bookings_%d.csv", *year)
	}

	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	if err := bc.bookingCSVService.WriteBookingsCSV(ctx, ctx.Writer, month, year); err != nil {
		ctx.Status(500)
		log.Error().Err(err).Msg("Failed to write bookings CSV")
		return
	}
}
