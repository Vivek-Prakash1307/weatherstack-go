package metrics

import (
	"sync"
	"time"
)

// MetricsManager handles application metrics
type MetricsManager struct {
	totalRequests     int64
	successRequests   int64
	cacheHits         int64
	cacheMisses       int64
	errors            int64
	responseTimes     []float64
	cityRequestCounts map[string]int64
	startTime         time.Time
	mu                sync.RWMutex
}

// NewMetricsManager creates a new metrics manager
func NewMetricsManager() *MetricsManager {
	return &MetricsManager{
		responseTimes:     make([]float64, 0, 1000),
		cityRequestCounts: make(map[string]int64),
		startTime:         time.Now(),
	}
}

// RecordRequest records a request with its duration and status
func (m *MetricsManager) RecordRequest(duration time.Duration, cacheHit bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests++

	if cacheHit {
		m.cacheHits++
	} else {
		m.cacheMisses++
	}

	if err != nil {
		m.errors++
	} else {
		m.successRequests++
	}

	// Record response time
	durationMs := float64(duration.Milliseconds())
	m.responseTimes = append(m.responseTimes, durationMs)

	// Keep only last 1000 response times
	if len(m.responseTimes) > 1000 {
		m.responseTimes = m.responseTimes[1:]
	}
}

// RecordCityRequest records a request for a specific city
func (m *MetricsManager) RecordCityRequest(city string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cityRequestCounts[city]++
}

// GetMetrics returns all metrics
func (m *MetricsManager) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	uptime := time.Since(m.startTime)
	avgResponseTime := m.calculateAverageResponseTime()
	p95ResponseTime := m.calculatePercentile(95)
	p99ResponseTime := m.calculatePercentile(99)

	// Get top 10 cities
	topCities := m.getTopCities(10)

	return map[string]interface{}{
		"total_requests":        m.totalRequests,
		"success_requests":      m.successRequests,
		"cache_hits":            m.cacheHits,
		"cache_misses":          m.cacheMisses,
		"cache_hit_rate":        m.calculateCacheHitRate(),
		"errors":                m.errors,
		"error_rate":            m.calculateErrorRate(),
		"average_response_ms":   avgResponseTime,
		"p95_response_ms":       p95ResponseTime,
		"p99_response_ms":       p99ResponseTime,
		"uptime_seconds":        uptime.Seconds(),
		"uptime":                uptime.String(),
		"requests_per_minute":   m.calculateRequestsPerMinute(uptime),
		"top_cities":            topCities,
		"total_unique_cities":   len(m.cityRequestCounts),
	}
}

// calculateAverageResponseTime calculates average response time
func (m *MetricsManager) calculateAverageResponseTime() float64 {
	if len(m.responseTimes) == 0 {
		return 0
	}

	var sum float64
	for _, rt := range m.responseTimes {
		sum += rt
	}
	return sum / float64(len(m.responseTimes))
}

// calculatePercentile calculates the nth percentile of response times
func (m *MetricsManager) calculatePercentile(percentile int) float64 {
	if len(m.responseTimes) == 0 {
		return 0
	}

	// Create a copy and sort it
	times := make([]float64, len(m.responseTimes))
	copy(times, m.responseTimes)

	// Simple bubble sort for small arrays
	for i := 0; i < len(times); i++ {
		for j := i + 1; j < len(times); j++ {
			if times[i] > times[j] {
				times[i], times[j] = times[j], times[i]
			}
		}
	}

	index := int(float64(len(times)) * float64(percentile) / 100.0)
	if index >= len(times) {
		index = len(times) - 1
	}
	return times[index]
}

// calculateCacheHitRate calculates cache hit rate percentage
func (m *MetricsManager) calculateCacheHitRate() float64 {
	total := m.cacheHits + m.cacheMisses
	if total == 0 {
		return 0
	}
	return float64(m.cacheHits) / float64(total) * 100
}

// calculateErrorRate calculates error rate percentage
func (m *MetricsManager) calculateErrorRate() float64 {
	if m.totalRequests == 0 {
		return 0
	}
	return float64(m.errors) / float64(m.totalRequests) * 100
}

// calculateRequestsPerMinute calculates requests per minute
func (m *MetricsManager) calculateRequestsPerMinute(uptime time.Duration) float64 {
	minutes := uptime.Minutes()
	if minutes == 0 {
		return 0
	}
	return float64(m.totalRequests) / minutes
}

// getTopCities returns top N most requested cities
func (m *MetricsManager) getTopCities(n int) []map[string]interface{} {
	type cityCount struct {
		city  string
		count int64
	}

	// Convert map to slice
	cities := make([]cityCount, 0, len(m.cityRequestCounts))
	for city, count := range m.cityRequestCounts {
		cities = append(cities, cityCount{city, count})
	}

	// Simple bubble sort by count (descending)
	for i := 0; i < len(cities); i++ {
		for j := i + 1; j < len(cities); j++ {
			if cities[i].count < cities[j].count {
				cities[i], cities[j] = cities[j], cities[i]
			}
		}
	}

	// Take top N
	if len(cities) > n {
		cities = cities[:n]
	}

	// Convert to map slice
	result := make([]map[string]interface{}, len(cities))
	for i, cc := range cities {
		result[i] = map[string]interface{}{
			"city":  cc.city,
			"count": cc.count,
		}
	}

	return result
}

// Reset resets all metrics (useful for testing)
func (m *MetricsManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests = 0
	m.successRequests = 0
	m.cacheHits = 0
	m.cacheMisses = 0
	m.errors = 0
	m.responseTimes = make([]float64, 0, 1000)
	m.cityRequestCounts = make(map[string]int64)
	m.startTime = time.Now()
}