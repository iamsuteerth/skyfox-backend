package response

type BulkCheckInResponse struct {
	CheckedIn   []int `json:"checked_in"`
	AlreadyDone []int `json:"already_done"`
	Invalid     []int `json:"invalid"`
}
