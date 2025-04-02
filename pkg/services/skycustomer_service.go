package services

import (
	"context"
	"fmt"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type SkyCustomerService interface {
	ValidateUserDetails(ctx context.Context, username, email, phoneNumber string) error
	CreateCustomer(ctx context.Context, customer *models.SkyCustomer, user *models.User, passwordHistory *models.PasswordHistory, securityQuestionID int, securityAnswer string) error
	ExistsByEmailOrMobile(ctx context.Context, email, mobileNumber string) (bool, string, error)
	GetUsernameByEmail(ctx context.Context, email string) (string, error)
	GetCustomerProfileImg(ctx context.Context, username string) ([]byte, error)
	CustomerDetails(ctx context.Context, username string) (*models.SkyCustomer, error)
}

type skyCustomerService struct {
	skyCustomerRepo      repositories.SkyCustomerRepository
	userRepo             repositories.UserRepository
	securityQuestionRepo repositories.SecurityQuestionRepository
}

func NewSkyCustomerService(
	skyCustomerRepo repositories.SkyCustomerRepository,
	userRepo repositories.UserRepository,
	securityQuestionRepo repositories.SecurityQuestionRepository,
) SkyCustomerService {
	return &skyCustomerService{
		skyCustomerRepo:      skyCustomerRepo,
		userRepo:             userRepo,
		securityQuestionRepo: securityQuestionRepo,
	}
}

func (s *skyCustomerService) ValidateUserDetails(ctx context.Context, username, email, phoneNumber string) error {
	existingUser, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error checking user existence", err)
	}
	if existingUser != nil {
		return utils.NewBadRequestError("USERNAME_EXISTS", "Username is already taken", nil)
	}

	exists, field, err := s.ExistsByEmailOrMobile(ctx, email, phoneNumber)
	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error checking customer data", err)
	}
	if exists {
		var errorMessage string
		switch field {
		case "email":
			errorMessage = "Email already exists"
		case "mobilenumber":
			errorMessage = "Mobile number already exists"
		default:
			errorMessage = fmt.Sprintf("%s already exists", field)
		}
		return utils.NewBadRequestError(field+"_EXISTS", errorMessage, nil)
	}

	return nil
}

func (s *skyCustomerService) CreateCustomer(
	ctx context.Context,
	customer *models.SkyCustomer,
	user *models.User,
	passwordHistory *models.PasswordHistory,
	securityQuestionID int,
	securityAnswer string,
) error {
	exists, err := s.securityQuestionRepo.QuestionExists(ctx, securityQuestionID)
	if err != nil {
		return err
	}

	if !exists {
		return utils.NewBadRequestError("INVALID_SECURITY_QUESTION", "The selected security question does not exist", nil)
	}

	securityAnswerHash, err := utils.HashPassword(securityAnswer)
	if err != nil {
		return utils.NewInternalServerError("SECURITY_ANSWER_HASH_ERROR", "Failed to hash security answer", err)
	}

	customer.SecurityQuestionID = securityQuestionID
	customer.SecurityAnswerHash = securityAnswerHash

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	if err := s.userRepo.SavePasswordHistory(ctx, passwordHistory); err != nil {
		return err
	}

	if err := s.skyCustomerRepo.Create(ctx, customer); err != nil {
		return err
	}

	return nil
}

func (s *skyCustomerService) ExistsByEmailOrMobile(ctx context.Context, email, mobileNumber string) (bool, string, error) {
	return s.skyCustomerRepo.ExistsByEmailOrMobile(ctx, email, mobileNumber)
}

func (s *skyCustomerService) GetUsernameByEmail(ctx context.Context, email string) (string, error) {
	customer, err := s.skyCustomerRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", utils.NewInternalServerError("DATABASE_ERROR", "Error fetching customer by email", err)
	}
	if customer == nil {
		return "", utils.NewNotFoundError("USER_NOT_FOUND", fmt.Sprintf("No user found with email: %s", email), nil)
	}
	return customer.Username, nil
}

func (s *skyCustomerService) GetCustomerProfileImg(ctx context.Context, username string) ([]byte, error) {
	customer, err := s.skyCustomerRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching customer profile", err)
	}
	if customer == nil || customer.ProfileImg == nil {
		return nil, utils.NewNotFoundError("PROFILE_IMG_NOT_FOUND", fmt.Sprintf("No profile image found for user: %s", username), nil)
	}
	return customer.ProfileImg, nil
}

func (s *skyCustomerService) CustomerDetails(ctx context.Context, username string) (*models.SkyCustomer, error) {
	customer, err := s.skyCustomerRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching customer details", err)
	}
	if customer == nil {
		return nil, utils.NewNotFoundError("USER_NOT_FOUND", fmt.Sprintf("No customer found with username: %s", username), nil)
	}
	return customer, nil
}
