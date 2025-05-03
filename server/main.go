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
	paymentservice "github.com/iamsuteerth/skyfox-backend/pkg/payment-service"
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
	paymentServiceConfig := config.GetPaymentServiceConfig()
	paymentService := paymentservice.NewPaymentService(paymentServiceConfig)
	s3Service := services.NewS3Service()

	userRepository := repositories.NewUserRepository(db)
	staffRepository := repositories.NewStaffRepository(db)
	skyCustomerRepository := repositories.NewSkyCustomerRepository(db)
	securityQuestionRepository := repositories.NewSecurityQuestionRepository(db)
	resetTokenRepository := repositories.NewResetTokenRepository(db)
	showRepository := repositories.NewShowRepository(db)
	bookingRepository := repositories.NewBookingRepository(db)
	slotRepository := repositories.NewSlotRepository(db)
	bookingSeatMappingRepository := repositories.NewBookingSeatMappingRepository(db)
	adminBookedCustomerRepository := repositories.NewAdminBookedCustomerRepository(db)
	pendingBookingRepository := repositories.NewPendingBookingRepository(db)
	paymentTransactionRepository := repositories.NewPaymentTransactionRepository(db)

	seed.SeedDB(userRepository, staffRepository)

	userService := services.NewUserService(userRepository)
	skyCustomerService := services.NewSkyCustomerService(skyCustomerRepository, userRepository, securityQuestionRepository, s3Service)
	securityQuestionService := services.NewSecurityQuestionService(securityQuestionRepository, skyCustomerRepository, resetTokenRepository)
	passwordResetService := services.NewPasswordResetService(resetTokenRepository, skyCustomerRepository, userRepository)
	showService := services.NewShowService(showRepository, bookingRepository, movieService, slotRepository)
	slotService := services.NewSlotService(slotRepository)
	adminStaffProfileService := services.NewAdminStaffProfileService(userRepository, staffRepository)
	bookingService := services.NewBookingService(showRepository, bookingRepository, bookingSeatMappingRepository, slotRepository, adminBookedCustomerRepository, skyCustomerRepository, movieService)
	adminBookingService := services.NewAdminBookingService(showRepository, bookingRepository, bookingSeatMappingRepository, adminBookedCustomerRepository, slotRepository)
	customerBookingService := services.NewCustomerBookingService(showRepository, bookingRepository, bookingSeatMappingRepository, pendingBookingRepository, paymentTransactionRepository, slotRepository, skyCustomerRepository, paymentService)
	checkInService := services.NewCheckInService(bookingRepository, showRepository)
	revenueService := services.NewRevenueService(bookingRepository, showRepository, slotRepository, movieService)

	authController := controllers.NewAuthController(userService)
	skyCustomerController := controllers.NewSkyCustomerController(userService, skyCustomerService, securityQuestionService)
	securityQuestionController := controllers.NewSecurityQuestionController(securityQuestionService)
	passwordResetController := controllers.NewPasswordResetController(passwordResetService, skyCustomerService)
	showController := controllers.NewShowController(showService)
	slotController := controllers.NewSlotController(slotService)
	adminStaffController := controllers.NewAdminStaffController(adminStaffProfileService)
	bookingController := controllers.NewBookingController(bookingService, adminBookingService, customerBookingService, checkInService)
	revenueController := controllers.NewDashboardRevenueController(revenueService)

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
		authAPIs.POST(constants.ChangePasswordEndPoint, passwordResetController.ChangePassword) // Change Password for User

		showsAPIs := authAPIs.Group(constants.ShowsEndPoint)
		{
			showsAPIs.GET("", showController.GetShows)                                    // Get Shows (RBAC-based)
			showsAPIs.GET(constants.BookingSeatMapEndPoint, bookingController.GetSeatMap) // Get Seat Map Data
		}

		showAPIs := authAPIs.Group(constants.ShowEndPoint)
		{
			showAPIs.GET("", showController.GetShowById)
		}

		bookingHelpers := authAPIs.Group(constants.BookingIdEndpoint)
		{
			bookingHelpers.GET(constants.QREndpoint, bookingController.GetQRCode) // Get QR Code Image as Base64
			bookingHelpers.GET(constants.PDFEndpoint, bookingController.GetPDF)   // // Get PDF as Base64
		}
	}

	customerAPIs := customeRouter.Group(constants.SkyCustomerEndPoint)
	{
		customerAPIs.GET(constants.ProfileEndPoint, skyCustomerController.GetCustomerProfile)               // Get Customer Profile
		customerAPIs.GET(constants.ProfileImageEndPoint, skyCustomerController.GetProfileImagePresignedURL) // Get Profile Image
		customerAPIs.POST(constants.UpdateProfileEndPoint, skyCustomerController.UpdateCustomerProfile)     // Update Customer Profile
		customerAPIs.POST(constants.UpdateProfileImageEndPoint, skyCustomerController.UpdateProfileImage)   // Update Customer Profile Image

		booking := customerAPIs.Group(constants.BookingEndpoint)
		{
			booking.POST(constants.BookingInitializeEndpoint, bookingController.InitializeCustomerBooking) // Initialize Booking
			booking.POST(constants.PaymentEndpoint, bookingController.ProcessPayment)                      // Handle Payment for Booking
			booking.DELETE(constants.CancelBookingEndpoint, bookingController.CancelBooking)               // Prematurely Cancel Pending Booking
		}

		bookings := customerAPIs.Group(constants.BookingsEndpoint)
		{
			bookings.GET("", bookingController.GetCustomerBookings)                                    // Get all customer bookings
			bookings.GET(constants.LatestBookingsEndpoint, bookingController.GetCustomerLatestBooking) // Get latest customer booking
		}
	}

	adminAPIs := adminRouter.Group("")
	{
		adminAPIs.GET(constants.SlotEndPoint, slotController.GetAvailableSlots) // Get Available Slots

		showCreation := adminAPIs.Group(constants.ShowEndPoint)
		{
			showCreation.GET(constants.MoviesEndPoint, showController.GetMovies) // Get Movies for Show creation
			showCreation.POST("", showController.CreateShow)                     // Create a Show
		}

		bookingAPIs := adminAPIs.Group(constants.AdminEndPoint)
		{
			bookingAPIs.POST(constants.CreateCustomerBookingEndpoint, bookingController.CreateAdminBooking) // Create Booking Through Admin
		}

		revenueAPIs := adminAPIs.Group(constants.RevenueEndpoint)
		{
			revenueAPIs.GET("", revenueController.GetRevenue) // Revenue API with necessary query params
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

		checkin := adminStaffAPIs.Group(constants.CheckinEndpoint)
		{
			checkin.GET(constants.BookingsEndpoint, bookingController.GetCheckInBookings)   // Get all confirmed bookings
			checkin.POST(constants.BookingsEndpoint, bookingController.BulkCheckInBookings) // Mark bookings as checked-in in bulk
			checkin.POST(constants.BookingEndpoint, bookingController.SingleCheckInBooking) // Mark a booking as checked-in
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
		port = "8080"
	}

	log.Info().Str("port", port).Msg("Server starting")
	if err := router.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
