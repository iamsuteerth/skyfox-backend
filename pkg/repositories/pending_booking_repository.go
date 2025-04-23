package repositories

import (
	"context"
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PendingBookingRepository interface {
	TrackPendingBooking(ctx context.Context, bookingId int, expirationTime time.Time) error
	GetExpirationTime(ctx context.Context, bookingId int) (*time.Time, error)
	RemoveTracker(ctx context.Context, bookingId int) error
	GetExpiredBookingIds(ctx context.Context, currentTime time.Time) ([]int, error)
}

type pendingBookingRepository struct {
	db *pgxpool.Pool
}

func NewPendingBookingRepository(db *pgxpool.Pool) PendingBookingRepository {
	return &pendingBookingRepository{db: db}
}

func (repo *pendingBookingRepository) TrackPendingBooking(ctx context.Context, bookingId int, expirationTime time.Time) error {
	query := `
        INSERT INTO pending_booking_tracker (booking_id, expiration_time)
        VALUES ($1, $2)
    `

	_, err := repo.db.Exec(ctx, query, bookingId, expirationTime)
	if err != nil {
		log.Error().Err(err).Int("bookingId", bookingId).Msg("Failed to track pending booking")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to track pending booking", err)
	}

	return nil
}

func (repo *pendingBookingRepository) GetExpirationTime(ctx context.Context, bookingId int) (*time.Time, error) {
	query := `
        SELECT expiration_time
        FROM pending_booking_tracker
        WHERE booking_id = $1
    `

	var expirationTime time.Time
	err := repo.db.QueryRow(ctx, query, bookingId).Scan(&expirationTime)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil 
		}
		log.Error().Err(err).Int("bookingId", bookingId).Msg("Failed to get expiration time")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve expiration time", err)
	}

	return &expirationTime, nil
}

func (repo *pendingBookingRepository) RemoveTracker(ctx context.Context, bookingId int) error {
	query := `
        DELETE FROM pending_booking_tracker
        WHERE booking_id = $1
    `

	_, err := repo.db.Exec(ctx, query, bookingId)
	if err != nil {
		log.Error().Err(err).Int("bookingId", bookingId).Msg("Failed to remove pending booking tracker")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to remove pending booking tracker", err)
	}

	return nil
}

func (repo *pendingBookingRepository) GetExpiredBookingIds(ctx context.Context, currentTime time.Time) ([]int, error) {
	query := `
        SELECT booking_id
        FROM pending_booking_tracker
        WHERE expiration_time < $1
    `

	rows, err := repo.db.Query(ctx, query, currentTime)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query expired bookings")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve expired bookings", err)
	}
	defer rows.Close()

	var bookingIds []int
	for rows.Next() {
		var bookingId int
		if err := rows.Scan(&bookingId); err != nil {
			log.Error().Err(err).Msg("Error scanning booking ID")
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to read booking data", err)
		}
		bookingIds = append(bookingIds, bookingId)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over booking IDs")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to process booking data", err)
	}

	return bookingIds, nil
}
