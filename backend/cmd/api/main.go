package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/your-org/pos-backend/internal/api/middlewares"
	"github.com/your-org/pos-backend/internal/auth"
	"github.com/your-org/pos-backend/internal/config"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	custommw "github.com/your-org/pos-backend/internal/middleware"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := logging.NewLogger(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting API server",
		zap.String("env", cfg.Server.Env),
		zap.String("port", cfg.Server.APIPort),
	)

	// Initialize database
	db, err := postgres.New(cfg, logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize Redis for rate limiting
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Warn("Redis connection failed, rate limiting will be disabled", zap.Error(err))
		redisClient = nil
	} else {
		logger.Info("Redis connection established")
	}

	// Initialize rate limiter
	rateLimiter := custommw.NewRateLimiter(redisClient, cfg.RateLimit, logger)

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(cfg.JWT.Secret)

	// Initialize global middleware wrapper for route modules
	middlewares.Initialize(authMiddleware, rateLimiter)
	logger.Info("middleware initialized for module routes")

	// Start system metrics collector
	systemCollector := metrics.NewSystemCollector(logger, 10*time.Second)
	go systemCollector.Start(context.Background())
	defer systemCollector.Stop()

	// Start database stats collector
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			stat := db.Stats()
			metrics.UpdateDatabaseConnectionStats(
				stat.AcquiredConns(),
				stat.IdleConns(),
				stat.MaxConns(),
			)
		}
	}()

	// Initialize router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// OBSERVABILITY: Collect metrics for all requests
	r.Use(metrics.MetricsMiddleware)

	// SECURITY: Apply rate limiting to all requests
	r.Use(rateLimiter.Limit())

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Health(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status": "unhealthy"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	})

	// Metrics endpoint (Prometheus scraping)
	if cfg.Metrics.Enabled {
		r.Handle("/metrics", promhttp.Handler())
		logger.Info("metrics endpoint enabled", zap.String("path", "/metrics"))
	}

	// API v1 routes - using new modular architecture (171 modules)
	// Legacy handlers removed - all functionality now provided by module routes below

	// Register all module routes (171 modules with full CRUD + Admin endpoints)
	RegisterAllModuleRoutes(r, db, logger)
	logger.Info("all module routes registered successfully")

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.APIPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("API server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited")
}
