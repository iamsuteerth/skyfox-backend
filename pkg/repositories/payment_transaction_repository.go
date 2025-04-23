package repositories

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PaymentTransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *models.PaymentTransaction) error
	GetTransactionByBookingId(ctx context.Context, bookingId int) (*models.PaymentTransaction, error)
}

type paymentTransactionRepository struct {
	db *pgxpool.Pool
}

func NewPaymentTransactionRepository(db *pgxpool.Pool) PaymentTransactionRepository {
	return &paymentTransactionRepository{db: db}
}

func (repo *paymentTransactionRepository) CreateTransaction(ctx context.Context, transaction *models.PaymentTransaction) error {
	query := `
        INSERT INTO payment_transaction (
            booking_id, transaction_id, payment_method, 
            amount, status
        )
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, processed_at
    `

	err := repo.db.QueryRow(ctx, query,
		transaction.BookingId,
		transaction.TransactionId,
		transaction.PaymentMethod,
		transaction.Amount,
		transaction.Status,
	).Scan(&transaction.Id, &transaction.ProcessedAt)

	if err != nil {
		log.Error().Err(err).Int("bookingId", transaction.BookingId).Msg("Failed to create payment transaction")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to record payment transaction", err)
	}

	return nil
}

func (repo *paymentTransactionRepository) GetTransactionByBookingId(ctx context.Context, bookingId int) (*models.PaymentTransaction, error) {
	query := `
        SELECT id, booking_id, transaction_id, payment_method, 
               amount, status, processed_at
        FROM payment_transaction
        WHERE booking_id = $1
    `

	var transaction models.PaymentTransaction
	err := repo.db.QueryRow(ctx, query, bookingId).Scan(
		&transaction.Id,
		&transaction.BookingId,
		&transaction.TransactionId,
		&transaction.PaymentMethod,
		&transaction.Amount,
		&transaction.Status,
		&transaction.ProcessedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil 
		}
		log.Error().Err(err).Int("bookingId", bookingId).Msg("Failed to get payment transaction")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve payment transaction", err)
	}

	return &transaction, nil
}
