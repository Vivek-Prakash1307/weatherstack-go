import React, { useState, useEffect } from 'react';
import './App.css';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080';

// Weather icon mapping
const weatherIcons = {
  'clear sky': '☀️',
  'few clouds': '⛅',
  'scattered clouds': '☁️',
  'broken clouds': '☁️',
  'overcast clouds': '☁️',
  'shower rain': '🌦️',
  'rain': '🌧️',
  'light rain': '🌧️',
  'moderate rain': '🌧️',
  'heavy rain': '⛈️',
  'thunderstorm': '⛈️',
  'snow': '❄️',
  'light snow': '🌨️',
  'mist': '🌫️',
  'fog': '🌫️',
  'haze': '🌫️',
  'smoke': '💨',
  'dust': '💨',
  'sand': '💨',
};

// Country flags mapping
const countryFlags = {
  'GB': '🇬🇧', 'US': '🇺🇸', 'JP': '🇯🇵', 'FR': '🇫🇷',
  'IN': '🇮🇳', 'AU': '🇦🇺', 'DE': '🇩🇪', 'CA': '🇨🇦',
  'IT': '🇮🇹', 'ES': '🇪🇸', 'BR': '🇧🇷', 'RU': '🇷🇺',
  'CN': '🇨🇳', 'MX': '🇲🇽', 'NL': '🇳🇱', 'SE': '🇸🇪',
  'CH': '🇨🇭', 'BE': '🇧🇪', 'AT': '🇦🇹', 'NO': '🇳🇴',
  'DK': '🇩🇰', 'FI': '🇫🇮', 'PL': '🇵🇱', 'PT': '🇵🇹',
  'GR': '🇬🇷', 'CZ': '🇨🇿', 'IE': '🇮🇪', 'NZ': '🇳🇿',
  'SG': '🇸🇬', 'TH': '🇹🇭', 'AE': '🇦🇪', 'SA': '🇸🇦',
};

function getWeatherIcon(description) {
  const lowerDesc = description.toLowerCase();
  return weatherIcons[lowerDesc] || '🌤️';
}

function getUVIndexColor(uv) {
  if (uv < 3) return '#00e400';
  if (uv < 6) return '#ffff00';
  if (uv < 8) return '#ff7e00';
  if (uv < 11) return '#ff0000';
  return '#b567a4';
}

function getAQIColor(aqi) {
  switch(aqi) {
    case 1: return '#00e400';
    case 2: return '#ffff00';
    case 3: return '#ff7e00';
    case 4: return '#ff0000';
    case 5: return '#8f3f97';
    default: return '#636e72';
  }
}

function formatTemperature(temp) {
  return Math.round(temp * 10) / 10;
}

function App() {
  const [cityInput, setCityInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [data, setData] = useState(null);
  const [serverHealth, setServerHealth] = useState(null);

  // Check server health on mount
  useEffect(() => {
    fetch(`${API_BASE_URL}/health`)
      .then((r) => {
        if (!r.ok) {
          throw new Error(`Server responded with status ${r.status}`);
        }
        return r.json();
      })
      .then((h) => {
        console.log("✅ Server health:", h);
        setServerHealth(h);
        setError("");
      })
      .catch((err) => {
        console.error("❌ Health check failed:", err);
        setError("⚠️ Unable to connect to weather server. Please ensure the server is running.");
      });
  }, []);

  async function fetchWeatherData(city) {
    setLoading(true);
    setError('');
    try {
      const res = await fetch(`${API_BASE_URL}/weather?city=${encodeURIComponent(city)}`);
      const json = await res.json();
      
      if (!res.ok) {
        throw new Error(json.error || json.message || 'Failed to fetch weather data');
      }
      
      setData(json);
    } catch (err) {
      console.error('Fetch error:', err);
      setError(err.message || 'Unable to fetch weather data. Please try again.');
      setData(null);
    } finally {
      setLoading(false);
    }
  }

  function handleSubmit(e) {
    e.preventDefault();
    const city = cityInput.trim();
    if (city) {
      fetchWeatherData(city);
    }
  }

  function searchCity(city) {
    setCityInput(city);
    fetchWeatherData(city);
  }

  return (
    <div className="app">
      <div className="container">
        <div className="header">
          <h1>🌤️ Weather Microservice</h1>
          <p>Real-time weather data with intelligent caching</p>
          {serverHealth && (
            <div className="server-status">
              <span className="status-dot"></span>
              Server Online • {serverHealth.cache_entries} cached cities
            </div>
          )}
        </div>

        <div className="search-container">
          <form className="search-form" onSubmit={handleSubmit}>
            <input
              type="text"
              id="cityInput"
              className="search-input"
              placeholder="Enter city name (e.g., London, New York, Tokyo)"
              value={cityInput}
              onChange={(e) => setCityInput(e.target.value)}
              required
            />
            <button type="submit" className="search-button" disabled={loading}>
              {loading ? '⏳ Loading...' : '🔍 Get Weather'}
            </button>
          </form>

          <div className="quick-cities">
            <span style={{color: '#2d3436', fontWeight: 'bold'}}>Quick search:</span>
            {['London', 'New York', 'Tokyo', 'Paris', 'Mumbai', 'Sydney', 'Dubai', 'Singapore'].map(c => (
              <button key={c} className="city-button" onClick={() => searchCity(c)}>
                {c}
              </button>
            ))}
          </div>
        </div>

        {error && (
          <div className="error-message">
            <strong>❌ Error:</strong> {error}
          </div>
        )}

        {loading && (
          <div className="loading">
            <div className="loading-spinner"></div>
            <p>Fetching weather data...</p>
          </div>
        )}

        {data && !loading && (
          <div className="weather-card show">
            <div className="city-header">
              <h2 className="city-name">{data.name}</h2>
              <span className="country-flag">
                {countryFlags[data.country] || '🌍'} {data.country}
              </span>
              {data.cache_hit && (
                <span className="cache-badge">💾 Cached</span>
              )}
            </div>

            <div className="main-weather">
              <div className="temperature-section">
                <div className="temperature-main">
                  {formatTemperature(data.main.temp_celsius)}°C
                </div>
                <div className="temperature-alternate">
                  {formatTemperature(data.main.temp_fahrenheit)}°F
                </div>
                <div className="temperature-feels">
                  Feels like {formatTemperature(data.main.feels_like.celsius)}°C
                </div>
                <div className="temperature-range">
                  <span>↓ {formatTemperature(data.main.temp_min.celsius)}°C</span>
                  <span>↑ {formatTemperature(data.main.temp_max.celsius)}°C</span>
                </div>
              </div>

              <div className="weather-description">
                <div className="weather-icon">
                  {data.weather && data.weather.length > 0 ? 
                    getWeatherIcon(data.weather[0].description) : '🌤️'}
                </div>
                <div className="weather-main">
                  {data.weather && data.weather.length > 0 ? data.weather[0].main : ''}
                </div>
                <div className="weather-desc">
                  {data.weather && data.weather.length > 0 ? 
                    data.weather[0].description.charAt(0).toUpperCase() + 
                    data.weather[0].description.slice(1) : ''}
                </div>
              </div>
            </div>

            <div className="details-grid">
              <div className="detail-item">
                <div className="detail-icon">💨</div>
                <div className="detail-label">Wind</div>
                <div className="detail-value">
                  {data.wind.speed_ms} m/s {data.wind.direction}
                </div>
                <div className="detail-sub">
                  {data.wind.speed_kmh.toFixed(1)} km/h
                </div>
              </div>

              <div className="detail-item">
                <div className="detail-icon">💧</div>
                <div className="detail-label">Humidity</div>
                <div className="detail-value">{data.main.humidity}%</div>
              </div>

              <div className="detail-item">
                <div className="detail-icon">🌡️</div>
                <div className="detail-label">Pressure</div>
                <div className="detail-value">{data.main.pressure}</div>
                <div className="detail-sub">hPa</div>
              </div>

              <div className="detail-item">
                <div className="detail-icon">☁️</div>
                <div className="detail-label">Cloudiness</div>
                <div className="detail-value">{data.clouds.all}%</div>
              </div>

              <div className="detail-item">
                <div className="detail-icon">👁️</div>
                <div className="detail-label">Visibility</div>
                <div className="detail-value">
                  {data.visibility_meters >= 1000 ? 
                    `${(data.visibility_meters / 1000).toFixed(1)} km` : 
                    `${data.visibility_meters} m`}
                </div>
              </div>

              <div className="detail-item">
                <div className="detail-icon">☀️</div>
                <div className="detail-label">UV Index</div>
                <div 
                  className="detail-value" 
                  style={{
                    color: data.uv_index >= 0 ? getUVIndexColor(data.uv_index) : '#636e72',
                    fontWeight: 'bold'
                  }}
                >
                  {data.uv_index >= 0 ? data.uv_index.toFixed(1) : 'N/A'}
                </div>
              </div>

              <div className="detail-item">
                <div className="detail-icon">🌅</div>
                <div className="detail-label">Sunrise</div>
                <div className="detail-value">{data.sunrise_time}</div>
              </div>

              <div className="detail-item">
                <div className="detail-icon">🌇</div>
                <div className="detail-label">Sunset</div>
                <div className="detail-value">{data.sunset_time}</div>
              </div>

              <div className="detail-item">
                <div className="detail-icon">🌬️</div>
                <div className="detail-label">Air Quality</div>
                <div 
                  className="detail-value"
                  style={{
                    color: data.aqi >= 0 ? getAQIColor(data.aqi) : '#636e72',
                    fontWeight: 'bold'
                  }}
                >
                  {data.aqi >= 0 ? data.air_quality : 'N/A'}
                </div>
                {data.aqi >= 0 && (
                  <div className="detail-sub">AQI: {data.aqi}</div>
                )}
              </div>

              <div className="detail-item">
                <div className="detail-icon">🕐</div>
                <div className="detail-label">Local Time</div>
                <div className="detail-value">{data.local_time}</div>
              </div>

              <div className="detail-item">
                <div className="detail-icon">📍</div>
                <div className="detail-label">Coordinates</div>
                <div className="detail-value">
                  {data.coordinates.latitude.toFixed(4)}°
                </div>
                <div className="detail-sub">
                  {data.coordinates.longitude.toFixed(4)}°
                </div>
              </div>

              <div className="detail-item">
                <div className="detail-icon">🔄</div>
                <div className="detail-label">Last Updated</div>
                <div className="detail-value" style={{fontSize: '0.85rem'}}>
                  {data.last_updated}
                </div>
              </div>
            </div>
          </div>
        )}

        <div className="footer-info">
          <div className="info-card">
            <div className="info-icon">💾</div>
            <div className="info-text">
              <strong>Smart Caching</strong>
              <p>Data cached for 10 minutes • 85% faster response time</p>
            </div>
          </div>
          <div className="info-card">
            <div className="info-icon">⚡</div>
            <div className="info-text">
              <strong>High Performance</strong>
              <p>Microservice architecture • 99.9% uptime</p>
            </div>
          </div>
          <div className="info-card">
            <div className="info-icon">🌍</div>
            <div className="info-text">
              <strong>Global Coverage</strong>
              <p>500+ cities • Real-time updates</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;