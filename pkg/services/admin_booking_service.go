package services

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/constants"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type AdminBookingService interface {
	CreateAdminBooking(ctx context.Context, req request.AdminBookingRequest) (*response.AdminBookingResponse, error)
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

func (s *adminBookingService) CreateAdminBooking(ctx context.Context, req request.AdminBookingRequest) (*response.AdminBookingResponse, error) {
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

	var amountPaid float64

	expectedPrice, err := s.calculateTotalPrice(ctx, show, req.SeatNumbers)
	if err != nil {
		log.Error().Err(err).Int("showID", req.ShowID).Strs("seatNumbers", req.SeatNumbers).Msg("Failed to calculate expected price")
		return nil, err
	}
	fmt.Printf("Expected Price %f Actual Price %f", expectedPrice, req.AmountPaid)

	if math.Abs(req.AmountPaid-expectedPrice) > 0.01 {
		log.Warn().Float64("submitted", req.AmountPaid).Float64("expected", expectedPrice).Msg("Amount mismatch")
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
		return nil, err
	}

	if err := s.bookingSeatMappingRepo.CreateMappings(ctx, booking.Id, req.SeatNumbers); err != nil {
		log.Error().Err(err).Int("bookingId", booking.Id).Strs("seatNumbers", req.SeatNumbers).Msg("Failed to create seat mappings")
		return nil, err
	}

	bookingResponse := &response.AdminBookingResponse{
		BookingID:    booking.Id,
		ShowID:       booking.ShowId,
		CustomerName: customer.Name,
		PhoneNumber:  customer.Number,
		SeatNumbers:  req.SeatNumbers,
		AmountPaid:   booking.AmountPaid,
		PaymentType:  booking.PaymentType,
		BookingTime:  booking.BookingTime,
		Status:       booking.Status,
	}

	return bookingResponse, nil
}

func (s *adminBookingService) calculateTotalPrice(ctx context.Context, show *models.Show, seatNumbers []string) (float64, error) {
	seatMap, err := s.showRepo.GetSeatMapForShow(ctx, show.Id)
	if err != nil {
		log.Error().Err(err).Int("showID", show.Id).Msg("Failed to get seat map for price calculation")
		return 0, err
	}

	seatPrices := make(map[string]float64)
	for _, seat := range seatMap {
		var price float64
		if seat.SeatType == "Deluxe" {
			price = show.Cost + constants.DELUXE_OFFSET
		} else {
			price = show.Cost
		}
		seatPrices[seat.SeatNumber] = price
	}

	var totalPrice float64
	for _, seatNumber := range seatNumbers {
		price, exists := seatPrices[seatNumber]
		if !exists {
			return 0, utils.NewBadRequestError("INVALID_SEAT", fmt.Sprintf("Seat %s not found in seat map", seatNumber), nil)
		}
		totalPrice += price
	}

	return totalPrice, nil
}
