package metrics

import (
	"context"
	"runtime"
	"time"

	"github.com/your-org/pos-backend/internal/logging"
	"go.uber.org/zap"
)

// SystemCollector periodically collects system metrics
type SystemCollector struct {
	logger   *logging.Logger
	interval time.Duration
	stopChan chan struct{}
}

// NewSystemCollector creates a new system metrics collector
func NewSystemCollector(logger *logging.Logger, interval time.Duration) *SystemCollector {
	return &SystemCollector{
		logger:   logger,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start begins collecting system metrics
func (sc *SystemCollector) Start(ctx context.Context) {
	ticker := time.NewTicker(sc.interval)
	defer ticker.Stop()

	sc.logger.Info("system metrics collector started", zap.Duration("interval", sc.interval))

	for {
		select {
		case <-ctx.Done():
			sc.logger.Info("system metrics collector stopped")
			return
		case <-sc.stopChan:
			sc.logger.Info("system metrics collector stopped")
			return
		case <-ticker.C:
			sc.collect()
		}
	}
}

// Stop stops the metrics collector
func (sc *SystemCollector) Stop() {
	close(sc.stopChan)
}

// collect gathers current system metrics
func (sc *SystemCollector) collect() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Update metrics
	GoRoutines.Set(float64(runtime.NumGoroutine()))
	MemoryUsage.Set(float64(m.Sys))
	MemoryAllocated.Set(float64(m.Alloc))
}

// CollectOnce collects metrics once (useful for testing)
func CollectOnce() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	GoRoutines.Set(float64(runtime.NumGoroutine()))
	MemoryUsage.Set(float64(m.Sys))
	MemoryAllocated.Set(float64(m.Alloc))
}
