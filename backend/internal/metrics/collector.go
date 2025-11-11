package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Collector provides business metrics for the application
type Collector struct {
	// Database metrics
	dbQueriesTotal       *prometheus.CounterVec
	dbQueryDuration      *prometheus.HistogramVec
	dbConnectionsOpen    prometheus.Gauge
	dbConnectionsIdle    prometheus.Gauge
	dbConnectionsInUse   prometheus.Gauge
	dbConnectionsWaitDuration *prometheus.HistogramVec

	// Business metrics
	salesTotal           *prometheus.CounterVec
	salesAmount          *prometheus.CounterVec
	authenticationTotal  *prometheus.CounterVec
	authenticationFailed *prometheus.CounterVec
	postingTotal         *prometheus.CounterVec
	postingFailed        *prometheus.CounterVec

	// Error metrics
	errorsTotal          *prometheus.CounterVec
}

// NewCollector creates a new metrics collector
func NewCollector(namespace string) *Collector {
	return &Collector{
		// Database metrics
		dbQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"operation", "table", "status"},
		),
		dbQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "db_query_duration_seconds",
				Help:      "Database query duration in seconds",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2, 5},
			},
			[]string{"operation", "table"},
		),
		dbConnectionsOpen: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_open",
				Help:      "Number of open database connections",
			},
		),
		dbConnectionsIdle: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_idle",
				Help:      "Number of idle database connections",
			},
		),
		dbConnectionsInUse: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_in_use",
				Help:      "Number of database connections in use",
			},
		),
		dbConnectionsWaitDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "db_connections_wait_duration_seconds",
				Help:      "Time spent waiting for a database connection",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
			},
			[]string{"acquired"},
		),

		// Business metrics
		salesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "sales_total",
				Help:      "Total number of sales transactions",
			},
			[]string{"organization_id", "status"},
		),
		salesAmount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "sales_amount_total",
				Help:      "Total sales amount",
			},
			[]string{"organization_id", "currency"},
		),
		authenticationTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "authentication_attempts_total",
				Help:      "Total number of authentication attempts",
			},
			[]string{"method", "status"},
		),
		authenticationFailed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "authentication_failures_total",
				Help:      "Total number of failed authentication attempts",
			},
			[]string{"method", "reason"},
		),
		postingTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "posting_operations_total",
				Help:      "Total number of posting operations",
			},
			[]string{"event_type", "status"},
		),
		postingFailed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "posting_failures_total",
				Help:      "Total number of failed posting operations",
			},
			[]string{"event_type", "reason"},
		),

		// Error metrics
		errorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "errors_total",
				Help:      "Total number of errors by type",
			},
			[]string{"error_type", "error_code"},
		),
	}
}

// Database metrics methods

// RecordDBQuery records a database query
func (c *Collector) RecordDBQuery(operation, table, status string, duration float64) {
	c.dbQueriesTotal.WithLabelValues(operation, table, status).Inc()
	c.dbQueryDuration.WithLabelValues(operation, table).Observe(duration)
}

// UpdateDBConnectionStats updates database connection pool statistics
func (c *Collector) UpdateDBConnectionStats(open, idle, inUse int) {
	c.dbConnectionsOpen.Set(float64(open))
	c.dbConnectionsIdle.Set(float64(idle))
	c.dbConnectionsInUse.Set(float64(inUse))
}

// RecordDBConnectionWait records time spent waiting for a connection
func (c *Collector) RecordDBConnectionWait(acquired bool, duration float64) {
	status := "false"
	if acquired {
		status = "true"
	}
	c.dbConnectionsWaitDuration.WithLabelValues(status).Observe(duration)
}

// Business metrics methods

// RecordSale records a sale transaction
func (c *Collector) RecordSale(orgID, status string) {
	c.salesTotal.WithLabelValues(orgID, status).Inc()
}

// RecordSaleAmount records a sale amount
func (c *Collector) RecordSaleAmount(orgID, currency string, amount float64) {
	// Note: For amounts, we use Add instead of Inc
	// The counter will track the total amount
	c.salesAmount.WithLabelValues(orgID, currency).Add(amount)
}

// RecordAuthenticationAttempt records an authentication attempt
func (c *Collector) RecordAuthenticationAttempt(method, status string) {
	c.authenticationTotal.WithLabelValues(method, status).Inc()
}

// RecordAuthenticationFailure records a failed authentication
func (c *Collector) RecordAuthenticationFailure(method, reason string) {
	c.authenticationFailed.WithLabelValues(method, reason).Inc()
}

// RecordPosting records a posting operation
func (c *Collector) RecordPosting(eventType, status string) {
	c.postingTotal.WithLabelValues(eventType, status).Inc()
}

// RecordPostingFailure records a failed posting operation
func (c *Collector) RecordPostingFailure(eventType, reason string) {
	c.postingFailed.WithLabelValues(eventType, reason).Inc()
}

// Error metrics methods

// RecordError records an error
func (c *Collector) RecordError(errorType, errorCode string) {
	c.errorsTotal.WithLabelValues(errorType, errorCode).Inc()
}

// GetCollectors returns all prometheus collectors for registration
func (c *Collector) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		c.dbQueriesTotal,
		c.dbQueryDuration,
		c.dbConnectionsOpen,
		c.dbConnectionsIdle,
		c.dbConnectionsInUse,
		c.dbConnectionsWaitDuration,
		c.salesTotal,
		c.salesAmount,
		c.authenticationTotal,
		c.authenticationFailed,
		c.postingTotal,
		c.postingFailed,
		c.errorsTotal,
	}
}
