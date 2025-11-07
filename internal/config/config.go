package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config represents the application configuration
type Config struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
	CacheExpiryMinutes   int    `json:"CacheExpiryMinutes"`
	RateLimitPerMinute   int    `json:"RateLimitPerMinute"`
	MaxConcurrentReqs    int    `json:"MaxConcurrentRequests"`
	ServerPort           string `json:"ServerPort"`
	LogLevel             string `json:"LogLevel"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file '%s' not found. Please create it with your OpenWeatherMap API key", filename)
	}

	// Read file
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Validate required fields
	if config.OpenWeatherMapApiKey == "" {
		return nil, fmt.Errorf("OpenWeatherMapApiKey is required in config file")
	}

	// Set default values
	if config.CacheExpiryMinutes == 0 {
		config.CacheExpiryMinutes = 10
	}
	if config.RateLimitPerMinute == 0 {
		config.RateLimitPerMinute = 100
	}
	if config.MaxConcurrentReqs == 0 {
		config.MaxConcurrentReqs = 50
	}
	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	return &config, nil
}

// SaveExampleConfig creates an example configuration file
func SaveExampleConfig(filename string) error {
	exampleConfig := Config{
		OpenWeatherMapApiKey: "your_api_key_here",
		CacheExpiryMinutes:   10,
		RateLimitPerMinute:   100,
		MaxConcurrentReqs:    50,
		ServerPort:           "8080",
		LogLevel:             "info",
	}

	bytes, err := json.MarshalIndent(exampleConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal example config: %v", err)
	}

	if err := ioutil.WriteFile(filename, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write example config: %v", err)
	}

	return nil
}