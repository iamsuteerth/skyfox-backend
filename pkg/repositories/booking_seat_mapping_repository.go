package repositories

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type BookingSeatMappingRepository interface {
	CreateMappings(ctx context.Context, bookingId int, seatNumbers []string) error
	GetSeatsByBookingId(ctx context.Context, bookingId int) ([]string, error)
	CheckSeatsAvailability(ctx context.Context, showId int, seatNumbers []string) (bool, error)
}

type bookingSeatMappingRepository struct {
	db *pgxpool.Pool
}

func NewBookingSeatMappingRepository(db *pgxpool.Pool) BookingSeatMappingRepository {
	return &bookingSeatMappingRepository{db: db}
}

func (repo *bookingSeatMappingRepository) CreateMappings(ctx context.Context, bookingId int, seatNumbers []string) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to start database transaction", err)
	}
	
	defer tx.Rollback(ctx)

	for _, seatNumber := range seatNumbers {
		query := `
			INSERT INTO booking_seat_mapping (booking_id, seat_number)
			VALUES ($1, $2)
		`

		_, err := tx.Exec(ctx, query, bookingId, seatNumber)
		if err != nil {
			log.Error().Err(err).Int("bookingId", bookingId).Str("seatNumber", seatNumber).Msg("Failed to create booking seat mapping")
			return utils.NewInternalServerError("DATABASE_ERROR", "Failed to map seats to booking", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to complete seat mapping", err)
	}

	return nil
}

func (repo *bookingSeatMappingRepository) GetSeatsByBookingId(ctx context.Context, bookingId int) ([]string, error) {
	query := `
		SELECT seat_number
		FROM booking_seat_mapping
		WHERE booking_id = $1
		ORDER BY
		  regexp_replace(seat_number, '[0-9]+', '', 'g'),
		  CAST(regexp_replace(seat_number, '[^0-9]+', '', 'g') AS INTEGER)
	`

	rows, err := repo.db.Query(ctx, query, bookingId)
	if err != nil {
		log.Error().Err(err).Int("bookingId", bookingId).Msg("Failed to get seats for booking")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve booking seats", err)
	}
	defer rows.Close()

	var seatNumbers []string
	for rows.Next() {
		var seatNumber string
		if err := rows.Scan(&seatNumber); err != nil {
			log.Error().Err(err).Msg("Error scanning seat number")
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to read seat data", err)
		}
		seatNumbers = append(seatNumbers, seatNumber)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over seat rows")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to process seat data", err)
	}

	return seatNumbers, nil
}

func (repo *bookingSeatMappingRepository) CheckSeatsAvailability(ctx context.Context, showId int, seatNumbers []string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM booking_seat_mapping bsm
		JOIN booking b ON bsm.booking_id = b.id
		WHERE b.show_id = $1 
		AND bsm.seat_number = ANY($2)
		AND b.status IN ('Pending', 'Confirmed', 'CheckedIn')
	`

	var count int
	err := repo.db.QueryRow(ctx, query, showId, seatNumbers).Scan(&count)
	if err != nil {
		log.Error().Err(err).Int("showId", showId).Strs("seatNumbers", seatNumbers).Msg("Failed to check seat availability")
		return false, utils.NewInternalServerError("DATABASE_ERROR", "Failed to check seat availability", err)
	}

	return count == 0, nil
}
