package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/go-playground/validator/v10"
	"github.com/iamsuteerth/skyfox-backend/pkg/config"
	"github.com/iamsuteerth/skyfox-backend/pkg/constants"
	"github.com/iamsuteerth/skyfox-backend/pkg/controllers"
	"github.com/iamsuteerth/skyfox-backend/pkg/database/seed"
	"github.com/iamsuteerth/skyfox-backend/pkg/middleware/cors"
	"github.com/iamsuteerth/skyfox-backend/pkg/middleware/security"
	customValidator "github.com/iamsuteerth/skyfox-backend/pkg/middleware/validator"
	movieservice "github.com/iamsuteerth/skyfox-backend/pkg/movie-service"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/services"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Warn().Msg("Warning: .env file not found")
	}

	utils.InitLogger()

	db := config.GetDBConnection()
	defer config.CloseDBConnection()

	movieServiceConfig := config.GetMovieServiceConfig()
	movieService := movieservice.NewMovieService(movieServiceConfig)

	userRepository := repositories.NewUserRepository(db)
	staffRepository := repositories.NewStaffRepository(db)
	skyCustomerRepository := repositories.NewSkyCustomerRepository(db)
	securityQuestionRepository := repositories.NewSecurityQuestionRepository(db)
	resetTokenRepository := repositories.NewResetTokenRepository(db)
	showRepository := repositories.NewShowRepository(db)
	bookingRepository := repositories.NewBookingRepository(db)
	slotRepository := repositories.NewSlotRepository(db)

	seed.SeedDB(userRepository, staffRepository)

	userService := services.NewUserService(userRepository)
	skyCustomerService := services.NewSkyCustomerService(skyCustomerRepository, userRepository, securityQuestionRepository)
	securityQuestionService := services.NewSecurityQuestionService(securityQuestionRepository, skyCustomerRepository, resetTokenRepository)
	forgotPasswordService := services.NewForgotPasswordService(resetTokenRepository, skyCustomerRepository, userRepository)
	showService := services.NewShowService(showRepository, bookingRepository, movieService, slotRepository)
	slotService := services.NewSlotService(slotRepository)

	authController := controllers.NewAuthController(userService)
	skyCustomerController := controllers.NewSkyCustomerController(userService, skyCustomerService, securityQuestionService)
	securityQuestionController := controllers.NewSecurityQuestionController(securityQuestionService)
	forgotPasswordController := controllers.NewForgotPasswordController(forgotPasswordService)
	showController := controllers.NewShowController(showService)
	slotController := controllers.NewSlotController(slotService)

	binding.Validator = new(customValidator.DtoValidator)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		customValidator.RegisterCustomValidations(v)
	}

	router := gin.Default()
	router.Use(cors.SetupCORS())

	noAuthRouter := router.Group("")

	authRouter := router.Group("")
	authRouter.Use(security.AuthMiddleware())

	adminRouter := router.Group("")
	adminRouter.Use(security.AuthMiddleware())
	adminRouter.Use(security.AdminMiddleware())

	staffRouter := router.Group("")
	staffRouter.Use(security.AuthMiddleware())
	staffRouter.Use(security.StaffMiddleware())

	// Login
	noAuthRouter.POST(constants.LoginEndPoint, authController.Login)
	// Customer Signup
	noAuthRouter.POST(constants.SkyCustomerSignUpEndPoint, skyCustomerController.Signup)
	// Get Security Questions
	noAuthRouter.GET(constants.SecurityQuestions, securityQuestionController.GetSecurityQuestions)
	noAuthRouter.GET(constants.SecurityQuestionByEmail, securityQuestionController.GetSecurityQuestionByEmail)
	// Verify Security Question Answer
	noAuthRouter.POST(constants.VerifySecurityAnswerEndpoint, securityQuestionController.VerifySecurityAnswer)
	// Forgot Password for Customer
	noAuthRouter.POST(constants.ForgotPasswordEndpoint, forgotPasswordController.ForgotPassword)
	// Get Shows (RBAC)
	authRouter.GET(constants.ShowEndPoint, showController.GetShows)
	// Show Creation - Admin
	showCreation := adminRouter.Group(constants.ShowEndPoint)
	{
		showCreation.GET(constants.MoviesEndPoint, showController.GetMovies)
		showCreation.POST("", showController.CreateShow)
	}
	// Get Available Slots - Admin
	adminRouter.GET(constants.SlotEndPoint, slotController.GetAvailableSlots)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("Server starting")
	if err := router.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
