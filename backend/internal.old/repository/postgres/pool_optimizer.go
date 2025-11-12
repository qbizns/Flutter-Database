package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/logging"
	"go.uber.org/zap"
)

// PoolOptimizer monitors and logs connection pool statistics
type PoolOptimizer struct {
	pool     *pgxpool.Pool
	logger   *logging.Logger
	interval time.Duration
}

// NewPoolOptimizer creates a new pool optimizer
func NewPoolOptimizer(pool *pgxpool.Pool, logger *logging.Logger, interval time.Duration) *PoolOptimizer {
	return &PoolOptimizer{
		pool:     pool,
		logger:   logger,
		interval: interval,
	}
}

// Start begins monitoring the connection pool
func (po *PoolOptimizer) Start(ctx context.Context) {
	ticker := time.NewTicker(po.interval)
	defer ticker.Stop()

	po.logger.Info("connection pool optimizer started", zap.Duration("interval", po.interval))

	for {
		select {
		case <-ctx.Done():
			po.logger.Info("connection pool optimizer stopped")
			return
		case <-ticker.C:
			po.logStats()
			po.checkHealth()
		}
	}
}

// logStats logs connection pool statistics
func (po *PoolOptimizer) logStats() {
	stat := po.pool.Stat()

	po.logger.Debug("connection pool stats",
		zap.Int32("acquired_conns", stat.AcquiredConns()),
		zap.Int32("constructing_conns", stat.ConstructingConns()),
		zap.Int32("idle_conns", stat.IdleConns()),
		zap.Int32("max_conns", stat.MaxConns()),
		zap.Int64("acquire_count", stat.AcquireCount()),
		zap.Duration("acquire_duration", stat.AcquireDuration()),
		zap.Int64("acquired_conns_count", stat.AcquiredConns()),
		zap.Int64("canceled_acquire_count", stat.CanceledAcquireCount()),
		zap.Int64("empty_acquire_count", stat.EmptyAcquireCount()),
		zap.Int32("total_conns", stat.TotalConns()),
	)
}

// checkHealth checks connection pool health and warns on issues
func (po *PoolOptimizer) checkHealth() {
	stat := po.pool.Stat()

	// Warn if pool is nearly exhausted
	utilization := float64(stat.AcquiredConns()) / float64(stat.MaxConns())
	if utilization > 0.9 {
		po.logger.Warn("connection pool nearly exhausted",
			zap.Float64("utilization", utilization),
			zap.Int32("acquired", stat.AcquiredConns()),
			zap.Int32("max", stat.MaxConns()),
		)
	}

	// Warn if there are many canceled acquires
	if stat.CanceledAcquireCount() > 100 {
		po.logger.Warn("high number of canceled connection acquires",
			zap.Int64("canceled_count", stat.CanceledAcquireCount()),
		)
	}

	// Warn if average acquire duration is high
	if stat.AcquireDuration() > 100*time.Millisecond {
		po.logger.Warn("high connection acquire duration",
			zap.Duration("avg_duration", stat.AcquireDuration()),
		)
	}

	// Warn if no idle connections but not at max
	if stat.IdleConns() == 0 && stat.TotalConns() < stat.MaxConns() {
		po.logger.Warn("no idle connections available, may need to increase pool size",
			zap.Int32("total_conns", stat.TotalConns()),
			zap.Int32("max_conns", stat.MaxConns()),
		)
	}
}

// GetRecommendations returns pool configuration recommendations
func (po *PoolOptimizer) GetRecommendations() PoolRecommendations {
	stat := po.pool.Stat()

	recommendations := PoolRecommendations{
		Current: PoolConfig{
			MaxConns:  int(stat.MaxConns()),
			IdleConns: int(stat.IdleConns()),
		},
	}

	// Calculate utilization
	utilization := float64(stat.AcquiredConns()) / float64(stat.MaxConns())

	// Recommend increasing max connections if utilization > 80%
	if utilization > 0.8 {
		recommendations.Recommended.MaxConns = int(stat.MaxConns()) + 10
		recommendations.Recommended.IdleConns = int(stat.MaxConns()/4) // 25% idle
		recommendations.Reason = "High utilization detected (>80%). Consider increasing pool size."
	}

	// Recommend decreasing if utilization consistently < 30%
	if utilization < 0.3 && stat.MaxConns() > 10 {
		recommendations.Recommended.MaxConns = int(stat.MaxConns()) - 5
		recommendations.Recommended.IdleConns = int(stat.MaxConns() / 5) // 20% idle
		recommendations.Reason = "Low utilization detected (<30%). Consider decreasing pool size to save resources."
	}

	// If no changes recommended
	if recommendations.Recommended.MaxConns == 0 {
		recommendations.Recommended = recommendations.Current
		recommendations.Reason = "Current pool configuration is optimal."
	}

	return recommendations
}

// PoolConfig represents connection pool configuration
type PoolConfig struct {
	MaxConns  int
	IdleConns int
}

// PoolRecommendations contains pool optimization recommendations
type PoolRecommendations struct {
	Current     PoolConfig
	Recommended PoolConfig
	Reason      string
}
