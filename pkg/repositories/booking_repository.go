package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository interface {
	BookedSeatsByShow(ctx context.Context, showId int) int
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
