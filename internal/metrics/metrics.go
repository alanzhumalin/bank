package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	TotalHttpRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "bank_http_requests_total",
	}, []string{"path", "method", "status"})

	TimeRequests = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "bank_http_request_duration_seconds",
	}, []string{"path", "method", "status"})
)

func Register() {
	prometheus.MustRegister(TotalHttpRequests, TimeRequests)
}
