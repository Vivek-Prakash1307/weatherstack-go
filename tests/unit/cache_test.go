package unit

import (
	"testing"
	"time"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/cache"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/models"
)

func TestCacheManager_SetAndGet(t *testing.T) {
	cm := cache.NewCacheManager(10 * time.Minute)

	// Create test data
	testData := models.WeatherData{
		Name:    "TestCity",
		Country: "TC",
	}

	// Set data
	cm.Set("testcity", testData)

	// Get data
	retrieved, found := cm.Get("testcity")
	if !found {
		t.Error("Expected to find cached data")
	}

	if retrieved.Name != testData.Name {
		t.Errorf("Expected name %s, got %s", testData.Name, retrieved.Name)
	}
}

func TestCacheManager_Expiry(t *testing.T) {
	// Use short expiry for testing
	cm := cache.NewCacheManager(100 * time.Millisecond)

	testData := models.WeatherData{
		Name: "TestCity",
	}

	cm.Set("testcity", testData)

	// Should be found immediately
	_, found := cm.Get("testcity")
	if !found {
		t.Error("Expected to find cached data immediately")
	}

	// Wait for expiry
	time.Sleep(150 * time.Millisecond)

	// Should not be found after expiry
	_, found = cm.Get("testcity")
	if found {
		t.Error("Expected data to be expired")
	}
}

func TestCacheManager_Clear(t *testing.T) {
	cm := cache.NewCacheManager(10 * time.Minute)

	testData := models.WeatherData{
		Name: "TestCity",
	}

	cm.Set("testcity", testData)

	// Verify data is there
	_, found := cm.Get("testcity")
	if !found {
		t.Error("Expected to find cached data before clear")
	}

	// Clear cache
	cm.Clear()

	// Verify data is gone
	_, found = cm.Get("testcity")
	if found {
		t.Error("Expected cache to be empty after clear")
	}
}

func TestCacheManager_Stats(t *testing.T) {
	cm := cache.NewCacheManager(10 * time.Minute)

	testData := models.WeatherData{
		Name: "TestCity",
	}

	cm.Set("testcity", testData)

	stats := cm.GetStats()

	if stats["total_entries"].(int) != 1 {
		t.Errorf("Expected 1 entry, got %v", stats["total_entries"])
	}

	// Test cache hit/miss
	cm.Get("testcity") // hit
	cm.Get("nonexistent") // miss

	stats = cm.GetStats()
	hitRate := stats["hit_rate"].(float64)

	if hitRate != 50.0 {
		t.Errorf("Expected hit rate 50%%, got %.2f%%", hitRate)
	}
}

func BenchmarkCacheManager_Set(b *testing.B) {
	cm := cache.NewCacheManager(10 * time.Minute)
	testData := models.WeatherData{Name: "TestCity"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Set("testcity", testData)
	}
}

func BenchmarkCacheManager_Get(b *testing.B) {
	cm := cache.NewCacheManager(10 * time.Minute)
	testData := models.WeatherData{Name: "TestCity"}
	cm.Set("testcity", testData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Get("testcity")
	}
}