package services

import (
	"context"
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/rs/zerolog/log"
)

type SlotService interface {
	GetAvailableSlots(ctx context.Context, date time.Time) ([]models.Slot, error)
}

type slotService struct {
	slotRepo repositories.SlotRepository
}

func NewSlotService(slotRepo repositories.SlotRepository) SlotService {
	return &slotService{
		slotRepo: slotRepo,
	}
}

func (s *slotService) GetAvailableSlots(ctx context.Context, date time.Time) ([]models.Slot, error) {
	slots, err := s.slotRepo.GetAvailableSlotsForDate(ctx, date)
	if err != nil {
		log.Error().Err(err).Time("date", date).Msg("Failed to get available slots for date")
		return nil, err
	}
	return slots, nil
}
