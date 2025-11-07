package openweathermap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/models"
)

// Client handles communication with OpenWeatherMap API
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new OpenWeatherMap API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.openweathermap.org/data/2.5",
	}
}

// GetWeather fetches weather data for a city
func (c *Client) GetWeather(city string) (*models.OpenWeatherResponse, error) {
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("%s/weather?q=%s&appid=%s&units=metric", c.baseURL, encodedCity, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("city '%s' not found", city)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status: %s", resp.Status)
	}

	var weatherResponse models.OpenWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %v", err)
	}

	return &weatherResponse, nil
}

// GetUVIndex fetches UV index for coordinates
func (c *Client) GetUVIndex(lat, lon float64) (float64, error) {
	url := fmt.Sprintf("%s/uvi?lat=%f&lon=%f&appid=%s", c.baseURL, lat, lon, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("UV API returned status: %s", resp.Status)
	}

	var uvData models.UVResponse
	if err := json.NewDecoder(resp.Body).Decode(&uvData); err != nil {
		return 0, err
	}

	return uvData.Value, nil
}

// GetAirQuality fetches air quality data for coordinates
func (c *Client) GetAirQuality(lat, lon float64) (int, string, error) {
	url := fmt.Sprintf("%s/air_pollution?lat=%f&lon=%f&appid=%s", c.baseURL, lat, lon, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("Air Quality API returned status: %s", resp.Status)
	}

	var aqData models.AirQualityResponse
	if err := json.NewDecoder(resp.Body).Decode(&aqData); err != nil {
		return 0, "", err
	}

	if len(aqData.List) == 0 {
		return 0, "Unknown", nil
	}

	aqi := aqData.List[0].Main.AQI
	quality := getAirQualityDescription(aqi)

	return aqi, quality, nil
}

// getAirQualityDescription converts AQI number to description
func getAirQualityDescription(aqi int) string {
	switch aqi {
	case 1:
		return "Good"
	case 2:
		return "Fair"
	case 3:
		return "Moderate"
	case 4:
		return "Poor"
	case 5:
		return "Very Poor"
	default:
		return "Unknown"
	}
}