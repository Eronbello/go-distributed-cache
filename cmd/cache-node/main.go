package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eronbello/distributed-cache/internal/application"
	"github.com/eronbello/distributed-cache/internal/domain/cache"
	"github.com/eronbello/distributed-cache/internal/domain/hashing"
	httpInfra "github.com/eronbello/distributed-cache/internal/infrastructure/http"
	"github.com/gorilla/mux"
)

// RemoteHTTPClient implements the RemoteClient interface required by DistributedCacheService.
type RemoteHTTPClient struct {
	BaseURL string
}

// Set satisfies the RemoteClient interface, making an HTTP POST request to store a key/value.
func (r *RemoteHTTPClient) Set(key string, value interface{}, ttl time.Duration) error {
	url := fmt.Sprintf("%s/cache", r.BaseURL)
	// Convert time.Duration to int64 seconds for the JSON payload
	ttlSeconds := int64(ttl.Seconds())
	payload := fmt.Sprintf(`{"key":"%s","value":%q,"ttl_seconds":%d}`, key, value, ttlSeconds)

	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote set failed: %v", resp.Status)
	}
	return nil
}

// Get satisfies the RemoteClient interface, making an HTTP GET request to retrieve a key.
func (r *RemoteHTTPClient) Get(key string) (interface{}, error) {
	url := fmt.Sprintf("%s/cache?key=%s", r.BaseURL, key)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote get failed: %v", resp.Status)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data["value"], nil
}

func main() {
	nodeID := getEnv("NODE_ID", "node1")
	host := getEnv("NODE_HOST", "localhost")
	portStr := getEnv("NODE_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("invalid port: %v", err)
	}

	// Parse cluster nodes from environment
	clusterStr := getEnv("CLUSTER_NODES", "")
	clusterNodes := strings.Split(clusterStr, ",") // e.g. "node1:8080,node2:8080"

	// Initialize the consistent hash ring
	hashRing := hashing.NewHashRing(3)

	// Register nodes (including this local one) in the hash ring
	remoteClients := make(map[string]application.RemoteClient)

	for _, n := range clusterNodes {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}

		parts := strings.Split(n, ":")
		if len(parts) != 2 {
			continue
		}

		nID := parts[0]
		nPort := parts[1]

		// Add node to the ring
		hashRing.AddNode(hashing.Node{
			ID:   nID,
			Host: nID, // simplistic assumption that Host = node ID
			Port: mustAtoi(nPort),
		})

		// If it's not the local node, we set up a RemoteHTTPClient
		if nID != nodeID {
			remoteClients[nID] = &RemoteHTTPClient{
				BaseURL: fmt.Sprintf("http://%s:%s", nID, nPort),
			}
		}
	}

	// Also ensure the local node is added (if not in CLUSTER_NODES)
	hashRing.AddNode(hashing.Node{
		ID:   nodeID,
		Host: host,
		Port: port,
	})

	// Prepare local cache
	localCache := cache.NewInMemoryCache()

	// DistributedCacheService handles routing logic
	distService := &application.DistributedCacheService{
		HashRing:    hashRing,
		LocalNodeID: nodeID,
		LocalCache:  localCache,
		RemoteCalls: remoteClients,
	}

	// Set up HTTP handlers
	cacheHandler := &httpInfra.CacheHandler{Service: distService}
	router := mux.NewRouter()
	router.HandleFunc("/cache", cacheHandler.Set).Methods("POST")
	router.HandleFunc("/cache", cacheHandler.Get).Methods("GET")

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("Starting cache node %s on %s", nodeID, addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

// getEnv reads environment variables with a fallback default.
func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// mustAtoi is a helper to parse integers, but will fatal if parsing fails.
func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("mustAtoi failed for %q: %v", s, err)
	}
	return i
}
