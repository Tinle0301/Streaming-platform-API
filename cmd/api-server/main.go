package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort        = "8080"
	defaultMetricsPort = "9090"
	shutdownTimeout    = 30 * time.Second
)

func main() {
	log.Println("Starting StreamHub API Server...")

	cfg := loadConfig()

	mux := http.NewServeMux()

	// Simple GraphQL endpoint
	mux.HandleFunc("/graphql", graphqlHandler)

	// GraphQL Playground
	if cfg.GraphQLPlayground {
		mux.HandleFunc("/playground", playgroundHandler)
		log.Println("GraphQL Playground enabled at /playground")
	}

	// Health check
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("/ready", readinessCheckHandler)

	// Metrics
	mux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Simple response - this is a demo project
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"hello":   "Hello from StreamHub API! 🚀 This is a portfolio demo project.",
			"message": "GraphQL resolvers would be implemented here in a production application.",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func playgroundHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>GraphQL Playground</title>
    <style>
        body { margin: 0; font-family: Arial, sans-serif; }
        .container { max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .query-box { background: #f5f5f5; padding: 20px; border-radius: 5px; margin: 20px 0; }
        pre { background: #282c34; color: #abb2bf; padding: 15px; border-radius: 5px; overflow-x: auto; }
        button { background: #0066cc; color: white; border: none; padding: 10px 20px; border-radius: 5px; cursor: pointer; font-size: 16px; }
        button:hover { background: #0052a3; }
        #result { margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎮 GraphQL Playground</h1>
        <p>StreamHub API - Portfolio Demo Project</p>
        
        <div class="query-box">
            <h3>Try this query:</h3>
            <pre>query {
  hello
  message
}</pre>
            <button onclick="runQuery()">Run Query</button>
        </div>
        
        <div id="result"></div>
    </div>
    
    <script>
        async function runQuery() {
            const response = await fetch('/graphql', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ query: 'query { hello message }' })
            });
            const data = await response.json();
            document.getElementById('result').innerHTML = 
                '<h3>Response:</h3><pre>' + JSON.stringify(data, null, 2) + '</pre>';
        }
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

func readinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ready","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}

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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
