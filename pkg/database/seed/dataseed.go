package seed

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

func SeedDB(userRepo repositories.UserRepository, staffRepo repositories.StaffRepository) {
	ctx := context.Background()
	log.Info().Msg("Starting database seeding...")

	seedUser(ctx, userRepo, staffRepo, "seed-user-1", "foobar", "admin", "Admin One", 101)
	seedUser(ctx, userRepo, staffRepo, "seed-user-2", "foobar", "admin", "Admin Two", 102)

	seedUser(ctx, userRepo, staffRepo, "staff-1", "foobar", "staff", "Staff One", 501)

	log.Info().Msg("Database seeding completed successfully")
}

func seedUser(ctx context.Context, userRepo repositories.UserRepository, staffRepo repositories.StaffRepository,
	username, password, role, staffName string, counterNumber int) {

	user, _ := userRepo.FindByUsername(ctx, username)
	staff, _ := staffRepo.FindByUsername(ctx, username)
	passwordHistory, _ := userRepo.FindByUsernameinPasswordHistory(ctx, username)

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to hash password during seeding")
		return
	}

	if user == nil {
		log.Info().Str("username", username).Msg("Creating new user")
		newUser := models.NewUser(username, hashedPassword, role)
		err = userRepo.Create(ctx, &newUser)
		if err != nil {
			log.Error().Err(err).Str("username", username).Msg("Failed to create user during seeding")
		} else {
			log.Info().Str("username", username).Str("role", role).Msg("User created successfully")
		}
	}

	if staff == nil {
		log.Info().Str("username", username).Msg("Creating new staff record")
		newStaff := models.NewStaff(username, staffName, counterNumber)
		err = staffRepo.Create(ctx, &newStaff)
		if err != nil {
			log.Error().Err(err).Str("username", username).Msg("Failed to create staff during seeding")
		} else {
			log.Info().Str("username", username).Str("name", staffName).Int("counterNumber", counterNumber).Msg("Staff created successfully")
		}
	}

	if passwordHistory == nil {
		log.Info().Str("username", username).Msg("Creating password history")
		newPasswordHistory := models.NewPasswordHistory(username, hashedPassword, "", "")
		err = userRepo.SavePasswordHistory(ctx, &newPasswordHistory)
		if err != nil {
			log.Error().Err(err).Str("username", username).Msg("Failed to create password history during seeding")
		} else {
			log.Info().Str("username", username).Msg("Password history created successfully")
		}
	}
}
