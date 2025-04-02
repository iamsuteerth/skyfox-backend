package repositories

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	SavePassword(ctx context.Context, username, password string) error
	FindByUsernameinPasswordHistory(ctx context.Context, username string) (*models.PasswordHistory, error)
	SavePasswordHistory(ctx context.Context, passwordHistory *models.PasswordHistory) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, password, role FROM usertable WHERE username = $1`

	var user models.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("username", username).Msg("Database error while finding user")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error querying database", err)
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO usertable (username, password, role) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(ctx, query, user.Username, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		log.Error().Err(err).Str("username", user.Username).Msg("Failed to create user")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to create user", err)
	}

	log.Info().Str("username", user.Username).Int("id", user.ID).Msg("User created successfully")
	return nil
}

func (r *userRepository) SavePassword(ctx context.Context, username, password string) error {
	query := `UPDATE usertable SET password = $1 WHERE username = $2`

	_, err := r.db.Exec(ctx, query, password, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to save password")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to save password", err)
	}

	return nil
}

func (r *userRepository) FindByUsernameinPasswordHistory(ctx context.Context, username string) (*models.PasswordHistory, error) {
	query := `SELECT id, username, previous_password_1, previous_password_2, previous_password_3 
              FROM password_history WHERE username = $1`

	var history models.PasswordHistory
	err := r.db.QueryRow(ctx, query, username).Scan(
		&history.ID,
		&history.Username,
		&history.PreviousPassword1,
		&history.PreviousPassword2,
		&history.PreviousPassword3,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("username", username).Msg("Error finding password history")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error finding password history", err)
	}

	return &history, nil
}

func (r *userRepository) SavePasswordHistory(ctx context.Context, history *models.PasswordHistory) error {
	query := `INSERT INTO password_history (username, previous_password_1, previous_password_2, previous_password_3) 
              VALUES ($1, $2, $3, $4) 
              ON CONFLICT (username) DO UPDATE 
              SET previous_password_1 = $2, previous_password_2 = $3, previous_password_3 = $4
              RETURNING id`

	err := r.db.QueryRow(ctx, query,
		history.Username,
		history.PreviousPassword1,
		history.PreviousPassword2,
		history.PreviousPassword3).Scan(&history.ID)

	if err != nil {
		log.Error().Err(err).Str("username", history.Username).Msg("Failed to save password history")
		return utils.NewInternalServerError("DATABASE_ERROR", "Error saving password history", err)
	}

	return nil
}
