package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/iamsuteerth/skyfox-backend/pkg/constants"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	paymentservice "github.com/iamsuteerth/skyfox-backend/pkg/payment-service"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type CustomerBookingService interface {
	InitializeBooking(ctx context.Context, username string, req request.InitializeBookingRequest) (*response.InitializeBookingResponse, error)
	ProcessPayment(ctx context.Context, username string, req request.ProcessPaymentRequest) (*response.BookingResponse, error)
	CancelPendingBooking(ctx context.Context, username string, bookingID int) error
	GetBookingsForCustomer(ctx context.Context, username string) ([]response.CustomerBookingInfo, error)
	GetLatestBookingForCustomer(ctx context.Context, username string) (*response.CustomerBookingInfo, error)
}

type customerBookingService struct {
	showRepo               repositories.ShowRepository
	bookingRepo            repositories.BookingRepository
	bookingSeatMappingRepo repositories.BookingSeatMappingRepository
	pendingBookingRepo     repositories.PendingBookingRepository
	paymentTransactionRepo repositories.PaymentTransactionRepository
	slotRepo               repositories.SlotRepository
	skyCustomerRepo        repositories.SkyCustomerRepository
	customerWalletRepo     repositories.CustomerWalletRepository
	walletTxdRepo          repositories.WalletTransactionRepository
	paymentService         paymentservice.PaymentService
}

func NewCustomerBookingService(
	showRepo repositories.ShowRepository,
	bookingRepo repositories.BookingRepository,
	bookingSeatMappingRepo repositories.BookingSeatMappingRepository,
	pendingBookingRepo repositories.PendingBookingRepository,
	paymentTransactionRepo repositories.PaymentTransactionRepository,
	slotRepo repositories.SlotRepository,
	skyCustomerRepo repositories.SkyCustomerRepository,
	customerWalletRepo repositories.CustomerWalletRepository,
	walletTxdRepo repositories.WalletTransactionRepository,
	paymentService paymentservice.PaymentService,
) CustomerBookingService {
	return &customerBookingService{
		showRepo:               showRepo,
		bookingRepo:            bookingRepo,
		bookingSeatMappingRepo: bookingSeatMappingRepo,
		pendingBookingRepo:     pendingBookingRepo,
		paymentTransactionRepo: paymentTransactionRepo,
		slotRepo:               slotRepo,
		skyCustomerRepo:        skyCustomerRepo,
		customerWalletRepo:     customerWalletRepo,
		walletTxdRepo:          walletTxdRepo,
		paymentService:         paymentService,
	}
}

func (s *customerBookingService) InitializeBooking(ctx context.Context, username string, req request.InitializeBookingRequest) (*response.InitializeBookingResponse, error) {
	if len(req.SeatNumbers) > constants.MAX_NO_OF_SEATS_PER_BOOKING {
		return nil, utils.NewBadRequestError("TOO_MANY_SEATS", fmt.Sprintf("Maximum %d seats can be booked per booking", constants.MAX_NO_OF_SEATS_PER_BOOKING), nil)
	}

	show, err := s.showRepo.FindById(ctx, req.ShowID)
	if err != nil {
		log.Error().Err(err).Int("showID", req.ShowID).Msg("Show not found for booking initialization")
		return nil, err
	}

	now := time.Now()
	slot, err := s.slotRepo.GetSlotById(ctx, show.SlotId)
	if err != nil {
		log.Error().Err(err).Int("slotId", show.SlotId).Msg("Failed to get slot details")
		return nil, err
	}

	startTimeParts := strings.Split(slot.StartTime, ":")
	hour, _ := strconv.Atoi(startTimeParts[0])
	minute, _ := strconv.Atoi(startTimeParts[1])

	showDateTime := time.Date(
		show.Date.Year(),
		show.Date.Month(),
		show.Date.Day(),
		hour,
		minute,
		0,
		0,
		now.Location(),
	)

	if now.After(showDateTime) {
		return nil, utils.NewBadRequestError("SHOW_ALREADY_STARTED", "Cannot book tickets for a show that has already started", nil)
	}

	areSeatsAvailable, err := s.bookingSeatMappingRepo.CheckSeatsAvailability(ctx, req.ShowID, req.SeatNumbers)
	if err != nil {
		log.Error().Err(err).Int("showID", req.ShowID).Strs("seatNumbers", req.SeatNumbers).Msg("Failed to check seat availability")
		return nil, err
	}

	if !areSeatsAvailable {
		return nil, utils.NewBadRequestError("SEATS_UNAVAILABLE", "One or more selected seats are not available", nil)
	}

	seatMap, err := s.showRepo.GetSeatMapForShow(ctx, req.ShowID)
	if err != nil {
		log.Error().Err(err).Int("showID", req.ShowID).Msg("Failed to get seat map for price calculation")
		return nil, err
	}

	deluxeOffset, _ := decimal.NewFromFloat64(constants.DELUXE_OFFSET)

	for i := range seatMap {
		if seatMap[i].SeatType == "Deluxe" {
			seatMap[i].Price, _ = show.Cost.Add(deluxeOffset)
		} else {
			seatMap[i].Price = show.Cost
		}
	}

	seatPrices := make(map[string]decimal.Decimal)
	for _, seat := range seatMap {
		seatPrices[seat.SeatNumber] = seat.Price
	}

	var totalPrice decimal.Decimal
	for _, seatNumber := range req.SeatNumbers {
		price, exists := seatPrices[seatNumber]
		if !exists {
			return nil, utils.NewBadRequestError("INVALID_SEAT", fmt.Sprintf("Seat %s not found in seat map", seatNumber), nil)
		}
		totalPrice, _ = totalPrice.Add(price)
	}

	booking := &models.Booking{
		Date:             show.Date,
		ShowId:           req.ShowID,
		CustomerUsername: &username,
		NoOfSeats:        len(req.SeatNumbers),
		AmountPaid:       totalPrice,
		Status:           "Pending",
		PaymentType:      "Card",
	}

	if err := s.bookingRepo.CreatePendingBooking(ctx, booking); err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to create pending booking")
		return nil, err
	}

	if err := s.bookingSeatMappingRepo.CreateMappings(ctx, booking.Id, req.SeatNumbers); err != nil {
		log.Error().Err(err).Int("bookingId", booking.Id).Strs("seatNumbers", req.SeatNumbers).Msg("Failed to map seats to booking")
		_ = s.bookingRepo.DeleteBookingsByIds(ctx, []int{booking.Id})
		return nil, err
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	if err := s.pendingBookingRepo.TrackPendingBooking(ctx, booking.Id, expirationTime); err != nil {
		log.Error().Err(err).Int("bookingId", booking.Id).Msg("Failed to track pending booking")
		_ = s.bookingRepo.DeleteBookingsByIds(ctx, []int{booking.Id})
		return nil, err
	}

	go s.monitorBookingExpiration(booking.Id, expirationTime)

	totalPriceFloat, _ := totalPrice.Float64()

	return &response.InitializeBookingResponse{
		BookingID:       booking.Id,
		ShowID:          booking.ShowId,
		SeatNumbers:     req.SeatNumbers,
		AmountDue:       totalPriceFloat,
		ExpirationTime:  expirationTime,
		TimeRemainingMs: int64(5 * time.Minute / time.Millisecond),
	}, nil
}

func (s *customerBookingService) ProcessPayment(ctx context.Context, username string, req request.ProcessPaymentRequest) (*response.BookingResponse, error) {
	booking, err := s.bookingRepo.GetBookingById(ctx, req.BookingID)
	if err != nil {
		log.Error().Err(err).Int("bookingID", req.BookingID).Msg("Failed to get booking")
		return nil, err
	}

	customer, err := s.skyCustomerRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to get customer details")
		return nil, err
	}

	if booking == nil {
		return nil, utils.NewNotFoundError("BOOKING_NOT_FOUND", "Booking not found", nil)
	}

	if booking.CustomerUsername == nil || *booking.CustomerUsername != username {
		log.Warn().Str("requestedBy", username).Str("owner", *booking.CustomerUsername).Int("bookingID", req.BookingID).Msg("Unauthorized booking access attempt")
		return nil, utils.NewForbiddenError("UNAUTHORIZED_ACCESS", "You don't have permission to access this booking", nil)
	}

	if booking.Status != "Pending" {
		return nil, utils.NewBadRequestError("INVALID_BOOKING_STATUS", "Payment can only be processed for bookings in pending state", nil)
	}

	expirationTime, err := s.pendingBookingRepo.GetExpirationTime(ctx, req.BookingID)
	if err != nil {
		log.Error().Err(err).Int("bookingID", req.BookingID).Msg("Failed to get expiration time")
		return nil, err
	}

	if expirationTime == nil || time.Now().After(*expirationTime) {
		_ = s.bookingRepo.DeleteBookingsByIds(ctx, []int{req.BookingID})
		return nil, utils.NewBadRequestError("BOOKING_EXPIRED", "This booking has expired. Please make a new booking", nil)
	}

	show, err := s.showRepo.FindById(ctx, booking.ShowId)
	if err != nil {
		log.Error().Err(err).Int("showID", booking.ShowId).Msg("Failed to get show details")
		return nil, err
	}

	slot, err := s.slotRepo.GetSlotById(ctx, show.SlotId)
	if err != nil {
		log.Error().Err(err).Int("slotId", show.SlotId).Msg("Failed to get slot details")
		return nil, err
	}

	var transactionID string

	if req.PaymentMethod == "Wallet" {
		wallet, err := s.customerWalletRepo.GetWalletByUsername(ctx, username)
		if err != nil {
			return nil, err
		}

		if wallet == nil {
			return nil, utils.NewNotFoundError("WALLET_NOT_FOUND", "Wallet not found", nil)
		}

		walletBalance := wallet.Balance
		bookingAmount := booking.AmountPaid

		switch {
		case walletBalance.Cmp(bookingAmount) != -1:
			deductTxnID := uuid.New().String()
			if err := s.customerWalletRepo.DeductFromWalletBalance(ctx, username, bookingAmount); err != nil {
				return nil, err
			}

			walletTxn := &models.WalletTransaction{
				WalletID:        wallet.ID,
				Username:        username,
				BookingID:       toPtr(int64(booking.Id)),
				TransactionID:   deductTxnID,
				Amount:          bookingAmount,
				TransactionType: "DEDUCT",
			}

			if err := s.walletTxdRepo.AddWalletTransaction(ctx, walletTxn); err != nil {
				if rollbackErr := s.customerWalletRepo.AddToWalletBalance(ctx, username, bookingAmount); rollbackErr != nil {
					log.Error().Err(rollbackErr).Str("username", username).Msg("Failed to rollback wallet deduction")
				}
				return nil, err
			}

			payTxn := &models.PaymentTransaction{
				BookingId:     booking.Id,
				TransactionId: deductTxnID,
				PaymentMethod: "Wallet",
				Amount:        bookingAmount,
				Status:        "Completed",
			}

			if err := s.paymentTransactionRepo.CreateTransaction(ctx, payTxn); err != nil {
				log.Error().Err(err).Int("bookingId", booking.Id).Msg("Failed to create wallet payment transaction")
				return nil, err
			}

			transactionID = deductTxnID
			booking.PaymentType = "Wallet"

			if err := s.bookingRepo.UpdateBookingPaymentType(ctx, booking.Id, booking.PaymentType); err != nil {
				log.Error().Err(err).Int("bookingId", booking.Id).Str("paymentType", booking.PaymentType).
					Msg("Failed to update booking payment type after processing payment")
				return nil, err
			}

		case walletBalance.Cmp(bookingAmount) == -1:
			if req.CardNumber == "" || req.CVV == "" || req.ExpiryMonth == "" || req.ExpiryYear == "" || req.CardholderName == "" {
				return nil, utils.NewBadRequestError("INSUFFICIENT_BALANCE", "Not enough wallet balance and no card details provided", nil)
			}
			requiredTopUp, _ := bookingAmount.Sub(walletBalance)

			expiry := fmt.Sprintf("%s/%s", req.ExpiryMonth, req.ExpiryYear)
			addFundsTxnID, err := s.paymentService.ProcessPayment(
				ctx,
				req.CardNumber,
				req.CVV,
				expiry,
				req.CardholderName,
				requiredTopUp,
			)

			if err != nil {
				log.Error().Err(err).Int("bookingID", req.BookingID).Msg("Partial payment card processing failed")
				return nil, err
			}

			if err := s.customerWalletRepo.AddToWalletBalance(ctx, username, requiredTopUp); err != nil {
				return nil, err
			}

			addFundsTxn := &models.WalletTransaction{
				WalletID:        wallet.ID,
				Username:        username,
				BookingID:       nil,
				TransactionID:   addFundsTxnID,
				Amount:          requiredTopUp,
				TransactionType: "ADD",
			}

			if err := s.walletTxdRepo.AddWalletTransaction(ctx, addFundsTxn); err != nil {
				if rollbackErr := s.customerWalletRepo.DeductFromWalletBalance(ctx, username, requiredTopUp); rollbackErr != nil {
					log.Error().Err(rollbackErr).Str("username", username).Msg("Failed to rollback wallet addition")
				}
				return nil, err
			}

			deductTxnID := uuid.New().String()
			if err := s.customerWalletRepo.DeductFromWalletBalance(ctx, username, bookingAmount); err != nil {
				rollbackErr := s.customerWalletRepo.DeductFromWalletBalance(ctx, username, requiredTopUp)
				if rollbackErr != nil {
					log.Error().Err(rollbackErr).Str("username", username).
						Str("operation", "rollback_addition").
						Str("amount", requiredTopUp.String()).
						Msg("Failed to rollback wallet addition after deduction failure")
				} else {
					log.Info().Str("username", username).
						Str("amount", requiredTopUp.String()).
						Msg("Successfully rolled back wallet addition after deduction failure")
				}
				return nil, err
			}

			walletTxn := &models.WalletTransaction{
				WalletID:        wallet.ID,
				Username:        username,
				BookingID:       toPtr(int64(booking.Id)),
				TransactionID:   deductTxnID,
				Amount:          bookingAmount,
				TransactionType: "DEDUCT",
			}
			if err := s.walletTxdRepo.AddWalletTransaction(ctx, walletTxn); err != nil {
				if rollbackErr := s.customerWalletRepo.AddToWalletBalance(ctx, username, bookingAmount); rollbackErr != nil {
					log.Error().Err(rollbackErr).Str("username", username).Msg("Failed to rollback wallet deduction")
				}
				if rollbackErr := s.customerWalletRepo.DeductFromWalletBalance(ctx, username, requiredTopUp); rollbackErr != nil {
					log.Error().Err(rollbackErr).Str("username", username).Msg("Failed to rollback wallet addition")
				}
				return nil, err
			}

			payTxn := &models.PaymentTransaction{
				BookingId:     booking.Id,
				TransactionId: deductTxnID,
				PaymentMethod: "Wallet",
				Amount:        bookingAmount,
				Status:        "Completed",
			}
			if err := s.paymentTransactionRepo.CreateTransaction(ctx, payTxn); err != nil {
				log.Error().Err(err).Int("bookingId", booking.Id).Msg("Failed to create partial wallet payment transaction")
				return nil, err
			}
			transactionID = deductTxnID
			booking.PaymentType = "Wallet"

			if err := s.bookingRepo.UpdateBookingPaymentType(ctx, booking.Id, booking.PaymentType); err != nil {
				log.Error().Err(err).Int("bookingId", booking.Id).Str("paymentType", booking.PaymentType).
					Msg("Failed to update booking payment type after processing payment")
				return nil, err
			}

		default:
			return nil, utils.NewBadRequestError("INSUFFICIENT_BALANCE", "Not enough wallet balance and no card details provided", nil)
		}

	} else if req.PaymentMethod == "Card" {
		expiry := fmt.Sprintf("%s/%s", req.ExpiryMonth, req.ExpiryYear)
		transactionID, err = s.paymentService.ProcessPayment(
			ctx,
			req.CardNumber,
			req.CVV,
			expiry,
			req.CardholderName,
			booking.AmountPaid,
		)
		if err != nil {
			log.Error().Err(err).Int("bookingID", req.BookingID).Msg("Card payment processing failed")
			return nil, err
		}

		transaction := &models.PaymentTransaction{
			BookingId:     booking.Id,
			TransactionId: transactionID,
			PaymentMethod: "Card",
			Amount:        booking.AmountPaid,
			Status:        "Completed",
		}
		if err := s.paymentTransactionRepo.CreateTransaction(ctx, transaction); err != nil {
			log.Error().Err(err).Int("bookingId", booking.Id).Msg("Failed to create card transaction record")
			return nil, err
		}
		booking.PaymentType = "Card"
	} else {
		return nil, utils.NewBadRequestError("INVALID_PAYMENT_METHOD", "Invalid payment method specified", nil)
	}

	if err := s.bookingRepo.UpdateBookingStatus(ctx, booking.Id, "Confirmed"); err != nil {
		log.Error().Err(err).Int("bookingId", booking.Id).Msg("Failed to update booking status")
		return nil, err
	}
	if err := s.pendingBookingRepo.RemoveTracker(ctx, booking.Id); err != nil {
		log.Error().Err(err).Int("bookingId", booking.Id).Msg("Failed to remove pending tracker")
	}

	seatNumbers, err := s.bookingSeatMappingRepo.GetSeatsByBookingId(ctx, booking.Id)
	if err != nil {
		log.Error().Err(err).Int("bookingId", booking.Id).Msg("Failed to get seat numbers")
		return nil, err
	}

	bookingAmtPaid, _ := booking.AmountPaid.Float64()

	response := &response.BookingResponse{
		BookingID:     booking.Id,
		ShowID:        booking.ShowId,
		ShowDate:      show.Date.Format("2006-01-02"),
		ShowTime:      slot.StartTime,
		CustomerName:  customer.Name,
		PhoneNumber:   customer.Number,
		SeatNumbers:   seatNumbers,
		AmountPaid:    bookingAmtPaid,
		PaymentType:   string(booking.PaymentType),
		BookingTime:   booking.BookingTime,
		Status:        "Confirmed",
		TransactionID: transactionID,
	}
	return response, nil
}

func toPtr[T any](v T) *T { return &v }

func (s *customerBookingService) CancelPendingBooking(ctx context.Context, username string, bookingID int) error {
	booking, err := s.bookingRepo.GetBookingById(ctx, bookingID)
	if err != nil {
		log.Error().Err(err).Int("bookingID", bookingID).Msg("Failed to get booking for cancellation")
		return err
	}

	if booking == nil {
		return utils.NewNotFoundError("BOOKING_NOT_FOUND", "Booking not found", nil)
	}

	if booking.CustomerUsername == nil || *booking.CustomerUsername != username {
		log.Warn().Str("requestedBy", username).Str("owner", *booking.CustomerUsername).Int("bookingID", bookingID).Msg("Unauthorized booking cancellation attempt")
		return utils.NewForbiddenError("UNAUTHORIZED_ACCESS", "You don't have permission to access this booking", nil)
	}

	if booking.Status != "Pending" {
		return utils.NewBadRequestError("INVALID_BOOKING_STATUS", "Only pending bookings can be cancelled", nil)
	}

	if err := s.bookingRepo.DeleteBookingsByIds(ctx, []int{bookingID}); err != nil {
		log.Error().Err(err).Int("bookingID", bookingID).Msg("Failed to delete booking during cancellation")
		return err
	}

	if err := s.pendingBookingRepo.RemoveTracker(ctx, bookingID); err != nil {
		log.Warn().Err(err).Int("bookingID", bookingID).Msg("Failed to remove pending tracker during cancellation")
	}

	log.Info().Int("bookingID", bookingID).Str("username", username).Msg("Booking successfully cancelled")
	return nil
}

func (s *customerBookingService) monitorBookingExpiration(bookingId int, expirationTime time.Time) {
	ctx := context.Background()

	sleepDuration := time.Until(expirationTime)

	log.Info().
		Int("bookingId", bookingId).
		Time("expirationTime", expirationTime).
		Dur("sleepDuration", sleepDuration).
		Msg("Started monitoring booking expiration")

	time.Sleep(sleepDuration)

	tracker, err := s.pendingBookingRepo.GetExpirationTime(ctx, bookingId)
	if err != nil {
		log.Error().Err(err).Int("bookingId", bookingId).Msg("Error checking booking expiration status")
		return
	}

	if tracker == nil {
		log.Debug().Int("bookingId", bookingId).Msg("Booking no longer pending, skipping expiration")
		return
	}

	booking, err := s.bookingRepo.GetBookingById(ctx, bookingId)
	if err != nil || booking == nil {
		log.Error().Err(err).Int("bookingId", bookingId).Msg("Error retrieving booking for expiration")
		return
	}

	if booking.Status == "Pending" {
		if err := s.bookingRepo.DeleteBookingsByIds(ctx, []int{bookingId}); err != nil {
			log.Error().Err(err).Int("bookingId", bookingId).Msg("Failed to delete expired booking")
		} else {
			log.Info().Int("bookingId", bookingId).Msg("Successfully deleted expired booking")
		}
	}
}

func (s *customerBookingService) GetBookingsForCustomer(ctx context.Context, username string) ([]response.CustomerBookingInfo, error) {
	bookings, err := s.bookingRepo.FindByCustomerUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	result := make([]response.CustomerBookingInfo, 0, len(bookings))
	for _, booking := range bookings {
		show, err := s.showRepo.FindById(ctx, booking.ShowId)
		if err != nil || show == nil {
			log.Warn().Int("booking_id", booking.Id).Msg("Skipping booking due to missing show data")
			continue
		}
		slot, err := s.slotRepo.GetSlotById(ctx, show.SlotId)
		if err != nil || slot == nil {
			log.Warn().Int("booking_id", booking.Id).Msg("Skipping booking due to missing slot data")
			continue
		}
		seats, err := s.bookingSeatMappingRepo.GetSeatsByBookingId(ctx, booking.Id)
		if err != nil {
			log.Warn().Int("booking_id", booking.Id).Msg("Failed to fetch seats, using empty list")
			seats = []string{}
		}

		bookingAmtPaid, _ := booking.AmountPaid.Float64()

		result = append(result, response.CustomerBookingInfo{
			BookingID:   booking.Id,
			ShowID:      booking.ShowId,
			ShowDate:    show.Date.Format("2006-01-02"),
			ShowTime:    slot.StartTime,
			SeatNumbers: seats,
			AmountPaid:  bookingAmtPaid,
			PaymentType: booking.PaymentType,
			BookingTime: booking.BookingTime,
			Status:      booking.Status,
		})
	}
	return result, nil
}

func (s *customerBookingService) GetLatestBookingForCustomer(ctx context.Context, username string) (*response.CustomerBookingInfo, error) {
	booking, err := s.bookingRepo.FindLatestByCustomerUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, nil
	}
	show, err := s.showRepo.FindById(ctx, booking.ShowId)
	if err != nil || show == nil {
		log.Warn().Int("booking_id", booking.Id).Msg("No show found for booking, returning nil")
		return nil, nil
	}
	slot, err := s.slotRepo.GetSlotById(ctx, show.SlotId)
	if err != nil || slot == nil {
		log.Warn().Int("booking_id", booking.Id).Msg("No slot found for booking, returning nil")
		return nil, nil
	}
	seats, err := s.bookingSeatMappingRepo.GetSeatsByBookingId(ctx, booking.Id)
	if err != nil {
		log.Warn().Int("booking_id", booking.Id).Msg("Failed to fetch seats, using empty list")
		seats = []string{}
	}

	bookingAmtPaid, _ := booking.AmountPaid.Float64()

	return &response.CustomerBookingInfo{
		BookingID:   booking.Id,
		ShowID:      booking.ShowId,
		ShowDate:    show.Date.Format("2006-01-02"),
		ShowTime:    slot.StartTime,
		SeatNumbers: seats,
		AmountPaid:  bookingAmtPaid,
		PaymentType: booking.PaymentType,
		BookingTime: booking.BookingTime,
		Status:      booking.Status,
	}, nil
}
