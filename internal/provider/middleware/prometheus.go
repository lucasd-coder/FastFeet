package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	httpRequestsCount    = "requests_total"
	httpRequestsDuration = "request_duration_seconds"
)

var cfg = struct {
	Namespace string
	Buckets   []float64
	Subsystem string
}{
	Namespace: "chi",
	Subsystem: "http",
	Buckets: []float64{
		0.0005,
		0.001, // 1ms
		0.002,
		0.005,
		0.01, // 10ms
		0.02,
		0.05,
		0.1, // 100 ms
		0.2,
		0.5,
		1.0, // 1s
		2.0,
		5.0,
		10.0, // 10s
		15.0,
		20.0,
		30.0,
	},
}

func PromMiddleware(next http.Handler) http.Handler {
	normalizeHTTPStatus := func(status int) string {
		switch {
		case status < http.StatusOK:
			return "1xx"
		case status < http.StatusMultipleChoices:
			return "2xx"
		case status < http.StatusBadRequest:
			return "3xx"
		case status < http.StatusInternalServerError:
			return "4xx"
		default:
			return "5xx"
		}
	}

	httpRequests := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      httpRequestsCount,
		Help:      "Number of HTTP operations",
	}, []string{"status", "method", "path"})

	httpDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      httpRequestsDuration,
		Help:      "Spend time by processing a route",
		Buckets:   cfg.Buckets,
	}, []string{"status", "method", "path"})

	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		path := chi.RouteContext(r.Context()).RoutePattern()

		status := normalizeHTTPStatus(ww.Status())

		httpRequests.WithLabelValues(status, r.Method, path).Inc()
		prometheus.NewTimer(httpDuration.WithLabelValues(status, r.Method, path)).ObserveDuration()
	}

	return http.HandlerFunc(fn)
}
