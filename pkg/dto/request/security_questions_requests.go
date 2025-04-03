package request

type GetSecurityQuestionByEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifySecurityAnswerRequest struct {
	Email          string `json:"email" binding:"required,email"`
	SecurityAnswer string `json:"security_answer" binding:"required,securityAnswer"`
}
