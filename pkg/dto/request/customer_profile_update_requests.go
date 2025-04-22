package request

type UpdateCustomerProfileRequest struct {
	Name           string `json:"name" binding:"required,customName"`
	Email          string `json:"email" binding:"required,email"`
	PhoneNumber    string `json:"phone_number" binding:"required,customPhone"`
	SecurityAnswer string `json:"security_answer" binding:"required"`
}

type UpdateProfileImageRequest struct {
    SecurityAnswer string `json:"security_answer" binding:"required,securityAnswer"`
    ProfileImg     string `json:"profile_img" binding:"required"`
    ProfileImgSHA  string `json:"profile_img_sha" binding:"required"`
}