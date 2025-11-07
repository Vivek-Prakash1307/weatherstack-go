package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vivek-Prakash1307/weather-Microservices/internal/config"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/handlers"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/middleware"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/services"
	"github.com/Vivek-Prakash1307/weather-Microservices/api/openweathermap"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/cache"
	"github.com/Vivek-Prakash1307/weather-Microservices/internal/metrics"

	"github.com/gorilla/mux"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", ".apiConfig", "Path to config file")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	// Override port with environment variable if set
	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	log.Println("✅ Configuration loaded successfully")

	// Initialize components
	cacheManager := cache.NewCacheManager(time.Duration(cfg.CacheExpiryMinutes) * time.Minute)
	metricsManager := metrics.NewMetricsManager()
	weatherClient := openweathermap.NewClient(cfg.OpenWeatherMapApiKey)
	weatherService := services.NewWeatherService(weatherClient, cacheManager, metricsManager)
	handler := handlers.NewHandler(weatherService, metricsManager, cacheManager)

	// Setup router with middleware
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimitPerMinute))

	// Register routes
	router.HandleFunc("/", handler.RootHandler).Methods("GET")
	router.HandleFunc("/health", handler.HealthHandler).Methods("GET")
	router.HandleFunc("/readiness", handler.ReadinessHandler).Methods("GET")
	router.HandleFunc("/weather", handler.WeatherHandler).Methods("GET")
	router.HandleFunc("/metrics", handler.MetricsHandler).Methods("GET")
	router.HandleFunc("/cache", handler.CacheHandler).Methods("GET")
	router.HandleFunc("/cache/clear", handler.CacheClearHandler).Methods("POST")

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + *port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🌍 Weather Microservice starting on port %s", *port)
		log.Printf("📖 API Documentation: http://localhost:%s", *port)
		log.Printf("🌤️  Weather endpoint: http://localhost:%s/weather?city=London", *port)
		log.Printf("❤️  Health check: http://localhost:%s/health", *port)
		log.Printf("📊 Metrics: http://localhost:%s/metrics", *port)
		log.Printf("💾 Cache status: http://localhost:%s/cache", *port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}