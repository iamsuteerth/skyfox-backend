package request

type BulkCheckInRequest struct {
	BookingIDs []int `json:"booking_ids" binding:"required,min=1,dive,required"`
}

type SingleCheckInRequest struct {
	BookingID int `json:"booking_id" binding:"required"`
}
