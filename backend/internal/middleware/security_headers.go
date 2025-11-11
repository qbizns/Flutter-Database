package middleware

import (
	"net/http"
	"strings"
)

// SecurityHeadersConfig holds configuration for security headers
type SecurityHeadersConfig struct {
	// ContentSecurityPolicy sets the Content-Security-Policy header
	// Default: "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'"
	ContentSecurityPolicy string

	// XContentTypeOptions sets X-Content-Type-Options header
	// Default: "nosniff"
	XContentTypeOptions string

	// XFrameOptions sets X-Frame-Options header
	// Default: "DENY"
	XFrameOptions string

	// XXSSProtection sets X-XSS-Protection header
	// Default: "1; mode=block"
	XXSSProtection string

	// StrictTransportSecurity sets Strict-Transport-Security header (HSTS)
	// Default: "max-age=31536000; includeSubDomains"
	StrictTransportSecurity string

	// ReferrerPolicy sets Referrer-Policy header
	// Default: "strict-origin-when-cross-origin"
	ReferrerPolicy string

	// PermissionsPolicy sets Permissions-Policy header
	// Default: "geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=(), accelerometer=()"
	PermissionsPolicy string

	// CustomHeaders allows adding custom security headers
	CustomHeaders map[string]string

	// EnableHSTS enables HTTP Strict Transport Security
	// Only set to true when using HTTPS in production
	EnableHSTS bool

	// IsDevelopment when true, uses relaxed CSP for development
	IsDevelopment bool

	// AllowedFrameOrigins specifies origins allowed to embed the application in iframes
	// Empty means no framing allowed (DENY)
	AllowedFrameOrigins []string
}

// DefaultSecurityHeadersConfig returns a production-ready security headers configuration
func DefaultSecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		ContentSecurityPolicy: "default-src 'self'; " +
			"script-src 'self'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'",
		XContentTypeOptions:     "nosniff",
		XFrameOptions:           "DENY",
		XXSSProtection:          "1; mode=block",
		StrictTransportSecurity: "max-age=31536000; includeSubDomains",
		ReferrerPolicy:          "strict-origin-when-cross-origin",
		PermissionsPolicy: "geolocation=(), microphone=(), camera=(), " +
			"payment=(), usb=(), magnetometer=(), gyroscope=(), accelerometer=()",
		EnableHSTS:    true,
		IsDevelopment: false,
	}
}

// DevelopmentSecurityHeadersConfig returns a relaxed configuration for development
func DevelopmentSecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https: http:; " +
			"font-src 'self' data:; " +
			"connect-src 'self' ws: wss: http: https:; " +
			"frame-ancestors 'self'; " +
			"base-uri 'self'; " +
			"form-action 'self'",
		XContentTypeOptions:     "nosniff",
		XFrameOptions:           "SAMEORIGIN",
		XXSSProtection:          "1; mode=block",
		StrictTransportSecurity: "",
		ReferrerPolicy:          "no-referrer-when-downgrade",
		PermissionsPolicy:       "",
		EnableHSTS:              false,
		IsDevelopment:           true,
	}
}

// SecurityHeaders middleware adds security headers to all responses
type SecurityHeaders struct {
	config SecurityHeadersConfig
}

// NewSecurityHeaders creates a new security headers middleware
func NewSecurityHeaders(config SecurityHeadersConfig) *SecurityHeaders {
	// Set defaults if not provided
	if config.XContentTypeOptions == "" {
		config.XContentTypeOptions = "nosniff"
	}
	if config.XFrameOptions == "" && len(config.AllowedFrameOrigins) == 0 {
		config.XFrameOptions = "DENY"
	}
	if config.XXSSProtection == "" {
		config.XXSSProtection = "1; mode=block"
	}
	if config.ReferrerPolicy == "" {
		config.ReferrerPolicy = "strict-origin-when-cross-origin"
	}

	// Handle frame options with allowed origins
	if len(config.AllowedFrameOrigins) > 0 {
		if len(config.AllowedFrameOrigins) == 1 && config.AllowedFrameOrigins[0] == "*" {
			config.XFrameOptions = "ALLOWALL"
		} else {
			// If specific origins are allowed, we'll use CSP frame-ancestors instead
			config.XFrameOptions = "SAMEORIGIN"
		}
	}

	// Build CSP with frame-ancestors if specified
	if config.ContentSecurityPolicy != "" && len(config.AllowedFrameOrigins) > 0 {
		origins := strings.Join(config.AllowedFrameOrigins, " ")

		// Check if frame-ancestors is already in CSP
		if strings.Contains(config.ContentSecurityPolicy, "frame-ancestors") {
			// Replace existing frame-ancestors directive
			parts := strings.Split(config.ContentSecurityPolicy, ";")
			var newParts []string
			for _, part := range parts {
				trimmed := strings.TrimSpace(part)
				if !strings.HasPrefix(trimmed, "frame-ancestors") {
					newParts = append(newParts, part)
				}
			}
			config.ContentSecurityPolicy = strings.Join(newParts, ";") + "; frame-ancestors " + origins
		} else {
			// Add frame-ancestors directive
			config.ContentSecurityPolicy += "; frame-ancestors " + origins
		}
	}

	return &SecurityHeaders{
		config: config,
	}
}

// Middleware returns the security headers middleware
func (s *SecurityHeaders) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add Content-Security-Policy
		if s.config.ContentSecurityPolicy != "" {
			w.Header().Set("Content-Security-Policy", s.config.ContentSecurityPolicy)
		}

		// Add X-Content-Type-Options
		if s.config.XContentTypeOptions != "" {
			w.Header().Set("X-Content-Type-Options", s.config.XContentTypeOptions)
		}

		// Add X-Frame-Options
		if s.config.XFrameOptions != "" {
			w.Header().Set("X-Frame-Options", s.config.XFrameOptions)
		}

		// Add X-XSS-Protection
		if s.config.XXSSProtection != "" {
			w.Header().Set("X-XSS-Protection", s.config.XXSSProtection)
		}

		// Add Strict-Transport-Security (only if HTTPS)
		if s.config.EnableHSTS && s.config.StrictTransportSecurity != "" {
			// Only set HSTS if the request is over HTTPS
			if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
				w.Header().Set("Strict-Transport-Security", s.config.StrictTransportSecurity)
			}
		}

		// Add Referrer-Policy
		if s.config.ReferrerPolicy != "" {
			w.Header().Set("Referrer-Policy", s.config.ReferrerPolicy)
		}

		// Add Permissions-Policy
		if s.config.PermissionsPolicy != "" {
			w.Header().Set("Permissions-Policy", s.config.PermissionsPolicy)
		}

		// Add custom headers
		for key, value := range s.config.CustomHeaders {
			w.Header().Set(key, value)
		}

		// Add additional security headers
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("X-Download-Options", "noopen")

		next.ServeHTTP(w, r)
	})
}

// WithCustomCSP creates a new SecurityHeaders with a custom CSP
func WithCustomCSP(base SecurityHeadersConfig, csp string) SecurityHeadersConfig {
	base.ContentSecurityPolicy = csp
	return base
}

// WithFrameOrigins creates a new SecurityHeaders with allowed frame origins
func WithFrameOrigins(base SecurityHeadersConfig, origins []string) SecurityHeadersConfig {
	base.AllowedFrameOrigins = origins
	return base
}

// WithCustomHeader adds a custom header to the configuration
func WithCustomHeader(base SecurityHeadersConfig, key, value string) SecurityHeadersConfig {
	if base.CustomHeaders == nil {
		base.CustomHeaders = make(map[string]string)
	}
	base.CustomHeaders[key] = value
	return base
}

// APISecurityHeadersConfig returns security headers optimized for API endpoints
func APISecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		ContentSecurityPolicy: "default-src 'none'; " +
			"frame-ancestors 'none'",
		XContentTypeOptions:     "nosniff",
		XFrameOptions:           "DENY",
		XXSSProtection:          "1; mode=block",
		StrictTransportSecurity: "max-age=31536000; includeSubDomains; preload",
		ReferrerPolicy:          "no-referrer",
		PermissionsPolicy:       "geolocation=(), microphone=(), camera=(), payment=(), usb=()",
		EnableHSTS:              true,
		IsDevelopment:           false,
		CustomHeaders: map[string]string{
			"X-API-Version": "1.0",
		},
	}
}

// GetSecureHeaders returns a map of all security headers that will be applied
// Useful for testing and debugging
func (s *SecurityHeaders) GetSecureHeaders() map[string]string {
	headers := make(map[string]string)

	if s.config.ContentSecurityPolicy != "" {
		headers["Content-Security-Policy"] = s.config.ContentSecurityPolicy
	}
	if s.config.XContentTypeOptions != "" {
		headers["X-Content-Type-Options"] = s.config.XContentTypeOptions
	}
	if s.config.XFrameOptions != "" {
		headers["X-Frame-Options"] = s.config.XFrameOptions
	}
	if s.config.XXSSProtection != "" {
		headers["X-XSS-Protection"] = s.config.XXSSProtection
	}
	if s.config.EnableHSTS && s.config.StrictTransportSecurity != "" {
		headers["Strict-Transport-Security"] = s.config.StrictTransportSecurity
	}
	if s.config.ReferrerPolicy != "" {
		headers["Referrer-Policy"] = s.config.ReferrerPolicy
	}
	if s.config.PermissionsPolicy != "" {
		headers["Permissions-Policy"] = s.config.PermissionsPolicy
	}

	for key, value := range s.config.CustomHeaders {
		headers[key] = value
	}

	headers["X-Permitted-Cross-Domain-Policies"] = "none"
	headers["X-Download-Options"] = "noopen"

	return headers
}
