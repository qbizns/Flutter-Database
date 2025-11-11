package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCSRFProtection_SafeMethods(t *testing.T) {
	config := DefaultCSRFConfig()
	config.CookieSecure = false // For testing
	csrf := NewCSRFProtection(config)

	handler := csrf.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	tests := []struct {
		name   string
		method string
	}{
		{"GET", "GET"},
		{"HEAD", "HEAD"},
		{"OPTIONS", "OPTIONS"},
		{"TRACE", "TRACE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Code)
			}

			// Should have CSRF token in header
			token := w.Header().Get(CSRFTokenHeader)
			if token == "" {
				t.Error("expected CSRF token in response header")
			}

			// Should have CSRF cookie
			cookies := w.Result().Cookies()
			found := false
			for _, cookie := range cookies {
				if cookie.Name == CSRFCookieName {
					found = true
					if cookie.Value == "" {
						t.Error("CSRF cookie value is empty")
					}
					if !cookie.HttpOnly {
						t.Error("CSRF cookie should be HttpOnly")
					}
				}
			}
			if !found {
				t.Error("CSRF cookie not found")
			}
		})
	}
}

func TestCSRFProtection_UnsafeMethods_NoToken(t *testing.T) {
	config := DefaultCSRFConfig()
	config.CookieSecure = false
	csrf := NewCSRFProtection(config)

	handler := csrf.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name   string
		method string
	}{
		{"POST", "POST"},
		{"PUT", "PUT"},
		{"DELETE", "DELETE"},
		{"PATCH", "PATCH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusForbidden {
				t.Errorf("expected status 403, got %d", w.Code)
			}

			body := w.Body.String()
			if !strings.Contains(body, "CSRF") {
				t.Errorf("expected CSRF error message, got: %s", body)
			}
		})
	}
}

func TestCSRFProtection_UnsafeMethods_WithValidToken(t *testing.T) {
	config := DefaultCSRFConfig()
	config.CookieSecure = false
	csrf := NewCSRFProtection(config)

	handler := csrf.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// First, get a CSRF token via GET request
	getReq := httptest.NewRequest("GET", "/api/test", nil)
	getW := httptest.NewRecorder()
	handler.ServeHTTP(getW, getReq)

	token := getW.Header().Get(CSRFTokenHeader)
	if token == "" {
		t.Fatal("failed to get CSRF token")
	}

	var csrfCookie *http.Cookie
	for _, cookie := range getW.Result().Cookies() {
		if cookie.Name == CSRFCookieName {
			csrfCookie = cookie
			break
		}
	}
	if csrfCookie == nil {
		t.Fatal("CSRF cookie not found")
	}

	tests := []struct {
		name   string
		method string
	}{
		{"POST", "POST"},
		{"PUT", "PUT"},
		{"DELETE", "DELETE"},
		{"PATCH", "PATCH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/test", nil)
			req.Header.Set(CSRFTokenHeader, token)
			req.AddCookie(csrfCookie)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestCSRFProtection_InvalidToken(t *testing.T) {
	config := DefaultCSRFConfig()
	config.CookieSecure = false
	csrf := NewCSRFProtection(config)

	handler := csrf.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name        string
		headerToken string
		cookieToken string
	}{
		{
			name:        "Mismatched tokens",
			headerToken: "token1",
			cookieToken: "token2",
		},
		{
			name:        "Invalid token format",
			headerToken: "invalid",
			cookieToken: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/test", nil)
			req.Header.Set(CSRFTokenHeader, tt.headerToken)
			req.AddCookie(&http.Cookie{
				Name:  CSRFCookieName,
				Value: tt.cookieToken,
			})

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusForbidden {
				t.Errorf("expected status 403, got %d", w.Code)
			}
		})
	}
}

func TestCSRFProtection_SkipPaths(t *testing.T) {
	config := DefaultCSRFConfig()
	config.CookieSecure = false
	config.SkipPaths = []string{"/health", "/metrics", "/auth/login", "/auth/register"}
	csrf := NewCSRFProtection(config)

	handler := csrf.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name string
		path string
	}{
		{"Health check", "/health"},
		{"Metrics", "/metrics"},
		{"Login", "/auth/login"},
		{"Register", "/auth/register"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200 for skip path, got %d", w.Code)
			}
		})
	}
}

func TestCSRFProtection_TokenExpiry(t *testing.T) {
	config := DefaultCSRFConfig()
	config.CookieSecure = false
	config.TokenLifetime = 100 * time.Millisecond // Very short lifetime for testing
	csrf := NewCSRFProtection(config)

	handler := csrf.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Get a token
	getReq := httptest.NewRequest("GET", "/api/test", nil)
	getW := httptest.NewRecorder()
	handler.ServeHTTP(getW, getReq)

	token := getW.Header().Get(CSRFTokenHeader)
	var csrfCookie *http.Cookie
	for _, cookie := range getW.Result().Cookies() {
		if cookie.Name == CSRFCookieName {
			csrfCookie = cookie
			break
		}
	}

	// Wait for token to expire
	time.Sleep(150 * time.Millisecond)

	// Try to use expired token
	req := httptest.NewRequest("POST", "/api/test", nil)
	req.Header.Set(CSRFTokenHeader, token)
	req.AddCookie(csrfCookie)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403 for expired token, got %d", w.Code)
	}
}

func TestCSRFProtection_FormToken(t *testing.T) {
	config := DefaultCSRFConfig()
	config.CookieSecure = false
	csrf := NewCSRFProtection(config)

	handler := csrf.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Get a token first
	getReq := httptest.NewRequest("GET", "/api/test", nil)
	getW := httptest.NewRecorder()
	handler.ServeHTTP(getW, getReq)

	token := getW.Header().Get(CSRFTokenHeader)
	var csrfCookie *http.Cookie
	for _, cookie := range getW.Result().Cookies() {
		if cookie.Name == CSRFCookieName {
			csrfCookie = cookie
			break
		}
	}

	// Submit token via form data
	formData := strings.NewReader("csrf_token=" + token + "&data=test")
	req := httptest.NewRequest("POST", "/api/test", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(csrfCookie)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 with form token, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSecureCompare(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"Equal strings", "token123", "token123", true},
		{"Different strings", "token123", "token456", false},
		{"Different lengths", "short", "longer", false},
		{"Empty strings", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := secureCompare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	config := DefaultCSRFConfig()
	csrf := NewCSRFProtection(config)

	token1, err := csrf.generateToken()
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	token2, err := csrf.generateToken()
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token1 == token2 {
		t.Error("generated tokens should be unique")
	}

	if len(token1) == 0 {
		t.Error("generated token is empty")
	}
}
