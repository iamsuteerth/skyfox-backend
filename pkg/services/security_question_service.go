package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type SecurityQuestionService interface {
	GetAllSecurityQuestions(ctx context.Context) ([]repositories.SecurityQuestion, error)
	ValidateSecurityQuestionExists(ctx context.Context, questionID int) error
	GetSecurityQuestionByEmail(ctx context.Context, email string) (*response.SecurityQuestionResponse, error)
	VerifySecurityAnswerAndGenerateToken(ctx context.Context, email, securityAnswer string) (*response.VerifySecurityAnswerResponse, error)
	VerifySecurityAnswer(ctx context.Context, email, securityAnswer string) (*response.VerifySecurityAnswerWithoutTokenResponse, error)
}

type securityQuestionService struct {
	securityQuestionRepo repositories.SecurityQuestionRepository
	skyCustomerRepo      repositories.SkyCustomerRepository
	resetTokenRepo       repositories.ResetTokenRepository
}

func NewSecurityQuestionService(securityQuestionRepo repositories.SecurityQuestionRepository, skyCustomerRepo repositories.SkyCustomerRepository, resetTokenRepo repositories.ResetTokenRepository) SecurityQuestionService {
	return &securityQuestionService{
		securityQuestionRepo: securityQuestionRepo,
		skyCustomerRepo:      skyCustomerRepo,
		resetTokenRepo:       resetTokenRepo,
	}
}

func (s *securityQuestionService) GetAllSecurityQuestions(ctx context.Context) ([]repositories.SecurityQuestion, error) {
	return s.securityQuestionRepo.FindAll(ctx)
}

func (s *securityQuestionService) ValidateSecurityQuestionExists(ctx context.Context, questionID int) error {
	exists, err := s.securityQuestionRepo.QuestionExists(ctx, questionID)
	if err != nil {
		return err
	}

	if !exists {
		return utils.NewBadRequestError("INVALID_SECURITY_QUESTION", "The selected security question does not exist", nil)
	}

	return nil
}

func (s *securityQuestionService) GetSecurityQuestionByEmail(ctx context.Context, email string) (*response.SecurityQuestionResponse, error) {
	customer, err := s.skyCustomerRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, utils.NewNotFoundError("USER_NOT_FOUND", "No user found with the provided email", nil)
	}

	question, err := s.securityQuestionRepo.FindByID(ctx, customer.SecurityQuestionID)
	if err != nil {
		return nil, err
	}

	return &response.SecurityQuestionResponse{
		QuestionID: question.ID,
		Question:   question.Question,
		Email:      email,
	}, nil
}

func (s *securityQuestionService) VerifySecurityAnswerAndGenerateToken(ctx context.Context, email, securityAnswer string) (*response.VerifySecurityAnswerResponse, error) {
	customer, err := s.skyCustomerRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, utils.NewNotFoundError("USER_NOT_FOUND", "No user found with the provided email", nil)
	}

	if !utils.CheckPasswordHash(securityAnswer, customer.SecurityAnswerHash) {
		return nil, utils.NewBadRequestError("INVALID_ANSWER", "The security answer provided is incorrect", nil)
	}

	token, expiresAt, exists, err := s.resetTokenRepo.GetValidToken(ctx, email)
	if err != nil {
		return nil, err
	}

	nowUTC := time.Now().UTC()

	if exists {
		expiresInSeconds := int(expiresAt.Sub(nowUTC).Seconds())
		return &response.VerifySecurityAnswerResponse{
			ResetToken: token,
			ExpiresIn:  expiresInSeconds,
		}, nil
	}

	if err := s.resetTokenRepo.DeletePreviousTokens(ctx, email); err != nil {
		return nil, err
	}

	token = uuid.New().String()
	expiresAt = nowUTC.Add(5 * time.Minute)
	expiresInSeconds := int(expiresAt.Sub(nowUTC).Seconds())

	if err := s.resetTokenRepo.StoreToken(ctx, email, token, expiresAt); err != nil {
		return nil, err
	}

	return &response.VerifySecurityAnswerResponse{
		ResetToken: token,
		ExpiresIn:  expiresInSeconds,
	}, nil
}

func (s *securityQuestionService) VerifySecurityAnswer(ctx context.Context, email, securityAnswer string) (*response.VerifySecurityAnswerWithoutTokenResponse, error) {
	customer, err := s.skyCustomerRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, utils.NewNotFoundError("USER_NOT_FOUND", "No user found with the provided email", nil)
	}

	if !utils.CheckPasswordHash(securityAnswer, customer.SecurityAnswerHash) {
		return &response.VerifySecurityAnswerWithoutTokenResponse{
			ValidAnswer: false,
		}, nil
	}

	return &response.VerifySecurityAnswerWithoutTokenResponse{
		ValidAnswer: true,
	}, nil
}
