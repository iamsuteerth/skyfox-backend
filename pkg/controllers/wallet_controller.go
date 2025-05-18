package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/request"
	"github.com/iamsuteerth/skyfox-backend/pkg/middleware/security"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type WalletController struct {
	walletService services.WalletService
}

func NewWalletController(walletService services.WalletService) *WalletController {
	return &WalletController{
		walletService: walletService,
	}
}

func (wc *WalletController) GetWalletBalance(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)
	
	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify user credentials", err), requestID)
		return
	}
	
	username, _ := claims["username"].(string)
	
	walletResponse, err := wc.walletService.GetWalletBalance(ctx.Request.Context(), username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to get wallet balance")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}
	
	utils.SendOKResponse(ctx, "Wallet balance retrieved successfully", requestID, walletResponse)
}

func (wc *WalletController) AddFunds(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)
	
	var addFundsRequest request.AddWalletFundsRequest
	if err := ctx.ShouldBindJSON(&addFundsRequest); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			utils.HandleErrorResponse(ctx, utils.NewValidationError(validationErrs), requestID)
			return
		}
		utils.HandleErrorResponse(ctx, utils.NewBadRequestError("INVALID_REQUEST", "Invalid request data", err), requestID)
		return
	}
	
	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify user credentials", err), requestID)
		return
	}
	
	username, _ := claims["username"].(string)
	
	walletResponse, err := wc.walletService.AddFunds(ctx.Request.Context(), username, &addFundsRequest)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to add funds to wallet")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}
	
	utils.SendOKResponse(ctx, "Funds added successfully", requestID, walletResponse)
}

func (wc *WalletController) GetTransactions(ctx *gin.Context) {
	requestID := utils.GetRequestID(ctx)
	
	claims, err := security.GetTokenClaims(ctx)
	if err != nil {
		utils.HandleErrorResponse(ctx, utils.NewUnauthorizedError("UNAUTHORIZED", "Unable to verify user credentials", err), requestID)
		return
	}
	
	username, _ := claims["username"].(string)
	
	transactions, err := wc.walletService.GetTransactions(ctx.Request.Context(), username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to get wallet transactions")
		utils.HandleErrorResponse(ctx, err, requestID)
		return
	}
	
	utils.SendOKResponse(ctx, "Wallet transactions retrieved successfully", requestID, transactions)
}
