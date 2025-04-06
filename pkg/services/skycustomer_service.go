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
	CreateCustomer(ctx context.Context, customer *models.SkyCustomer, user *models.User, passwordHistory *models.PasswordHistory, securityQuestionID int, securityAnswer string, profileImgBytes []byte, profileImgSHA string) error
	ExistsByEmailOrMobile(ctx context.Context, email, mobileNumber string) (bool, string, error)
	GetUsernameByEmail(ctx context.Context, email string) (string, error)
	GetProfileImageURL(ctx context.Context, username string) (string, error)
	CustomerDetails(ctx context.Context, username string) (*models.SkyCustomer, error)
	UpdateProfileImage(ctx context.Context, username string, imageBytes []byte, imageSHA string) error
}

type skyCustomerService struct {
	skyCustomerRepo      repositories.SkyCustomerRepository
	userRepo             repositories.UserRepository
	securityQuestionRepo repositories.SecurityQuestionRepository
	s3Service            S3Service
}

func NewSkyCustomerService(
	skyCustomerRepo repositories.SkyCustomerRepository,
	userRepo repositories.UserRepository,
	securityQuestionRepo repositories.SecurityQuestionRepository,
	s3Service S3Service,
) SkyCustomerService {
	return &skyCustomerService{
		skyCustomerRepo:      skyCustomerRepo,
		userRepo:             userRepo,
		securityQuestionRepo: securityQuestionRepo,
		s3Service:            s3Service,
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
	profileImgBytes []byte,
	profileImgSHA string,
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

	if len(profileImgBytes) > 0 && profileImgSHA != "" {
		imageURL, err := s.s3Service.UploadProfileImage(ctx, profileImgBytes, user.Username, profileImgSHA)
		if err != nil {
			return err
		}
		customer.ProfileImg = imageURL
	} else {
		return utils.NewBadRequestError("INVALID_IMAGE_HASH", "The hash or the image bytes provided are invalid.", nil)
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
		if customer.ProfileImg != "" {
			_ = s.s3Service.DeleteProfileImage(ctx, customer.ProfileImg)
		}
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

func (s *skyCustomerService) GetProfileImageURL(ctx context.Context, username string) (string, error) {
	return s.skyCustomerRepo.GetCustomerProfileImg(ctx, username)
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

func (s *skyCustomerService) UpdateProfileImage(ctx context.Context, username string, imageBytes []byte, imageSHA string) error {
	customer, err := s.skyCustomerRepo.FindByUsername(ctx, username)
	if err != nil {
		return err
	}
	if customer == nil {
		return utils.NewNotFoundError("USER_NOT_FOUND", fmt.Sprintf("No user found with username: %s", username), nil)
	}

	newImageURL, err := s.s3Service.UploadProfileImage(ctx, imageBytes, username, imageSHA)
	if err != nil {
		return err
	}

	if err := s.skyCustomerRepo.UpdateProfileImageURL(ctx, username, newImageURL); err != nil {
		_ = s.s3Service.DeleteProfileImage(ctx, newImageURL)
		return utils.NewInternalServerError("DB_UPDATE_FAILED", "Failed to update profile image URL in database", err)
	}

	if oldImageURL := customer.ProfileImg; oldImageURL != "" {
		if delErr := s.s3Service.DeleteProfileImage(ctx, oldImageURL); delErr != nil {
			return utils.NewInternalServerError(
				"PARTIAL_UPDATE",
				"Profile image updated, but failed to delete old image from S3",
				delErr,
			)
		}
	}

	return nil
}
