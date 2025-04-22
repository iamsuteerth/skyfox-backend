package repositories

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type BookingRepository interface {
	CreateAdminBooking(ctx context.Context, booking *models.Booking) error
	GetBookingById(ctx context.Context, id int) (*models.Booking, error)
	BookedSeatsByShow(ctx context.Context, showId int) int
}

type bookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) BookingRepository {
	return &bookingRepository{db: db}
}

func (repo *bookingRepository) CreateAdminBooking(ctx context.Context, booking *models.Booking) error {
	query := `
		INSERT INTO booking (
			date, show_id, customer_id, no_of_seats, 
			amount_paid, status, payment_type
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, booking_time
	`

	err := repo.db.QueryRow(ctx, query,
		booking.Date,
		booking.ShowId,
		booking.CustomerId,
		booking.NoOfSeats,
		booking.AmountPaid,
		booking.Status,
		booking.PaymentType,
	).Scan(&booking.Id, &booking.BookingTime)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create booking")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to create booking", err)
	}

	return nil
}

func (repo *bookingRepository) GetBookingById(ctx context.Context, id int) (*models.Booking, error) {
	query := `
		SELECT 
			id, date, show_id, customer_id, customer_username, 
			no_of_seats, amount_paid, status, booking_time, payment_type
		FROM booking
		WHERE id = $1
	`

	var booking models.Booking
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&booking.Id,
		&booking.Date,
		&booking.ShowId,
		&booking.CustomerId,
		&booking.CustomerUsername,
		&booking.NoOfSeats,
		&booking.AmountPaid,
		&booking.Status,
		&booking.BookingTime,
		&booking.PaymentType,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.NewNotFoundError("BOOKING_NOT_FOUND", "Booking not found", nil)
		}
		log.Error().Err(err).Int("id", id).Msg("Failed to get booking by ID")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve booking", err)
	}

	return &booking, nil
}

func (repo *bookingRepository) BookedSeatsByShow(ctx context.Context, showId int) int {
	query := `
		SELECT COALESCE(SUM(no_of_seats), 0)
		FROM booking
		WHERE show_id = $1 AND status IN ('Confirmed', 'CheckedIn')
	`

	var count int
	err := repo.db.QueryRow(ctx, query, showId).Scan(&count)

	if err != nil {
		log.Error().Err(err).Int("showId", showId).Msg("Failed to count booked seats")
		return 0
	}

	return count
}
