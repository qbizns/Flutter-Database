package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/your-org/pos-backend/internal/config"
)

// SecurityHeaders adds security headers to all responses
type SecurityHeaders struct {
	config *config.Config
}

// NewSecurityHeaders creates a new security headers middleware
func NewSecurityHeaders(cfg *config.Config) *SecurityHeaders {
	return &SecurityHeaders{
		config: cfg,
	}
}

// Handler returns middleware that adds security headers
func (sh *SecurityHeaders) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Strict-Transport-Security (HSTS)
			// Force HTTPS for 1 year, include subdomains
			if sh.config.Server.Env == "production" {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			}

			// X-Content-Type-Options
			// Prevent MIME type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// X-Frame-Options
			// Prevent clickjacking attacks
			w.Header().Set("X-Frame-Options", "DENY")

			// X-XSS-Protection
			// Enable XSS filter in older browsers
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Content-Security-Policy
			// Restrict resource loading to prevent XSS
			csp := sh.buildCSP()
			w.Header().Set("Content-Security-Policy", csp)

			// Referrer-Policy
			// Control referrer information
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Permissions-Policy (formerly Feature-Policy)
			// Disable unnecessary browser features
			permissions := sh.buildPermissionsPolicy()
			w.Header().Set("Permissions-Policy", permissions)

			// X-Permitted-Cross-Domain-Policies
			// Restrict Adobe Flash and PDF cross-domain requests
			w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")

			// Remove server identification
			w.Header().Del("Server")
			w.Header().Del("X-Powered-By")

			next.ServeHTTP(w, r)
		})
	}
}

// buildCSP builds Content Security Policy header
func (sh *SecurityHeaders) buildCSP() string {
	policies := []string{
		"default-src 'self'",
		"script-src 'self' 'unsafe-inline' 'unsafe-eval'", // TODO: Remove unsafe-* in production
		"style-src 'self' 'unsafe-inline'",
		"img-src 'self' data: https:",
		"font-src 'self' data:",
		"connect-src 'self'",
		"frame-ancestors 'none'",
		"base-uri 'self'",
		"form-action 'self'",
		"upgrade-insecure-requests",
	}

	// In development, be more permissive
	if sh.config.IsDevelopment() {
		policies = []string{
			"default-src 'self'",
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'",
			"style-src 'self' 'unsafe-inline'",
			"img-src 'self' data: https: http:",
			"font-src 'self' data:",
			"connect-src 'self' http://localhost:* ws://localhost:*",
		}
	}

	return strings.Join(policies, "; ")
}

// buildPermissionsPolicy builds Permissions-Policy header
func (sh *SecurityHeaders) buildPermissionsPolicy() string {
	permissions := []string{
		"accelerometer=()",
		"ambient-light-sensor=()",
		"autoplay=()",
		"battery=()",
		"camera=()",
		"display-capture=()",
		"document-domain=()",
		"encrypted-media=()",
		"execution-while-not-rendered=()",
		"execution-while-out-of-viewport=()",
		"fullscreen=()",
		"geolocation=()",
		"gyroscope=()",
		"layout-animations=()",
		"legacy-image-formats=()",
		"magnetometer=()",
		"microphone=()",
		"midi=()",
		"navigation-override=()",
		"oversized-images=()",
		"payment=()",
		"picture-in-picture=()",
		"publickey-credentials-get=()",
		"sync-xhr=()",
		"usb=()",
		"vr=()",
		"wake-lock=()",
		"screen-wake-lock=()",
		"web-share=()",
		"xr-spatial-tracking=()",
	}

	return strings.Join(permissions, ", ")
}

// CORS middleware with proper configuration
func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
				if origin == allowedOrigin || allowedOrigin == "*" {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")

				// Handle preflight requests
				if r.Method == http.MethodOptions {
					methods := strings.Join(cfg.CORS.AllowedMethods, ", ")
					headers := strings.Join(cfg.CORS.AllowedHeaders, ", ")

					w.Header().Set("Access-Control-Allow-Methods", methods)
					w.Header().Set("Access-Control-Allow-Headers", headers)
					w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CSRF protection middleware
type CSRF struct {
	config *config.Config
}

// NewCSRF creates a new CSRF middleware
func NewCSRF(cfg *config.Config) *CSRF {
	return &CSRF{config: cfg}
}

// Protect returns middleware that validates CSRF tokens
func (c *CSRF) Protect() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip CSRF for safe methods
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			// Skip CSRF for API endpoints (they use JWT auth)
			if strings.HasPrefix(r.URL.Path, "/api/") {
				next.ServeHTTP(w, r)
				return
			}

			// For web forms, validate CSRF token
			token := r.Header.Get("X-CSRF-Token")
			if token == "" {
				token = r.FormValue("csrf_token")
			}

			// TODO: Implement proper CSRF token validation
			// For now, just check if token exists
			if token == "" {
				http.Error(w, "CSRF token missing", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestID adds a unique request ID to each request
func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				// Generate a simple request ID
				requestID = fmt.Sprintf("%d", time.Now().UnixNano())
			}

			// Add to response headers
			w.Header().Set("X-Request-ID", requestID)

			// Add to context
			ctx := context.WithValue(r.Context(), "request_id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
