package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alanzhumalin/bank/internal/metrics"
)

type recorderWithStatus struct {
	statusCode int
	http.ResponseWriter
}

func (r *recorderWithStatus) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func MetricsMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := &recorderWithStatus{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(rec, r)

			statusCode := strconv.Itoa(rec.statusCode)

			metrics.TotalHttpRequests.WithLabelValues(r.URL.Path, r.Method, statusCode).Inc()
			metrics.TimeRequests.WithLabelValues(r.URL.Path, r.Method, statusCode).Observe(time.Since(start).Seconds())
		})
	}
}
