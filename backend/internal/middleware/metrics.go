package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsMiddleware provides Prometheus metrics for HTTP requests
type MetricsMiddleware struct {
	requestsTotal      *prometheus.CounterVec
	requestDuration    *prometheus.HistogramVec
	requestSize        *prometheus.HistogramVec
	responseSize       *prometheus.HistogramVec
	requestsInFlight   prometheus.Gauge
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(namespace string) *MetricsMiddleware {
	m := &MetricsMiddleware{
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request latency in seconds",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path", "status"},
		),
		requestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_size_bytes",
				Help:      "HTTP request size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 7), // 100 bytes to ~100MB
			},
			[]string{"method", "path"},
		),
		responseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 7), // 100 bytes to ~100MB
			},
			[]string{"method", "path"},
		),
		requestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "http_requests_in_flight",
				Help:      "Current number of HTTP requests being served",
			},
		),
	}

	return m
}

// Handler returns the middleware handler
func (m *MetricsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Increment in-flight requests
		m.requestsInFlight.Inc()
		defer m.requestsInFlight.Dec()

		// Record request size
		if r.ContentLength > 0 {
			m.requestSize.WithLabelValues(r.Method, r.URL.Path).Observe(float64(r.ContentLength))
		}

		// Wrap response writer to capture status code and response size
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Serve the request
		next.ServeHTTP(ww, r)

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get status code
		status := strconv.Itoa(ww.Status())

		// Get path pattern (for better cardinality)
		path := r.URL.Path

		// Try to get route pattern from chi context
		// This provides better cardinality (e.g., /users/{id} instead of /users/123)
		if routePattern := middleware.GetReqID(r.Context()); routePattern != "" {
			// For chi router, we can use the route pattern
			// In production, you might want to use chi.RouteContext(r.Context()).RoutePattern()
			// For now, we'll use the actual path (can be improved later)
		}

		// Record metrics
		m.requestsTotal.WithLabelValues(r.Method, path, status).Inc()
		m.requestDuration.WithLabelValues(r.Method, path, status).Observe(duration)
		m.responseSize.WithLabelValues(r.Method, path).Observe(float64(ww.BytesWritten()))
	})
}

// GetMetrics returns the metrics for testing/inspection
func (m *MetricsMiddleware) GetMetrics() map[string]prometheus.Collector {
	return map[string]prometheus.Collector{
		"requests_total":      m.requestsTotal,
		"request_duration":    m.requestDuration,
		"request_size":        m.requestSize,
		"response_size":       m.responseSize,
		"requests_in_flight":  m.requestsInFlight,
	}
}
