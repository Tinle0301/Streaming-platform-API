package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yourusername/streaming-platform-api/internal/websocket"
)

const (
	defaultWSPort = "8081"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking in production
		return true
	},
}

func main() {
	log.Println("Starting StreamHub WebSocket Server...")

	// Create WebSocket hub
	hub := websocket.NewHub()

	// Start hub in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)

	// Setup HTTP server
	mux := http.NewServeMux()

	// WebSocket endpoint
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Metrics endpoint
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := hub.GetMetrics()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// TODO: Format metrics properly
		log.Printf("Current metrics: %+v", metrics)
	})

	port := getEnv("WS_PORT", defaultWSPort)
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("🔌 WebSocket Server listening on port %s", port)
		log.Printf("   Connect at: ws://localhost:%s/ws", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start WebSocket server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down WebSocket server...")

	// Cancel hub context
	cancel()

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("WebSocket server forced to shutdown: %v", err)
	}

	log.Println("WebSocket server exited")
}

// serveWs handles websocket requests from clients
func serveWs(hub *websocket.Hub, w http.ResponseWriter, r *http.Request) {
	// Extract user ID from query params or JWT token
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		// TODO: Extract from JWT token in production
		userID = "anonymous"
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create new client
	client := websocket.NewClient(hub, conn, userID)

	// Register client with hub
	hub.Register <- client

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()

	log.Printf("New WebSocket connection: userID=%s", userID)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
