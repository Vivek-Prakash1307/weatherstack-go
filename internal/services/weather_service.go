package services

import (
	"fmt"
	"log"
	"strings"
	"time"
	"github.com/Vivek-Prakash1307/weather-Microservices/api/openweathermap"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/cache"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/metrics"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/models"
	"github.com/Vivek-Prakash1307/weather-Microservices/pkg/utils"
)

// WeatherService handles weather-related business logic
type WeatherService struct {
	weatherClient  *openweathermap.Client
	cacheManager   *cache.CacheManager
	metricsManager *metrics.MetricsManager
}

// NewWeatherService creates a new weather service
func NewWeatherService(
	weatherClient *openweathermap.Client,
	cacheManager *cache.CacheManager,
	metricsManager *metrics.MetricsManager,
) *WeatherService {
	return &WeatherService{
		weatherClient:  weatherClient,
		cacheManager:   cacheManager,
		metricsManager: metricsManager,
	}
}

// GetWeatherData fetches weather data for a city
func (ws *WeatherService) GetWeatherData(city string) (*models.WeatherData, error) {
	startTime := time.Now()
	city = strings.ToLower(strings.TrimSpace(city))

	if city == "" {
		return nil, fmt.Errorf("city name cannot be empty")
	}

	// Record city request
	ws.metricsManager.RecordCityRequest(city)

	// Check cache first
	if cachedData, found := ws.cacheManager.Get(city); found {
		cachedData.CacheHit = true
		duration := time.Since(startTime)
		ws.metricsManager.RecordRequest(duration, true, nil)
		log.Printf("✅ Cache hit for city: %s (took %v)", city, duration)
		return &cachedData, nil
	}

	log.Printf("🔄 Cache miss for city: %s, fetching from API...", city)

	// Fetch from API
	apiResponse, err := ws.weatherClient.GetWeather(city)
	if err != nil {
		duration := time.Since(startTime)
		ws.metricsManager.RecordRequest(duration, false, err)
		return nil, err
	}

	// Transform API response to our weather data model
	weatherData := ws.transformWeatherData(apiResponse)

	// Fetch UV Index and Air Quality in parallel
	uvChan := make(chan float64, 1)
	aqiChan := make(chan struct {
		aqi     int
		quality string
	}, 1)

	go func() {
		uv, err := ws.weatherClient.GetUVIndex(apiResponse.Coord.Lat, apiResponse.Coord.Lon)
		if err != nil {
			log.Printf("⚠️  Failed to get UV index: %v", err)
			uv = -1
		}
		uvChan <- uv
	}()

	go func() {
		aqi, quality, err := ws.weatherClient.GetAirQuality(apiResponse.Coord.Lat, apiResponse.Coord.Lon)
		if err != nil {
			log.Printf("⚠️  Failed to get air quality: %v", err)
			aqi = -1
			quality = "Unknown"
		}
		aqiChan <- struct {
			aqi     int
			quality string
		}{aqi, quality}
	}()

	// Wait for parallel requests
	weatherData.UVIndex = <-uvChan
	aqData := <-aqiChan
	weatherData.AQI = aqData.aqi
	weatherData.AirQuality = aqData.quality

	weatherData.CacheHit = false
	weatherData.LastUpdated = time.Now().Format("2006-01-02 15:04:05 MST")

	// Cache the result
	ws.cacheManager.Set(city, *weatherData)

	duration := time.Since(startTime)
	ws.metricsManager.RecordRequest(duration, false, nil)
	log.Printf("✅ Successfully fetched and cached weather data for: %s (took %v)", weatherData.Name, duration)

	return weatherData, nil
}

// transformWeatherData converts OpenWeatherMap response to our WeatherData model
func (ws *WeatherService) transformWeatherData(apiResponse *models.OpenWeatherResponse) *models.WeatherData {
	var data models.WeatherData

	// Basic info
	data.Name = apiResponse.Name
	data.Country = apiResponse.Sys.Country
	data.Timezone = apiResponse.Timezone
	data.Coordinates.Latitude = apiResponse.Coord.Lat
	data.Coordinates.Longitude = apiResponse.Coord.Lon

	// Time calculations
	location := time.FixedZone("Local", apiResponse.Timezone)
	data.LocalTime = time.Now().In(location).Format("15:04:05 MST")
	data.Sunrise = apiResponse.Sys.Sunrise
	data.Sunset = apiResponse.Sys.Sunset
	data.SunriseTime = utils.UnixToTime(apiResponse.Timezone, apiResponse.Sys.Sunrise)
	data.SunsetTime = utils.UnixToTime(apiResponse.Timezone, apiResponse.Sys.Sunset)

	// Temperature conversions
	celsius := apiResponse.Main.Temp
	kelvin, fahrenheit := utils.ConvertTemperatures(celsius)

	data.Main.Celsius = celsius
	data.Main.Kelvin = kelvin
	data.Main.Fahrenheit = fahrenheit

	// Feels like temperature
	feelsLikeKelvin, feelsLikeFahrenheit := utils.ConvertTemperatures(apiResponse.Main.FeelsLike)
	data.Main.FeelsLike.Celsius = apiResponse.Main.FeelsLike
	data.Main.FeelsLike.Kelvin = feelsLikeKelvin
	data.Main.FeelsLike.Fahrenheit = feelsLikeFahrenheit

	// Min temperature
	minKelvin, minFahrenheit := utils.ConvertTemperatures(apiResponse.Main.TempMin)
	data.Main.MinTemp.Celsius = apiResponse.Main.TempMin
	data.Main.MinTemp.Kelvin = minKelvin
	data.Main.MinTemp.Fahrenheit = minFahrenheit

	// Max temperature
	maxKelvin, maxFahrenheit := utils.ConvertTemperatures(apiResponse.Main.TempMax)
	data.Main.MaxTemp.Celsius = apiResponse.Main.TempMax
	data.Main.MaxTemp.Kelvin = maxKelvin
	data.Main.MaxTemp.Fahrenheit = maxFahrenheit

	// Other main data
	data.Main.Humidity = apiResponse.Main.Humidity
	data.Main.Pressure = apiResponse.Main.Pressure

	// Wind data
	data.Wind.Speed = apiResponse.Wind.Speed
	data.Wind.SpeedKmh = apiResponse.Wind.Speed * 3.6 // Convert m/s to km/h
	data.Wind.Degrees = apiResponse.Wind.Deg
	data.Wind.Direction = utils.WindDirection(apiResponse.Wind.Deg)

	// Cloud data
	data.Clouds.Cloudiness = apiResponse.Clouds.All

	// Weather description
	data.Weather = make([]struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}, len(apiResponse.Weather))

	for i, w := range apiResponse.Weather {
		data.Weather[i].Main = w.Main
		data.Weather[i].Description = w.Description
		data.Weather[i].Icon = w.Icon
	}

	// Visibility
	data.Visibility = apiResponse.Visibility

	return &data
}