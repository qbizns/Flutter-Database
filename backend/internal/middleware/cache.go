package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/your-org/pos-backend/internal/cache"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
)

// CacheMiddleware provides HTTP response caching
type CacheMiddleware struct {
	cache *cache.Cache
}

// NewCacheMiddleware creates a new cache middleware
func NewCacheMiddleware(c *cache.Cache) *CacheMiddleware {
	return &CacheMiddleware{cache: c}
}

// CachedResponse represents a cached HTTP response
type CachedResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
}

// CacheGET caches GET requests
func (cm *CacheMiddleware) CacheGET(ttl time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only cache GET requests
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()

			// Get organization ID from context
			orgID, ok := appctx.GetOrganizationID(ctx)
			if !ok {
				// No org context, don't cache
				next.ServeHTTP(w, r)
				return
			}

			// Create cache key from URL + org ID + query params
			cacheKey := generateCacheKey(orgID.String(), r)

			// Try to get from cache
			var cachedResponse CachedResponse
			if err := cm.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
				// Cache hit
				w.Header().Set("X-Cache", "HIT")
				w.Header().Set("X-Cache-Key", cacheKey)

				// Restore headers
				for k, v := range cachedResponse.Headers {
					w.Header().Set(k, v)
				}

				// Write response
				w.WriteHeader(cachedResponse.StatusCode)
				w.Write(cachedResponse.Body)
				return
			}

			// Cache miss - capture response
			w.Header().Set("X-Cache", "MISS")
			w.Header().Set("X-Cache-Key", cacheKey)

			rec := &responseCaptureWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           &bytes.Buffer{},
				headers:        make(map[string]string),
			}

			next.ServeHTTP(rec, r)

			// Cache successful responses only (2xx status codes)
			if rec.statusCode >= 200 && rec.statusCode < 300 {
				cachedResp := CachedResponse{
					StatusCode: rec.statusCode,
					Headers:    rec.headers,
					Body:       rec.body.Bytes(),
				}

				// Store in cache (ignore errors to not block request)
				_ = cm.cache.Set(ctx, cacheKey, cachedResp, ttl)
			}
		})
	}
}

// responseCaptureWriter captures the response for caching
type responseCaptureWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
	headers    map[string]string
}

func (rcw *responseCaptureWriter) WriteHeader(code int) {
	rcw.statusCode = code

	// Capture important headers
	for _, h := range []string{"Content-Type", "Content-Encoding", "ETag", "Last-Modified"} {
		if v := rcw.ResponseWriter.Header().Get(h); v != "" {
			rcw.headers[h] = v
		}
	}

	rcw.ResponseWriter.WriteHeader(code)
}

func (rcw *responseCaptureWriter) Write(b []byte) (int, error) {
	// Write to both response and buffer
	rcw.body.Write(b)
	return rcw.ResponseWriter.Write(b)
}

// generateCacheKey creates a cache key from request
func generateCacheKey(orgID string, r *http.Request) string {
	// Include: org, path, query params
	h := sha256.New()
	io.WriteString(h, orgID)
	io.WriteString(h, r.URL.Path)
	io.WriteString(h, r.URL.RawQuery)

	// Include relevant headers that affect response
	if accept := r.Header.Get("Accept"); accept != "" {
		io.WriteString(h, accept)
	}
	if acceptLang := r.Header.Get("Accept-Language"); acceptLang != "" {
		io.WriteString(h, acceptLang)
	}

	return fmt.Sprintf("http:%s:%x", orgID, h.Sum(nil))
}

// InvalidateOrgCache invalidates all cached responses for an organization
func (cm *CacheMiddleware) InvalidateOrgCache(ctx context.Context, orgID string) error {
	pattern := fmt.Sprintf("http:%s:*", orgID)
	return cm.cache.InvalidatePattern(ctx, pattern)
}

// InvalidateResourceCache invalidates cache for a specific resource
func (cm *CacheMiddleware) InvalidateResourceCache(ctx context.Context, orgID, resourceType, resourceID string) error {
	// This is a simple approach - in production you might want more sophisticated cache invalidation
	return cm.InvalidateOrgCache(ctx, orgID)
}
