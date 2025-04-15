package request

type UpdateCustomerProfileRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	PhoneNumber    string `json:"phone_number" binding:"required,customPhone"`
	SecurityAnswer string `json:"security_answer" binding:"required"`
}
