package utils

import (
	"time"
)

// WindDirection converts wind degrees to cardinal direction
func WindDirection(degrees int) string {
	directions := []string{
		"N", "NNE", "NE", "ENE",
		"E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW",
		"W", "WNW", "NW", "NNW",
	}

	// Normalize degrees to 0-360 range
	degrees = degrees % 360
	if degrees < 0 {
		degrees += 360
	}

	// Each direction covers 22.5 degrees (360/16)
	index := int((float64(degrees)+11.25)/22.5) % 16
	return directions[index]
}

// UnixToTime converts Unix timestamp to local time string
func UnixToTime(timezoneOffset int, timestamp int64) string {
	location := time.FixedZone("Local", timezoneOffset)
	return time.Unix(timestamp, 0).In(location).Format("15:04:05")
}

// ConvertTemperatures converts Celsius to Kelvin and Fahrenheit
func ConvertTemperatures(celsius float64) (kelvin, fahrenheit float64) {
	kelvin = celsius + 273.15
	fahrenheit = (celsius * 9 / 5) + 32
	return kelvin, fahrenheit
}

// RoundToDecimal rounds a float to specified decimal places
func RoundToDecimal(value float64, decimals int) float64 {
	multiplier := 1.0
	for i := 0; i < decimals; i++ {
		multiplier *= 10
	}
	return float64(int(value*multiplier+0.5)) / multiplier
}