package response

type SecurityQuestionResponse struct {
	QuestionID int    `json:"question_id"`
	Question   string `json:"question"`
	Email      string `json:"email,omitempty"`
}

type VerifySecurityAnswerResponse struct {
	ResetToken string `json:"reset_token"`
	ExpiresIn  int    `json:"expires_in_seconds"`
}

type VerifySecurityAnswerWithoutTokenResponse struct {
	ValidAnswer bool `json:"security_answer_valid"`
}