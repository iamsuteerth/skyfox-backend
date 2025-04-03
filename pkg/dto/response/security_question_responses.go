package response

type SecurityQuestionResponse struct {
	QuestionID int    `json:"question_id"`
	Question   string `json:"question"`
	Email      string `json:"email,omitempty"`
}
