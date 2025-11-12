package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics
var (
	// HTTP Metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	HTTPRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path"},
	)

	HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path"},
	)

	// Database Metrics
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query latency in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2, 5},
		},
		[]string{"query_type", "table"},
	)

	DatabaseConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	DatabaseConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	DatabaseConnectionsMax = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_max",
			Help: "Maximum number of database connections",
		},
	)

	DatabaseQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"query_type", "status"},
	)

	// Posting Engine Metrics
	PostingEngineExecutions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "posting_engine_executions_total",
			Help: "Total number of posting engine executions",
		},
		[]string{"document_type", "event", "status"},
	)

	PostingEngineDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "posting_engine_duration_seconds",
			Help:    "Posting engine execution duration in seconds",
			Buckets: []float64{.01, .05, .1, .25, .5, 1, 2, 5, 10},
		},
		[]string{"document_type", "event"},
	)

	PostingEngineValidationErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "posting_engine_validation_errors_total",
			Help: "Total number of posting validation errors",
		},
		[]string{"document_type", "error_type"},
	)

	// Cache Metrics
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type", "key_prefix"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type", "key_prefix"},
	)

	CacheOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_operation_duration_seconds",
			Help:    "Cache operation duration in seconds",
			Buckets: []float64{.0001, .0005, .001, .005, .01, .05, .1},
		},
		[]string{"operation", "cache_type"},
	)

	// Rate Limiter Metrics
	RateLimitExceeded = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Total number of rate limit exceeded events",
		},
		[]string{"endpoint", "limit_type"},
	)

	// Business Metrics
	SalesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sales_total",
			Help: "Total number of sales transactions",
		},
		[]string{"organization_id", "payment_method", "status"},
	)

	SalesAmount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sales_amount_total",
			Help: "Total sales amount",
		},
		[]string{"organization_id", "currency"},
	)

	InvoicesCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "invoices_created_total",
			Help: "Total number of invoices created",
		},
		[]string{"organization_id", "type"},
	)

	// System Metrics
	GoRoutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_count",
			Help: "Number of goroutines",
		},
	)

	MemoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)

	MemoryAllocated = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_allocated_bytes",
			Help: "Memory allocated in bytes",
		},
	)

	// Background Job Metrics
	JobExecutionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "job_executions_total",
			Help: "Total number of background job executions",
		},
		[]string{"job_type", "status"},
	)

	JobExecutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "job_execution_duration_seconds",
			Help:    "Background job execution duration in seconds",
			Buckets: []float64{.1, .5, 1, 5, 10, 30, 60, 300, 600},
		},
		[]string{"job_type"},
	)

	JobQueueSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "job_queue_size",
			Help: "Number of jobs waiting in queue",
		},
		[]string{"queue_name"},
	)

	JobRetries = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "job_retries_total",
			Help: "Total number of job retries",
		},
		[]string{"job_type"},
	)
)

// MetricsMiddleware wraps HTTP handlers to collect metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code and size
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			size:           0,
		}

		// Track request size
		HTTPRequestSize.WithLabelValues(r.Method, r.URL.Path).Observe(float64(r.ContentLength))

		// Execute request
		next.ServeHTTP(wrapped, r)

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Record metrics
		HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		HTTPRequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			fmt.Sprintf("%d", wrapped.statusCode),
		).Inc()
		HTTPResponseSize.WithLabelValues(r.Method, r.URL.Path).Observe(float64(wrapped.size))
	})
}

// responseWriter wraps http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// RecordDatabaseQuery records a database query metric
func RecordDatabaseQuery(queryType, table string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	DatabaseQueryDuration.WithLabelValues(queryType, table).Observe(duration.Seconds())
	DatabaseQueriesTotal.WithLabelValues(queryType, status).Inc()
}

// RecordPostingExecution records a posting engine execution
func RecordPostingExecution(docType, event, status string, duration time.Duration) {
	PostingEngineExecutions.WithLabelValues(docType, event, status).Inc()
	PostingEngineDuration.WithLabelValues(docType, event).Observe(duration.Seconds())
}

// RecordCacheOperation records a cache operation
func RecordCacheOperation(operation, cacheType, keyPrefix string, duration time.Duration, hit bool) {
	CacheOperationDuration.WithLabelValues(operation, cacheType).Observe(duration.Seconds())

	if operation == "get" {
		if hit {
			CacheHits.WithLabelValues(cacheType, keyPrefix).Inc()
		} else {
			CacheMisses.WithLabelValues(cacheType, keyPrefix).Inc()
		}
	}
}

// UpdateDatabaseConnectionStats updates database connection pool stats
func UpdateDatabaseConnectionStats(active, idle, max int32) {
	DatabaseConnectionsActive.Set(float64(active))
	DatabaseConnectionsIdle.Set(float64(idle))
	DatabaseConnectionsMax.Set(float64(max))
}

// RecordJobExecution records a background job execution
func RecordJobExecution(jobType, status string, duration time.Duration) {
	JobExecutionsTotal.WithLabelValues(jobType, status).Inc()
	JobExecutionDuration.WithLabelValues(jobType).Observe(duration.Seconds())
}

// UpdateJobQueueSize updates the job queue size metric
func UpdateJobQueueSize(queueName string, size int) {
	JobQueueSize.WithLabelValues(queueName).Set(float64(size))
}

// RecordJobRetry records a job retry
func RecordJobRetry(jobType string) {
	JobRetries.WithLabelValues(jobType).Inc()
}
