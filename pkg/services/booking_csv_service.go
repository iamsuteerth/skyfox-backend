package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type BookingCSVService interface {
	WriteBookingsCSV(ctx context.Context, w io.Writer, month, year *int) error
}

type bookingCSVService struct {
	bookingRepo             repositories.BookingRepository
	showRepo                repositories.ShowRepository
	adminBookedCustomerRepo repositories.AdminBookedCustomerRepository
	skyCustomerRepo         repositories.SkyCustomerRepository
}

func NewBookingCSVService(
	bookingRepo repositories.BookingRepository,
	showRepo repositories.ShowRepository,
	adminBookedCustomerRepo repositories.AdminBookedCustomerRepository,
	skyCustomerRepo repositories.SkyCustomerRepository,
) BookingCSVService {
	return &bookingCSVService{
		bookingRepo:             bookingRepo,
		showRepo:                showRepo,
		adminBookedCustomerRepo: adminBookedCustomerRepo,
		skyCustomerRepo:         skyCustomerRepo,
	}
}

func (s *bookingCSVService) WriteBookingsCSV(ctx context.Context, w io.Writer, month, year *int) error {
	bookings, err := s.bookingRepo.FindBookingsByStatusAndDate(ctx, []string{"Confirmed", "CheckedIn"}, month, year)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch bookings for CSV export")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve bookings", err)
	}

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	header := []string{
		"Booking ID",
		"Show ID",
		"Show Date",
		"Customer Name",
		"Phone Number",
		"Number of Seats",
		"Amount Paid",
		"Payment Type",
		"Booking Time",
		"Status",
	}

	if err := csvWriter.Write(header); err != nil {
		log.Error().Err(err).Msg("Failed to write CSV header")
		return utils.NewInternalServerError("CSV_ERROR", "Failed to write CSV header", err)
	}

	if len(bookings) == 0 {
		return nil
	}

	showCache := make(map[int]*models.Show)
	adminCustomerCache := make(map[int]*models.AdminBookedCustomer)
	skyCustomerCache := make(map[string]*models.SkyCustomer)

	showIDs := make(map[int]bool)
	customerIDs := make(map[int]bool)
	customerUsernames := make(map[string]bool)

	for _, booking := range bookings {
		showIDs[booking.ShowId] = true

		if booking.CustomerId != nil {
			customerIDs[*booking.CustomerId] = true
		} else if booking.CustomerUsername != nil {
			customerUsernames[*booking.CustomerUsername] = true
		}
	}

	for showID := range showIDs {
		show, err := s.showRepo.FindById(ctx, showID)
		if err == nil && show != nil {
			showCache[showID] = show
		}
	}

	for customerID := range customerIDs {
		customer, err := s.adminBookedCustomerRepo.FindById(ctx, customerID)
		if err == nil && customer != nil {
			adminCustomerCache[customerID] = customer
		}
	}

	for username := range customerUsernames {
		customer, err := s.skyCustomerRepo.FindByUsername(ctx, username)
		if err == nil && customer != nil {
			skyCustomerCache[username] = customer
		}
	}

	for _, booking := range bookings {
		var customerName, phoneNumber string
		if booking.CustomerId != nil {
			if customer, exists := adminCustomerCache[*booking.CustomerId]; exists {
				customerName = customer.Name
				phoneNumber = customer.Number
			}
		} else if booking.CustomerUsername != nil {
			if customer, exists := skyCustomerCache[*booking.CustomerUsername]; exists {
				customerName = customer.Name
				phoneNumber = customer.Number
			}
		}

		row := []string{
			fmt.Sprintf("%d", booking.Id),
			fmt.Sprintf("%d", booking.ShowId),
			booking.Date.Format("2006-01-02"),
			customerName,
			phoneNumber,
			fmt.Sprintf("%d", booking.NoOfSeats),
			fmt.Sprintf("%.2f", booking.AmountPaid),
			booking.PaymentType,
			booking.BookingTime.Format("2006-01-02 15:04:05"),
			booking.Status,
		}

		if err := csvWriter.Write(row); err != nil {
			log.Error().Err(err).Int("booking_id", booking.Id).Msg("Failed to write CSV row")
			return utils.NewInternalServerError("CSV_ERROR", "Failed to write booking data", err)
		}
	}

	return nil
}
