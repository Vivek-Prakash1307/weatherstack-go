package cache

import (
	"log"
	"sync"
	"time"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/models"
)

// CacheManager handles caching of weather data
type CacheManager struct {
	data      map[string]models.WeatherData
	expiry    map[string]time.Time
	mu        sync.RWMutex
	cacheTime time.Duration
	hitCount  int64
	missCount int64
}

// NewCacheManager creates a new cache manager
func NewCacheManager(cacheTime time.Duration) *CacheManager {
	cm := &CacheManager{
		data:      make(map[string]models.WeatherData),
		expiry:    make(map[string]time.Time),
		cacheTime: cacheTime,
	}

	// Start cleanup goroutine
	go cm.cleanupExpired()

	log.Printf("✅ Cache initialized with %v expiry time", cacheTime)
	return cm
}

// Get retrieves data from cache
func (cm *CacheManager) Get(key string) (models.WeatherData, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if data, found := cm.data[key]; found && time.Now().Before(cm.expiry[key]) {
		cm.hitCount++
		return data, true
	}

	cm.missCount++
	return models.WeatherData{}, false
}

// Set stores data in cache
func (cm *CacheManager) Set(key string, data models.WeatherData) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.data[key] = data
	cm.expiry[key] = time.Now().Add(cm.cacheTime)
}

// Clear removes all cached data
func (cm *CacheManager) Clear() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.data = make(map[string]models.WeatherData)
	cm.expiry = make(map[string]time.Time)
	log.Println("🗑️  Cache cleared")
}

// GetStats returns cache statistics
func (cm *CacheManager) GetStats() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	entries := make(map[string]string)
	validEntries := 0
	now := time.Now()

	for key, expiry := range cm.expiry {
		if now.Before(expiry) {
			entries[key] = expiry.Format("2006-01-02 15:04:05")
			validEntries++
		}
	}

	return map[string]interface{}{
		"total_entries":  validEntries,
		"hit_count":      cm.hitCount,
		"miss_count":     cm.missCount,
		"hit_rate":       cm.calculateHitRate(),
		"cache_duration": cm.cacheTime.String(),
		"entries":        entries,
	}
}

// calculateHitRate calculates the cache hit rate percentage
func (cm *CacheManager) calculateHitRate() float64 {
	total := cm.hitCount + cm.missCount
	if total == 0 {
		return 0
	}
	return float64(cm.hitCount) / float64(total) * 100
}

// cleanupExpired removes expired entries periodically
func (cm *CacheManager) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cm.mu.Lock()
		now := time.Now()
		cleaned := 0

		for key, expiry := range cm.expiry {
			if now.After(expiry) {
				delete(cm.data, key)
				delete(cm.expiry, key)
				cleaned++
			}
		}

		if cleaned > 0 {
			log.Printf("🧹 Cleaned %d expired cache entries", cleaned)
		}
		cm.mu.Unlock()
	}
}

// GetSize returns the number of cached entries
func (cm *CacheManager) GetSize() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	validEntries := 0
	now := time.Now()

	for _, expiry := range cm.expiry {
		if now.Before(expiry) {
			validEntries++
		}
	}

	return validEntries
}