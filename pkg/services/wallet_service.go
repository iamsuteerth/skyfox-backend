package services

import (
	"context"
	"fmt"
	"time"

	"github.com/govalues/decimal"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	paymentservice "github.com/iamsuteerth/skyfox-backend/pkg/payment-service"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type WalletService interface {
	AddFunds(ctx context.Context, username string, req *request.AddWalletFundsRequest) (*response.WalletResponse, error)
	GetWalletBalance(ctx context.Context, username string) (*response.WalletResponse, error)
	GetTransactions(ctx context.Context, username string) (*response.WalletTransactionsResponse, error)
}

type walletService struct {
	customerWalletRepo     repositories.CustomerWalletRepository
	walletTxdRepo          repositories.WalletTransactionRepository
	paymentTransactionRepo repositories.PaymentTransactionRepository
	paymentService         paymentservice.PaymentService
}

func NewWalletService(
	customerWalletRepo repositories.CustomerWalletRepository,
	walletTxdRepo repositories.WalletTransactionRepository,
	paymentTransactionRepo repositories.PaymentTransactionRepository,
	paymentService paymentservice.PaymentService,
) WalletService {
	return &walletService{
		customerWalletRepo:     customerWalletRepo,
		walletTxdRepo:          walletTxdRepo,
		paymentTransactionRepo: paymentTransactionRepo,
		paymentService:         paymentService,
	}
}

func (s *walletService) AddFunds(ctx context.Context, username string, req *request.AddWalletFundsRequest) (*response.WalletResponse, error) {
	if req.Amount.Cmp(decimal.Zero) <= 0 {
		return nil, utils.NewBadRequestError("INVALID_AMOUNT", "Amount must be greater than zero", nil)
	}

	maxAmount, _ := decimal.NewFromFloat64(10000)

	if req.Amount.Cmp(maxAmount) > 0 {
		return nil, utils.NewBadRequestError("AMOUNT_TOO_LARGE", "Maximum amount allowed is 10000", nil)
	}
	expiry := fmt.Sprintf("%s/%s", req.ExpiryMonth, req.ExpiryYear)
	transactionID, err := s.paymentService.ProcessPayment(
		ctx,
		req.CardNumber,
		req.CVV,
		expiry,
		req.CardholderName,
		req.Amount,
	)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Card payment failed when adding funds to wallet")
		return nil, err
	}

	wallet, err := s.customerWalletRepo.GetWalletByUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to retrieve wallet")
		return nil, err
	}

	if wallet == nil {
		return nil, utils.NewNotFoundError("WALLET_NOT_FOUND", "Wallet not found for user", nil)
	}

	if err := s.customerWalletRepo.AddToWalletBalance(ctx, username, req.Amount); err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to update wallet balance")
		return nil, err
	}

	walletTxn := &models.WalletTransaction{
		WalletID:        wallet.ID,
		Username:        username,
		BookingID:       nil,
		TransactionID:   transactionID,
		Amount:          req.Amount,
		TransactionType: "ADD",
	}

	if err := s.walletTxdRepo.AddWalletTransaction(ctx, walletTxn); err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to record wallet transaction")
		return nil, err
	}

	updatedWallet, err := s.customerWalletRepo.GetWalletByUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to retrieve updated wallet")
		return nil, err
	}

	updatedWalletBalance, _ := updatedWallet.Balance.Float64()

	return &response.WalletResponse{
		Username:  username,
		Balance:   updatedWalletBalance,
		UpdatedAt: updatedWallet.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *walletService) GetWalletBalance(ctx context.Context, username string) (*response.WalletResponse, error) {
	wallet, err := s.customerWalletRepo.GetWalletByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, utils.NewNotFoundError("WALLET_NOT_FOUND", "Wallet not found for user", nil)
	}

	balance, _ := wallet.Balance.Float64()

	return &response.WalletResponse{
		Username:  username,
		Balance:   balance,
		UpdatedAt: wallet.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *walletService) GetTransactions(ctx context.Context, username string) (*response.WalletTransactionsResponse, error) {
	wallet, err := s.customerWalletRepo.GetWalletByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, utils.NewNotFoundError("WALLET_NOT_FOUND", "Wallet not found for user", nil)
	}

	transactions, err := s.walletTxdRepo.GetWalletTransactionsForUser(ctx, username)
	if err != nil {
		return nil, err
	}

	var transactionResponses []response.WalletTransactionResponse
	for _, t := range transactions {
		amount, _ := t.Amount.Float64()
		response := response.WalletTransactionResponse{
			ID:              t.ID,
			Amount:          amount,
			TransactionType: t.TransactionType,
			TransactionID:   t.TransactionID,
			Timestamp:       t.Timestamp.Format(time.RFC3339),
		}

		if t.BookingID != nil {
			response.BookingID = t.BookingID
		}

		transactionResponses = append(transactionResponses, response)
	}

	return &response.WalletTransactionsResponse{
		Transactions: transactionResponses,
	}, nil
}
