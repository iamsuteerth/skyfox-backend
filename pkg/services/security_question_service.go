package services

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type SecurityQuestionService interface {
	GetAllSecurityQuestions(ctx context.Context) ([]repositories.SecurityQuestion, error)
	ValidateSecurityQuestionExists(ctx context.Context, questionID int) error
}

type securityQuestionService struct {
	securityQuestionRepo repositories.SecurityQuestionRepository
}

func NewSecurityQuestionService(securityQuestionRepo repositories.SecurityQuestionRepository) SecurityQuestionService {
	return &securityQuestionService{
		securityQuestionRepo: securityQuestionRepo,
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
