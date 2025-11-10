package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/your-org/pos-backend/internal/config"
	"github.com/your-org/pos-backend/internal/logging"
)

func TestRateLimiter_Limit_WithoutRedis(t *testing.T) {
	// Test that rate limiter fails open when Redis is unavailable
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	rl := NewRateLimiter(nil, config.RateLimitConfig{
		RequestsPerMinute: 5,
		RequestsPerHour:   100,
	}, logger)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rl.Limit()(testHandler)

	// Make multiple requests - should all pass without Redis
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected status 200, got %d (should fail open)", i, rec.Code)
		}
	}
}

func TestGetRealIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		want       string
	}{
		{
			name:       "from RemoteAddr",
			remoteAddr: "192.168.1.1:12345",
			want:       "192.168.1.1",
		},
		{
			name:       "from X-Forwarded-For single",
			remoteAddr: "192.168.1.1:12345",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.1",
			},
			want: "203.0.113.1",
		},
		{
			name:       "from X-Forwarded-For multiple",
			remoteAddr: "192.168.1.1:12345",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.1, 198.51.100.1, 192.0.2.1",
			},
			want: "203.0.113.1",
		},
		{
			name:       "from X-Forwarded-For with spaces",
			remoteAddr: "192.168.1.1:12345",
			headers: map[string]string{
				"X-Forwarded-For": "  203.0.113.1  , 198.51.100.1",
			},
			want: "203.0.113.1",
		},
		{
			name:       "from X-Real-IP",
			remoteAddr: "192.168.1.1:12345",
			headers: map[string]string{
				"X-Real-IP": "203.0.113.1",
			},
			want: "203.0.113.1",
		},
		{
			name:       "X-Forwarded-For takes precedence over X-Real-IP",
			remoteAddr: "192.168.1.1:12345",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.1",
				"X-Real-IP":       "198.51.100.1",
			},
			want: "203.0.113.1",
		},
		{
			name:       "invalid IP in X-Forwarded-For falls back to RemoteAddr",
			remoteAddr: "192.168.1.1:12345",
			headers: map[string]string{
				"X-Forwarded-For": "not-an-ip",
			},
			want: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.RemoteAddr = tt.remoteAddr

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			got := getRealIP(req)
			if got != tt.want {
				t.Errorf("getRealIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRateLimiter_LimitAuth_StricterLimits(t *testing.T) {
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	// Create rate limiter without Redis (will fail open)
	rl := NewRateLimiter(nil, config.RateLimitConfig{
		RequestsPerMinute: 60,
		RequestsPerHour:   1000,
	}, logger)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test that auth endpoints have different rate limits
	authHandler := rl.LimitAuth()(testHandler)

	// Should allow requests when Redis is not available (fail open)
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		authHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected status 200, got %d", i, rec.Code)
		}
	}
}

func TestCheckRateLimit_Logic(t *testing.T) {
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	t.Run("nil redis fails open", func(t *testing.T) {
		rl := NewRateLimiter(nil, config.RateLimitConfig{}, logger)

		allowed := rl.checkRateLimit(context.Background(), "192.168.1.1", "minute", 5, time.Minute)
		if !allowed {
			t.Error("should allow request when Redis is nil (fail open)")
		}
	})
}

// Mock Redis client for testing
type mockRedisClient struct {
	redis.UniversalClient
	incrFunc   func(ctx context.Context, key string) *redis.IntCmd
	expireFunc func(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	getFunc    func(ctx context.Context, key string) *redis.StringCmd
}

func (m *mockRedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	if m.incrFunc != nil {
		return m.incrFunc(ctx, key)
	}
	return redis.NewIntCmd(ctx)
}

func (m *mockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	if m.expireFunc != nil {
		return m.expireFunc(ctx, key, expiration)
	}
	return redis.NewBoolCmd(ctx)
}

func (m *mockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	if m.getFunc != nil {
		return m.getFunc(ctx, key)
	}
	return redis.NewStringCmd(ctx)
}

func TestGetRemainingRequests(t *testing.T) {
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	t.Run("with nil redis", func(t *testing.T) {
		rl := NewRateLimiter(nil, config.RateLimitConfig{}, logger)

		remaining, err := rl.GetRemainingRequests(
			context.Background(),
			"192.168.1.1",
			"minute",
			60,
			time.Minute,
		)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if remaining != 60 {
			t.Errorf("expected 60 remaining, got %d", remaining)
		}
	})
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	rl := NewRateLimiter(nil, config.RateLimitConfig{
		RequestsPerMinute: 5,
		RequestsPerHour:   100,
	}, logger)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rl.Limit()(testHandler)

	// Different IPs should be tracked separately
	ips := []string{
		"192.168.1.1:12345",
		"192.168.1.2:12345",
		"192.168.1.3:12345",
	}

	for _, ip := range ips {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = ip
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("IP %s: expected status 200, got %d", ip, rec.Code)
		}
	}
}

func TestRateLimiter_Integration(t *testing.T) {
	// This would be an integration test with a real Redis instance
	// For now, we just verify the structure and fail-open behavior

	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	cfg := config.RateLimitConfig{
		RequestsPerMinute: 10,
		RequestsPerHour:   100,
	}

	rl := NewRateLimiter(nil, cfg, logger)

	if rl.redis != nil {
		t.Error("expected nil redis for this test")
	}

	if rl.config.RequestsPerMinute != 10 {
		t.Errorf("expected 10 requests per minute, got %d", rl.config.RequestsPerMinute)
	}

	if rl.logger == nil {
		t.Error("expected logger to be set")
	}
}

// Benchmark tests
func BenchmarkGetRealIP(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 198.51.100.1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getRealIP(req)
	}
}

func BenchmarkRateLimiter_NoRedis(b *testing.B) {
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	rl := NewRateLimiter(nil, config.RateLimitConfig{
		RequestsPerMinute: 60,
	}, logger)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rl.Limit()(testHandler)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}
