package repositories

import (
	"context"
	"fmt"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository interface {
	BookedSeatsByShow(ctx context.Context, showId int) int
	GetBookedSeats(ctx context.Context, showId int) ([]models.BookingSeatMapping, error)
}

type bookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) BookingRepository {
	return &bookingRepository{db: db}
}

func (repo *bookingRepository) BookedSeatsByShow(ctx context.Context, showId int) int {
	query := `
        SELECT COALESCE(SUM(no_of_seats), 0)
        FROM booking
        WHERE show_id = $1 AND status IN ('Confirmed', 'CheckedIn')
    `

	var bookedSeats int
	err := repo.db.QueryRow(ctx, query, showId).Scan(&bookedSeats)
	if err != nil {
		if err != pgx.ErrNoRows {
			fmt.Printf("Error getting booked seats count for show %d: %v\n", showId, err)
		}
		return 0
	}

	return bookedSeats
}

func (repo *bookingRepository) GetBookedSeats(ctx context.Context, showId int) ([]models.BookingSeatMapping, error) {
	query := `
        SELECT bsm.seat_number
        FROM booking_seat_mapping bsm
        JOIN booking b ON b.id = bsm.booking_id
        WHERE b.show_id = $1
    `

	rows, err := repo.db.Query(ctx, query, showId)
	if err != nil {
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching booked seats", err)
	}
	defer rows.Close()

	var seats []models.BookingSeatMapping
	for rows.Next() {
		var seat models.BookingSeatMapping
		if err := rows.Scan(&seat.SeatNumber); err != nil {
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error scanning seat data", err)
		}
		seats = append(seats, seat)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error iterating over seats", err)
	}

	return seats, nil
}
