package response

type RevenueGroupStats struct {
	Label            string  `json:"label"` // E.g. date, week, month, year, movie name, genre, slot label
	TotalRevenue     float64 `json:"total_revenue"`
	MeanRevenue      float64 `json:"mean_revenue"`
	MedianRevenue    float64 `json:"median_revenue"`
	TotalBookings    int     `json:"total_bookings"`
	TotalSeatsBooked int     `json:"total_seats_booked"`
}

type RevenueDashboardResponse struct {
	// For selected/filtered context (not just "all")
	TotalRevenue     float64 `json:"total_revenue"`
	MeanRevenue      float64 `json:"mean_revenue"`
	MedianRevenue    float64 `json:"median_revenue"`
	TotalBookings    int     `json:"total_bookings"`
	TotalSeatsBooked int     `json:"total_seats_booked"`
	// Breakdown per group (timeframe, movie, slot, genre, etc.)
	Groups []RevenueGroupStats `json:"groups"`
	// For comparative dashboards-if e.g. no filter is sent on movie, will have one for each movie, etc.
}
