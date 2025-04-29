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
	CreatePendingBooking(ctx context.Context, booking *models.Booking) error
	UpdateBookingStatus(ctx context.Context, bookingID int, status string) error
	DeleteBookingsByIds(ctx context.Context, bookingIds []int) error
	FindByCustomerUsername(ctx context.Context, username string) ([]*models.Booking, error)
    FindLatestByCustomerUsername(ctx context.Context, username string) (*models.Booking, error)
	FindConfirmedBookings(ctx context.Context) ([]*models.Booking, error)
	MarkBookingsCheckedIn(ctx context.Context, bookingIDs []int) (int, error)
	MarkBookingCheckedIn(ctx context.Context, bookingID int) (bool, error)
	FindBookingsByIds(ctx context.Context, bookingIDs []int) ([]*models.Booking, error)
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
		WHERE show_id = $1 AND status IN ('Confirmed', 'CheckedIn', 'Pending')
	`

	var count int
	err := repo.db.QueryRow(ctx, query, showId).Scan(&count)

	if err != nil {
		log.Error().Err(err).Int("showId", showId).Msg("Failed to count booked seats")
		return 0
	}

	return count
}

func (repo *bookingRepository) CreatePendingBooking(ctx context.Context, booking *models.Booking) error {
	query := `
		INSERT INTO booking (
			date, show_id, customer_username, no_of_seats, 
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
		booking.CustomerUsername,
		booking.NoOfSeats,
		booking.AmountPaid,
		booking.Status,
		booking.PaymentType,
	).Scan(&booking.Id, &booking.BookingTime)

	if err != nil {
		log.Error().Err(err).Str("username", *booking.CustomerUsername).Msg("Failed to create pending booking")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to create booking", err)
	}

	return nil
}

func (repo *bookingRepository) UpdateBookingStatus(ctx context.Context, bookingID int, status string) error {
	query := `
		UPDATE booking
		SET status = $1
		WHERE id = $2
	`

	_, err := repo.db.Exec(ctx, query, status, bookingID)
	if err != nil {
		log.Error().Err(err).Int("bookingID", bookingID).Str("status", status).Msg("Failed to update booking status")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to update booking status", err)
	}

	return nil
}

func (repo *bookingRepository) DeleteBookingsByIds(ctx context.Context, bookingIds []int) error {
	if len(bookingIds) == 0 {
		return nil
	}

	query := `
		DELETE FROM booking
		WHERE id = ANY($1)
	`

	_, err := repo.db.Exec(ctx, query, bookingIds)
	if err != nil {
		log.Error().Err(err).Interface("bookingIds", bookingIds).Msg("Failed to delete bookings")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to delete expired bookings", err)
	}

	return nil
}

func (repo *bookingRepository) FindByCustomerUsername(ctx context.Context, username string) ([]*models.Booking, error) {
    const query = `
        SELECT id, show_id, customer_id, customer_username, amount_paid, payment_type, status, booking_time
        FROM booking
        WHERE customer_username = $1
        ORDER BY booking_time DESC
    `
    rows, err := repo.db.Query(ctx, query, username)
    if err != nil {
        log.Error().Err(err).Msg("Failed to query bookings for given username")
        return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve bookings", err)
    }
    defer rows.Close()

    var bookings []*models.Booking
    for rows.Next() {
        var booking models.Booking
        err := rows.Scan(
            &booking.Id,
            &booking.ShowId,
            &booking.CustomerId,
            &booking.CustomerUsername,
            &booking.AmountPaid,
            &booking.PaymentType,
            &booking.Status,
            &booking.BookingTime,
        )
        if err != nil {
            log.Error().Err(err).Msg("Error scanning booking row")
            return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to scan booking data", err)
        }
        bookings = append(bookings, &booking)
    }

    if err := rows.Err(); err != nil {
        log.Error().Err(err).Msg("Error iterating over booking rows")
        return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to iterate over bookings", err)
    }

    return bookings, nil
}

func (repo *bookingRepository) FindLatestByCustomerUsername(ctx context.Context, username string) (*models.Booking, error) {
    const query = `
        SELECT id, show_id, customer_id, customer_username, amount_paid, payment_type, status, booking_time
        FROM booking
        WHERE customer_username = $1
        ORDER BY booking_time DESC
        LIMIT 1
    `
    var booking models.Booking
    err := repo.db.QueryRow(ctx, query, username).Scan(
        &booking.Id,
        &booking.ShowId,
        &booking.CustomerId,
        &booking.CustomerUsername,
        &booking.AmountPaid,
        &booking.PaymentType,
        &booking.Status,
        &booking.BookingTime,
    )
    if err == pgx.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        log.Error().Err(err).Msg("Failed to query latest booking for given username")
        return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve latest booking", err)
    }
    return &booking, nil
}

func (repo *bookingRepository) FindConfirmedBookings(ctx context.Context) ([]*models.Booking, error) {
	const query = `
		SELECT 
			id, date, show_id, customer_id, customer_username, no_of_seats, amount_paid, status, booking_time, payment_type
		FROM booking
		WHERE status = 'Confirmed'
		ORDER BY booking_time DESC
	`
	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch confirmed bookings")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve confirmed bookings", err)
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
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
			log.Error().Err(err).Msg("Error scanning confirmed booking row")
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to scan confirmed booking data", err)
		}
		bookings = append(bookings, &booking)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over confirmed bookings")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to iterate over confirmed bookings", err)
	}

	return bookings, nil
}

func (repo *bookingRepository) MarkBookingsCheckedIn(ctx context.Context, bookingIDs []int) (int, error) {
	if len(bookingIDs) == 0 {
		return 0, nil
	}
	query := `
		UPDATE booking
		SET status = 'CheckedIn'
		WHERE id = ANY($1) AND status = 'Confirmed'
	`
	cmdTag, err := repo.db.Exec(ctx, query, bookingIDs)
	if err != nil {
		log.Error().Err(err).Interface("bookingIDs", bookingIDs).Msg("Failed to bulk update bookings to CheckedIn")
		return 0, utils.NewInternalServerError("DATABASE_ERROR", "Failed to update booking statuses", err)
	}
	return int(cmdTag.RowsAffected()), nil
}

func (repo *bookingRepository) MarkBookingCheckedIn(ctx context.Context, bookingID int) (bool, error) {
	query := `
		UPDATE booking
		SET status = 'CheckedIn'
		WHERE id = $1 AND status = 'Confirmed'
	`
	cmdTag, err := repo.db.Exec(ctx, query, bookingID)
	if err != nil {
		log.Error().Err(err).Int("bookingID", bookingID).Msg("Failed to update booking to CheckedIn")
		return false, utils.NewInternalServerError("DATABASE_ERROR", "Failed to update booking status", err)
	}
	return cmdTag.RowsAffected() == 1, nil
}

func (repo *bookingRepository) FindBookingsByIds(ctx context.Context, bookingIDs []int) ([]*models.Booking, error) {
	if len(bookingIDs) == 0 {
		return []*models.Booking{}, nil
	}
	query := `
		SELECT id, date, show_id, customer_id, customer_username, no_of_seats, amount_paid, status, booking_time, payment_type
		FROM booking
		WHERE id = ANY($1)
	`
	rows, err := repo.db.Query(ctx, query, bookingIDs)
	if err != nil {
		log.Error().Err(err).Interface("bookingIDs", bookingIDs).Msg("Failed to fetch bookings by IDs")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to fetch bookings by IDs", err)
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
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
			log.Error().Err(err).Msg("Error scanning bulk booking row")
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to scan bookings data", err)
		}
		bookings = append(bookings, &booking)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over bookings by IDs")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to iterate over bookings by IDs", err)
	}
	return bookings, nil
}
