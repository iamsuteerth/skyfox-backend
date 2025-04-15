package services

import (
	"context"
	"fmt"
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/dto/response"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type AdminStaffProfileService interface {
	GetProfile(ctx context.Context, username string) (*response.AdminStaffProfileResponse, error)
}

type adminStaffProfileService struct {
	userRepo  repositories.UserRepository
	staffRepo repositories.StaffRepository
}

func NewAdminStaffProfileService(userRepo repositories.UserRepository, staffRepo repositories.StaffRepository) AdminStaffProfileService {
	return &adminStaffProfileService{
		userRepo:  userRepo,
		staffRepo: staffRepo,
	}
}

func (s *adminStaffProfileService) GetProfile(ctx context.Context, username string) (*response.AdminStaffProfileResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching user details", err)
	}
	
	if user == nil {
		return nil, utils.NewNotFoundError("USER_NOT_FOUND", fmt.Sprintf("No user found with username: %s", username), nil)
	}
	
	staff, err := s.staffRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching user details", err)
	}
	
	if staff == nil {
		return nil, utils.NewNotFoundError("USER_NOT_FOUND", fmt.Sprintf("No staff or admin details found for username: %s", username), nil)
	}
	
	createdAtStr := user.CreatedAt.Format(time.RFC3339)
	
	profile := &response.AdminStaffProfileResponse{
		Username:  user.Username,
		Name:      staff.Name,
		CounterNo: staff.CounterNumber,
		CreatedAt: createdAtStr,
	}
	
	return profile, nil
}
