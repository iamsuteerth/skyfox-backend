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
	GetWalletTransactionsByUsername(ctx context.Context, username string) ([]models.PaymentTransaction, error)
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

func (repo *paymentTransactionRepository) GetWalletTransactionsByUsername(ctx context.Context, username string) ([]models.PaymentTransaction, error) {
	query := `
		SELECT pt.id, pt.booking_id, pt.transaction_id, pt.payment_method,
			pt.amount, pt.status, pt.processed_at
		FROM payment_transaction pt
		JOIN booking b ON pt.booking_id = b.id
		WHERE b.customer_username = $1 AND pt.payment_method = 'Wallet'
		ORDER BY pt.processed_at DESC
	`
	
	rows, err := repo.db.Query(ctx, query, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to get wallet transactions")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to retrieve wallet transactions", err)
	}
	defer rows.Close()
	
	var transactions []models.PaymentTransaction
	for rows.Next() {
		var transaction models.PaymentTransaction
		err := rows.Scan(
			&transaction.Id,
			&transaction.BookingId,
			&transaction.TransactionId,
			&transaction.PaymentMethod,
			&transaction.Amount,
			&transaction.Status,
			&transaction.ProcessedAt,
		)
		if err != nil {
			log.Error().Err(err).Str("username", username).Msg("Failed to scan wallet transaction")
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to scan wallet transaction", err)
		}
		transactions = append(transactions, transaction)
	}
	
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Str("username", username).Msg("Error iterating wallet transactions")
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error iterating wallet transactions", err)
	}
	
	return transactions, nil
}
