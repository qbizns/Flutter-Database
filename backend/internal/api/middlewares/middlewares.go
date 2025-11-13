package middlewares

import (
	"net/http"

	"github.com/your-org/pos-backend/internal/auth"
	custommw "github.com/your-org/pos-backend/internal/middleware"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
)

// Global middleware instances - initialized in main.go
var (
	authMW       *auth.Middleware
	rateLimiterMW *custommw.RateLimiter
)

// Initialize sets up global middleware instances
// Call this from main.go after creating middleware instances
func Initialize(authMiddleware *auth.Middleware, rateLimiter *custommw.RateLimiter) {
	authMW = authMiddleware
	rateLimiterMW = rateLimiter
}

// AuthRequired is middleware that requires authentication
func AuthRequired(next http.Handler) http.Handler {
	if authMW == nil {
		// Fallback - should not happen in production
		return next
	}
	return authMW.Authenticate(next)
}

// OrganizationContext middleware validates organization access
// This is a placeholder - can be enhanced with additional org-level checks
func OrganizationContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Organization context is already set by AuthRequired middleware
		// This middleware can add additional organization-level validation
		// For now, it's a pass-through
		next.ServeHTTP(w, r)
	})
}

// RateLimiter applies rate limiting
func RateLimiter(next http.Handler) http.Handler {
	if rateLimiterMW == nil {
		// Fallback - rate limiting disabled
		return next
	}
	return rateLimiterMW.Limit()(next)
}

// AdminOnly restricts access to admin users only
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user has admin role
		if !appctx.HasRole(r.Context(), "admin") {
			http.Error(w, `{"error":"admin access required"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// OptionalAuth allows requests with or without authentication
func OptionalAuth(next http.Handler) http.Handler {
	if authMW == nil {
		return next
	}
	return authMW.OptionalAuthenticate(next)
}
