package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
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
	// Production-safe CSP - NO unsafe-inline or unsafe-eval
	if sh.config.IsProduction() {
		policies := []string{
			"default-src 'self'",
			"script-src 'self'", // FIXED: Removed unsafe-inline and unsafe-eval
			"style-src 'self'",  // FIXED: Removed unsafe-inline
			"img-src 'self' data: https:",
			"font-src 'self' data:",
			"connect-src 'self'",
			"frame-ancestors 'none'",
			"base-uri 'self'",
			"form-action 'self'",
			"upgrade-insecure-requests",
			"object-src 'none'", // Block plugins
		}
		return strings.Join(policies, "; ")
	}

	// In development, be more permissive for hot reload and debugging
	policies := []string{
		"default-src 'self'",
		"script-src 'self' 'unsafe-inline' 'unsafe-eval'", // Allow for development
		"style-src 'self' 'unsafe-inline'",
		"img-src 'self' data: https: http:",
		"font-src 'self' data:",
		"connect-src 'self' http://localhost:* ws://localhost:* wss://localhost:*",
		"frame-ancestors 'self'",
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
	store  *CSRFTokenStore
}

// CSRFTokenStore manages CSRF tokens with expiration
type CSRFTokenStore struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

// NewCSRFTokenStore creates a new CSRF token store
func NewCSRFTokenStore() *CSRFTokenStore {
	store := &CSRFTokenStore{
		tokens: make(map[string]time.Time),
	}
	// Cleanup expired tokens every 10 minutes
	go store.cleanupExpired()
	return store
}

// Generate creates a new CSRF token
func (s *CSRFTokenStore) Generate() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(bytes)

	s.mu.Lock()
	s.tokens[token] = time.Now().Add(2 * time.Hour)
	s.mu.Unlock()

	return token, nil
}

// Validate checks if a CSRF token is valid
func (s *CSRFTokenStore) Validate(token string) bool {
	if token == "" {
		return false
	}

	s.mu.RLock()
	expiry, exists := s.tokens[token]
	s.mu.RUnlock()

	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		s.mu.Lock()
		delete(s.tokens, token)
		s.mu.Unlock()
		return false
	}

	return true
}

// Invalidate removes a CSRF token (after successful use)
func (s *CSRFTokenStore) Invalidate(token string) {
	s.mu.Lock()
	delete(s.tokens, token)
	s.mu.Unlock()
}

// cleanupExpired removes expired tokens periodically
func (s *CSRFTokenStore) cleanupExpired() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for token, expiry := range s.tokens {
			if now.After(expiry) {
				delete(s.tokens, token)
			}
		}
		s.mu.Unlock()
	}
}

// NewCSRF creates a new CSRF middleware
func NewCSRF(cfg *config.Config) *CSRF {
	return &CSRF{
		config: cfg,
		store:  NewCSRFTokenStore(),
	}
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

			if !c.store.Validate(token) {
				http.Error(w, "CSRF token invalid or expired", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GenerateToken creates a new CSRF token for a session
func (c *CSRF) GenerateToken() (string, error) {
	return c.store.Generate()
}

// RequestID adds a unique request ID to each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// FIXED: Use UUID instead of timestamp for proper uniqueness
			// This ensures uniqueness across distributed systems
			requestID = generateRequestID()
		}

		// Add to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Add to context
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// generateRequestID creates a cryptographically secure unique request ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp if random fails (should never happen)
		return fmt.Sprintf("req_%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("req_%s", base64.URLEncoding.EncodeToString(bytes)[:22])
}
