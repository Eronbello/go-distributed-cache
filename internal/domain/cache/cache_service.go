package cache

import (
	"time"
)

// InMemoryCache is a basic in-memory implementation of the Cache interface.
type InMemoryCache struct {
	storage map[string]*CacheItem
}

// NewInMemoryCache creates a new instance of our basic in-memory cache.
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		storage: make(map[string]*CacheItem),
	}
}

// Get retrieves an item from the cache.
func (c *InMemoryCache) Get(key string) (interface{}, error) {
	item, found := c.storage[key]
	if !found {
		return nil, ErrKeyNotFound
	}
	if time.Now().After(item.ExpiresAt) {
		delete(c.storage, key)
		return nil, ErrKeyExpired
	}
	return item.Value, nil
}

// Set adds or updates an item in the cache with a TTL.
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	expiration := time.Now().Add(ttl)
	c.storage[key] = &CacheItem{
		Key:       key,
		Value:     value,
		ExpiresAt: expiration,
	}
	return nil
}

// Delete removes an item from the cache.
func (c *InMemoryCache) Delete(key string) error {
	delete(c.storage, key)
	return nil
}
