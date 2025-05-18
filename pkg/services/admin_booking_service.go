package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/govalues/decimal"
	"github.com/iamsuteerth/skyfox-backend/pkg/constants"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type AdminBookingService interface {
	CreateAdminBooking(ctx context.Context, req request.AdminBookingRequest) (*response.BookingResponse, error)
}

type adminBookingService struct {
	showRepo                repositories.ShowRepository
	bookingRepo             repositories.BookingRepository
	bookingSeatMappingRepo  repositories.BookingSeatMappingRepository
	adminBookedCustomerRepo repositories.AdminBookedCustomerRepository
	slotRepo                repositories.SlotRepository
}

func NewAdminBookingService(
	showRepo repositories.ShowRepository,
	bookingRepo repositories.BookingRepository,
	bookingSeatMappingRepo repositories.BookingSeatMappingRepository,
	adminBookedCustomerRepo repositories.AdminBookedCustomerRepository,
	slotRepo repositories.SlotRepository,
) AdminBookingService {
	return &adminBookingService{
		showRepo:                showRepo,
		bookingRepo:             bookingRepo,
		bookingSeatMappingRepo:  bookingSeatMappingRepo,
		adminBookedCustomerRepo: adminBookedCustomerRepo,
		slotRepo:                slotRepo,
	}
}

func (s *adminBookingService) CreateAdminBooking(ctx context.Context, req request.AdminBookingRequest) (*response.BookingResponse, error) {
	if len(req.SeatNumbers) > constants.MAX_NO_OF_SEATS_PER_BOOKING {
		return nil, utils.NewBadRequestError("TOO_MANY_SEATS", fmt.Sprintf("Maximum %d seats can be booked per booking", constants.MAX_NO_OF_SEATS_PER_BOOKING), nil)
	}

	show, err := s.showRepo.FindById(ctx, req.ShowID)
	if err != nil {
		log.Error().Err(err).Int("showID", req.ShowID).Msg("Show not found for admin booking")
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

	var amountPaid decimal.Decimal

	expectedPrice, err := s.calculateTotalPrice(ctx, show, req.SeatNumbers)
	if err != nil {
		log.Error().Err(err).Int("showID", req.ShowID).Strs("seatNumbers", req.SeatNumbers).Msg("Failed to calculate expected price")
		return nil, err
	}

	if !req.AmountPaid.Equal(expectedPrice) {
		log.Warn().
			Str("submitted", req.AmountPaid.String()).
			Str("expected", expectedPrice.String()).
			Msg("Amount mismatch")
		return nil, utils.NewBadRequestError("INVALID_AMOUNT", "The amount paid does not match the expected price", nil)
	}

	amountPaid = req.AmountPaid

	customer := &models.AdminBookedCustomer{
		Name:   req.CustomerName,
		Number: req.PhoneNumber,
	}

	if err := s.adminBookedCustomerRepo.Create(ctx, customer); err != nil {
		log.Error().Err(err).Interface("customer", customer).Msg("Failed to create admin booked customer")
		return nil, err
	}

	booking := &models.Booking{
		Date:        show.Date,
		ShowId:      req.ShowID,
		CustomerId:  &customer.Id,
		NoOfSeats:   len(req.SeatNumbers),
		AmountPaid:  amountPaid,
		Status:      "Confirmed",
		PaymentType: "Cash",
	}

	if err := s.bookingRepo.CreateAdminBooking(ctx, booking); err != nil {
		log.Error().Err(err).Interface("booking", booking).Msg("Failed to create admin booking")
		_ = s.adminBookedCustomerRepo.DeleteById(ctx, customer.Id)
		return nil, err
	}

	if err := s.adminBookedCustomerRepo.UpdateBookingId(ctx, customer.Id, booking.Id); err != nil {
		log.Error().Err(err).Int("customerId", customer.Id).Int("bookingId", booking.Id).Msg("Failed to update admin booked customer with booking ID")
		return nil, err
	}

	if err := s.bookingSeatMappingRepo.CreateMappings(ctx, booking.Id, req.SeatNumbers); err != nil {
		log.Error().Err(err).Int("bookingId", booking.Id).Strs("seatNumbers", req.SeatNumbers).Msg("Failed to create seat mappings")
		_ = s.bookingRepo.DeleteBookingsByIds(ctx, []int{booking.Id})
		_ = s.adminBookedCustomerRepo.DeleteById(ctx, customer.Id)
		return nil, err
	}

	bookingAmtPaid, _ := booking.AmountPaid.Float64()

	bookingResponse := &response.BookingResponse{
		BookingID:    booking.Id,
		ShowID:       booking.ShowId,
		ShowDate:     show.Date.Format("2006-01-02"),
		ShowTime:     slot.StartTime,
		CustomerName: customer.Name,
		PhoneNumber:  customer.Number,
		SeatNumbers:  req.SeatNumbers,
		AmountPaid:   bookingAmtPaid,
		PaymentType:  booking.PaymentType,
		BookingTime:  booking.BookingTime,
		Status:       booking.Status,
	}

	return bookingResponse, nil
}

func (s *adminBookingService) calculateTotalPrice(ctx context.Context, show *models.Show, seatNumbers []string) (decimal.Decimal, error) {
	seatMap, err := s.showRepo.GetSeatMapForShow(ctx, show.Id)
	if err != nil {
		log.Error().Err(err).Int("showID", show.Id).Msg("Failed to get seat map for price calculation")
		return decimal.Zero, err
	}

	seatPrices := make(map[string]decimal.Decimal)
	deluxeOffset, _ := decimal.NewFromFloat64(constants.DELUXE_OFFSET)
	for _, seat := range seatMap {
		var price decimal.Decimal
		if seat.SeatType == "Deluxe" {
			price, _ = show.Cost.Add(deluxeOffset)
		} else {
			price = show.Cost
		}
		seatPrices[seat.SeatNumber] = price
	}

	var totalPrice decimal.Decimal
	for _, seatNumber := range seatNumbers {
		price, exists := seatPrices[seatNumber]
		if !exists {
			return decimal.Zero, utils.NewBadRequestError("INVALID_SEAT", fmt.Sprintf("Seat %s not found in seat map", seatNumber), nil)
		}
		totalPrice, _ = totalPrice.Add(price)
	}

	return totalPrice, nil
}
