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
	"github.com/your-org/pos-backend/internal/auth"
	"github.com/your-org/pos-backend/internal/config"
	domainauth "github.com/your-org/pos-backend/internal/domain/auth"
	"github.com/your-org/pos-backend/internal/http/rest"
	"github.com/your-org/pos-backend/internal/logging"
	custommiddleware "github.com/your-org/pos-backend/internal/middleware"
	"github.com/your-org/pos-backend/internal/metrics"
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

	// Initialize auth repositories
	userRepo := postgres.NewUserRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	userRoleRepo := postgres.NewUserRoleRepository(db)

	// Initialize auth services
	userService := domainauth.NewUserService(userRepo, roleRepo, userRoleRepo, logger)
	tokenService := auth.NewTokenService(cfg)

	// Initialize metrics collectors
	metricsMiddleware := custommiddleware.NewMetricsMiddleware("pos_backend")
	// Business metrics collector (for future use in domain services)
	// metricsCollector := metrics.NewCollector("pos_backend")
	dbMetricsCollector := metrics.NewDBCollector("pos_backend", db.Pool)

	// Start collecting database pool metrics every 15 seconds
	metricsCtx, metricsCancel := context.WithCancel(context.Background())
	defer metricsCancel()
	dbMetricsCollector.StartCollecting(metricsCtx, 15*time.Second)

	logger.Info("metrics collection initialized",
		zap.String("namespace", "pos_backend"),
		zap.String("metrics_port", cfg.Metrics.MetricsPort))

	// Initialize router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Metrics middleware (after logger, before business logic)
	r.Use(metricsMiddleware.Handler)

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
	r.Handle("/metrics", promhttp.Handler())

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (no auth)
		r.Group(func(r chi.Router) {
			// Authentication
			r.Post("/auth/login", rest.LoginHandler(cfg, userService, tokenService, logger))
			r.Post("/auth/register", rest.RegisterHandler(cfg, userService, tokenService, logger))
			r.Post("/auth/refresh", rest.RefreshTokenHandler(tokenService, logger))
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
				r.Get("/customers/{id}", rest.GetCustomerHandler(db, logger))
				r.Patch("/customers/{id}", rest.UpdateCustomerHandler(db, logger))
				r.Delete("/customers/{id}", rest.DeleteCustomerHandler(db, logger))

				// Suppliers
				r.Get("/suppliers", rest.ListSuppliersHandler(db, logger))
				r.Post("/suppliers", rest.CreateSupplierHandler(db, logger))
				r.Get("/suppliers/{id}", rest.GetSupplierHandler(db, logger))
				r.Patch("/suppliers/{id}", rest.UpdateSupplierHandler(db, logger))
				r.Delete("/suppliers/{id}", rest.DeleteSupplierHandler(db, logger))

				// Categories
				r.Get("/categories", rest.ListCategoriesHandler(db, logger))
				r.Post("/categories", rest.CreateCategoryHandler(db, logger))
				r.Get("/categories/{id}", rest.GetCategoryHandler(db, logger))
				r.Patch("/categories/{id}", rest.UpdateCategoryHandler(db, logger))
				r.Delete("/categories/{id}", rest.DeleteCategoryHandler(db, logger))

				// Locations
				r.Get("/locations", rest.ListLocationsHandler(db, logger))
				r.Post("/locations", rest.CreateLocationHandler(db, logger))
				r.Get("/locations/{id}", rest.GetLocationHandler(db, logger))
				r.Patch("/locations/{id}", rest.UpdateLocationHandler(db, logger))
				r.Delete("/locations/{id}", rest.DeleteLocationHandler(db, logger))

				// Sales
				r.Get("/sales", rest.ListSalesHandler(db, logger))
				r.Post("/sales", rest.CreateSaleHandler(db, logger))
				r.Get("/sales/{id}", rest.GetSaleHandler(db, logger))
				r.Patch("/sales/{id}", rest.UpdateSaleHandler(db, logger))
				r.Delete("/sales/{id}", rest.DeleteSaleHandler(db, logger))

				// Posting Engine (Phase 3B)
				r.Post("/posting/post", rest.PostDocumentHandler(db, logger))
				r.Get("/posting/rules", rest.GetPostingRulesHandler(db, logger))
				r.Get("/posting/audit", rest.GetPostingAuditHandler(db, logger))

				// Journal Entries (Phase 3B)
				r.Get("/journal-entries", rest.ListJournalEntriesHandler(db, logger))
				r.Post("/journal-entries", rest.CreateJournalEntryHandler(db, logger))
				r.Get("/journal-entries/{id}", rest.GetJournalEntryHandler(db, logger))
				r.Post("/journal-entries/{id}/post", rest.PostJournalEntryHandler(db, logger))

				// Financial Reports (Phase 3B)
				r.Get("/reports/balance-sheet", rest.BalanceSheetHandler(db, logger))
				r.Get("/reports/income-statement", rest.IncomeStatementHandler(db, logger))
				r.Get("/reports/trial-balance", rest.TrialBalanceHandler(db, logger))

				// Chart of Accounts (Phase 3B)
				r.Get("/chart-of-accounts", rest.ListChartOfAccountsHandler(db, logger))
				r.Post("/chart-of-accounts", rest.CreateChartOfAccountHandler(db, logger))
				r.Get("/chart-of-accounts/{id}", rest.GetChartOfAccountHandler(db, logger))
				r.Patch("/chart-of-accounts/{id}", rest.UpdateChartOfAccountHandler(db, logger))
				r.Delete("/chart-of-accounts/{id}", rest.DeleteChartOfAccountHandler(db, logger))

				// Fiscal Years (Phase 3B)
				r.Get("/fiscal-years", rest.ListFiscalYearsHandler(db, logger))
				r.Post("/fiscal-years", rest.CreateFiscalYearHandler(db, logger))
				r.Get("/fiscal-years/{id}", rest.GetFiscalYearHandler(db, logger))
				r.Patch("/fiscal-years/{id}", rest.UpdateFiscalYearHandler(db, logger))
				r.Post("/fiscal-years/{id}/close", rest.CloseFiscalYearHandler(db, logger))

				// Accounting Periods (Phase 3B)
				r.Get("/accounting-periods", rest.ListAccountingPeriodsHandler(db, logger))
				r.Post("/accounting-periods", rest.CreateAccountingPeriodHandler(db, logger))

				// Customer Invoices (Receivables - Phase 3B)
				r.Get("/customer-invoices", rest.ListCustomerInvoicesHandler(db, logger))
				r.Post("/customer-invoices", rest.CreateCustomerInvoiceHandler(db, logger))
				r.Get("/customer-invoices/{id}", rest.GetCustomerInvoiceHandler(db, logger))
				r.Patch("/customer-invoices/{id}", rest.UpdateCustomerInvoiceHandler(db, logger))
				r.Delete("/customer-invoices/{id}", rest.DeleteCustomerInvoiceHandler(db, logger))

				// Customer Payments (Receivables - Phase 3B)
				r.Get("/customer-payments", rest.ListCustomerPaymentsHandler(db, logger))
				r.Post("/customer-payments", rest.CreateCustomerPaymentHandler(db, logger))
				r.Get("/customer-payments/{id}", rest.GetCustomerPaymentHandler(db, logger))

				// Vendor Bills (Payables - Phase 3B)
				r.Get("/vendor-bills", rest.ListVendorBillsHandler(db, logger))
				r.Post("/vendor-bills", rest.CreateVendorBillHandler(db, logger))
				r.Get("/vendor-bills/{id}", rest.GetVendorBillHandler(db, logger))
				r.Patch("/vendor-bills/{id}", rest.UpdateVendorBillHandler(db, logger))
				r.Delete("/vendor-bills/{id}", rest.DeleteVendorBillHandler(db, logger))

				// Vendor Payments (Payables - Phase 3B)
				r.Get("/vendor-payments", rest.ListVendorPaymentsHandler(db, logger))
				r.Post("/vendor-payments", rest.CreateVendorPaymentHandler(db, logger))
				r.Get("/vendor-payments/{id}", rest.GetVendorPaymentHandler(db, logger))
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
