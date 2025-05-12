// wallet_repository.go
package repositories

import (
	"context"
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *models.CustomerWallet) error
	GetWalletByUsername(ctx context.Context, username string) (*models.CustomerWallet, error)
	UpdateWalletBalance(ctx context.Context, username string, amount float64) error
}

type walletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) WalletRepository {
	return &walletRepository{db: db}
}

func (repo *walletRepository) CreateWallet(ctx context.Context, wallet *models.CustomerWallet) error {
	query := `
        INSERT INTO customer_wallet (username, balance, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	now := time.Now()
	wallet.CreatedAt = now
	wallet.UpdatedAt = now

	err := repo.db.QueryRow(ctx, query,
		wallet.Username,
		wallet.Balance,
		wallet.CreatedAt,
		wallet.UpdatedAt).Scan(&wallet.ID)

	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error creating wallet", err)
	}

	return nil
}

func (repo *walletRepository) GetWalletByUsername(ctx context.Context, username string) (*models.CustomerWallet, error) {
	query := `SELECT id, username, balance, created_at, updated_at FROM customer_wallet WHERE username = $1`

	var wallet models.CustomerWallet
	err := repo.db.QueryRow(ctx, query, username).Scan(
		&wallet.ID,
		&wallet.Username,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching wallet", err)
	}

	return &wallet, nil
}

func (repo *walletRepository) UpdateWalletBalance(ctx context.Context, username string, amount float64) error {
	query := `
        UPDATE customer_wallet 
        SET balance = balance + $1, updated_at = $2
        WHERE username = $3
        RETURNING balance
    `

	now := time.Now()
	var newBalance float64
	err := repo.db.QueryRow(ctx, query, amount, now, username).Scan(&newBalance)

	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error updating wallet balance", err)
	}

	if newBalance < 0 {
		rollbackQuery := `
            UPDATE customer_wallet 
            SET balance = balance - $1, updated_at = $2
            WHERE username = $3
        `
		_, rollbackErr := repo.db.Exec(ctx, rollbackQuery, amount, now, username)
		if rollbackErr != nil {
			return utils.NewInternalServerError("DATABASE_ERROR", "Error rolling back negative balance", rollbackErr)
		}

		return utils.NewBadRequestError("INSUFFICIENT_BALANCE", "Insufficient wallet balance", nil)
	}

	return nil
}
