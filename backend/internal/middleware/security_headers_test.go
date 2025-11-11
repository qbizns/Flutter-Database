package middleware

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSecurityHeaders_DefaultConfig(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tests := []struct {
		header   string
		expected string
	}{
		{"X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", "DENY"},
		{"X-XSS-Protection", "1; mode=block"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
		{"X-Permitted-Cross-Domain-Policies", "none"},
		{"X-Download-Options", "noopen"},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			value := w.Header().Get(tt.header)
			if value != tt.expected {
				t.Errorf("header %s: expected '%s', got '%s'", tt.header, tt.expected, value)
			}
		})
	}

	// Check CSP is set
	csp := w.Header().Get("Content-Security-Policy")
	if csp == "" {
		t.Error("Content-Security-Policy header not set")
	}
	if !strings.Contains(csp, "default-src 'self'") {
		t.Error("CSP should contain default-src 'self'")
	}

	// Check Permissions-Policy is set
	pp := w.Header().Get("Permissions-Policy")
	if pp == "" {
		t.Error("Permissions-Policy header not set")
	}
}

func TestSecurityHeaders_HSTS_WithHTTPS(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	config.EnableHSTS = true
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "https://example.com/api/test", nil)
	req.TLS = &tls.ConnectionState{} // Simulate HTTPS
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	hsts := w.Header().Get("Strict-Transport-Security")
	if hsts == "" {
		t.Error("HSTS header should be set for HTTPS requests")
	}
	if !strings.Contains(hsts, "max-age=31536000") {
		t.Error("HSTS should contain max-age=31536000")
	}
}

func TestSecurityHeaders_HSTS_WithoutHTTPS(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	config.EnableHSTS = true
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "http://example.com/api/test", nil)
	// No TLS connection state - HTTP request
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	hsts := w.Header().Get("Strict-Transport-Security")
	if hsts != "" {
		t.Error("HSTS header should not be set for HTTP requests")
	}
}

func TestSecurityHeaders_HSTS_XForwardedProto(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	config.EnableHSTS = true
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	hsts := w.Header().Get("Strict-Transport-Security")
	if hsts == "" {
		t.Error("HSTS header should be set when X-Forwarded-Proto is https")
	}
}

func TestSecurityHeaders_DevelopmentConfig(t *testing.T) {
	config := DevelopmentSecurityHeadersConfig()
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Development should have relaxed CSP
	csp := w.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "unsafe-inline") || !strings.Contains(csp, "unsafe-eval") {
		t.Error("Development CSP should allow unsafe-inline and unsafe-eval")
	}

	// Development should use SAMEORIGIN for X-Frame-Options
	xfo := w.Header().Get("X-Frame-Options")
	if xfo != "SAMEORIGIN" {
		t.Errorf("Development X-Frame-Options should be SAMEORIGIN, got %s", xfo)
	}

	// HSTS should not be set in development
	hsts := w.Header().Get("Strict-Transport-Security")
	if hsts != "" {
		t.Error("HSTS should not be set in development")
	}
}

func TestSecurityHeaders_APIConfig(t *testing.T) {
	config := APISecurityHeadersConfig()
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.TLS = &tls.ConnectionState{} // Simulate HTTPS
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// API config should have strict CSP
	csp := w.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "default-src 'none'") {
		t.Error("API CSP should have default-src 'none'")
	}

	// Check for preload in HSTS
	hsts := w.Header().Get("Strict-Transport-Security")
	if !strings.Contains(hsts, "preload") {
		t.Error("API HSTS should include preload")
	}

	// Check for no-referrer
	rp := w.Header().Get("Referrer-Policy")
	if rp != "no-referrer" {
		t.Errorf("API Referrer-Policy should be no-referrer, got %s", rp)
	}

	// Check custom header
	apiVersion := w.Header().Get("X-API-Version")
	if apiVersion != "1.0" {
		t.Error("Custom X-API-Version header should be set")
	}
}

func TestSecurityHeaders_CustomHeaders(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	config.CustomHeaders = map[string]string{
		"X-Custom-Header-1": "value1",
		"X-Custom-Header-2": "value2",
	}
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Header().Get("X-Custom-Header-1") != "value1" {
		t.Error("X-Custom-Header-1 should be set")
	}
	if w.Header().Get("X-Custom-Header-2") != "value2" {
		t.Error("X-Custom-Header-2 should be set")
	}
}

func TestSecurityHeaders_AllowedFrameOrigins(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	config.AllowedFrameOrigins = []string{"https://trusted.com"}
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// When specific origins are allowed, X-Frame-Options should be SAMEORIGIN
	xfo := w.Header().Get("X-Frame-Options")
	if xfo != "SAMEORIGIN" {
		t.Errorf("expected X-Frame-Options SAMEORIGIN, got %s", xfo)
	}

	// CSP should include frame-ancestors
	csp := w.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "frame-ancestors") {
		t.Error("CSP should include frame-ancestors when AllowedFrameOrigins is set")
	}
	if !strings.Contains(csp, "https://trusted.com") {
		t.Error("CSP frame-ancestors should include the allowed origin")
	}
}

func TestSecurityHeaders_AllowAllFrames(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	config.AllowedFrameOrigins = []string{"*"}
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	xfo := w.Header().Get("X-Frame-Options")
	if xfo != "ALLOWALL" {
		t.Errorf("expected X-Frame-Options ALLOWALL, got %s", xfo)
	}
}

func TestWithCustomCSP(t *testing.T) {
	base := DefaultSecurityHeadersConfig()
	config := WithCustomCSP(base, "default-src 'none'; script-src 'self'")

	if config.ContentSecurityPolicy != "default-src 'none'; script-src 'self'" {
		t.Error("Custom CSP not applied")
	}
}

func TestWithFrameOrigins(t *testing.T) {
	base := DefaultSecurityHeadersConfig()
	config := WithFrameOrigins(base, []string{"https://example.com"})

	if len(config.AllowedFrameOrigins) != 1 {
		t.Error("Frame origins not set")
	}
	if config.AllowedFrameOrigins[0] != "https://example.com" {
		t.Error("Frame origin value incorrect")
	}
}

func TestWithCustomHeader(t *testing.T) {
	base := DefaultSecurityHeadersConfig()
	config := WithCustomHeader(base, "X-Test", "test-value")

	if config.CustomHeaders["X-Test"] != "test-value" {
		t.Error("Custom header not set")
	}
}

func TestGetSecureHeaders(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	config.CustomHeaders = map[string]string{
		"X-Custom": "custom-value",
	}
	sh := NewSecurityHeaders(config)

	headers := sh.GetSecureHeaders()

	expectedHeaders := []string{
		"Content-Security-Policy",
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-XSS-Protection",
		"Strict-Transport-Security",
		"Referrer-Policy",
		"Permissions-Policy",
		"X-Permitted-Cross-Domain-Policies",
		"X-Download-Options",
		"X-Custom",
	}

	for _, header := range expectedHeaders {
		if _, exists := headers[header]; !exists {
			t.Errorf("expected header %s not found in GetSecureHeaders", header)
		}
	}
}

func TestSecurityHeaders_EmptyConfig(t *testing.T) {
	// Test with minimal config - should use defaults
	config := SecurityHeadersConfig{}
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should have default headers
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Default X-Content-Type-Options not set")
	}
	if w.Header().Get("X-Frame-Options") != "DENY" {
		t.Error("Default X-Frame-Options not set")
	}
	if w.Header().Get("X-XSS-Protection") != "1; mode=block" {
		t.Error("Default X-XSS-Protection not set")
	}
}

func TestSecurityHeaders_MultipleRequests(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	sh := NewSecurityHeaders(config)

	handler := sh.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Make multiple requests to ensure headers are consistent
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Errorf("Request %d: X-Content-Type-Options header not consistent", i+1)
		}
	}
}
