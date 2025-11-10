package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/pos-backend/internal/config"
)

func TestSecurityHeaders_Handler(t *testing.T) {
	tests := []struct {
		name          string
		env           string
		checkHeaders  map[string]string
		checkContains map[string]string
	}{
		{
			name: "production environment",
			env:  "production",
			checkHeaders: map[string]string{
				"Strict-Transport-Security":        "max-age=31536000; includeSubDomains; preload",
				"X-Content-Type-Options":           "nosniff",
				"X-Frame-Options":                  "DENY",
				"X-XSS-Protection":                 "1; mode=block",
				"Referrer-Policy":                  "strict-origin-when-cross-origin",
				"X-Permitted-Cross-Domain-Policies": "none",
			},
			checkContains: map[string]string{
				"Content-Security-Policy": "default-src 'self'",
				"Permissions-Policy":      "camera=()",
			},
		},
		{
			name: "development environment",
			env:  "development",
			checkHeaders: map[string]string{
				"X-Content-Type-Options": "nosniff",
				"X-Frame-Options":        "DENY",
				"X-XSS-Protection":       "1; mode=block",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test config
			cfg := &config.Config{
				Server: config.ServerConfig{
					Env: tt.env,
				},
			}

			// Create security headers middleware
			sh := NewSecurityHeaders(cfg)

			// Create test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with security headers
			handler := sh.Handler()(testHandler)

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			// Execute request
			handler.ServeHTTP(rec, req)

			// Check status
			if rec.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rec.Code)
			}

			// Check exact headers
			for key, expected := range tt.checkHeaders {
				actual := rec.Header().Get(key)
				if actual != expected {
					t.Errorf("header %s = %q, want %q", key, actual, expected)
				}
			}

			// Check headers contain substring
			for key, contains := range tt.checkContains {
				actual := rec.Header().Get(key)
				if actual == "" {
					t.Errorf("header %s not set", key)
					continue
				}
				if !containsString(actual, contains) {
					t.Errorf("header %s = %q, should contain %q", key, actual, contains)
				}
			}
		})
	}
}

func TestSecurityHeaders_RemovesServerHeader(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Env: "production"},
	}

	sh := NewSecurityHeaders(cfg)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "BadServer/1.0")
		w.Header().Set("X-Powered-By", "PHP/7.0")
		w.WriteHeader(http.StatusOK)
	})

	handler := sh.Handler()(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Verify server identification headers are removed
	if rec.Header().Get("Server") != "" {
		t.Error("Server header should be removed")
	}

	if rec.Header().Get("X-Powered-By") != "" {
		t.Error("X-Powered-By header should be removed")
	}
}

func TestBuildCSP_Production(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Env: "production"},
	}

	sh := NewSecurityHeaders(cfg)
	csp := sh.buildCSP()

	// Check for required directives
	requiredDirectives := []string{
		"default-src 'self'",
		"frame-ancestors 'none'",
		"base-uri 'self'",
		"form-action 'self'",
		"upgrade-insecure-requests",
	}

	for _, directive := range requiredDirectives {
		if !containsString(csp, directive) {
			t.Errorf("CSP should contain %q, got: %s", directive, csp)
		}
	}
}

func TestBuildCSP_Development(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Env: "development"},
	}

	sh := NewSecurityHeaders(cfg)
	csp := sh.buildCSP()

	// Development should allow localhost and websockets
	developmentDirectives := []string{
		"http://localhost:*",
		"ws://localhost:*",
	}

	for _, directive := range developmentDirectives {
		if !containsString(csp, directive) {
			t.Errorf("Development CSP should contain %q, got: %s", directive, csp)
		}
	}

	// Should NOT have upgrade-insecure-requests in development
	if containsString(csp, "upgrade-insecure-requests") {
		t.Error("Development CSP should not have upgrade-insecure-requests")
	}
}

func TestBuildPermissionsPolicy(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Env: "production"},
	}

	sh := NewSecurityHeaders(cfg)
	policy := sh.buildPermissionsPolicy()

	// Check for some restricted features
	restrictedFeatures := []string{
		"camera=()",
		"microphone=()",
		"geolocation=()",
		"payment=()",
		"usb=()",
	}

	for _, feature := range restrictedFeatures {
		if !containsString(policy, feature) {
			t.Errorf("Permissions-Policy should contain %q, got: %s", feature, policy)
		}
	}
}

func TestCORS_AllowedOrigin(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000", "https://example.com"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
	}

	corsMiddleware := CORS(cfg)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := corsMiddleware(testHandler)

	// Test allowed origin
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("expected CORS origin header, got %q", rec.Header().Get("Access-Control-Allow-Origin"))
	}

	if rec.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("expected credentials to be allowed")
	}
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
		},
	}

	corsMiddleware := CORS(cfg)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := corsMiddleware(testHandler)

	// Test disallowed origin
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should not set CORS headers for disallowed origin
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("should not set CORS headers for disallowed origin")
	}
}

func TestCORS_PreflightRequest(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
	}

	corsMiddleware := CORS(cfg)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for preflight")
	})

	handler := corsMiddleware(testHandler)

	// Test preflight request
	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}

	// Check preflight headers
	if rec.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("preflight should set allowed methods")
	}

	if rec.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Error("preflight should set allowed headers")
	}

	if rec.Header().Get("Access-Control-Max-Age") != "86400" {
		t.Error("preflight should set max age")
	}
}

func TestRequestID(t *testing.T) {
	middleware := RequestID()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check context has request ID
		requestID := r.Context().Value("request_id")
		if requestID == nil {
			t.Error("request_id not found in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(testHandler)

	t.Run("generates request ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		// Check response header
		requestID := rec.Header().Get("X-Request-ID")
		if requestID == "" {
			t.Error("X-Request-ID header not set")
		}
	})

	t.Run("preserves existing request ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-ID", "existing-id")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		// Should preserve existing ID
		requestID := rec.Header().Get("X-Request-ID")
		if requestID != "existing-id" {
			t.Errorf("expected existing-id, got %q", requestID)
		}
	})
}

func TestCSRF_SkipSafeMethods(t *testing.T) {
	cfg := &config.Config{}
	csrf := NewCSRF(cfg)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := csrf.Protect()(testHandler)

	safeMethods := []string{http.MethodGet, http.MethodHead, http.MethodOptions}

	for _, method := range safeMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/test", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("safe method %s should be allowed, got status %d", method, rec.Code)
			}
		})
	}
}

func TestCSRF_SkipAPIEndpoints(t *testing.T) {
	cfg := &config.Config{}
	csrf := NewCSRF(cfg)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := csrf.Protect()(testHandler)

	// POST to API endpoint without CSRF token should be allowed
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("API endpoint should skip CSRF, got status %d", rec.Code)
	}
}

func TestCSRF_RequireTokenForForms(t *testing.T) {
	cfg := &config.Config{}
	csrf := NewCSRF(cfg)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := csrf.Protect()(testHandler)

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/form/submit", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status 403, got %d", rec.Code)
		}
	})

	t.Run("with token in header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/form/submit", nil)
		req.Header.Set("X-CSRF-Token", "test-token")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		// For now it just checks if token exists
		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
	})
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
