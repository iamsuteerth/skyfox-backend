package response

type AdminStaffProfileResponse struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	CounterNo int    `json:"counter_no"`
	CreatedAt string `json:"created_at"`
}
