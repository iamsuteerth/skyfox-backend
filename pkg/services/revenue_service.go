package services

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/govalues/decimal"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	movieservice "github.com/iamsuteerth/skyfox-backend/pkg/movie-service"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/rs/zerolog/log"
)

type RevenueService interface {
	GetRevenue(ctx context.Context, req request.RevenueDashboardRequest) (*response.RevenueDashboardResponse, error)
}

type revenueService struct {
	bookingRepo  repositories.BookingRepository
	showRepo     repositories.ShowRepository
	slotRepo     repositories.SlotRepository
	movieService movieservice.MovieService
}

func NewRevenueService(
	bookingRepo repositories.BookingRepository,
	showRepo repositories.ShowRepository,
	slotRepo repositories.SlotRepository,
	movieService movieservice.MovieService,
) RevenueService {
	return &revenueService{
		bookingRepo:  bookingRepo,
		showRepo:     showRepo,
		slotRepo:     slotRepo,
		movieService: movieService,
	}
}

func (s *revenueService) GetRevenue(ctx context.Context, req request.RevenueDashboardRequest) (*response.RevenueDashboardResponse, error) {
	if req.Timeframe == "all" {
		req.Timeframe = ""
	}

	bookings, err := s.bookingRepo.FindBookingsByStatus(ctx, []string{"Confirmed", "CheckedIn"})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch bookings for revenue calculation")
		return nil, err
	}

	filteredBookings := s.filterBookings(ctx, bookings, req)

	totalRevenue, meanRevenue, medianRevenue, totalBookings, totalSeats := s.calculateOverallStats(filteredBookings)

	groups, err := s.groupBookingData(ctx, filteredBookings, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to group booking data for revenue analysis")
		return nil, err
	}

	totalRevenueFloat, _ := totalRevenue.Float64()
	meanRevenueFloat, _ := meanRevenue.Float64()
	medianRevenueFloat, _ := medianRevenue.Float64()

	return &response.RevenueDashboardResponse{
		TotalRevenue:     totalRevenueFloat,
		MeanRevenue:      meanRevenueFloat,
		MedianRevenue:    medianRevenueFloat,
		TotalBookings:    totalBookings,
		TotalSeatsBooked: totalSeats,
		Groups:           groups,
	}, nil
}

func (s *revenueService) filterBookings(ctx context.Context, bookings []*models.Booking, req request.RevenueDashboardRequest) []*models.Booking {
	var filtered []*models.Booking

	var startLimit time.Time
	now := time.Now()

	if req.Timeframe != "" {
		switch req.Timeframe {
		case "daily":
			startLimit = now.AddDate(0, 0, -30)
		case "weekly":
			startLimit = now.AddDate(0, 0, -16*7)
		case "monthly":
			startLimit = now.AddDate(0, -12, 0)
		}
	}

	for _, booking := range bookings {
		if booking.Status != "Confirmed" && booking.Status != "CheckedIn" {
			continue
		}

		include := true

		if req.Timeframe != "" && !startLimit.IsZero() {
			if booking.BookingTime.Before(startLimit) {
				include = false
			}
		}

		if req.Month != nil {
			bookingMonth := int(booking.BookingTime.Month())
			if bookingMonth != *req.Month {
				include = false
			}
		}

		if req.Year != nil {
			bookingYear := booking.BookingTime.Year()
			if bookingYear != *req.Year {
				include = false
			}
		}

		if req.MovieID != nil {
			show, err := s.showRepo.FindById(ctx, booking.ShowId)
			if err != nil || show == nil || show.MovieId != *req.MovieID {
				include = false
			}
		}

		if req.SlotID != nil {
			show, err := s.showRepo.FindById(ctx, booking.ShowId)
			if err != nil || show == nil || show.SlotId != *req.SlotID {
				include = false
			}
		}

		if req.Genre != nil && len(*req.Genre) > 0 {
			show, err := s.showRepo.FindById(ctx, booking.ShowId)
			if err != nil || show == nil {
				include = false
			} else {
				movie, err := s.movieService.GetMovieById(ctx, show.MovieId)
				if err != nil || movie == nil {
					include = false
				} else {
					if !strings.Contains(movie.Genre, *req.Genre) {
						include = false
					}
				}
			}
		}

		if include {
			filtered = append(filtered, booking)
		}
	}

	return filtered
}

func (s *revenueService) calculateOverallStats(bookings []*models.Booking) (decimal.Decimal, decimal.Decimal, decimal.Decimal, int, int) {
	if len(bookings) == 0 {
		return decimal.Zero, decimal.Zero, decimal.Zero, 0, 0
	}

	var totalRevenue decimal.Decimal
	var amounts []decimal.Decimal
	totalSeats := 0

	for _, booking := range bookings {
		totalRevenue, _ = totalRevenue.Add(booking.AmountPaid)
		amounts = append(amounts, booking.AmountPaid)
		totalSeats += booking.NoOfSeats
	}

	lenOfBookings, _ := decimal.NewFromInt64(int64(len(bookings)), 0, 0)

	meanRevenue, _ := totalRevenue.Quo(lenOfBookings)
	medianRevenue := calculateMedian(amounts)

	return totalRevenue, meanRevenue, medianRevenue, len(bookings), totalSeats
}

func calculateMedian(values []decimal.Decimal) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].Cmp(values[j]) < 0
	})

	middle := len(values) / 2
	if len(values)%2 == 0 {
		sum, _ := values[middle-1].Add(values[middle])
		two := decimal.MustParse("2")
		median, _ := sum.Quo(two)
		return median
	}

	return values[middle]
}

func (s *revenueService) groupBookingData(ctx context.Context, bookings []*models.Booking, req request.RevenueDashboardRequest) ([]response.RevenueGroupStats, error) {
	result := []response.RevenueGroupStats{}

	filterMap := make(map[string]struct {
		name  string
		value string
	})

	if req.Genre != nil {
		filterMap["genre"] = struct {
			name  string
			value string
		}{name: "genre", value: *req.Genre}
	}

	if req.MovieID != nil {
		filterMap["movie_id"] = struct {
			name  string
			value string
		}{name: "movie", value: ""}
	}

	if req.SlotID != nil {
		filterMap["slot_id"] = struct {
			name  string
			value string
		}{name: "slot", value: ""}
	}

	if req.Timeframe != "" {
		filterMap["timeframe"] = struct {
			name  string
			value string
		}{name: "timeframe", value: req.Timeframe}
	} else {
		if req.Month != nil {
			monthName := time.Month(*req.Month).String()
			filterMap["month"] = struct {
				name  string
				value string
			}{name: "month", value: monthName}
		}

		if req.Year != nil {
			yearStr := strconv.Itoa(*req.Year)
			filterMap["year"] = struct {
				name  string
				value string
			}{name: "year", value: yearStr}
		}
	}

	if len(filterMap) == 0 {
		filterMap["all"] = struct {
			name  string
			value string
		}{name: "all", value: "All"}
		if len(req.ParamOrder) == 0 {
			req.ParamOrder = []string{"all"}
		}
	}

	groups := make(map[string][]*models.Booking)

	showCache := make(map[int]*models.Show)
	movieCache := make(map[string]*models.Movie)
	slotCache := make(map[int]*models.Slot)

	showIDs := make(map[int]bool)
	for _, booking := range bookings {
		showIDs[booking.ShowId] = true
	}

	for showID := range showIDs {
		show, err := s.showRepo.FindById(ctx, showID)
		if err != nil || show == nil {
			continue
		}

		showCache[showID] = show

		if req.SlotID != nil || contains(req.ParamOrder, "slot_id") {
			slot, err := s.slotRepo.GetSlotById(ctx, show.SlotId)
			if err == nil && slot != nil {
				slotCache[show.SlotId] = slot
			}
		}

		if req.MovieID != nil || req.Genre != nil || contains(req.ParamOrder, "movie_id") {
			movie, err := s.movieService.GetMovieById(ctx, show.MovieId)
			if err == nil && movie != nil {
				movieCache[show.MovieId] = movie
			}
		}
	}

	for _, booking := range bookings {
		show, exists := showCache[booking.ShowId]
		if !exists {
			continue
		}

		var keyParts []string

		for _, paramName := range req.ParamOrder {
			filter, exists := filterMap[paramName]
			if !exists {
				continue
			}

			var labelValue string
			switch filter.name {
			case "genre":
				labelValue = filter.value

			case "movie":
				movie, exists := movieCache[show.MovieId]
				if !exists {
					continue
				}
				labelValue = movie.Name

			case "slot":
				slot, exists := slotCache[show.SlotId]
				if !exists {
					continue
				}
				labelValue = slot.Name

			case "timeframe":
				switch filter.value {
				case "daily":
					labelValue = booking.BookingTime.Format("2006-01-02")
				case "weekly":
					year, week := booking.BookingTime.ISOWeek()
					labelValue = fmt.Sprintf("%d-W%02d", year, week)
				case "monthly":
					labelValue = booking.BookingTime.Format("2006-01")
				case "yearly":
					labelValue = booking.BookingTime.Format("2006")
				default:
					labelValue = "All"
				}

			case "month":
				labelValue = filter.value

			case "year":
				labelValue = filter.value

			case "all":
				labelValue = "All"
			}

			keyParts = append(keyParts, labelValue)
		}

		if len(keyParts) == 0 {
			keyParts = append(keyParts, "All")
		}

		key := strings.Join(keyParts, ";")
		groups[key] = append(groups[key], booking)
	}

	for label, groupBookings := range groups {
		totalRev, meanRev, medianRev, totalBook, totalSeats := s.calculateOverallStats(groupBookings)

		totalRevenueFloat, _ := totalRev.Float64()
		meanRevenueFloat, _ := meanRev.Float64()
		medianRevenueFloat, _ := medianRev.Float64()

		result = append(result, response.RevenueGroupStats{
			Label:            label,
			TotalRevenue:     totalRevenueFloat,
			MeanRevenue:      meanRevenueFloat,
			MedianRevenue:    medianRevenueFloat,
			TotalBookings:    totalBook,
			TotalSeatsBooked: totalSeats,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Label < result[j].Label
	})

	return result, nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
