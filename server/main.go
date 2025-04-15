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
	s3Service := services.NewS3Service()

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
	skyCustomerService := services.NewSkyCustomerService(skyCustomerRepository, userRepository, securityQuestionRepository, s3Service)
	securityQuestionService := services.NewSecurityQuestionService(securityQuestionRepository, skyCustomerRepository, resetTokenRepository)
	forgotPasswordService := services.NewForgotPasswordService(resetTokenRepository, skyCustomerRepository, userRepository)
	showService := services.NewShowService(showRepository, bookingRepository, movieService, slotRepository)
	slotService := services.NewSlotService(slotRepository)
	adminStaffProfileService := services.NewAdminStaffProfileService(userRepository, staffRepository)

	authController := controllers.NewAuthController(userService)
	skyCustomerController := controllers.NewSkyCustomerController(userService, skyCustomerService, securityQuestionService)
	securityQuestionController := controllers.NewSecurityQuestionController(securityQuestionService)
	forgotPasswordController := controllers.NewForgotPasswordController(forgotPasswordService)
	showController := controllers.NewShowController(showService)
	slotController := controllers.NewSlotController(slotService)
	adminStaffController := controllers.NewAdminStaffController(adminStaffProfileService)

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

	adminStaffRouter := router.Group("")
	adminStaffRouter.Use(security.AuthMiddleware())
	adminStaffRouter.Use(security.AdminStaffMiddleware())

	noAuthAPIs := noAuthRouter.Group("")
	{
		login := noAuthAPIs.Group("")
		{
			login.POST(constants.LoginEndPoint, authController.Login)                             // Login
			login.POST(constants.ForgotPasswordEndPoint, forgotPasswordController.ForgotPassword) // Forgot Password
		}
		signup := noAuthAPIs.Group("")
		{
			signup.POST(constants.SkyCustomerSignUpEndPoint, skyCustomerController.Signup) // Customer Signup
		}
		securityQuestions := noAuthAPIs.Group(constants.SecurityQuestionsEndPoint)
		{
			securityQuestions.GET(constants.SecurityQuestionsEndPoint, securityQuestionController.GetSecurityQuestions)     // Get all Security Questions
			securityQuestions.GET(constants.ByEmailEndPoint, securityQuestionController.GetSecurityQuestionByEmail)         // Get Security Question by Email
			securityQuestions.POST(constants.VerifySecurityAnswerEndPoint, securityQuestionController.VerifySecurityAnswer) // Verify Security Answer
		}
	}

	authAPIs := authRouter.Group("")
	{
		shows := authAPIs.Group("")
		{
			shows.GET(constants.ShowEndPoint, showController.GetShows) // Get Shows (RBAC-based)
		}
		customer := authAPIs.Group(constants.SkyCustomerEndPoint)
		{
			customer.GET(constants.ProfileImageEndPoint, skyCustomerController.GetProfileImagePresignedURL) // Get Profile Image
			customer.GET(constants.ProfileEndPoint, skyCustomerController.GetCustomerProfile)               // Get customer profile details
		}
	}

	adminAPIs := adminRouter.Group("")
	{
		// Slot Management
		adminAPIs.GET(constants.SlotEndPoint, slotController.GetAvailableSlots) // Get Available Slots

		// Show Management
		showCreation := adminAPIs.Group(constants.ShowEndPoint)
		{
			showCreation.GET(constants.MoviesEndPoint, showController.GetMovies) // Get Movies for Show creation
			showCreation.POST("", showController.CreateShow)                     // Create a Show
		}
	}

	adminStaffAPIs := adminStaffRouter.Group("")
	{
		admin := adminStaffAPIs.Group(constants.AdminEndPoint)
		{
			admin.GET(constants.ProfileEndPoint, adminStaffController.GetAdminProfile) // Get Admin Profile
		}
		staff := adminStaffAPIs.Group(constants.StaffEndPoint)
		{
			staff.GET(constants.ProfileEndPoint, adminStaffController.GetStaffProfile) // Get Staff Profile
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("Server starting")
	if err := router.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
