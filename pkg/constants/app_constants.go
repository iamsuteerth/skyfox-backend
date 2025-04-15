package constants

const (
	// No Auth Routes
	LoginEndPoint                = "/login"
	SkyCustomerSignUpEndPoint    = "/customer/signup"
	SecurityQuestionsEndPoint    = "/security-questions"
	SecurityQuestionByEmail      = "/security-questions/by-email"
	VerifySecurityAnswerEndPoint = "/verify-security-answer"
	ForgotPasswordEndPoint       = "/forgot-password"
	// Shows Page
	ShowEndPoint   = "/shows"
	SlotEndPoint   = "/slot"
	MoviesEndPoint = "/movies"
	// Role Endpoints
	CustomerEndPoint = "/customer"
	AdminEndPoint    = "/admin"
	StaffEndPoint    = "/staff"
	// Profile Page
	ChangePasswordEndPoint     = "/change-password"
	ProfileImageEndPoint       = "/profile-image"
	ProfileEndPoint            = "/profile"
	UpdateProfileEndPoint      = "/update-profile"
	UpdateProfileImageEndPoint = "/update-profile-image"
)

const (
	TOTAL_NO_OF_SEATS           = 100
	MAX_NO_OF_SEATS_PER_BOOKING = 10
)
