package v1

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/takama/bit"
)

var (
	totalDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of http requests duration in seconds.",
			Buckets: []float64{0.0001, 0.001, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 7.5, 10},
		}, 
		[]string{"status"},
	)

	totalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of http requests.",
		},
		[]string{"status"},
	)
)

// MetricsFunc returns a func for work with Prometheus
func (h *Handler) MetricsFunc() func(c bit.Control) {
	handler := promhttp.Handler()

	return func(c bit.Control) {
		c.Code(http.StatusOK)
		handler.ServeHTTP(c, c.Request())
	}
}

func init() {
	prometheus.MustRegister(totalDuration)
	prometheus.MustRegister(totalCounter)
}