package metrics

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// DBCollector collects database connection pool metrics
type DBCollector struct {
	pool *pgxpool.Pool

	// Metrics
	acquireCount      *prometheus.CounterVec
	acquireDuration   *prometheus.HistogramVec
	acquiredConns     prometheus.Gauge
	canceledAcquires  prometheus.Counter
	constructingConns prometheus.Gauge
	emptyAcquires     prometheus.Counter
	idleConns         prometheus.Gauge
	maxConns          prometheus.Gauge
	totalConns        prometheus.Gauge
	maxLifetimeDestroy prometheus.Counter
	maxIdleDestroy     prometheus.Counter
}

// NewDBCollector creates a new database metrics collector
func NewDBCollector(namespace string, pool *pgxpool.Pool) *DBCollector {
	collector := &DBCollector{
		pool: pool,
		acquireCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_pool_acquire_count",
				Help:      "Cumulative count of successful acquires from the pool",
			},
			[]string{"status"},
		),
		acquireDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "db_pool_acquire_duration_seconds",
				Help:      "Duration of successful acquire from the pool",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
			},
			[]string{},
		),
		acquiredConns: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_pool_acquired_conns",
				Help:      "Number of currently acquired connections in the pool",
			},
		),
		canceledAcquires: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_pool_canceled_acquires_total",
				Help:      "Cumulative count of acquires canceled by context",
			},
		),
		constructingConns: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_pool_constructing_conns",
				Help:      "Number of conns with construction in progress in the pool",
			},
		),
		emptyAcquires: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_pool_empty_acquires_total",
				Help:      "Cumulative count of successful acquires that waited for a resource to be released or constructed",
			},
		),
		idleConns: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_pool_idle_conns",
				Help:      "Number of currently idle conns in the pool",
			},
		),
		maxConns: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_pool_max_conns",
				Help:      "Maximum size of the pool",
			},
		),
		totalConns: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_pool_total_conns",
				Help:      "Total number of resources currently in the pool (acquired, idle, constructing)",
			},
		),
		maxLifetimeDestroy: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_pool_max_lifetime_destroy_total",
				Help:      "Cumulative count of connections destroyed for exceeding max lifetime",
			},
		),
		maxIdleDestroy: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_pool_max_idle_destroy_total",
				Help:      "Cumulative count of connections destroyed for exceeding max idle time",
			},
		),
	}

	return collector
}

// Collect updates all metrics with current pool stats
func (c *DBCollector) Collect() {
	if c.pool == nil {
		return
	}

	stat := c.pool.Stat()

	// Update gauges
	c.acquiredConns.Set(float64(stat.AcquiredConns()))
	c.constructingConns.Set(float64(stat.ConstructingConns()))
	c.idleConns.Set(float64(stat.IdleConns()))
	c.maxConns.Set(float64(stat.MaxConns()))
	c.totalConns.Set(float64(stat.TotalConns()))

	// Update counters (only if they increased)
	c.canceledAcquires.Add(float64(stat.CanceledAcquireCount()))
	c.emptyAcquires.Add(float64(stat.EmptyAcquireCount()))
	c.maxLifetimeDestroy.Add(float64(stat.MaxLifetimeDestroyCount()))
	c.maxIdleDestroy.Add(float64(stat.MaxIdleDestroyCount()))

	// Note: AcquireCount and AcquireDuration are recorded per-operation
	// via the wrapper methods below
}

// StartCollecting starts a goroutine that periodically collects metrics
func (c *DBCollector) StartCollecting(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.Collect()
			}
		}
	}()
}

// RecordAcquire records a connection acquire operation
func (c *DBCollector) RecordAcquire(duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	c.acquireCount.WithLabelValues(status).Inc()
	if success {
		c.acquireDuration.WithLabelValues().Observe(duration.Seconds())
	}
}

// GetCollectors returns all prometheus collectors
func (c *DBCollector) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		c.acquireCount,
		c.acquireDuration,
		c.acquiredConns,
		c.canceledAcquires,
		c.constructingConns,
		c.emptyAcquires,
		c.idleConns,
		c.maxConns,
		c.totalConns,
		c.maxLifetimeDestroy,
		c.maxIdleDestroy,
	}
}
