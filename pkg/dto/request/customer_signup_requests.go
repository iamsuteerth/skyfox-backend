package request

type SignupRequest struct {
	Name        string `json:"name" binding:"required,customName"`
	Username    string `json:"username" binding:"required,customUsername"`
	Password    string `json:"password" binding:"required,customPassword"`
	PhoneNumber string `json:"number" binding:"required,customPhone"`
	Email       string `json:"email" binding:"required,email"`
	ProfileImg  []byte `json:"profile_img"`
}
