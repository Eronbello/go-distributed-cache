package application

import (
	"time"

	dCache "github.com/eronbello/distributed-cache/internal/domain/cache"
	"github.com/eronbello/distributed-cache/internal/domain/hashing"
)

// DistributedCacheService manages operations that might involve multiple nodes.
type DistributedCacheService struct {
	HashRing    *hashing.HashRing
	LocalNodeID string
	LocalCache  dCache.Cache
	RemoteCalls map[string]RemoteClient
}

// RemoteClient represents a minimal contract for calling remote nodes.
type RemoteClient interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
}

// Set routes a Set request to the node responsible for a given key.
func (d *DistributedCacheService) Set(key string, value interface{}, ttl time.Duration) error {
	node := d.HashRing.GetNode(key)
	if node.ID == d.LocalNodeID {
		// Store in local cache
		return d.LocalCache.Set(key, value, ttl)
	}
	// Otherwise, forward the request to the remote node
	client := d.RemoteCalls[node.ID]
	if client == nil {
		return d.LocalCache.Set(key, value, ttl) // fallback local if remote is unknown
	}
	return client.Set(key, value, ttl)
}

// Get routes a Get request to the node responsible for the given key.
func (d *DistributedCacheService) Get(key string) (interface{}, error) {
	node := d.HashRing.GetNode(key)
	if node.ID == d.LocalNodeID {
		return d.LocalCache.Get(key)
	}
	client := d.RemoteCalls[node.ID]
	if client == nil {
		return d.LocalCache.Get(key) // fallback local if remote is unknown
	}
	return client.Get(key)
}
