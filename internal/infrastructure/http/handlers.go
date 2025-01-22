package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eronbello/distributed-cache/internal/application"
	dCache "github.com/eronbello/distributed-cache/internal/domain/cache"
)

type CacheHandler struct {
	Service *application.DistributedCacheService
}

type setRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl_seconds"`
}

// Set handles POST requests to store a key/value in the cache.
func (h *CacheHandler) Set(w http.ResponseWriter, r *http.Request) {
	var req setRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.Service.Set(req.Key, req.Value, time.Duration(req.TTL)*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Get handles GET requests to retrieve a value from the cache.
func (h *CacheHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key parameter", http.StatusBadRequest)
		return
	}

	value, err := h.Service.Get(key)
	if err != nil {
		switch err {
		case dCache.ErrKeyNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case dCache.ErrKeyExpired:
			http.Error(w, err.Error(), http.StatusGone)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp := map[string]interface{}{
		"key":   key,
		"value": value,
	}
	json.NewEncoder(w).Encode(resp)
}
