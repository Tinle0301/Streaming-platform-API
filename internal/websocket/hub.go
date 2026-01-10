package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Room-based subscriptions (e.g., stream-specific rooms)
	rooms map[string]map[*Client]bool

	// Inbound messages from the clients
	Broadcast chan *Message

	// Register requests from the clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Metrics
	metrics *HubMetrics
}

// HubMetrics tracks hub statistics
type HubMetrics struct {
	TotalConnections  int64
	ActiveConnections int32
	TotalMessagesSent int64
	TotalMessagesRecv int64
	LastMessageTime   time.Time
	RoomCounts        map[string]int
	mu                sync.RWMutex
}

// Message represents a WebSocket message
type Message struct {
	Type      string                 `json:"type"`
	Room      string                 `json:"room,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// Client represents a single WebSocket connection
type Client struct {
	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// User ID associated with this client
	userID string

	// Rooms this client is subscribed to
	rooms map[string]bool

	// Client metadata
	metadata map[string]string

	// Mutex for client operations
	mu sync.RWMutex
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512KB

	// Send buffer size
	sendBufferSize = 256
)

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan *Message, 1000),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		metrics: &HubMetrics{
			RoomCounts: make(map[string]int),
		},
	}
}

// Run starts the hub's main event loop
func (h *Hub) Run(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Hub shutting down...")
			h.shutdown()
			return

		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastMessage(message)

		case <-ticker.C:
			// Periodic maintenance tasks
			h.logMetrics()
		}
	}
}

// registerClient registers a new client connection
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	h.metrics.ActiveConnections++
	h.metrics.TotalConnections++

	log.Printf("Client registered: userID=%s, total=%d", client.userID, len(h.clients))
}

// unregisterClient removes a client connection
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		// Remove from all rooms
		for room := range client.rooms {
			h.removeFromRoom(room, client)
		}

		delete(h.clients, client)
		close(client.send)
		h.metrics.ActiveConnections--

		log.Printf("Client unregistered: userID=%s, total=%d", client.userID, len(h.clients))
	}
}

// broadcastMessage sends a message to all clients in a room or all clients
func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	var targetClients []*Client

	if message.Room != "" {
		// Send to specific room
		if roomClients, ok := h.rooms[message.Room]; ok {
			targetClients = make([]*Client, 0, len(roomClients))
			for client := range roomClients {
				targetClients = append(targetClients, client)
			}
		}
	} else {
		// Broadcast to all clients
		targetClients = make([]*Client, 0, len(h.clients))
		for client := range h.clients {
			targetClients = append(targetClients, client)
		}
	}

	// Send messages asynchronously
	for _, client := range targetClients {
		select {
		case client.send <- messageBytes:
			h.metrics.TotalMessagesSent++
		default:
			// Client's send buffer is full, close the connection
			log.Printf("Client send buffer full, closing connection: userID=%s", client.userID)
			h.Unregister <- client
		}
	}

	h.metrics.LastMessageTime = time.Now()
}

// JoinRoom adds a client to a room
func (h *Hub) JoinRoom(room string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*Client]bool)
	}

	h.rooms[room][client] = true
	client.rooms[room] = true
	h.metrics.RoomCounts[room]++

	log.Printf("Client joined room: userID=%s, room=%s, count=%d",
		client.userID, room, len(h.rooms[room]))
}

// LeaveRoom removes a client from a room
func (h *Hub) LeaveRoom(room string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.removeFromRoom(room, client)
}

// removeFromRoom is an internal helper (caller must hold lock)
func (h *Hub) removeFromRoom(room string, client *Client) {
	if roomClients, ok := h.rooms[room]; ok {
		delete(roomClients, client)
		delete(client.rooms, room)
		h.metrics.RoomCounts[room]--

		if len(roomClients) == 0 {
			delete(h.rooms, room)
			delete(h.metrics.RoomCounts, room)
		}

		log.Printf("Client left room: userID=%s, room=%s", client.userID, room)
	}
}

// BroadcastToRoom sends a message to all clients in a specific room
func (h *Hub) BroadcastToRoom(room string, messageType string, data map[string]interface{}) {
	message := &Message{
		Type:      messageType,
		Room:      room,
		Data:      data,
		Timestamp: time.Now(),
	}

	h.Broadcast <- message
}

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(messageType string, data map[string]interface{}) {
	message := &Message{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
	}

	h.Broadcast <- message
}

// GetMetrics returns current hub metrics
func (h *Hub) GetMetrics() *HubMetrics {
	h.metrics.mu.RLock()
	defer h.metrics.mu.RUnlock()

	// Create a copy to avoid race conditions
	metricsCopy := &HubMetrics{
		TotalConnections:  h.metrics.TotalConnections,
		ActiveConnections: h.metrics.ActiveConnections,
		TotalMessagesSent: h.metrics.TotalMessagesSent,
		TotalMessagesRecv: h.metrics.TotalMessagesRecv,
		LastMessageTime:   h.metrics.LastMessageTime,
		RoomCounts:        make(map[string]int),
	}

	for room, count := range h.metrics.RoomCounts {
		metricsCopy.RoomCounts[room] = count
	}

	return metricsCopy
}

// logMetrics logs current hub metrics
func (h *Hub) logMetrics() {
	metrics := h.GetMetrics()
	log.Printf("Hub Metrics - Active: %d, Total: %d, Messages Sent: %d, Rooms: %d",
		metrics.ActiveConnections,
		metrics.TotalConnections,
		metrics.TotalMessagesSent,
		len(metrics.RoomCounts))
}

// shutdown gracefully shuts down the hub
func (h *Hub) shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Println("Closing all client connections...")

	for client := range h.clients {
		close(client.send)
		client.conn.Close()
	}

	h.clients = make(map[*Client]bool)
	h.rooms = make(map[string]map[*Client]bool)

	log.Println("Hub shutdown complete")
}

// GetRoomCount returns the number of clients in a specific room
func (h *Hub) GetRoomCount(room string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if roomClients, ok := h.rooms[room]; ok {
		return len(roomClients)
	}
	return 0
}

// GetTotalClients returns the total number of connected clients
func (h *Hub) GetTotalClients() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.clients)
}
