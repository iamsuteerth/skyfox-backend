// server/main.go
package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/iamsuteerth/skyfox-backend/pkg/config"
	"github.com/iamsuteerth/skyfox-backend/pkg/controllers"
	"github.com/iamsuteerth/skyfox-backend/pkg/database/seed"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Warn().Msg("Warning: .env file not found")
	}

	// Initialize logger
	utils.InitLogger()

	// Get database connection
	db := config.GetDBConnection()
	defer config.CloseDBConnection()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	staffRepo := repositories.NewStaffRepository(db)

	// Seed the database
	seed.SeedDB(userRepo, staffRepo)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(userService)

	// Set up router
	router := gin.Default()

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		api.POST("/login", authController.Login)

		// Protected routes will come later
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("Server starting")
	if err := router.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
