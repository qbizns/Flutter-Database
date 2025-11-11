package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	// RequestsPerWindow is the maximum number of requests allowed in the time window
	RequestsPerWindow int

	// Window is the time window for rate limiting
	Window time.Duration

	// KeyFunc extracts the key to rate limit by (e.g., IP address, user ID)
	// Default: uses IP address
	KeyFunc func(r *http.Request) string

	// OnLimitReached is called when rate limit is exceeded
	// Default: returns 429 Too Many Requests
	OnLimitReached func(w http.ResponseWriter, r *http.Request)

	// SkipPaths is a list of paths to skip rate limiting
	SkipPaths []string

	// AllowBurst allows burst requests up to this number above the limit
	AllowBurst int

	// UseXRateLimitHeaders adds X-RateLimit-* headers to responses
	UseXRateLimitHeaders bool
}

// DefaultRateLimitConfig returns a default rate limit configuration
// Default: 100 requests per minute per IP
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerWindow:    100,
		Window:               1 * time.Minute,
		KeyFunc:              defaultKeyFunc,
		OnLimitReached:       defaultOnLimitReached,
		SkipPaths:            []string{"/health", "/metrics"},
		AllowBurst:           10,
		UseXRateLimitHeaders: true,
	}
}

// StrictAuthRateLimitConfig returns a strict rate limit config for auth endpoints
// 5 requests per minute per IP to prevent brute force
func StrictAuthRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerWindow:    5,
		Window:               1 * time.Minute,
		KeyFunc:              defaultKeyFunc,
		OnLimitReached:       defaultOnLimitReached,
		SkipPaths:            []string{},
		AllowBurst:           2,
		UseXRateLimitHeaders: true,
	}
}

// bucket represents a token bucket for rate limiting
type bucket struct {
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	config  RateLimitConfig
	buckets sync.Map // map[string]*bucket
}

// NewRateLimiter creates a new rate limiter middleware
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	// Set defaults
	if config.RequestsPerWindow == 0 {
		config.RequestsPerWindow = 100
	}
	if config.Window == 0 {
		config.Window = 1 * time.Minute
	}
	if config.KeyFunc == nil {
		config.KeyFunc = defaultKeyFunc
	}
	if config.OnLimitReached == nil {
		config.OnLimitReached = defaultOnLimitReached
	}

	rl := &RateLimiter{
		config: config,
	}

	// Start cleanup goroutine
	go rl.cleanupExpiredBuckets()

	return rl
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip rate limiting for configured paths
		if rl.shouldSkipPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get the rate limit key (e.g., IP address)
		key := rl.config.KeyFunc(r)
		if key == "" {
			// If no key, allow request (fail open)
			next.ServeHTTP(w, r)
			return
		}

		// Check rate limit
		allowed, remaining, resetTime := rl.allow(key)

		// Add rate limit headers if configured
		if rl.config.UseXRateLimitHeaders {
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.RequestsPerWindow))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
		}

		if !allowed {
			rl.config.OnLimitReached(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// allow checks if a request is allowed for the given key
func (rl *RateLimiter) allow(key string) (allowed bool, remaining int, resetTime time.Time) {
	now := time.Now()

	// Get or create bucket
	bucketInterface, _ := rl.buckets.LoadOrStore(key, &bucket{
		tokens:     rl.config.RequestsPerWindow + rl.config.AllowBurst,
		lastRefill: now,
	})

	b := bucketInterface.(*bucket)
	b.mu.Lock()
	defer b.mu.Unlock()

	// Calculate tokens to add based on time elapsed
	elapsed := now.Sub(b.lastRefill)
	if elapsed >= rl.config.Window {
		// Full refill
		b.tokens = rl.config.RequestsPerWindow + rl.config.AllowBurst
		b.lastRefill = now
	} else {
		// Partial refill based on elapsed time
		tokensToAdd := int(float64(rl.config.RequestsPerWindow) * (float64(elapsed) / float64(rl.config.Window)))
		b.tokens = min(b.tokens+tokensToAdd, rl.config.RequestsPerWindow+rl.config.AllowBurst)
		if tokensToAdd > 0 {
			b.lastRefill = now
		}
	}

	// Check if request is allowed
	if b.tokens > 0 {
		b.tokens--
		remaining = b.tokens
		resetTime = b.lastRefill.Add(rl.config.Window)
		return true, remaining, resetTime
	}

	// Rate limit exceeded
	remaining = 0
	resetTime = b.lastRefill.Add(rl.config.Window)
	return false, remaining, resetTime
}

// shouldSkipPath checks if a path should skip rate limiting
func (rl *RateLimiter) shouldSkipPath(path string) bool {
	for _, skipPath := range rl.config.SkipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// cleanupExpiredBuckets periodically removes old buckets
func (rl *RateLimiter) cleanupExpiredBuckets() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		rl.buckets.Range(func(key, value interface{}) bool {
			b := value.(*bucket)
			b.mu.Lock()
			// Remove buckets that haven't been used in 2x the window duration
			if now.Sub(b.lastRefill) > 2*rl.config.Window {
				rl.buckets.Delete(key)
			}
			b.mu.Unlock()
			return true
		})
	}
}

// defaultKeyFunc extracts IP address from request
func defaultKeyFunc(r *http.Request) string {
	// Try to get real IP from X-Forwarded-For or X-Real-IP headers
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, use the first one
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	// RemoteAddr format: "IP:port"
	parts := strings.Split(r.RemoteAddr, ":")
	if len(parts) > 0 {
		return parts[0]
	}

	return r.RemoteAddr
}

// defaultOnLimitReached is the default handler when rate limit is exceeded
func defaultOnLimitReached(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// KeyByIPAndUser creates a rate limit key combining IP and user ID
// Useful for authenticated endpoints
func KeyByIPAndUser(r *http.Request) string {
	ip := defaultKeyFunc(r)
	// Try to get user ID from context (assuming it's set by auth middleware)
	userID := r.Context().Value("user_id")
	if userID != nil {
		return fmt.Sprintf("%s:%v", ip, userID)
	}
	return ip
}

// KeyByUser creates a rate limit key based on user ID only
// Useful for per-user rate limiting
func KeyByUser(r *http.Request) string {
	userID := r.Context().Value("user_id")
	if userID != nil {
		return fmt.Sprintf("user:%v", userID)
	}
	// Fall back to IP if no user ID
	return defaultKeyFunc(r)
}
