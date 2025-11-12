package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	"go.uber.org/zap"
)

// Cache wraps Redis client with metrics and logging
type Cache struct {
	client *redis.Client
	logger *logging.Logger
}

// New creates a new cache instance
func New(client *redis.Client, logger *logging.Logger) *Cache {
	return &Cache{
		client: client,
		logger: logger,
	}
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		hit := ctx.Err() == nil
		metrics.RecordCacheOperation("get", "redis", extractKeyPrefix(key), duration, hit)
	}()

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		c.logger.Debug("cache miss", zap.String("key", key))
		return fmt.Errorf("cache miss")
	}
	if err != nil {
		c.logger.Error("cache get error", zap.Error(err), zap.String("key", key))
		return err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		c.logger.Error("cache unmarshal error", zap.Error(err), zap.String("key", key))
		return err
	}

	c.logger.Debug("cache hit", zap.String("key", key))
	return nil
}

// Set stores a value in cache
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordCacheOperation("set", "redis", extractKeyPrefix(key), duration, false)
	}()

	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("cache marshal error", zap.Error(err), zap.String("key", key))
		return err
	}

	err = c.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		c.logger.Error("cache set error", zap.Error(err), zap.String("key", key))
		return err
	}

	c.logger.Debug("cache set", zap.String("key", key), zap.Duration("ttl", ttl))
	return nil
}

// Delete removes a key from cache
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if len(keys) > 0 {
			metrics.RecordCacheOperation("delete", "redis", extractKeyPrefix(keys[0]), duration, false)
		}
	}()

	err := c.client.Del(ctx, keys...).Err()
	if err != nil {
		c.logger.Error("cache delete error", zap.Error(err), zap.Strings("keys", keys))
		return err
	}

	c.logger.Debug("cache delete", zap.Strings("keys", keys))
	return nil
}

// InvalidatePattern invalidates all keys matching a pattern
func (c *Cache) InvalidatePattern(ctx context.Context, pattern string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordCacheOperation("invalidate", "redis", extractKeyPrefix(pattern), duration, false)
	}()

	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		c.logger.Error("cache scan error", zap.Error(err), zap.String("pattern", pattern))
		return err
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			c.logger.Error("cache delete pattern error", zap.Error(err), zap.String("pattern", pattern))
			return err
		}
		c.logger.Info("cache pattern invalidated",
			zap.String("pattern", pattern),
			zap.Int("keys_deleted", len(keys)),
		)
	}

	return nil
}

// Exists checks if a key exists in cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SetNX sets a value only if the key doesn't exist (distributed lock)
func (c *Cache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	set, err := c.client.SetNX(ctx, key, data, ttl).Result()
	if err != nil {
		c.logger.Error("cache setnx error", zap.Error(err), zap.String("key", key))
		return false, err
	}

	return set, nil
}

// GetSet gets the old value and sets a new value atomically
func (c *Cache) GetSet(ctx context.Context, key string, newValue interface{}, dest interface{}) error {
	data, err := json.Marshal(newValue)
	if err != nil {
		return err
	}

	oldVal, err := c.client.GetSet(ctx, key, data).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss")
	}
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(oldVal), dest)
}

// Increment atomically increments a counter
func (c *Cache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	return c.client.IncrBy(ctx, key, delta).Result()
}

// Expire sets an expiration on an existing key
func (c *Cache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl).Err()
}

// TTL returns the remaining time to live of a key
func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// extractKeyPrefix extracts the prefix from a cache key for metrics
func extractKeyPrefix(key string) string {
	// Extract first part before ":"
	for i, c := range key {
		if c == ':' {
			return key[:i]
		}
	}
	return key
}

// CacheKey generates a cache key for common patterns
type CacheKey struct{}

// Product generates a product cache key
func (CacheKey) Product(orgID, productID string) string {
	return fmt.Sprintf("product:%s:%s", orgID, productID)
}

// ProductList generates a products list cache key
func (CacheKey) ProductList(orgID string, page, pageSize int) string {
	return fmt.Sprintf("products:%s:page:%d:size:%d", orgID, page, pageSize)
}

// Customer generates a customer cache key
func (CacheKey) Customer(orgID, customerID string) string {
	return fmt.Sprintf("customer:%s:%s", orgID, customerID)
}

// Sale generates a sale cache key
func (CacheKey) Sale(orgID, saleID string) string {
	return fmt.Sprintf("sale:%s:%s", orgID, saleID)
}

// Dashboard generates a dashboard cache key
func (CacheKey) Dashboard(orgID, period string) string {
	return fmt.Sprintf("dashboard:%s:%s", orgID, period)
}

// Report generates a report cache key
func (CacheKey) Report(orgID, reportType, startDate, endDate string) string {
	return fmt.Sprintf("report:%s:%s:%s:%s", orgID, reportType, startDate, endDate)
}

// OrganizationPattern generates a pattern to invalidate all org-specific caches
func (CacheKey) OrganizationPattern(orgID string) string {
	return fmt.Sprintf("*:%s:*", orgID)
}
