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
	"github.com/your-org/pos-backend/internal/auth"
	"github.com/your-org/pos-backend/internal/config"
	"github.com/your-org/pos-backend/internal/http/rest"
	"github.com/your-org/pos-backend/internal/logging"
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

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(cfg.JWT.Secret)

	// Initialize router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

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

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (no auth)
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", rest.LoginHandler(cfg, db, logger))
			r.Post("/auth/register", rest.RegisterHandler(cfg, db, logger))
		})

		// Protected routes (require auth)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			// Organizations
			r.Route("/organizations/{org_id}", func(r chi.Router) {
				// Products
				r.Get("/products", rest.ListProductsHandler(db, logger))
				r.Post("/products", rest.CreateProductHandler(db, logger))
				r.Get("/products/{id}", rest.GetProductHandler(db, logger))
				r.Patch("/products/{id}", rest.UpdateProductHandler(db, logger))
				r.Delete("/products/{id}", rest.DeleteProductHandler(db, logger))

				// Customers
				r.Get("/customers", rest.ListCustomersHandler(db, logger))
				r.Post("/customers", rest.CreateCustomerHandler(db, logger))

				// Sales
				r.Get("/sales", rest.ListSalesHandler(db, logger))
				r.Post("/sales", rest.CreateSaleHandler(db, logger))
				r.Get("/sales/{id}", rest.GetSaleHandler(db, logger))

				// Posting Engine (CRITICAL)
				r.Post("/posting/post", rest.PostDocumentHandler(db, logger))
				r.Get("/posting/rules", rest.GetPostingRulesHandler(db, logger))
				r.Get("/posting/audit", rest.GetPostingAuditHandler(db, logger))

				// Journal Entries
				r.Get("/journal-entries", rest.ListJournalEntriesHandler(db, logger))
				r.Post("/journal-entries", rest.CreateJournalEntryHandler(db, logger))

				// Reports
				r.Get("/reports/balance-sheet", rest.BalanceSheetHandler(db, logger))
				r.Get("/reports/income-statement", rest.IncomeStatementHandler(db, logger))
				r.Get("/reports/trial-balance", rest.TrialBalanceHandler(db, logger))
			})
		})
	})

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
