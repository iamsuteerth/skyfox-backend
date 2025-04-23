package constants

const (
	// No Auth Routes
	LoginEndPoint                = "/login"
	SkyCustomerSignUpEndPoint    = "/customer/signup"
	SecurityQuestionsEndPoint    = "/security-questions"
	ByEmailEndPoint              = "/by-email"
	VerifySecurityAnswerEndPoint = "/verify-security-answer"
	ForgotPasswordEndPoint       = "/forgot-password"
	// Shows Page
	ShowEndPoint           = "/shows"
	SlotEndPoint           = "/slot"
	MoviesEndPoint         = "/movies"
	BookingSeatMapEndPoint = "/:show_id/seat-map"
	// Role Endpoints
	SkyCustomerEndPoint = "/customer"
	AdminEndPoint       = "/admin"
	StaffEndPoint       = "/staff"
	// Profile Page
	ChangePasswordEndPoint     = "/change-password"
	ProfileImageEndPoint       = "/profile-image"
	ProfileEndPoint            = "/profile"
	UpdateProfileEndPoint      = "/update-profile"
	UpdateProfileImageEndPoint = "/update-profile-image"
	// Booking Related Endpoints
	GetBookingsEndpoint           = "/bookings"
	CreateCustomerBookingEndpoint = "/create-customer-booking"
	BookingEndpoint               = "/booking"
	BookingInitializeEndpoint     = "/initialize"
	PaymentEndpoint               = "/payment"
)

const (
	TOTAL_NO_OF_SEATS           = 100
	MAX_NO_OF_SEATS_PER_BOOKING = 10
	DELUXE_OFFSET               = 150.0
)
