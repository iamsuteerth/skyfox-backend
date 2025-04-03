package services

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type SecurityQuestionService interface {
	GetAllSecurityQuestions(ctx context.Context) ([]repositories.SecurityQuestion, error)
	ValidateSecurityQuestionExists(ctx context.Context, questionID int) error
	GetSecurityQuestionByEmail(ctx context.Context, email string) (*response.SecurityQuestionResponse, error)
}

type securityQuestionService struct {
	securityQuestionRepo repositories.SecurityQuestionRepository
	skyCustomerRepo      repositories.SkyCustomerRepository
}

func NewSecurityQuestionService(securityQuestionRepo repositories.SecurityQuestionRepository, skyCustomerRepo repositories.SkyCustomerRepository) SecurityQuestionService {
	return &securityQuestionService{
		securityQuestionRepo: securityQuestionRepo,
		skyCustomerRepo:      skyCustomerRepo,
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
