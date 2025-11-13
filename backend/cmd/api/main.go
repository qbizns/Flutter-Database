package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq" // postgres driver
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	// Internal packages
	"github.com/your-org/pos-backend/internal/auth"
	"github.com/your-org/pos-backend/internal/config"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	custommw "github.com/your-org/pos-backend/internal/middleware"

	// Module imports
	"github.com/your-org/pos-backend/internal/category"
	"github.com/your-org/pos-backend/internal/customer"
	"github.com/your-org/pos-backend/internal/kitchen_station"
	"github.com/your-org/pos-backend/internal/order"
	"github.com/your-org/pos-backend/internal/payment"
	"github.com/your-org/pos-backend/internal/product"
	"github.com/your-org/pos-backend/internal/restaurant_table"
	"github.com/your-org/pos-backend/internal/role"
	"github.com/your-org/pos-backend/internal/shift"
	"github.com/your-org/pos-backend/internal/user"
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

	// Initialize database connection
	db, err := connectDatabase(cfg, logger)
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

	// Start system metrics collector
	systemCollector := metrics.NewSystemCollector(logger, 10*time.Second)
	go systemCollector.Start(context.Background())
	defer systemCollector.Stop()

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
	r.Get("/health", healthCheckHandler(db))

	// Metrics endpoint (Prometheus scraping)
	if cfg.Metrics.Enabled {
		r.Handle("/metrics", promhttp.Handler())
		logger.Info("metrics endpoint enabled", zap.String("path", "/metrics"))
	}

	// Initialize all module handlers
	initializeModules(r, db, logger, authMiddleware, cfg)

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

func connectDatabase(cfg *config.Config, logger *logging.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	// Test connection
	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}

	logger.Info("database connection established")
	return db, nil
}

func healthCheckHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status": "unhealthy"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}
}

func initializeModules(r chi.Router, db *sql.DB, logger *logging.Logger, authMiddleware *auth.Middleware, cfg *config.Config) {
	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (no auth, but stricter rate limiting for auth endpoints)
		r.Group(func(r chi.Router) {
			// Auth endpoints would go here
			// r.Post("/auth/login", authHandler.Login)
			// r.Post("/auth/register", authHandler.Register)
		})

		// Protected routes (require auth)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			// Organizations
			r.Route("/organizations/{org_id}", func(r chi.Router) {
				// Initialize all module repositories, services, and handlers
				initializeOrdersModule(r, db, logger)
				initializeProductsModule(r, db, logger)
				initializeCategoriesModule(r, db, logger)
				initializeTablesModule(r, db, logger)
				initializePaymentsModule(r, db, logger)
				initializeCustomersModule(r, db, logger)
				initializeStaffModule(r, db, logger)
				initializeRolesModule(r, db, logger)
				initializeShiftsModule(r, db, logger)
				initializeKitchenStationsModule(r, db, logger)
				initializeZonesModule(r, db, logger)
			})
		})
	})
}

// Orders Module - Priority 1
func initializeOrdersModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	repo := order.NewRepository(db)
	service := order.NewService(repo, logger)
	handler := order.NewHandler(service, logger)

	r.Route("/orders", func(r chi.Router) {
		// Basic CRUD
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Patch("/{id}", handler.Update)

		// Custom endpoints per NEW_APIS.md
		r.Patch("/{id}/status", handler.UpdateStatus)
		r.Post("/{id}/cancel", handler.Cancel)
		r.Get("/statistics", handler.GetStatistics)
	})
}

// Products Module - Priority 1
func initializeProductsModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	repo := product.NewRepository(db)
	service := product.NewService(repo, logger)
	handler := product.NewHandler(service, logger)

	r.Route("/products", func(r chi.Router) {
		// Basic CRUD
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Patch("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)

		// Custom endpoints per NEW_APIS.md
		r.Get("/search", handler.Search)
		r.Post("/batch", handler.GetBatch)
		r.Get("/featured", handler.GetFeatured)
		r.Get("/low-stock", handler.GetLowStock)
		r.Patch("/{id}/stock", handler.UpdateStock)
		r.Patch("/{id}/availability", handler.UpdateAvailability)
	})
}

// Categories Module - Priority 1
func initializeCategoriesModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	repo := category.NewRepository(db)
	service := category.NewService(repo, logger)
	handler := category.NewHandler(service, logger)

	r.Route("/categories", func(r chi.Router) {
		// Basic CRUD
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Patch("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})

	// Products count by category (separate endpoint)
	r.Get("/products/count-by-category", handler.GetProductsCountByCategory)
}

// Tables Module - Priority 1
func initializeTablesModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	repo := restaurant_table.NewRepository(db)
	service := restaurant_table.NewService(repo, logger)
	handler := restaurant_table.NewHandler(service, logger)

	r.Route("/tables", func(r chi.Router) {
		// Basic CRUD
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Patch("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)

		// Custom endpoints per NEW_APIS.md
		r.Patch("/{id}/status", handler.UpdateStatus)
		r.Post("/{id}/assign-order", handler.AssignOrder)
		r.Post("/{id}/clear", handler.Clear)
		r.Get("/statistics", handler.GetStatistics)
		r.Get("/count-by-zone", handler.GetCountByZone)
	})
}

// Payments Module - Priority 1
func initializePaymentsModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	repo := payment.NewRepository(db)
	service := payment.NewService(repo, logger)
	handler := payment.NewHandler(service, logger)

	r.Route("/payments", func(r chi.Router) {
		// Create payment (process payment)
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)

		// Custom endpoints per NEW_APIS.md
		r.Post("/{id}/cancel", handler.Cancel)
		r.Get("/statistics", handler.GetStatistics)
	})

	// Refunds
	r.Route("/refunds", func(r chi.Router) {
		r.Post("/", handler.CreateRefund)
		r.Get("/", handler.ListRefunds)
		r.Get("/{id}", handler.GetRefund)
	})
}

// Customers Module - Priority 3
func initializeCustomersModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	repo := customer.NewRepository(db)
	service := customer.NewService(repo, logger)
	handler := customer.NewHandler(service, logger)

	r.Route("/customers", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Patch("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

// Staff Module - Priority 2
func initializeStaffModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	repo := user.NewRepository(db)
	service := user.NewService(repo, logger)
	handler := user.NewHandler(service, logger)

	r.Route("/staff", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Patch("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

// Roles Module - Priority 2
func initializeRolesModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	repo := role.NewRepository(db)
	service := role.NewService(repo, logger)
	handler := role.NewHandler(service, logger)

	r.Route("/roles", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Patch("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

// Shifts Module - Priority 2
func initializeShiftsModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	repo := shift.NewRepository(db)
	service := shift.NewService(repo, logger)
	handler := shift.NewHandler(service, logger)

	r.Route("/shifts", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Patch("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

// Kitchen Stations Module - Priority 2
func initializeKitchenStationsModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	repo := kitchen_station.NewRepository(db)
	service := kitchen_station.NewService(repo, logger)
	handler := kitchen_station.NewHandler(service, logger)

	r.Route("/kitchen-stations", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Patch("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

// Zones Module - Priority 2
func initializeZonesModule(r chi.Router, db *sql.DB, logger *logging.Logger) {
	// TODO: Create zone module
	// For now, zones will be handled through table_section module
	r.Route("/zones", func(r chi.Router) {
		// r.Post("/", handler.Create)
		// r.Get("/", handler.List)
		// r.Get("/{id}", handler.GetByID)
	})
}
