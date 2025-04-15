package response

type CustomerProfileResponse struct {
	Username               string `json:"username"`
	Name                   string `json:"name"`
	Email                  string `json:"email"`
	PhoneNumber            string `json:"phone_number"`
	SecurityQuestionExists bool   `json:"security_question_exists"`
	CreatedAt              string `json:"created_at"`
}

type UpdateCustomerProfileResponse struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
