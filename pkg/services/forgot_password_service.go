package services

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type ForgotPasswordService interface {
	ForgotPassword(ctx context.Context, email, token, newPassword string) error
}

type forgotPasswordService struct {
	resetTokenRepo  repositories.ResetTokenRepository
	skyCustomerRepo repositories.SkyCustomerRepository
	userRepo        repositories.UserRepository
}

func NewForgotPasswordService(resetTokenRepo repositories.ResetTokenRepository, skyCustomerRepo repositories.SkyCustomerRepository, userRepo repositories.UserRepository) ForgotPasswordService {
	return &forgotPasswordService{
		resetTokenRepo:  resetTokenRepo,
		skyCustomerRepo: skyCustomerRepo,
		userRepo:        userRepo,
	}
}

func (s *forgotPasswordService) ForgotPassword(ctx context.Context, email, token, newPassword string) error {
	customer, err := s.skyCustomerRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if customer == nil {
		return utils.NewNotFoundError("USER_NOT_FOUND", "No user found with the provided email", nil)
	}

	valid, err := s.resetTokenRepo.ValidateToken(ctx, email, token)
	if err != nil {
		return err
	}

	if !valid {
		if err := s.resetTokenRepo.DeletePreviousTokens(ctx, email); err != nil {
		}
		return utils.NewBadRequestError("INVALID_RESET_TOKEN", "The reset token is invalid, expired, or has already been used", nil)
	}

	passwordHistory, err := s.userRepo.FindByUsernameinPasswordHistory(ctx, customer.Username)
	if err != nil {
		return err
	}

	if passwordHistory != nil {
		if passwordHistory.PreviousPassword1 != "" && utils.CheckPasswordHash(newPassword, passwordHistory.PreviousPassword1) {
			return utils.NewBadRequestError("PASSWORD_REUSE", "New password cannot match any of your previous passwords", nil)
		}
		if passwordHistory.PreviousPassword2 != "" && utils.CheckPasswordHash(newPassword, passwordHistory.PreviousPassword2) {
			return utils.NewBadRequestError("PASSWORD_REUSE", "New password cannot match any of your previous passwords", nil)
		}
		if passwordHistory.PreviousPassword3 != "" && utils.CheckPasswordHash(newPassword, passwordHistory.PreviousPassword3) {
			return utils.NewBadRequestError("PASSWORD_REUSE", "New password cannot match any of your previous passwords", nil)
		}
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return utils.NewInternalServerError("PASSWORD_HASH_ERROR", "Failed to hash password", err)
	}

	if passwordHistory == nil {
		passwordHistory = &models.PasswordHistory{
			Username:          customer.Username,
			PreviousPassword1: hashedPassword,
		}
	} else {
		passwordHistory.PreviousPassword3 = passwordHistory.PreviousPassword2
		passwordHistory.PreviousPassword2 = passwordHistory.PreviousPassword1
		passwordHistory.PreviousPassword1 = hashedPassword
	}

	if err := s.userRepo.SavePassword(ctx, customer.Username, hashedPassword); err != nil {
		return err
	}

	if err := s.userRepo.SavePasswordHistory(ctx, passwordHistory); err != nil {
		return err
	}

	if err := s.resetTokenRepo.InvalidateToken(ctx, email, token); err != nil {
		return err
	}

	return nil
}
