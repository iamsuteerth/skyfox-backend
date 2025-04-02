package repositories

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type StaffRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.Staff, error)
	Create(ctx context.Context, staff *models.Staff) error
}

type staffRepository struct {
	db *pgxpool.Pool
}

func NewStaffRepository(db *pgxpool.Pool) StaffRepository {
	return &staffRepository{
		db: db,
	}
}

func (r *staffRepository) FindByUsername(ctx context.Context, username string) (*models.Staff, error) {
	query := `SELECT id, username, name, counter_no FROM stafftable WHERE username = $1`

	var staff models.Staff
	err := r.db.QueryRow(ctx, query, username).Scan(
		&staff.ID,
		&staff.Username,
		&staff.Name,
		&staff.CounterNumber,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("username", username).Msg("Database error while finding staff")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error querying database", err)
	}

	return &staff, nil
}

func (r *staffRepository) Create(ctx context.Context, staff *models.Staff) error {
	query := `INSERT INTO stafftable (username, name, counter_no) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(ctx, query, staff.Username, staff.Name, staff.CounterNumber).Scan(&staff.ID)
	if err != nil {
		log.Error().Err(err).Str("username", staff.Username).Msg("Failed to create staff")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to create staff", err)
	}

	log.Info().Str("username", staff.Username).Int("id", staff.ID).Msg("Staff created successfully")
	return nil
}
