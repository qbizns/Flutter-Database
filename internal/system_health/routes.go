package system_health

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/your-org/pos-backend/internal/api/middlewares"
)

// RegisterRoutes registers all SystemHealth routes
func RegisterRoutes(r chi.Router, handler *Handler) {
	
	// Organization-scoped routes
	r.Route("/api/v1/organizations/{orgID}/system_health", func(r chi.Router) {
		// Apply middleware
		r.Use(middlewares.AuthRequired)
		r.Use(middlewares.OrganizationContext)
		r.Use(middlewares.RateLimiter)
		r.Use(middleware.Compress(5))

		// CRUD endpoints
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)

		// Custom endpoints (add as needed)
		// r.Get("/search", handler.Search)
		// r.Post("/bulk", handler.BulkCreate)
		// r.Delete("/bulk", handler.BulkDelete)
	})
	
}

// RegisterPublicRoutes registers public SystemHealth routes (if any)
func RegisterPublicRoutes(r chi.Router, handler *Handler) {
	r.Route("/api/v1/public/system_health", func(r chi.Router) {
		r.Use(middlewares.RateLimiter)
		r.Use(middleware.Compress(5))

		// Add public endpoints here
		// Example: r.Get("/", handler.ListPublic)
	})
}

// RegisterAdminRoutes registers admin-only SystemHealth routes
func RegisterAdminRoutes(r chi.Router, handler *Handler) {
	r.Route("/api/v1/admin/system_health", func(r chi.Router) {
		// Apply middleware
		r.Use(middlewares.AuthRequired)
		r.Use(middlewares.AdminOnly)
		r.Use(middlewares.RateLimiter)
		r.Use(middleware.Compress(5))

		// Admin endpoints
		r.Get("/", handler.AdminList)
		r.Get("/stats", handler.GetStats)
		r.Post("/export", handler.Export)
		r.Post("/import", handler.Import)
		
	})
}

// RouteConfig holds route configuration
type RouteConfig struct {
	Prefix      string
	Middlewares []func(http.Handler) http.Handler
	Handler     *Handler
}

// RegisterCustomRoutes allows flexible route registration
func RegisterCustomRoutes(r chi.Router, config RouteConfig) {
	r.Route(config.Prefix, func(r chi.Router) {
		// Apply custom middleware
		for _, mw := range config.Middlewares {
			r.Use(mw)
		}

		// Standard CRUD
		r.Post("/", config.Handler.Create)
		r.Get("/", config.Handler.List)
		r.Get("/{id}", config.Handler.GetByID)
		r.Put("/{id}", config.Handler.Update)
		r.Delete("/{id}", config.Handler.Delete)
	})
}
