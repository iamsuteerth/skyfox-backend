package request

type GetSecurityQuestionByEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}
