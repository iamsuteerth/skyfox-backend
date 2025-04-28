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
	ShowsEndPoint          = "/shows"
	ShowEndPoint           = "show"
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
	BookingsEndpoint              = "/bookings"
	LatestBookingsEndpoint        = "/latest"
	CreateCustomerBookingEndpoint = "/create-customer-booking"
	BookingEndpoint               = "/booking"
	BookingIdEndpoint             = "/booking/:id"
	QREndpoint                    = "/qr"
	PDFEndpoint                   = "/pdf"
	BookingInitializeEndpoint     = "/initialize"
	CancelBookingEndpoint         = "/:id/cancel"
	PaymentEndpoint               = "/payment"
)

const (
	TOTAL_NO_OF_SEATS           = 100
	MAX_NO_OF_SEATS_PER_BOOKING = 10
	DELUXE_OFFSET               = 150.0
)
