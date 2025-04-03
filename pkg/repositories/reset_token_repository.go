package repositories

import (
	"context"
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResetTokenRepository interface {
	StoreToken(ctx context.Context, email, token string, expiresAt time.Time) error
	ValidateToken(ctx context.Context, email, token string) (bool, error)
	InvalidateToken(ctx context.Context, email, token string) error
	GetValidToken(ctx context.Context, email string) (string, time.Time, bool, error)
	DeletePreviousTokens(ctx context.Context, email string) error
}

type resetTokenRepository struct {
	db *pgxpool.Pool
}

func NewResetTokenRepository(db *pgxpool.Pool) ResetTokenRepository {
	return &resetTokenRepository{db: db}
}

func (repo *resetTokenRepository) StoreToken(ctx context.Context, email, token string, expiresAt time.Time) error {
	now := time.Now().UTC()
	expiresAt = expiresAt.UTC()
	query := `
        INSERT INTO password_reset_tokens (email, token, created_at, expires_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err := repo.db.Exec(ctx, query, email, token, now, expiresAt)
	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error storing reset token", err)
	}

	return nil
}

func (repo *resetTokenRepository) ValidateToken(ctx context.Context, email, token string) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM password_reset_tokens 
            WHERE email = $1 AND token = $2 AND expires_at > $3 AND used = false
        )
    `

	nowUTC := time.Now().UTC()
	var valid bool
	err := repo.db.QueryRow(ctx, query, email, token, nowUTC).Scan(&valid)
	if err != nil {
		return false, utils.NewInternalServerError("DATABASE_ERROR", "Error validating reset token", err)
	}
	return valid, nil
}

func (repo *resetTokenRepository) InvalidateToken(ctx context.Context, email, token string) error {
	query := `
		UPDATE password_reset_tokens 
		SET used = true 
		WHERE email = $1 AND token = $2
	`

	_, err := repo.db.Exec(ctx, query, email, token)
	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error invalidating reset token", err)
	}

	return nil
}

func (repo *resetTokenRepository) GetValidToken(ctx context.Context, email string) (string, time.Time, bool, error) {
	query := `
        SELECT token, expires_at FROM password_reset_tokens 
        WHERE email = $1 AND expires_at > $2 AND used = false
        ORDER BY created_at DESC 
        LIMIT 1
    `

	var token string
	var expiresAt time.Time
	nowUTC := time.Now().UTC()

	err := repo.db.QueryRow(ctx, query, email, nowUTC).Scan(&token, &expiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", time.Time{}, false, nil
		}
		return "", time.Time{}, false, utils.NewInternalServerError("DATABASE_ERROR", "Error getting reset token", err)
	}

	expiresAt = expiresAt.UTC()
	return token, expiresAt, true, nil
}

func (repo *resetTokenRepository) DeletePreviousTokens(ctx context.Context, email string) error {
	query := `DELETE FROM password_reset_tokens WHERE email = $1`

	_, err := repo.db.Exec(ctx, query, email)
	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error deleting previous tokens", err)
	}

	return nil
}
