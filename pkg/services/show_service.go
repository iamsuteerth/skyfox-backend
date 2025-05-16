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
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	movieservice "github.com/iamsuteerth/skyfox-backend/pkg/movie-service"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type ShowService interface {
	GetShows(ctx context.Context, date time.Time) ([]models.Show, error)
	GetShowById(ctx context.Context, id int) (*models.Show, error)
	GetMovieById(ctx context.Context, id string) (*models.Movie, error)
	GetMovies(ctx context.Context) ([]*models.Movie, error)
	CreateShow(ctx context.Context, showRequest request.ShowRequest) (*models.Show, error)
	AvailableSeats(ctx context.Context, showId int) int
}

type showService struct {
	showRepo     repositories.ShowRepository
	bookingRepo  repositories.BookingRepository
	movieService movieservice.MovieService
	slotRepo     repositories.SlotRepository
}

func NewShowService(
	showRepo repositories.ShowRepository,
	bookingRepo repositories.BookingRepository,
	movieService movieservice.MovieService,
	slotRepo repositories.SlotRepository,
) ShowService {
	return &showService{
		showRepo:     showRepo,
		bookingRepo:  bookingRepo,
		movieService: movieService,
		slotRepo:     slotRepo,
	}
}

func (s *showService) GetShows(ctx context.Context, date time.Time) ([]models.Show, error) {
	shows, err := s.showRepo.GetAllShowsOn(ctx, date)
	if err != nil {
		log.Error().Err(err).Str("date", date.String()).Msg("Failed to get shows for date")
		return nil, err
	}
	return shows, nil
}

func (s *showService) GetShowById(ctx context.Context, id int) (*models.Show, error) {
	show, err := s.showRepo.FindById(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Failed to get shows for id : %d", id))
		return nil, err
	}
	return show, nil
}

func (s *showService) GetMovies(ctx context.Context) ([]*models.Movie, error) {
	movies, err := s.movieService.GetAllMovies(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get movies from movie service")
		return nil, err
	}
	return movies, nil
}

func (s *showService) GetMovieById(ctx context.Context, id string) (*models.Movie, error) {
	movie, err := s.movieService.GetMovieById(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("movieId", id).Msg("Failed to get movie by ID")
		return nil, err
	}
	return movie, nil
}

func (s *showService) CreateShow(ctx context.Context, showRequest request.ShowRequest) (*models.Show, error) {
	if showRequest.Cost.Cmp(decimal.Zero) != 1 {
		return nil, utils.NewBadRequestError("INVALID_COST", "The cost must be greater than 0", nil)
	}

	maxCost, _ := decimal.NewFromInt64(3000, 0, 0)

	if showRequest.Cost.Cmp(maxCost) == 1  {
		return nil, utils.NewBadRequestError("INVALID_COST", "The cost must be less than or equal to 3000", nil)
	}

	showDate, err := time.Parse("2006-01-02", showRequest.Date)
	if err != nil {
		return nil, utils.NewBadRequestError("INVALID_DATE_FORMAT", "The date format is not valid", err)
	}

	now := time.Now()

	slot, err := s.slotRepo.GetSlotById(ctx, showRequest.SlotId)
	if err != nil {
		log.Error().Err(err).Int("slotId", showRequest.SlotId).Msg("Error fetching slot details")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching slot details", err)
	}

	if slot == nil {
		return nil, utils.NewBadRequestError("INVALID_SLOT", "The selected slot does not exist", nil)
	}

	startTimeParts := strings.Split(slot.StartTime, ":")
	if len(startTimeParts) < 2 {
		return nil, utils.NewInternalServerError("INVALID_SLOT_TIME", "Invalid slot start time format", nil)
	}

	hour, _ := strconv.Atoi(startTimeParts[0])
	minute, _ := strconv.Atoi(startTimeParts[1])

	showDateTime := time.Date(
		showDate.Year(),
		showDate.Month(),
		showDate.Day(),
		hour,
		minute,
		0,
		0,
		now.Location(),
	)

	if showDateTime.Before(now) {
		return nil, utils.NewBadRequestError("PAST_DATETIME", "The show cannot be scheduled for a time in the past", nil)
	}

	isSlotAvailable, err := s.slotRepo.IsSlotAvailableForDate(ctx, showRequest.SlotId, showDate)
	if err != nil {
		log.Error().Err(err).Msg("Error checking slot availability")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error checking slot availability", err)
	}

	if !isSlotAvailable {
		return nil, utils.NewBadRequestError(
			"SLOT_NOT_AVAILABLE",
			fmt.Sprintf("The selected slot is not available on %s", showDate.Format("2006-01-02")),
			nil,
		)
	}

	movie, err := s.movieService.GetMovieById(ctx, showRequest.MovieId)
	if err != nil {
		log.Error().Err(err).Str("movieId", showRequest.MovieId).Msg("Error verifying movie existence")
		return nil, utils.NewBadRequestError("INVALID_MOVIE", "The selected movie does not exist", err)
	}

	if movie == nil {
		return nil, utils.NewBadRequestError("INVALID_MOVIE", "The selected movie does not exist", nil)
	}

	show := models.Show{
		MovieId: showRequest.MovieId,
		Date:    showDate,
		SlotId:  showRequest.SlotId,
		Cost:    showRequest.Cost,
	}

	if err := s.showRepo.Create(ctx, &show); err != nil {
		log.Error().Err(err).Msg("Failed to create show")
		return nil, err
	}

	completeShow, err := s.showRepo.FindById(ctx, show.Id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch show information for created show. Request not processed.")
		return nil, err
	}

	return completeShow, nil
}

func (s *showService) AvailableSeats(ctx context.Context, showId int) int {
	bookedSeats := s.bookingRepo.BookedSeatsByShow(ctx, showId)
	return constants.TOTAL_NO_OF_SEATS - bookedSeats
}
