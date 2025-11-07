package models

// OpenWeatherResponse represents the response from OpenWeatherMap API
type OpenWeatherResponse struct {
	Name  string `json:"name"`
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int64 `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Timezone int `json:"timezone"`
}

// UVResponse represents UV index response
type UVResponse struct {
	Value float64 `json:"value"`
}

// AirQualityResponse represents air quality response
type AirQualityResponse struct {
	List []struct {
		Main struct {
			AQI int `json:"aqi"`
		} `json:"main"`
	} `json:"list"`
}

// WeatherData represents the enhanced weather data structure
type WeatherData struct {
	Name      string `json:"name"`
	Country   string `json:"country"`
	Timezone  int    `json:"timezone"`
	LocalTime string `json:"local_time"`
	Main      struct {
		Kelvin     float64 `json:"temp_kelvin"`
		Celsius    float64 `json:"temp_celsius"`
		Fahrenheit float64 `json:"temp_fahrenheit"`
		FeelsLike  struct {
			Kelvin     float64 `json:"kelvin"`
			Celsius    float64 `json:"celsius"`
			Fahrenheit float64 `json:"fahrenheit"`
		} `json:"feels_like"`
		MinTemp struct {
			Kelvin     float64 `json:"kelvin"`
			Celsius    float64 `json:"celsius"`
			Fahrenheit float64 `json:"fahrenheit"`
		} `json:"temp_min"`
		MaxTemp struct {
			Kelvin     float64 `json:"kelvin"`
			Celsius    float64 `json:"celsius"`
			Fahrenheit float64 `json:"fahrenheit"`
		} `json:"temp_max"`
		Humidity int `json:"humidity"`
		Pressure int `json:"pressure"`
	} `json:"main"`
	Wind struct {
		Speed     float64 `json:"speed_ms"`
		SpeedKmh  float64 `json:"speed_kmh"`
		Direction string  `json:"direction"`
		Degrees   int     `json:"degrees"`
	} `json:"wind"`
	Clouds struct {
		Cloudiness int `json:"all"`
	} `json:"clouds"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Visibility  int     `json:"visibility_meters"`
	Sunrise     int64   `json:"sunrise"`
	Sunset      int64   `json:"sunset"`
	SunriseTime string  `json:"sunrise_time"`
	SunsetTime  string  `json:"sunset_time"`
	UVIndex     float64 `json:"uv_index"`
	AirQuality  string  `json:"air_quality"`
	AQI         int     `json:"aqi"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	LastUpdated string `json:"last_updated"`
	CacheHit    bool   `json:"cache_hit"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}