package services

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	movieservice "github.com/iamsuteerth/skyfox-backend/pkg/movie-service"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/rs/zerolog/log"
)

// RevenueService defines methods for obtaining revenue analytics
type RevenueService interface {
	// GetRevenue retrieves revenue data with optional filtering
	GetRevenue(ctx context.Context, req request.RevenueDashboardRequest) (*response.RevenueDashboardResponse, error)
}

type revenueService struct {
	bookingRepo repositories.BookingRepository
	showRepo    repositories.ShowRepository
	slotRepo    repositories.SlotRepository
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
	// Default timeframe to "all" if not specified
	if req.Timeframe == "" {
		req.Timeframe = "all"
	}

	// Step 1: Get all confirmed and checked-in bookings
	// We'll need to extend BookingRepository to support this
	bookings, err := s.bookingRepo.FindBookingsByStatus(ctx, []string{"Confirmed", "CheckedIn"})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch bookings for revenue calculation")
		return nil, err
	}

	// Filter bookings based on request parameters
	filteredBookings := s.filterBookings(ctx, bookings, req)

	// Step 2: Calculate overall statistics
	totalRevenue, meanRevenue, medianRevenue, totalBookings, totalSeats := s.calculateOverallStats(filteredBookings)

	// Step 3: Group data by the requested dimension (timeframe, movie, slot, genre)
	groups, err := s.groupBookingData(ctx, filteredBookings, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to group booking data for revenue analysis")
		return nil, err
	}

	return &response.RevenueDashboardResponse{
		TotalRevenue:     totalRevenue,
		MeanRevenue:      meanRevenue,
		MedianRevenue:    medianRevenue,
		TotalBookings:    totalBookings,
		TotalSeatsBooked: totalSeats,
		Groups:           groups,
	}, nil
}

// filterBookings applies request filters to the list of bookings
func (s *revenueService) filterBookings(ctx context.Context, bookings []*models.Booking, req request.RevenueDashboardRequest) []*models.Booking {
	var filtered []*models.Booking

	for _, booking := range bookings {
		// Skip if status is not Confirmed or CheckedIn (defensive check)
		if booking.Status != "Confirmed" && booking.Status != "CheckedIn" {
			continue
		}

		include := true

		// Filter by month if specified
		if req.Month != nil {
			bookingMonth := int(booking.BookingTime.Month())
			if bookingMonth != *req.Month {
				include = false
			}
		}

		// Filter by year if specified
		if req.Year != nil {
			bookingYear := booking.BookingTime.Year()
			if bookingYear != *req.Year {
				include = false
			}
		}

		// Filter by movie if specified
		if req.MovieID != nil {
			show, err := s.showRepo.FindById(ctx, booking.ShowId)
			if err != nil || show == nil || show.MovieId != *req.MovieID {
				include = false
			}
		}

		// Filter by slot if specified
		if req.SlotID != nil {
			show, err := s.showRepo.FindById(ctx, booking.ShowId)
			if err != nil || show == nil || show.SlotId != *req.SlotID {
				include = false
			}
		}

		// Filter by genre requires movie details
		if req.Genre != nil && len(*req.Genre) > 0 {
			show, err := s.showRepo.FindById(ctx, booking.ShowId)
			if err != nil || show == nil {
				include = false
			} else {
				movie, err := s.movieService.GetMovieById(ctx, show.MovieId)
				if err != nil || movie == nil {
					include = false
				} else {
					// Check if genre is in movie's genre list
					// Assuming genre is comma-separated
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

// calculateOverallStats computes aggregate statistics for the filtered bookings
func (s *revenueService) calculateOverallStats(bookings []*models.Booking) (float64, float64, float64, int, int) {
	if len(bookings) == 0 {
		return 0, 0, 0, 0, 0
	}

	var totalRevenue float64
	var amounts []float64
	totalSeats := 0

	for _, booking := range bookings {
		totalRevenue += booking.AmountPaid
		amounts = append(amounts, booking.AmountPaid)
		totalSeats += booking.NoOfSeats
	}

	meanRevenue := totalRevenue / float64(len(bookings))
	medianRevenue := calculateMedian(amounts)

	return totalRevenue, meanRevenue, medianRevenue, len(bookings), totalSeats
}

// calculateMedian computes the median value from a slice of float64 values
func calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Sort the values
	sort.Float64s(values)

	// Find the median
	middle := len(values) / 2
	if len(values)%2 == 0 {
		// Even number of elements
		return (values[middle-1] + values[middle]) / 2
	}
	// Odd number of elements
	return values[middle]
}

// groupBookingData groups bookings by the requested dimension
func (s *revenueService) groupBookingData(ctx context.Context, bookings []*models.Booking, req request.RevenueDashboardRequest) ([]response.RevenueGroupStats, error) {
	// Define grouping function based on timeframe
	var getGroupKey func(*models.Booking) (string, error)

	switch req.Timeframe {
	case "daily":
		getGroupKey = func(b *models.Booking) (string, error) {
			return b.BookingTime.Format("2006-01-02"), nil
		}
	case "weekly":
		getGroupKey = func(b *models.Booking) (string, error) {
			year, week := b.BookingTime.ISOWeek()
			return fmt.Sprintf("%d-W%02d", year, week), nil
		}
	case "monthly":
		getGroupKey = func(b *models.Booking) (string, error) {
			return b.BookingTime.Format("2006-01"), nil
		}
	case "yearly":
		getGroupKey = func(b *models.Booking) (string, error) {
			return b.BookingTime.Format("2006"), nil
		}
	default:
		// Handle grouping by other dimensions (movie, slot, genre)
		if req.MovieID != nil {
			// If movie_id is specified, we're looking at one movie only
			getGroupKey = func(b *models.Booking) (string, error) {
				show, err := s.showRepo.FindById(ctx, b.ShowId)
				if err != nil {
					return "", err
				}
				movie, err := s.movieService.GetMovieById(ctx, show.MovieId)
				if err != nil {
					return "", err
				}
				return movie.Name, nil
			}
		} else if req.SlotID != nil {
			// If slot_id is specified, we're looking at one slot only
			getGroupKey = func(b *models.Booking) (string, error) {
				show, err := s.showRepo.FindById(ctx, b.ShowId)
				if err != nil {
					return "", err
				}
				slot, err := s.slotRepo.GetSlotById(ctx, show.SlotId)
				if err != nil {
					return "", err
				}
				return slot.Name, nil
			}
		} else if req.Genre != nil {
			// If genre is specified, we're looking at one genre only
			getGroupKey = func(b *models.Booking) (string, error) {
				return *req.Genre, nil
			}
		} else {
			// Default case - aggregate all
			getGroupKey = func(b *models.Booking) (string, error) {
				return "All", nil
			}
		}
	}

	// Group bookings
	groups := make(map[string][]*models.Booking)
	for _, booking := range bookings {
		key, err := getGroupKey(booking)
		if err != nil {
			log.Warn().Err(err).Int("booking_id", booking.Id).Msg("Failed to get group key, skipping")
			continue
		}
		groups[key] = append(groups[key], booking)
	}

	// Calculate stats for each group
	var result []response.RevenueGroupStats
	for label, groupBookings := range groups {
		totalRev, meanRev, medianRev, totalBook, totalSeats := s.calculateOverallStats(groupBookings)

		result = append(result, response.RevenueGroupStats{
			Label:            label,
			TotalRevenue:     totalRev,
			MeanRevenue:      meanRev,
			MedianRevenue:    medianRev,
			TotalBookings:    totalBook,
			TotalSeatsBooked: totalSeats,
		})
	}

	// Sort groups by label for consistency
	sort.Slice(result, func(i, j int) bool {
		return result[i].Label < result[j].Label
	})

	return result, nil
}
