package cache

import (
	"sync"
	"time"
)

// CacheManager is the global cache manager instance
var CacheManager = NewCacheManager(time.Duration(60) * time.Minute)

// Cache is the cache manager struct
type Cache struct {
	caches          sync.Map
	cleanupInterval time.Duration
}

// CacheData is the struct that holds the cache data
type CacheData struct {
	Data      []byte
	Endpoint  string
	CreatedAt time.Time
	ExpiredAt time.Time
}

// NewCacheManager creates a new cache manager with the given cleanup interval
func NewCacheManager(cleanupInterval time.Duration) *Cache {
	cm := &Cache{
		cleanupInterval: cleanupInterval,
	}
	go cm.startCleanup()
	return cm
}

// Set sets the cache data for the given endpoint
func (cm *Cache) Set(endpoint string, data CacheData) {
	cm.caches.Store(endpoint, data)
}

// Get gets the cache data for the given endpoint
func (cm *Cache) Get(endpoint string) (CacheData, bool) {
	data, ok := cm.caches.Load(endpoint)
	if !ok {
		return CacheData{}, false
	}
	return data.(CacheData), true
}

// DeleteAll deletes the cache data for the given endpoint
func (cm *Cache) DeleteAll() {
	cm.caches.Range(func(key, value interface{}) bool {
		cm.caches.Delete(key)
		return true
	})
}

// Cleanup deletes the expired cache data
func (cm *Cache) Cleanup() {
	cm.caches.Range(func(key, value interface{}) bool {
		cacheData := value.(CacheData)
		if time.Now().After(cacheData.ExpiredAt) {
			cm.caches.Delete(key)
		}
		return true
	})
}

// startCleanup starts the cleanup process in the background
func (cm *Cache) startCleanup() {
	for {
		time.Sleep(cm.cleanupInterval)
		cm.Cleanup()
	}
}
