package main

import (
	"net/http"
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
	passwordResetService := services.NewPasswordResetService(resetTokenRepository, skyCustomerRepository, userRepository)
	showService := services.NewShowService(showRepository, bookingRepository, movieService, slotRepository)
	slotService := services.NewSlotService(slotRepository)
	adminStaffProfileService := services.NewAdminStaffProfileService(userRepository, staffRepository)

	authController := controllers.NewAuthController(userService)
	skyCustomerController := controllers.NewSkyCustomerController(userService, skyCustomerService, securityQuestionService)
	securityQuestionController := controllers.NewSecurityQuestionController(securityQuestionService)
	passwordResetController := controllers.NewPasswordResetController(passwordResetService, skyCustomerService)
	showController := controllers.NewShowController(showService)
	slotController := controllers.NewSlotController(slotService)
	adminStaffController := controllers.NewAdminStaffController(adminStaffProfileService)

	binding.Validator = new(customValidator.DtoValidator)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		customValidator.RegisterCustomValidations(v)
	}

	router := gin.Default()
	router.Use(cors.SetupCORS())
	router.Use(security.APIKeyAuthMiddleware())

	noAuthRouter := router.Group("")

	authRouter := router.Group("")
	authRouter.Use(security.AuthMiddleware())

	customeRouter := router.Group("")
	customeRouter.Use(security.AuthMiddleware())
	customeRouter.Use(security.CustomerMiddleware())

	adminRouter := router.Group("")
	adminRouter.Use(security.AuthMiddleware())
	adminRouter.Use(security.AdminMiddleware())

	adminStaffRouter := router.Group("")
	adminStaffRouter.Use(security.AuthMiddleware())
	adminStaffRouter.Use(security.AdminStaffMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	noAuthAPIs := noAuthRouter.Group("")
	{
		login := noAuthAPIs.Group("")
		{
			login.POST(constants.LoginEndPoint, authController.Login)                            // Login
			login.POST(constants.ForgotPasswordEndPoint, passwordResetController.ForgotPassword) // Forgot Password
		}
		signup := noAuthAPIs.Group("")
		{
			signup.POST(constants.SkyCustomerSignUpEndPoint, skyCustomerController.Signup) // Customer Signup
		}
		securityQuestions := noAuthAPIs.Group(constants.SecurityQuestionsEndPoint)
		{
			securityQuestions.GET("", securityQuestionController.GetSecurityQuestions)                                      // Get all Security Questions
			securityQuestions.GET(constants.ByEmailEndPoint, securityQuestionController.GetSecurityQuestionByEmail)         // Get Security Question by Email
			securityQuestions.POST(constants.VerifySecurityAnswerEndPoint, securityQuestionController.VerifySecurityAnswer) // Verify Security Answer
		}
	}

	authAPIs := authRouter.Group("")
	{
		authAPIs.GET(constants.ShowEndPoint, showController.GetShows)                           // Get Shows (RBAC-based)
		authAPIs.POST(constants.ChangePasswordEndPoint, passwordResetController.ChangePassword) // Change Password for User
	}

	customerAPIs := customeRouter.Group(constants.SkyCustomerEndPoint)
	{
		customerAPIs.GET(constants.ProfileEndPoint, skyCustomerController.GetCustomerProfile)               // Get Customer Profile
		customerAPIs.GET(constants.ProfileImageEndPoint, skyCustomerController.GetProfileImagePresignedURL) // Get Profile Image
		customerAPIs.POST(constants.UpdateProfileEndPoint, skyCustomerController.UpdateCustomerProfile)     // Update Customer Profile
		customerAPIs.POST(constants.UpdateProfileImageEndPoint, skyCustomerController.UpdateProfileImage)   // Update Customer Profile Image
	}

	adminAPIs := adminRouter.Group("")
	{
		adminAPIs.GET(constants.SlotEndPoint, slotController.GetAvailableSlots) // Get Available Slots

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

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "ERROR",
			"code":    "ROUTE_NOT_FOUND",
			"message": "This route does not exist",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Info().Str("port", port).Msg("Server starting")
	if err := router.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
