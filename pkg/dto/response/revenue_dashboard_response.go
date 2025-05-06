package response

type RevenueGroupStats struct {
	Label            string  `json:"label"` 
	TotalRevenue     float64 `json:"total_revenue"`
	MeanRevenue      float64 `json:"mean_revenue"`
	MedianRevenue    float64 `json:"median_revenue"`
	TotalBookings    int     `json:"total_bookings"`
	TotalSeatsBooked int     `json:"total_seats_booked"`
}

type RevenueDashboardResponse struct {
	TotalRevenue     float64 `json:"total_revenue"`
	MeanRevenue      float64 `json:"mean_revenue"`
	MedianRevenue    float64 `json:"median_revenue"`
	TotalBookings    int     `json:"total_bookings"`
	TotalSeatsBooked int     `json:"total_seats_booked"`
	Groups []RevenueGroupStats `json:"groups"`
}
