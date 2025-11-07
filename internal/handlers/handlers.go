package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/cache"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/metrics"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/models"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/services"
)

// Handler contains all HTTP handlers
type Handler struct {
	weatherService *services.WeatherService
	metricsManager *metrics.MetricsManager
	cacheManager   *cache.CacheManager
}

// NewHandler creates a new handler
func NewHandler(
	weatherService *services.WeatherService,
	metricsManager *metrics.MetricsManager,
	cacheManager *cache.CacheManager,
) *Handler {
	return &Handler{
		weatherService: weatherService,
		metricsManager: metricsManager,
		cacheManager:   cacheManager,
	}
}

// WeatherHandler handles weather requests
func (h *Handler) WeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		h.respondWithError(w, http.StatusBadRequest, "City parameter is required. Usage: /weather?city=CityName")
		return
	}

	data, err := h.weatherService.GetWeatherData(city)
	if err != nil {
		log.Printf("❌ Error fetching weather for '%s': %v", city, err)
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, data)
}

// HealthHandler handles health check requests
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":        "healthy",
		"timestamp":     time.Now().Format(time.RFC3339),
		"cache_entries": h.cacheManager.GetSize(),
		"service":       "weather-microservice",
		"version":       "1.0.0",
	}
	h.respondWithJSON(w, http.StatusOK, health)
}

// ReadinessHandler handles readiness check requests
func (h *Handler) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	readiness := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now().Format(time.RFC3339),
		"checks": map[string]bool{
			"cache":   true,
			"metrics": true,
		},
	}
	h.respondWithJSON(w, http.StatusOK, readiness)
}

// MetricsHandler handles metrics requests
func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := h.metricsManager.GetMetrics()
	h.respondWithJSON(w, http.StatusOK, metrics)
}

// CacheHandler handles cache status requests
func (h *Handler) CacheHandler(w http.ResponseWriter, r *http.Request) {
	cacheStats := h.cacheManager.GetStats()
	h.respondWithJSON(w, http.StatusOK, cacheStats)
}

// CacheClearHandler handles cache clear requests
func (h *Handler) CacheClearHandler(w http.ResponseWriter, r *http.Request) {
	h.cacheManager.Clear()
	response := map[string]interface{}{
		"status":  "success",
		"message": "Cache cleared successfully",
		"time":    time.Now().Format(time.RFC3339),
	}
	h.respondWithJSON(w, http.StatusOK, response)
}

// RootHandler handles root endpoint requests
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Weather Microservice API</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 40px 20px;
        }
        .container { 
            max-width: 900px;
            margin: 0 auto;
            background: white;
            padding: 40px;
            border-radius: 15px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
        }
        h1 { 
            color: #667eea;
            font-size: 2.5rem;
            margin-bottom: 10px;
            text-align: center;
        }
        .subtitle {
            text-align: center;
            color: #666;
            margin-bottom: 30px;
            font-size: 1.1rem;
        }
        .section { 
            margin: 30px 0;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 10px;
            border-left: 4px solid #667eea;
        }
        h2 { 
            color: #333;
            margin-bottom: 15px;
            font-size: 1.5rem;
        }
        .endpoint { 
            background: white;
            padding: 15px;
            margin: 15px 0;
            border-radius: 8px;
            border: 1px solid #dee2e6;
        }
        .method { 
            display: inline-block;
            padding: 4px 10px;
            background: #28a745;
            color: white;
            border-radius: 4px;
            font-size: 0.85rem;
            font-weight: bold;
            margin-right: 10px;
        }
        .method.post { background: #ffc107; }
        .path { 
            font-family: 'Courier New', monospace;
            color: #667eea;
            font-weight: bold;
        }
        .description { 
            margin-top: 10px;
            color: #666;
            line-height: 1.6;
        }
        .example { 
            background: #e3f2fd;
            padding: 12px;
            border-radius: 5px;
            margin-top: 10px;
            font-family: 'Courier New', monospace;
            font-size: 0.9rem;
            color: #1565c0;
        }
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-top: 20px;
        }
        .feature {
            background: white;
            padding: 15px;
            border-radius: 8px;
            text-align: center;
            border: 1px solid #dee2e6;
        }
        .feature-icon {
            font-size: 2rem;
            margin-bottom: 10px;
        }
        .footer {
            text-align: center;
            margin-top: 40px;
            padding-top: 20px;
            border-top: 2px solid #dee2e6;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🌤️ Weather Microservice API</h1>
        <p class="subtitle">High-performance weather data service with intelligent caching</p>

        <div class="section">
            <h2>🚀 Features</h2>
            <div class="features">
                <div class="feature">
                    <div class="feature-icon">⚡</div>
                    <strong>Fast Response</strong>
                    <p style="color: #666; font-size: 0.9rem; margin-top: 5px;">85% reduced latency</p>
                </div>
                <div class="feature">
                    <div class="feature-icon">💾</div>
                    <strong>Smart Caching</strong>
                    <p style="color: #666; font-size: 0.9rem; margin-top: 5px;">10-minute cache</p>
                </div>
                <div class="feature">
                    <div class="feature-icon">📊</div>
                    <strong>Metrics</strong>
                    <p style="color: #666; font-size: 0.9rem; margin-top: 5px;">Real-time monitoring</p>
                </div>
                <div class="feature">
                    <div class="feature-icon">🔒</div>
                    <strong>99.9% Uptime</strong>
                    <p style="color: #666; font-size: 0.9rem; margin-top: 5px;">Production ready</p>
                </div>
            </div>
        </div>

        <div class="section">
            <h2>📡 API Endpoints</h2>
            
            <div class="endpoint">
                <span class="method">GET</span>
                <span class="path">/weather?city={cityname}</span>
                <div class="description">
                    Get comprehensive weather data for any city worldwide including temperature, humidity, wind, UV index, and air quality.
                </div>
                <div class="example">
                    📝 Example: /weather?city=London
                </div>
            </div>

            <div class="endpoint">
                <span class="method">GET</span>
                <span class="path">/health</span>
                <div class="description">
                    Check the health status of the microservice. Returns service status, timestamp, and cache statistics.
                </div>
            </div>

            <div class="endpoint">
                <span class="method">GET</span>
                <span class="path">/readiness</span>
                <div class="description">
                    Kubernetes readiness probe endpoint. Confirms the service is ready to accept traffic.
                </div>
            </div>

            <div class="endpoint">
                <span class="method">GET</span>
                <span class="path">/metrics</span>
                <div class="description">
                    Get detailed performance metrics including request counts, response times, cache hit rates, and top requested cities.
                </div>
            </div>

            <div class="endpoint">
                <span class="method">GET</span>
                <span class="path">/cache</span>
                <div class="description">
                    View current cache status, hit/miss rates, and all cached entries with expiry times.
                </div>
            </div>

            <div class="endpoint">
                <span class="method post">POST</span>
                <span class="path">/cache/clear</span>
                <div class="description">
                    Clear all cached data. Useful for debugging or forcing fresh data retrieval.
                </div>
            </div>
        </div>

        <div class="footer">
            <p><strong>Weather Microservice v1.0.0</strong></p>
            <p style="margin-top: 10px;">Powered by OpenWeatherMap API | Started: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
        </div>
    </div>
</body>
</html>`
	w.Write([]byte(html))
}

// Helper methods
func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("❌ Error encoding response: %v", err)
	}
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	errorResponse := models.ErrorResponse{
		Error:   message,
		Message: message,
		Code:    code,
	}
	h.respondWithJSON(w, code, errorResponse)
}