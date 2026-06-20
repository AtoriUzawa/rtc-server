// Package redis provides a simple wrapper for the Redis client.
//
// Design goals:
//   - Wrap go-redis basic operations (Set / Get / Del etc.)
//   - Provide unified JSON serialization
//   - Avoid business layer depending directly on low-level redis API
//
// Usage:
//  1. Call NewClient in main to initialize
//  2. Inject Client into various services for use
//
// Note:
//   - This package does not handle key design (e.g., prefixes, naming conventions)
//   - Does not include caching strategies (e.g., breakdown prevention, expiration policies)
//   - Recommended to wrap further at the business layer
package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client is a wrapped Redis client.
//
// Internally holds a go-redis Client instance and provides common operations.
type Client struct {
	rdb *redis.Client
}

// NewClient creates a new Redis client.
//
// Parameters:
//   - addr: Redis address (e.g., "localhost:6379")
//   - password: Redis password (pass empty string if no password)
//   - db: database number (default 0)
//
// Returns:
//   - *Client wrapped client instance
//
// Note:
//   - This method only creates a client, does not guarantee connectivity
//   - Recommended to call the Ping method to check connectivity
func NewClient(addr, password string, db int) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Client{
		rdb: rdb,
	}
}

// Ping checks if the Redis connection is available.
//
// Parameters:
//   - ctx: context (recommended to pass a context with timeout)
//
// Returns:
//   - error: returns error if connection fails
//
// Use cases:
//   - Health check on application startup
//   - Periodic Redis status check
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Set sets a cache value (string or basic type).
//
// Parameters:
//   - ctx: context
//   - key: key
//   - val: value (will be automatically converted to string)
//   - exp: expiration duration (0 means no expiration)
//
// Returns:
//   - error
//
// Note:
//   - For complex structs, use SetJSON instead
func (c *Client) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	return c.rdb.Set(ctx, key, val, exp).Err()
}

// Get retrieves a cache value (string).
//
// Parameters:
//   - ctx: context
//   - key: key
//
// Returns:
//   - string: value
//   - error: returns redis.Nil if key does not exist
//
// Note:
//   - Caller should check for redis.Nil (indicates cache miss)
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// Del deletes one or more keys.
//
// Parameters:
//   - ctx: context
//   - key: one or more keys
//
// Returns:
//   - error
//
// Use cases:
//   - Delete cache
//   - Invalidate specific data
func (c *Client) Del(ctx context.Context, key ...string) error {
	return c.rdb.Del(ctx, key...).Err()
}

// SetJSON sets a cache value in JSON format.
//
// Parameters:
//   - ctx: context
//   - key: key
//   - val: any serializable object
//   - exp: expiration duration
//
// Returns:
//   - error
//
// Note:
//   - val must be serializable by json.Marshal
//   - Commonly used for caching struct data
func (c *Client) SetJSON(ctx context.Context, key string, val any, exp time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, exp).Err()
}

// GetJSON retrieves JSON and deserializes into the target object.
//
// Parameters:
//   - ctx: context
//   - key: key
//   - dest: target object (must be a pointer)
//
// Returns:
//   - error
//
// Note:
//   - Returns redis.Nil if key does not exist
//   - dest must be a pointer type, otherwise deserialization fails
func (c *Client) GetJSON(ctx context.Context, key string, dest any) error {
	str, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(str), dest)
}
