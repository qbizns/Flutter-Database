package cache

import (
	"context"
	"time"
)

// Cache provides a simple in-memory cache interface
// This is a stub implementation - can be enhanced with Redis/Memcached
type Cache struct {
	// Add cache implementation here
}

// New creates a new cache instance
func New() *Cache {
	return &Cache{}
}

// Get retrieves a value from the cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	// Stub implementation - always returns cache miss
	return false, nil
}

// Set stores a value in the cache
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Stub implementation - does nothing
	return nil
}

// Delete removes a value from the cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	// Stub implementation - does nothing
	return nil
}

// Exists checks if a key exists in the cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	// Stub implementation - always returns false
	return false, nil
}

// InvalidatePattern invalidates all keys matching a pattern
func (c *Cache) InvalidatePattern(ctx context.Context, pattern string) error {
	// Stub implementation - does nothing
	return nil
}

// ErrCacheMiss is returned when a key is not found
var ErrCacheMiss = &CacheError{message: "cache miss"}

// CacheError represents a cache error
type CacheError struct {
	message string
}

func (e *CacheError) Error() string {
	return e.message
}
