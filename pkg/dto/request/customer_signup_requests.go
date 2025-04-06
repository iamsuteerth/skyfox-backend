package request

type SignupRequest struct {
	Name               string `json:"name" binding:"required,customName"`
	Username           string `json:"username" binding:"required,customUsername"`
	Password           string `json:"password" binding:"required,customPassword"`
	PhoneNumber        string `json:"number" binding:"required,customPhone"`
	Email              string `json:"email" binding:"required,email"`
	ProfileImg         string `json:"profile_img"`
	ProfileImgSHA      string `json:"profile_img_sha"`
	SecurityQuestionID int    `json:"security_question_id" binding:"required"`
	SecurityAnswer     string `json:"security_answer" binding:"required,securityAnswer"`
}
