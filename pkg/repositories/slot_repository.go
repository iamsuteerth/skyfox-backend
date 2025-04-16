package repositories

import (
	"context"
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type SlotRepository interface {
	GetAvailableSlotsForDate(ctx context.Context, date time.Time) ([]models.Slot, error)
	GetSlotById(ctx context.Context, slotId int) (*models.Slot, error)
	IsSlotAvailableForDate(ctx context.Context, slotId int, date time.Time) (bool, error)
}

type slotRepository struct {
	db *pgxpool.Pool
}

func NewSlotRepository(db *pgxpool.Pool) SlotRepository {
	return &slotRepository{db: db}
}


func (repo *slotRepository) GetAvailableSlotsForDate(ctx context.Context, date time.Time) ([]models.Slot, error) {
	query := `
		SELECT s.id, s.name, s.start_time, s.end_time
		FROM slot s
		WHERE NOT EXISTS (
			SELECT 1 FROM show sh
			WHERE sh.date = $1 AND sh.slot_id = s.id
		)
		ORDER BY s.id
	`

	rows, err := repo.db.Query(ctx, query, date)
	if err != nil {
		log.Error().Err(err).Time("date", date).Msg("Failed to query available slots for date")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve available slots", err)
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var slot models.Slot
		if err := rows.Scan(&slot.Id, &slot.Name, &slot.StartTime, &slot.EndTime); err != nil {
			log.Error().Err(err).Msg("Error scanning available slot row")
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to scan slot data", err)
		}
		slots = append(slots, slot)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over available slot rows")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to iterate over available slots", err)
	}

	return slots, nil
}

func (repo *slotRepository) GetSlotById(ctx context.Context, slotId int) (*models.Slot, error) {
	query := `SELECT id, name, start_time, end_time FROM slot WHERE id = $1`

	var slot models.Slot
	err := repo.db.QueryRow(ctx, query, slotId).Scan(
		&slot.Id,
		&slot.Name,
		&slot.StartTime,
		&slot.EndTime,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Int("slotId", slotId).Msg("Error fetching slot by ID")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching slot details", err)
	}

	return &slot, nil
}

func (repo *slotRepository) IsSlotAvailableForDate(ctx context.Context, slotId int, date time.Time) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM show 
        WHERE slot_id = $1 AND date = $2
    `

	var count int
	err := repo.db.QueryRow(ctx, query, slotId, date).Scan(&count)
	if err != nil {
		log.Error().Err(err).Int("slotId", slotId).Time("date", date).Msg("Error checking slot availability")
		return false, utils.NewInternalServerError("DATABASE_ERROR", "Error checking slot availability", err)
	}

	return count == 0, nil
}
