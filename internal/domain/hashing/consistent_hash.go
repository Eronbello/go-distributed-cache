package hashing

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Node represents a node in the cluster.
type Node struct {
	ID   string
	Host string
	Port int
}

// HashRing implements consistent hashing for distributing keys across nodes.
type HashRing struct {
	replicas int
	keys     []int
	hashMap  map[int]Node
}

// NewHashRing creates a hash ring with the specified number of virtual replicas.
func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		replicas: replicas,
		hashMap:  make(map[int]Node),
	}
}

// AddNode adds a new node to the hash ring, creating replicas for better distribution.
func (h *HashRing) AddNode(node Node) {
	for i := 0; i < h.replicas; i++ {
		combinedID := node.ID + strconv.Itoa(i)
		hashVal := int(crc32.ChecksumIEEE([]byte(combinedID)))
		h.keys = append(h.keys, hashVal)
		h.hashMap[hashVal] = node
	}
	sort.Ints(h.keys)
}

// GetNode finds the appropriate node for a given key using consistent hashing.
func (h *HashRing) GetNode(key string) Node {
	if len(h.keys) == 0 {
		return Node{}
	}

	hashVal := int(crc32.ChecksumIEEE([]byte(key)))
	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hashVal
	})
	if idx == len(h.keys) {
		idx = 0
	}
	return h.hashMap[h.keys[idx]]
}
