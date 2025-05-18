package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type ShowRepository interface {
	Create(ctx context.Context, show *models.Show) error
	GetAllShowsOn(ctx context.Context, date time.Time) ([]models.Show, error)
	FindById(ctx context.Context, id int) (*models.Show, error)
	GetSeatMapForShow(ctx context.Context, showID int) ([]models.SeatMapEntry, error)
}

type showRepository struct {
	db *pgxpool.Pool
}

func NewShowRepository(db *pgxpool.Pool) ShowRepository {
	return &showRepository{db: db}
}

func (repo *showRepository) Create(ctx context.Context, show *models.Show) error {
	query := `
        INSERT INTO show (movie_id, date, slot_id, cost)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	err := repo.db.QueryRow(ctx, query,
		show.MovieId,
		show.Date,
		show.SlotId,
		show.Cost,
	).Scan(&show.Id)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create show")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to create show", err)
	}

	return nil
}

func (repo *showRepository) GetAllShowsOn(ctx context.Context, date time.Time) ([]models.Show, error) {
	query := `
        SELECT s.id, s.movie_id, s.date, s.slot_id, s.cost, 
               sl.id, sl.name, sl.start_time, sl.end_time
        FROM show s
        JOIN slot sl ON s.slot_id = sl.id
        WHERE s.date = $1
    `

	rows, err := repo.db.Query(ctx, query, date)
	if err != nil {
		log.Error().Err(err).Str("date", date.String()).Msg("Failed to query shows for date")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve shows", err)
	}
	defer rows.Close()

	var shows []models.Show
	for rows.Next() {
		var show models.Show
		var slot models.Slot

		err := rows.Scan(
			&show.Id,
			&show.MovieId,
			&show.Date,
			&show.SlotId,
			&show.Cost,
			&slot.Id,
			&slot.Name,
			&slot.StartTime,
			&slot.EndTime,
		)

		if err != nil {
			log.Error().Err(err).Msg("Error scanning show row")
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to scan show data", err)
		}

		show.Slot = slot
		shows = append(shows, show)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over show rows")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to iterate over shows", err)
	}

	return shows, nil
}

func (repo *showRepository) FindById(ctx context.Context, id int) (*models.Show, error) {
	query := `
        SELECT s.id, s.movie_id, s.date, s.slot_id, s.cost, 
               sl.id, sl.name, sl.start_time, sl.end_time
        FROM show s
        JOIN slot sl ON s.slot_id = sl.id
        WHERE s.id = $1
    `

	var show models.Show
	var slot models.Slot

	err := repo.db.QueryRow(ctx, query, id).Scan(
		&show.Id,
		&show.MovieId,
		&show.Date,
		&show.SlotId,
		&show.Cost,
		&slot.Id,
		&slot.Name,
		&slot.StartTime,
		&slot.EndTime,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.NewNotFoundError("SHOW_NOT_FOUND", fmt.Sprintf("Show not found for id: %d", id), nil)
		}
		log.Error().Err(err).Int("id", id).Msg("Failed to find show by ID")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error retrieving show", err)
	}

	show.Slot = slot
	return &show, nil
}

func (repo *showRepository) GetSeatMapForShow(ctx context.Context, showID int) ([]models.SeatMapEntry, error) {
	query := `
        SELECT 
            s.seat_number,
            SUBSTRING(s.seat_number, 1, 1) AS seat_row,
            SUBSTRING(s.seat_number, 2) AS seat_column,
            s.seat_type,
            0.0 AS price, 
            EXISTS (
                SELECT 1
                FROM booking_seat_mapping bsm
                JOIN booking b ON bsm.booking_id = b.id
                WHERE b.show_id = $1 
                AND bsm.seat_number = s.seat_number
            ) AS occupied
        FROM 
            seat s
        ORDER BY 
            SUBSTRING(s.seat_number, 1, 1) ASC,
            CAST(SUBSTRING(s.seat_number, 2) AS INTEGER) ASC
    `

	rows, err := repo.db.Query(ctx, query, showID)
	if err != nil {
		log.Error().Err(err).Int("show_id", showID).Msg("Failed to query seat map for show")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve seat map", err)
	}
	defer rows.Close()

	var seatMap []models.SeatMapEntry
	for rows.Next() {
		var seat models.SeatMapEntry
		err := rows.Scan(
			&seat.SeatNumber,
			&seat.SeatRow,
			&seat.SeatColumn,
			&seat.SeatType,
			&seat.Price,
			&seat.Occupied,
		)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning seat map row")
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to scan seat data", err)
		}
		seatMap = append(seatMap, seat)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over seat map rows")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to iterate over seat map", err)
	}

	if len(seatMap) == 0 {
		log.Warn().Int("show_id", showID).Msg("No seats found for show")
		return []models.SeatMapEntry{}, nil
	}

	return seatMap, nil
}
