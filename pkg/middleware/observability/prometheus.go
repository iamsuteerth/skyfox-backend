package observability

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamsuteerth/skyfox-backend/pkg/metrics"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		metrics.HttpRequestsInFlight.Inc()
		
		c.Next()
	
		endpointGroup := getEndpointGroup(c.FullPath())

		duration := time.Since(start).Seconds()
		method := c.Request.Method
		statusCode := strconv.Itoa(c.Writer.Status())

		if endpointGroup == "" {
            metrics.HttpRequestsInFlight.Dec()
            return
        }
		
		metrics.HttpRequestsTotal.WithLabelValues(method, endpointGroup, statusCode).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, endpointGroup).Observe(duration)
		
		metrics.HttpRequestsInFlight.Dec()
	}
}

func getEndpointGroup(path string) string {
	switch {
	// Authentication & Security
	case path == "/login":
		return "auth"
	case path == "/forgot-password":
		return "auth"
	case path == "/change-password":
		return "auth"
	case strings.HasPrefix(path, "/security-questions"):
		return "auth"
		
	// Customer Management
	case path == "/customer/signup":
		return "customer_mgmt"
	case strings.HasPrefix(path, "/customer/profile"):
		return "customer_mgmt"
	case strings.HasPrefix(path, "/customer/update-profile"):
		return "customer_mgmt"
		
	// Customer Wallet
	case strings.HasPrefix(path, "/customer/wallet"):
		return "wallet"
		
	// Show & Movie Management
	case path == "/shows" || path == "/show":
		return "shows"
	case path == "/show/movies":
		return "shows"
	case strings.HasPrefix(path, "/slot"):
		return "shows"
		
	// Booking Operations
	case strings.HasPrefix(path, "/customer/booking"):
		return "booking"
	case strings.HasPrefix(path, "/admin/create-customer-booking"):
		return "booking"
	case strings.HasPrefix(path, "/shows/") && strings.Contains(path, "seat-map"):
		return "booking"
	case strings.HasPrefix(path, "/booking/") && (strings.Contains(path, "/qr") || strings.Contains(path, "/pdf")):
		return "booking"
	case strings.HasPrefix(path, "/customer/bookings"):
		return "booking"
		
	// Check-in Operations
	case strings.HasPrefix(path, "/check-in"):
		return "checkin"
		
	// Admin Operations
	case strings.HasPrefix(path, "/admin/profile") || strings.HasPrefix(path, "/staff/profile"):
		return "admin"
	case path == "/revenue":
		return "admin"
	case path == "/booking-csv":
		return "admin"
		
	// System endpoints
	case path == "/health":
		return ""
	case path == "/metrics":
		return ""
		
	default:
		return "other"
	}
}
