package response

type LoginResponse struct {
	User  UserInfo `json:"user"`
	Token string   `json:"token"`
}

type UserInfo struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
