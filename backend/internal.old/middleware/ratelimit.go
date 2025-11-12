package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/your-org/pos-backend/internal/config"
	"github.com/your-org/pos-backend/internal/logging"
	"go.uber.org/zap"
)

// RateLimiter provides rate limiting middleware using Redis
type RateLimiter struct {
	redis  *redis.Client
	config config.RateLimitConfig
	logger *logging.Logger
}

// NewRateLimiter creates a new rate limiter middleware
func NewRateLimiter(redis *redis.Client, cfg config.RateLimitConfig, logger *logging.Logger) *RateLimiter {
	return &RateLimiter{
		redis:  redis,
		config: cfg,
		logger: logger,
	}
}

// Limit returns middleware that limits requests per IP
func (rl *RateLimiter) Limit() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getRealIP(r)

			// Check per-minute rate limit
			if !rl.checkRateLimit(r.Context(), ip, "minute", rl.config.RequestsPerMinute, time.Minute) {
				rl.logger.Warn("rate limit exceeded",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
					zap.String("limit_type", "per_minute"),
				)
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			// Check per-hour rate limit
			if !rl.checkRateLimit(r.Context(), ip, "hour", rl.config.RequestsPerHour, time.Hour) {
				rl.logger.Warn("rate limit exceeded",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
					zap.String("limit_type", "per_hour"),
				)
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// LimitAuth returns stricter rate limiting for authentication endpoints
func (rl *RateLimiter) LimitAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getRealIP(r)

			// Strict limit: 5 requests per minute for auth endpoints
			if !rl.checkRateLimit(r.Context(), ip, "auth_minute", 5, time.Minute) {
				rl.logger.Warn("auth rate limit exceeded",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
				)
				http.Error(w, "Too many authentication attempts. Please try again later.", http.StatusTooManyRequests)
				return
			}

			// 20 requests per hour for auth
			if !rl.checkRateLimit(r.Context(), ip, "auth_hour", 20, time.Hour) {
				rl.logger.Warn("auth hourly rate limit exceeded",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
				)
				http.Error(w, "Too many authentication attempts. Please try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// checkRateLimit checks if request is within rate limit
func (rl *RateLimiter) checkRateLimit(ctx context.Context, ip, limitType string, maxRequests int, window time.Duration) bool {
	// If Redis is unavailable, allow request (fail open for availability)
	if rl.redis == nil {
		rl.logger.Warn("rate limiter: redis not available, allowing request")
		return true
	}

	// Create key for this IP and time window
	now := time.Now()
	windowKey := fmt.Sprintf("ratelimit:%s:%s:%d", ip, limitType, now.Unix()/int64(window.Seconds()))

	// Increment counter
	count, err := rl.redis.Incr(ctx, windowKey).Result()
	if err != nil {
		rl.logger.Error("rate limiter: redis error", zap.Error(err))
		// Fail open - allow request if Redis has issues
		return true
	}

	// Set expiry on first request in window
	if count == 1 {
		rl.redis.Expire(ctx, windowKey, window+time.Second)
	}

	return count <= int64(maxRequests)
}

// getRealIP extracts the real client IP address
func getRealIP(r *http.Request) string {
	// Check X-Forwarded-For header (behind proxy/load balancer)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs (client, proxy1, proxy2, ...)
		// Take the first one (client IP)
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header (some proxies use this)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if net.ParseIP(xri) != nil {
			return xri
		}
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

// GetRemainingRequests returns how many requests remain in the current window
func (rl *RateLimiter) GetRemainingRequests(ctx context.Context, ip string, limitType string, maxRequests int, window time.Duration) (int, error) {
	if rl.redis == nil {
		return maxRequests, nil
	}

	now := time.Now()
	windowKey := fmt.Sprintf("ratelimit:%s:%s:%d", ip, limitType, now.Unix()/int64(window.Seconds()))

	count, err := rl.redis.Get(ctx, windowKey).Int()
	if err == redis.Nil {
		return maxRequests, nil
	}
	if err != nil {
		return 0, err
	}

	remaining := maxRequests - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}
