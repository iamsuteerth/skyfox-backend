package repositories

import (
	"context"
	"time"

	"github.com/govalues/decimal"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerWalletRepository interface {
	CreateWallet(ctx context.Context, wallet *models.CustomerWallet) error
	GetWalletByUsername(ctx context.Context, username string) (*models.CustomerWallet, error)
	GetWalletById(ctx context.Context, walletId int64) (*models.CustomerWallet, error)
	AddToWalletBalance(ctx context.Context, username string, amount decimal.Decimal) error
	DeductFromWalletBalance(ctx context.Context, username string, amount decimal.Decimal) error
}

type customerWalletRepository struct {
	db *pgxpool.Pool
}

func NewCustomerWalletRepository(db *pgxpool.Pool) CustomerWalletRepository {
	return &customerWalletRepository{db: db}
}

func (r *customerWalletRepository) CreateWallet(ctx context.Context, wallet *models.CustomerWallet) error {
	query := `
		INSERT INTO customer_wallet (username, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	now := time.Now()
	wallet.CreatedAt = now
	wallet.UpdatedAt = now
	return r.db.QueryRow(ctx, query, wallet.Username, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt).Scan(&wallet.ID)
}

func (r *customerWalletRepository) GetWalletByUsername(ctx context.Context, username string) (*models.CustomerWallet, error) {
	query := `SELECT id, username, balance, created_at, updated_at FROM customer_wallet WHERE username = $1`
	var wallet models.CustomerWallet
	err := r.db.QueryRow(ctx, query, username).Scan(&wallet.ID, &wallet.Username, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching wallet", err)
	}
	return &wallet, nil
}

func (r *customerWalletRepository) GetWalletById(ctx context.Context, walletId int64) (*models.CustomerWallet, error) {
	query := `SELECT id, username, balance, created_at, updated_at FROM customer_wallet WHERE id = $1`
	var wallet models.CustomerWallet
	err := r.db.QueryRow(ctx, query, walletId).Scan(&wallet.ID, &wallet.Username, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching wallet by id", err)
	}
	return &wallet, nil
}

func (r *customerWalletRepository) AddToWalletBalance(ctx context.Context, username string, amount decimal.Decimal) error {
	if amount.Cmp(decimal.Zero) != 1 {
		return utils.NewBadRequestError("INVALID_AMOUNT", "Amount must be positive", nil)
	}

	query := `
        UPDATE customer_wallet 
        SET balance = balance + $1, updated_at = $2
        WHERE username = $3
    `
	now := time.Now()
	_, err := r.db.Exec(ctx, query, amount, now, username)
	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error adding to wallet balance", err)
	}
	return nil
}

func (r *customerWalletRepository) DeductFromWalletBalance(ctx context.Context, username string, amount decimal.Decimal) error {
	if amount.Cmp(decimal.Zero) != 1 {
		return utils.NewBadRequestError("INVALID_AMOUNT", "Amount must be positive", nil)
	}

	query := `
        UPDATE customer_wallet 
        SET balance = balance - $1, updated_at = $2
        WHERE username = $3
        RETURNING balance
    `
	now := time.Now()
	var newBalance decimal.Decimal
	err := r.db.QueryRow(ctx, query, amount, now, username).Scan(&newBalance)
	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error deducting from wallet balance", err)
	}

	if newBalance.Cmp(decimal.Zero) == -1 {
		rollbackQuery := `
            UPDATE customer_wallet 
            SET balance = balance + $1, updated_at = $2
            WHERE username = $3
        `
		_, rollbackErr := r.db.Exec(ctx, rollbackQuery, amount, now, username)
		if rollbackErr != nil {
			return utils.NewInternalServerError("DATABASE_ERROR", "Error rolling back negative balance", rollbackErr)
		}
		return utils.NewBadRequestError("INSUFFICIENT_BALANCE", "Insufficient wallet balance", nil)
	}

	return nil
}
