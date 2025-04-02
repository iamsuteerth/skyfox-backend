// pkg/services/user_service.go
package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type UserService interface {
	Login(ctx context.Context, username, password string) (*models.User, string, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Login(ctx context.Context, username, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", utils.NewUnauthorizedError("INVALID_CREDENTIALS", "Invalid username or password", nil)
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", utils.NewUnauthorizedError("INVALID_CREDENTIALS", "Invalid username or password", nil)
	}

	token, err := generateToken(user)
	if err != nil {
		return nil, "", utils.NewInternalServerError("TOKEN_GENERATION_FAILED", "Failed to generate token", err)
	}

	return user, token, nil
}

func generateToken(user *models.User) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("JWT secret key not set")
	}

	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
