package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort        = "8080"
	defaultMetricsPort = "9090"
	shutdownTimeout    = 30 * time.Second
)

func main() {
	log.Println("Starting StreamHub API Server...")

	// Load configuration
	cfg := loadConfig()

	// Initialize dependencies
	// db := initDatabase(cfg.DatabaseURL)
	// cache := initCache(cfg.RedisURL)
	// eventPublisher := initEventPublisher(cfg)

	// Create GraphQL server
	srv := handler.NewDefaultServer(nil) // TODO: Add generated GraphQL schema

	// Setup HTTP server with middleware
	mux := http.NewServeMux()

	// GraphQL endpoint
	mux.Handle("/graphql", corsMiddleware(authMiddleware(rateLimitMiddleware(srv))))

	// GraphQL Playground (development only)
	if cfg.GraphQLPlayground {
		mux.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
		log.Println("GraphQL Playground enabled at /playground")
	}

	// Health check endpoint
	mux.HandleFunc("/health", healthCheckHandler)

	// Readiness check endpoint
	mux.HandleFunc("/ready", readinessCheckHandler)

	// Metrics endpoint (Prometheus)
	mux.Handle("/metrics", promhttp.Handler())

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 API Server listening on port %s", cfg.Port)
		log.Printf("📊 Metrics available at http://localhost:%s/metrics", cfg.MetricsPort)
		if cfg.GraphQLPlayground {
			log.Printf("🎮 GraphQL Playground at http://localhost:%s/playground", cfg.Port)
		}

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// Config holds application configuration
type Config struct {
	Port              string
	MetricsPort       string
	DatabaseURL       string
	RedisURL          string
	RabbitMQURL       string
	GraphQLPlayground bool
	JWTSecret         string
	Environment       string
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
	return Config{
		Port:              getEnv("API_PORT", defaultPort),
		MetricsPort:       getEnv("METRICS_PORT", defaultMetricsPort),
		DatabaseURL:       getEnv("DATABASE_URL", "postgresql://localhost:5432/streamhub"),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379"),
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		GraphQLPlayground: getEnv("GRAPHQL_PLAYGROUND", "true") == "true",
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Environment:       getEnv("ENVIRONMENT", "development"),
	}
}

// Middleware implementations

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement JWT authentication
		// Extract token from Authorization header
		// Validate token
		// Add user context to request

		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement rate limiting
		// Check request count per user/IP
		// Return 429 if limit exceeded

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log request details
		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// Health check handlers

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

func readinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Check database connection, Redis connection, etc.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ready","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
