package unit

import (
	"testing"
	"time"
	"github.com/Vivek-Prakash1307/weather-Microservices/api/openweathermap"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/cache"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/metrics"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/services"
)

func TestWeatherService_GetWeatherData(t *testing.T) {
	// Setup
	apiKey := "test_api_key"
	cacheManager := cache.NewCacheManager(10 * time.Minute)
	metricsManager := metrics.NewMetricsManager()
	weatherClient := openweathermap.NewClient(apiKey)
	service := services.NewWeatherService(weatherClient, cacheManager, metricsManager)

	tests := []struct {
		name    string
		city    string
		wantErr bool
	}{
		{
			name:    "Empty city name",
			city:    "",
			wantErr: true,
		},
		{
			name:    "Valid city name",
			city:    "London",
			wantErr: false,
		},
		{
			name:    "City with spaces",
			city:    "New York",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetWeatherData(tt.city)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWeatherData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWeatherService_CacheHit(t *testing.T) {
	// Setup
	apiKey := "test_api_key"
	cacheManager := cache.NewCacheManager(10 * time.Minute)
	metricsManager := metrics.NewMetricsManager()
	weatherClient := openweathermap.NewClient(apiKey)
	service := services.NewWeatherService(weatherClient, cacheManager, metricsManager)

	city := "TestCity"

	// First call should be a cache miss
	data1, err1 := service.GetWeatherData(city)
	if err1 == nil && !data1.CacheHit {
		t.Log("First call - cache miss as expected")
	}

	// Second call should be a cache hit (if first was successful)
	if err1 == nil {
		data2, err2 := service.GetWeatherData(city)
		if err2 != nil {
			t.Errorf("Second call failed: %v", err2)
		}
		if !data2.CacheHit {
			t.Error("Expected cache hit on second call")
		}
	}
}

func BenchmarkWeatherService_GetWeatherData(b *testing.B) {
	apiKey := "test_api_key"
	cacheManager := cache.NewCacheManager(10 * time.Minute)
	metricsManager := metrics.NewMetricsManager()
	weatherClient := openweathermap.NewClient(apiKey)
	service := services.NewWeatherService(weatherClient, cacheManager, metricsManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetWeatherData("London")
	}
}