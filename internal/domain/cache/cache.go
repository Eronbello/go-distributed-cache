package cache

import (
	"time"
)

// CacheItem represents a single item in the cache with an optional expiration.
type CacheItem struct {
	Key       string
	Value     interface{}
	ExpiresAt time.Time
}

// Cache is the interface that any cache implementation must fulfill.
type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
}
