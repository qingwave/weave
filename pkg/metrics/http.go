package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"method", "path", "code"},
	)

	HTTPInflightRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_inflight_requests",
			Help: "Status of HTTP response",
		},
		[]string{"method", "path"},
	)

	HTTPRequestsDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_requests_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"method", "path"})
)

func init() {
	prometheus.Register(HTTPRequestsTotal)
	prometheus.Register(HTTPInflightRequests)
	prometheus.Register(HTTPRequestsDuration)
}
