package services

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/constants"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
)

type BookingService interface {
	GetSeatMapForShow(ctx context.Context, showID int) ([]models.SeatMapEntry, error)
}

type bookingService struct {
	showRepo repositories.ShowRepository
}

func NewBookingService(
	showRepo repositories.ShowRepository,
) BookingService {
	return &bookingService{
		showRepo: showRepo,
	}
}

func (s *bookingService) GetSeatMapForShow(ctx context.Context, showID int) ([]models.SeatMapEntry, error) {
	show, err := s.showRepo.FindById(ctx, showID)
	if err != nil {
		return nil, err
	}

	seatMap, err := s.showRepo.GetSeatMapForShow(ctx, showID)
	if err != nil {
		return nil, err
	}

	for i := range seatMap {
		if seatMap[i].SeatType == "Deluxe" {
			seatMap[i].Price = show.Cost + constants.DELUXE_OFFSET
		} else {
			seatMap[i].Price = show.Cost
		}
	}

	return seatMap, nil
}
