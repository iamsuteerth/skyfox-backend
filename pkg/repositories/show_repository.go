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
	GetShowByMovieAndSlot(ctx context.Context, movieId string, showDate string, slotId int) (*models.Show, error)
	GetAllShows(ctx context.Context) ([]models.Show, error)
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

func (repo *showRepository) GetShowByMovieAndSlot(ctx context.Context, movieId string, showDate string, slotId int) (*models.Show, error) {
	query := `
        SELECT s.id, s.movie_id, s.date, s.slot_id, s.cost, 
               sl.id, sl.name, sl.start_time, sl.end_time
        FROM show s
        JOIN slot sl ON s.slot_id = sl.id
        WHERE s.movie_id = $1 AND s.date = $2 AND s.slot_id = $3
    `

	var show models.Show
	var slot models.Slot

	err := repo.db.QueryRow(ctx, query, movieId, showDate, slotId).Scan(
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
			return nil, nil
		}
		log.Error().Err(err).Msg("Error querying for show by movie and slot")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Could not check existing show", err)
	}

	show.Slot = slot
	return &show, nil
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

func (repo *showRepository) GetAllShows(ctx context.Context) ([]models.Show, error) {
	query := `
        SELECT s.id, s.movie_id, s.date, s.slot_id, s.cost, 
               sl.id, sl.name, sl.start_time, sl.end_time
        FROM show s
        JOIN slot sl ON s.slot_id = sl.id
    `

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query all shows")
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
