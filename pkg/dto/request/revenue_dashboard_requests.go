package request

type RevenueDashboardRequest struct {
	Timeframe  string   `form:"timeframe" binding:"omitempty,oneof=daily weekly monthly yearly all"`
	Month      *int     `form:"month" binding:"omitempty,min=1,max=12"`
	Year       *int     `form:"year" binding:"omitempty"`
	MovieID    *string  `form:"movie_id" binding:"omitempty"`
	SlotID     *int     `form:"slot_id" binding:"omitempty"`
	Genre      *string  `form:"genre" binding:"omitempty"`
	ParamOrder []string `form:"-"`
}
