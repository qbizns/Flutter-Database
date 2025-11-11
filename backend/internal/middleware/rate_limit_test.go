package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerWindow:    10,
		Window:               1 * time.Second,
		KeyFunc:              defaultKeyFunc,
		OnLimitReached:       defaultOnLimitReached,
		UseXRateLimitHeaders: true,
	}
	rl := NewRateLimiter(config)

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Make 10 requests (should all succeed)
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected status 200, got %d", i+1, w.Code)
		}

		// Check rate limit headers
		limit := w.Header().Get("X-RateLimit-Limit")
		if limit == "" {
			t.Error("X-RateLimit-Limit header not set")
		}

		remaining := w.Header().Get("X-RateLimit-Remaining")
		if remaining == "" {
			t.Error("X-RateLimit-Remaining header not set")
		}
	}

	// 11th request should be rate limited
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", w.Code)
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerWindow: 5,
		Window:            1 * time.Second,
		KeyFunc:           defaultKeyFunc,
		OnLimitReached:    defaultOnLimitReached,
	}
	rl := NewRateLimiter(config)

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// IP1 makes 5 requests
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("IP1 request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// IP2 should still be able to make requests
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.2:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("IP2 request %d: expected 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimiter_Refill(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerWindow: 3,
		Window:            500 * time.Millisecond,
		KeyFunc:           defaultKeyFunc,
		OnLimitReached:    defaultOnLimitReached,
	}
	rl := NewRateLimiter(config)

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Use up the limit
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// Next request should be rate limited
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}

	// Wait for refill
	time.Sleep(600 * time.Millisecond)

	// Should be able to make requests again
	req = httptest.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("after refill: expected 200, got %d", w.Code)
	}
}

func TestRateLimiter_SkipPaths(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerWindow: 2,
		Window:            1 * time.Second,
		KeyFunc:           defaultKeyFunc,
		OnLimitReached:    defaultOnLimitReached,
		SkipPaths:         []string{"/health", "/metrics"},
	}
	rl := NewRateLimiter(config)

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Make 10 requests to /health (should not be rate limited)
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("health check %d: expected 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimiter_XForwardedFor(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerWindow: 5,
		Window:            1 * time.Second,
		KeyFunc:           defaultKeyFunc,
		OnLimitReached:    defaultOnLimitReached,
	}
	rl := NewRateLimiter(config)

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Make requests with X-Forwarded-For header
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("X-Forwarded-For", "10.0.0.1")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// 6th request should be rate limited
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}

func TestRateLimiter_Burst(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerWindow: 5,
		Window:            1 * time.Second,
		AllowBurst:        3,
		KeyFunc:           defaultKeyFunc,
		OnLimitReached:    defaultOnLimitReached,
	}
	rl := NewRateLimiter(config)

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Should be able to make 5 + 3 = 8 requests initially
	for i := 0; i < 8; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// 9th request should be rate limited
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}

func TestRateLimiter_CustomKeyFunc(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerWindow: 3,
		Window:            1 * time.Second,
		KeyFunc: func(r *http.Request) string {
			// Rate limit by custom header
			return r.Header.Get("X-API-Key")
		},
		OnLimitReached: defaultOnLimitReached,
	}
	rl := NewRateLimiter(config)

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// API Key 1 makes 3 requests
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("X-API-Key", "key1")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("key1 request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// 4th request should be rate limited
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("X-API-Key", "key1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}

	// API Key 2 should still work
	req = httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("X-API-Key", "key2")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("key2 request: expected 200, got %d", w.Code)
	}
}

func TestStrictAuthRateLimitConfig(t *testing.T) {
	config := StrictAuthRateLimitConfig()
	rl := NewRateLimiter(config)

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Should allow 5 + 2 burst = 7 requests
	for i := 0; i < 7; i++ {
		req := httptest.NewRequest("POST", "/auth/login", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// 8th request should be rate limited
	req := httptest.NewRequest("POST", "/auth/login", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}

func TestKeyByUser(t *testing.T) {
	// Test with user ID in context
	req := httptest.NewRequest("GET", "/api/test", nil)
	ctx := context.WithValue(req.Context(), "user_id", "user123")
	req = req.WithContext(ctx)

	key := KeyByUser(req)
	if key != "user:user123" {
		t.Errorf("expected key 'user:user123', got '%s'", key)
	}

	// Test without user ID (should fall back to IP)
	req2 := httptest.NewRequest("GET", "/api/test", nil)
	req2.RemoteAddr = "192.168.1.1:1234"
	key2 := KeyByUser(req2)
	if key2 != "192.168.1.1" {
		t.Errorf("expected key '192.168.1.1', got '%s'", key2)
	}
}

func TestKeyByIPAndUser(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	ctx := context.WithValue(req.Context(), "user_id", "user123")
	req = req.WithContext(ctx)

	key := KeyByIPAndUser(req)
	if key != "192.168.1.1:user123" {
		t.Errorf("expected key '192.168.1.1:user123', got '%s'", key)
	}

	// Test without user ID (should use IP only)
	req2 := httptest.NewRequest("GET", "/api/test", nil)
	req2.RemoteAddr = "192.168.1.1:1234"
	key2 := KeyByIPAndUser(req2)
	if key2 != "192.168.1.1" {
		t.Errorf("expected key '192.168.1.1', got '%s'", key2)
	}
}

func TestDefaultKeyFunc_XRealIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("X-Real-IP", "203.0.113.1")
	req.RemoteAddr = "192.168.1.1:1234"

	key := defaultKeyFunc(req)
	if key != "203.0.113.1" {
		t.Errorf("expected key '203.0.113.1', got '%s'", key)
	}
}

func TestDefaultKeyFunc_XForwardedForMultiple(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 192.168.1.1, 10.0.0.1")

	key := defaultKeyFunc(req)
	if key != "203.0.113.1" {
		t.Errorf("expected key '203.0.113.1', got '%s'", key)
	}
}
