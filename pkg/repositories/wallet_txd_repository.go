package repositories

import (
	"context"
	"time"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"
)

type WalletTransactionRepository interface {
	AddWalletTransaction(ctx context.Context, txn *models.WalletTransaction) error
	GetWalletTransactionsByWalletId(ctx context.Context, walletId int64) ([]*models.WalletTransaction, error)
	GetWalletTransactionsForUser(ctx context.Context, username string) ([]*models.WalletTransaction, error)
	GetWalletTransactionsForBooking(ctx context.Context, bookingId int64) ([]*models.WalletTransaction, error)
	GetLatestWalletTransactionForUser(ctx context.Context, username string) (*models.WalletTransaction, error)
}

type walletTransactionRepository struct {
	db *pgxpool.Pool
}

func NewWalletTransactionRepository(db *pgxpool.Pool) WalletTransactionRepository {
	return &walletTransactionRepository{db: db}
}

func (r *walletTransactionRepository) AddWalletTransaction(ctx context.Context, txn *models.WalletTransaction) error {
	query := `
		INSERT INTO wallet_transaction
			(wallet_id, username, booking_id, transaction_id, amount, timestamp, transaction_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	now := time.Now()
	txn.Timestamp = now
	return r.db.QueryRow(ctx, query,
		txn.WalletID, txn.Username, txn.BookingID, txn.TransactionID, txn.Amount, txn.Timestamp, txn.TransactionType,
	).Scan(&txn.ID)
}

func (r *walletTransactionRepository) GetWalletTransactionsByWalletId(ctx context.Context, walletId int64) ([]*models.WalletTransaction, error) {
    query := `
        SELECT id, wallet_id, username, booking_id, transaction_id, amount, timestamp, transaction_type
        FROM wallet_transaction
        WHERE wallet_id = $1
        ORDER BY timestamp DESC
    `
    rows, err := r.db.Query(ctx, query, walletId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var transactions []*models.WalletTransaction
    for rows.Next() {
        var txn models.WalletTransaction
        var bookingID pgtype.Int8
        err := rows.Scan(&txn.ID, &txn.WalletID, &txn.Username, &bookingID, &txn.TransactionID, &txn.Amount, &txn.Timestamp, &txn.TransactionType)
        if err != nil {
            return nil, err
        }
        if bookingID.Valid {
            id := bookingID.Int64
            txn.BookingID = &id
        }
        transactions = append(transactions, &txn)
    }
    return transactions, nil
}

func (r *walletTransactionRepository) GetWalletTransactionsForUser(ctx context.Context, username string) ([]*models.WalletTransaction, error) {
    query := `
        SELECT id, wallet_id, username, booking_id, transaction_id, amount, timestamp, transaction_type
        FROM wallet_transaction
        WHERE username = $1
        ORDER BY timestamp DESC
    `
    rows, err := r.db.Query(ctx, query, username)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var transactions []*models.WalletTransaction
    for rows.Next() {
        var txn models.WalletTransaction
        var bookingID pgtype.Int8
        err := rows.Scan(&txn.ID, &txn.WalletID, &txn.Username, &bookingID, &txn.TransactionID, &txn.Amount, &txn.Timestamp, &txn.TransactionType)
        if err != nil {
            return nil, err
        }
        if bookingID.Valid {
            id := bookingID.Int64
            txn.BookingID = &id
        }
        transactions = append(transactions, &txn)
    }
    return transactions, nil
}


func (r *walletTransactionRepository) GetWalletTransactionsForBooking(ctx context.Context, bookingId int64) ([]*models.WalletTransaction, error) {
	query := `
		SELECT id, wallet_id, username, booking_id, transaction_id, amount, timestamp, transaction_type
		FROM wallet_transaction
		WHERE booking_id = $1
		ORDER BY timestamp ASC
	`
	rows, err := r.db.Query(ctx, query, bookingId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []*models.WalletTransaction
	for rows.Next() {
		var txn models.WalletTransaction
		var bookingID  pgtype.Int8
		err := rows.Scan(&txn.ID, &txn.WalletID, &txn.Username, &bookingID, &txn.TransactionID, &txn.Amount, &txn.Timestamp, &txn.TransactionType)
		if err != nil {
			return nil, err
		}
		if bookingID.Valid {
			id := bookingID.Int64
			txn.BookingID = &id
		}
		transactions = append(transactions, &txn)
	}
	return transactions, nil
}

func (r *walletTransactionRepository) GetLatestWalletTransactionForUser(ctx context.Context, username string) (*models.WalletTransaction, error) {
	query := `
		SELECT id, wallet_id, username, booking_id, transaction_id, amount, timestamp, transaction_type
		FROM wallet_transaction
		WHERE username = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`
	var txn models.WalletTransaction
	var bookingID  pgtype.Int8
	err := r.db.QueryRow(ctx, query, username).Scan(
		&txn.ID, &txn.WalletID, &txn.Username, &bookingID, &txn.TransactionID, &txn.Amount, &txn.Timestamp, &txn.TransactionType,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching latest wallet transaction", err)
	}
	if bookingID.Valid {
		id := bookingID.Int64
		txn.BookingID = &id
	}
	return &txn, nil
}
