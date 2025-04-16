package response

type ProfileImageResponse struct {
	PresignedURL string `json:"presigned_url"`
	ExpiresAt    string `json:"expires_at"`
}

type UpdateProfileImageResponse struct {
    Username string `json:"username"`
}