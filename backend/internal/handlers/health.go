package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/logging"
	"go.uber.org/zap"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db     *pgxpool.Pool
	redis  *redis.Client
	logger *logging.Logger
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *pgxpool.Pool, redis *redis.Client, logger *logging.Logger) *HealthHandler {
	return &HealthHandler{
		db:        db,
		redis:     redis,
		logger:    logger,
		startTime: time.Now(),
	}
}

// HealthStatus represents the overall health status
type HealthStatus struct {
	Status      string                 `json:"status"` // ok, degraded, error
	Timestamp   time.Time              `json:"timestamp"`
	Version     string                 `json:"version"`
	Uptime      string                 `json:"uptime"`
	Checks      map[string]CheckResult `json:"checks"`
	System      SystemInfo             `json:"system"`
}

// CheckResult represents a health check result
type CheckResult struct {
	Status    string                 `json:"status"` // ok, degraded, error
	Message   string                 `json:"message,omitempty"`
	Duration  string                 `json:"duration"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// SystemInfo represents system information
type SystemInfo struct {
	Goroutines    int     `json:"goroutines"`
	MemoryAlloc   uint64  `json:"memory_alloc_bytes"`
	MemoryTotal   uint64  `json:"memory_total_bytes"`
	MemorySys     uint64  `json:"memory_sys_bytes"`
	NumGC         uint32  `json:"num_gc"`
	CPUCores      int     `json:"cpu_cores"`
}

// Health returns basic health status (for quick checks)
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0", // TODO: Get from build info
		Uptime:    time.Since(h.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// HealthDetailed returns detailed health status with all checks
func (h *HealthHandler) HealthDetailed(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Run all health checks concurrently
	checks := make(map[string]CheckResult)
	var mu sync.Mutex
	var wg sync.WaitGroup

	checkFuncs := map[string]func(context.Context) CheckResult{
		"database":     h.checkDatabase,
		"redis":        h.checkRedis,
		"disk_space":   h.checkDiskSpace,
		"memory":       h.checkMemory,
	}

	for name, checkFunc := range checkFuncs {
		wg.Add(1)
		go func(name string, fn func(context.Context) CheckResult) {
			defer wg.Done()
			result := fn(ctx)
			mu.Lock()
			checks[name] = result
			mu.Unlock()
		}(name, checkFunc)
	}

	wg.Wait()

	// Determine overall status
	overallStatus := "ok"
	for _, check := range checks {
		if check.Status == "error" {
			overallStatus = "error"
			break
		}
		if check.Status == "degraded" {
			overallStatus = "degraded"
		}
	}

	// Get system info
	systemInfo := h.getSystemInfo()

	status := HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0", // TODO: Get from build info
		Uptime:    time.Since(h.startTime).String(),
		Checks:    checks,
		System:    systemInfo,
	}

	// Set HTTP status based on health
	httpStatus := http.StatusOK
	if overallStatus == "error" {
		httpStatus = http.StatusServiceUnavailable
	} else if overallStatus == "degraded" {
		httpStatus = http.StatusOK // Still return 200 for degraded
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(status)
}

// Readiness returns readiness probe status (for Kubernetes)
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check critical dependencies
	dbCheck := h.checkDatabase(ctx)
	redisCheck := h.checkRedis(ctx)

	if dbCheck.Status == "error" || redisCheck.Status == "error" {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "not ready",
			"reason": "critical dependencies unavailable",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

// Liveness returns liveness probe status (for Kubernetes)
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	// Simple liveness check - just return 200
	// If the service can respond, it's alive
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
}

// checkDatabase checks database connectivity and performance
func (h *HealthHandler) checkDatabase(ctx context.Context) CheckResult {
	start := time.Now()
	result := CheckResult{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	// Check connection
	err := h.db.Ping(ctx)
	if err != nil {
		result.Status = "error"
		result.Message = "Database ping failed: " + err.Error()
		result.Duration = time.Since(start).String()
		return result
	}

	// Get pool stats
	stats := h.db.Stat()
	details := map[string]interface{}{
		"acquired_conns":   stats.AcquiredConns(),
		"idle_conns":       stats.IdleConns(),
		"max_conns":        stats.MaxConns(),
		"total_conns":      stats.TotalConns(),
	}

	// Check if pool is exhausted
	utilization := float64(stats.AcquiredConns()) / float64(stats.MaxConns())
	if utilization > 0.9 {
		result.Status = "degraded"
		result.Message = "Database connection pool nearly exhausted"
	}

	details["utilization"] = utilization
	result.Details = details
	result.Duration = time.Since(start).String()

	return result
}

// checkRedis checks Redis connectivity and performance
func (h *HealthHandler) checkRedis(ctx context.Context) CheckResult {
	start := time.Now()
	result := CheckResult{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	// Check connection
	_, err := h.redis.Ping(ctx).Result()
	if err != nil {
		result.Status = "error"
		result.Message = "Redis ping failed: " + err.Error()
		result.Duration = time.Since(start).String()
		return result
	}

	// Get Redis info
	info, err := h.redis.Info(ctx, "memory").Result()
	if err != nil {
		result.Status = "degraded"
		result.Message = "Could not get Redis info"
	} else {
		result.Details = map[string]interface{}{
			"info": "available",
		}
	}

	// Get database size
	dbSize, err := h.redis.DBSize(ctx).Result()
	if err == nil {
		if result.Details == nil {
			result.Details = make(map[string]interface{})
		}
		result.Details["keys"] = dbSize
	}

	result.Duration = time.Since(start).String()
	return result
}

// checkDiskSpace checks available disk space
func (h *HealthHandler) checkDiskSpace(ctx context.Context) CheckResult {
	start := time.Now()
	result := CheckResult{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	// Note: In production, you'd want to actually check disk space
	// For now, we'll just return ok
	result.Message = "Disk space check not implemented"
	result.Duration = time.Since(start).String()

	return result
}

// checkMemory checks memory usage
func (h *HealthHandler) checkMemory(ctx context.Context) CheckResult {
	start := time.Now()
	result := CheckResult{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	details := map[string]interface{}{
		"alloc_bytes":      m.Alloc,
		"total_alloc_bytes": m.TotalAlloc,
		"sys_bytes":        m.Sys,
		"num_gc":           m.NumGC,
	}

	// Check if memory usage is high (> 1GB allocated)
	if m.Alloc > 1*1024*1024*1024 {
		result.Status = "degraded"
		result.Message = "High memory usage detected"
	}

	result.Details = details
	result.Duration = time.Since(start).String()

	return result
}

// getSystemInfo returns current system information
func (h *HealthHandler) getSystemInfo() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemInfo{
		Goroutines:  runtime.NumGoroutine(),
		MemoryAlloc: m.Alloc,
		MemoryTotal: m.TotalAlloc,
		MemorySys:   m.Sys,
		NumGC:       m.NumGC,
		CPUCores:    runtime.NumCPU(),
	}
}
