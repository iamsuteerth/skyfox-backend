package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "skyfox_http_requests_total",
			Help: "Total HTTP requests by endpoint group and status",
		},
		[]string{"method", "endpoint_group", "status_code"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "skyfox_http_request_duration_seconds",
			Help: "HTTP request duration in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "endpoint_group"},
	)

	HttpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "skyfox_http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)
)

func InitMetrics() {
}
