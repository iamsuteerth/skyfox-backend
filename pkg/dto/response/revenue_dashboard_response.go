package response

import "github.com/govalues/decimal"

type RevenueGroupStats struct {
	Label            string          `json:"label"`
	TotalRevenue     decimal.Decimal `json:"total_revenue"`
	MeanRevenue      decimal.Decimal `json:"mean_revenue"`
	MedianRevenue    decimal.Decimal `json:"median_revenue"`
	TotalBookings    int             `json:"total_bookings"`
	TotalSeatsBooked int             `json:"total_seats_booked"`
}

type RevenueDashboardResponse struct {
	TotalRevenue     decimal.Decimal     `json:"total_revenue"`
	MeanRevenue      decimal.Decimal     `json:"mean_revenue"`
	MedianRevenue    decimal.Decimal     `json:"median_revenue"`
	TotalBookings    int                 `json:"total_bookings"`
	TotalSeatsBooked int                 `json:"total_seats_booked"`
	Groups           []RevenueGroupStats `json:"groups"`
}
