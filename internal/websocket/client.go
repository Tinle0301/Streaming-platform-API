package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// NewClient creates a new Client instance
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, sendBufferSize),
		userID:   userID,
		rooms:    make(map[string]bool),
		metadata: make(map[string]string),
	}
}

// ReadPump pumps messages from the websocket connection to the hub.
//
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse the incoming message
		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Handle the message based on type
		c.handleMessage(&message)

		// Update metrics
		c.hub.metrics.TotalMessagesRecv++
	}
}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages from the client
func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case "subscribe":
		// Subscribe to a room (e.g., stream-specific notifications)
		if room, ok := msg.Data["room"].(string); ok {
			c.hub.JoinRoom(room, c)
			c.sendAck("subscribed", room)
		}

	case "unsubscribe":
		// Unsubscribe from a room
		if room, ok := msg.Data["room"].(string); ok {
			c.hub.LeaveRoom(room, c)
			c.sendAck("unsubscribed", room)
		}

	case "ping":
		// Respond to ping with pong
		c.sendMessage("pong", map[string]interface{}{
			"timestamp": time.Now().Unix(),
		})

	case "message":
		// Handle custom messages (e.g., chat messages)
		// This could be forwarded to a message queue or processed directly
		log.Printf("Received message from client %s: %+v", c.userID, msg.Data)

	default:
		log.Printf("Unknown message type from client %s: %s", c.userID, msg.Type)
	}
}

// sendAck sends an acknowledgment message to the client
func (c *Client) sendAck(action, room string) {
	c.sendMessage("ack", map[string]interface{}{
		"action": action,
		"room":   room,
	})
}

// sendMessage sends a message to the client
func (c *Client) sendMessage(messageType string, data map[string]interface{}) {
	message := Message{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case c.send <- messageBytes:
	default:
		log.Printf("Client send buffer full, message dropped: userID=%s", c.userID)
	}
}

// SendNotification sends a notification to this specific client
func (c *Client) SendNotification(notificationType string, data map[string]interface{}) {
	c.sendMessage("notification", map[string]interface{}{
		"type": notificationType,
		"data": data,
	})
}

// GetUserID returns the user ID associated with this client
func (c *Client) GetUserID() string {
	return c.userID
}

// SetMetadata sets custom metadata for the client
func (c *Client) SetMetadata(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metadata[key] = value
}

// GetMetadata retrieves metadata for the client
func (c *Client) GetMetadata(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.metadata[key]
	return value, ok
}

// IsInRoom checks if the client is subscribed to a specific room
func (c *Client) IsInRoom(room string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.rooms[room]
	return ok
}

// GetRooms returns all rooms the client is subscribed to
func (c *Client) GetRooms() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	rooms := make([]string, 0, len(c.rooms))
	for room := range c.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}
