package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Event represents a domain event in the system
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	UserID    string                 `json:"user_id,omitempty"`
	StreamID  string                 `json:"stream_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
}

// Publisher defines the interface for publishing events
type Publisher interface {
	Publish(ctx context.Context, event Event) error
	PublishBatch(ctx context.Context, events []Event) error
	Close() error
}

// RedisPublisher implements Publisher using Redis Pub/Sub
type RedisPublisher struct {
	client *redis.Client
}

// NewRedisPublisher creates a new Redis-based event publisher
func NewRedisPublisher(redisURL string) (*RedisPublisher, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Connected to Redis for event publishing")

	return &RedisPublisher{
		client: client,
	}, nil
}

// Publish publishes a single event to Redis
func (p *RedisPublisher) Publish(ctx context.Context, event Event) error {
	// Set timestamp if not set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Set version if not set
	if event.Version == "" {
		event.Version = "1.0"
	}

	// Marshal event to JSON
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to Redis channel based on event type
	channel := fmt.Sprintf("events:%s", event.Type)
	if err := p.client.Publish(ctx, channel, eventBytes).Err(); err != nil {
		return fmt.Errorf("failed to publish event to Redis: %w", err)
	}

	log.Printf("Published event: type=%s, id=%s, channel=%s", event.Type, event.ID, channel)
	return nil
}

// PublishBatch publishes multiple events in a batch
func (p *RedisPublisher) PublishBatch(ctx context.Context, events []Event) error {
	pipe := p.client.Pipeline()

	for _, event := range events {
		// Set defaults
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}
		if event.Version == "" {
			event.Version = "1.0"
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		channel := fmt.Sprintf("events:%s", event.Type)
		pipe.Publish(ctx, channel, eventBytes)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute batch publish: %w", err)
	}

	log.Printf("Published %d events in batch", len(events))
	return nil
}

// Close closes the Redis connection
func (p *RedisPublisher) Close() error {
	return p.client.Close()
}

// RabbitMQPublisher implements Publisher using RabbitMQ
type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQPublisher creates a new RabbitMQ-based event publisher
func NewRabbitMQPublisher(amqpURL string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange for events
	err = channel.ExchangeDeclare(
		"events",  // name
		"topic",   // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	log.Println("Connected to RabbitMQ for event publishing")

	return &RabbitMQPublisher{
		conn:    conn,
		channel: channel,
	}, nil
}

// Publish publishes a single event to RabbitMQ
func (p *RabbitMQPublisher) Publish(ctx context.Context, event Event) error {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if event.Version == "" {
		event.Version = "1.0"
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create routing key from event type (e.g., "stream.live")
	routingKey := event.Type

	err = p.channel.PublishWithContext(
		ctx,
		"events",    // exchange
		routingKey,  // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         eventBytes,
			DeliveryMode: amqp.Persistent,
			Timestamp:    event.Timestamp,
			MessageId:    event.ID,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event to RabbitMQ: %w", err)
	}

	log.Printf("Published event: type=%s, id=%s, routingKey=%s", event.Type, event.ID, routingKey)
	return nil
}

// PublishBatch publishes multiple events
func (p *RabbitMQPublisher) PublishBatch(ctx context.Context, events []Event) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}

	log.Printf("Published %d events in batch", len(events))
	return nil
}

// Close closes the RabbitMQ connection
func (p *RabbitMQPublisher) Close() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}

// MultiPublisher publishes to multiple backends (Redis and RabbitMQ)
type MultiPublisher struct {
	publishers []Publisher
}

// NewMultiPublisher creates a publisher that publishes to multiple backends
func NewMultiPublisher(publishers ...Publisher) *MultiPublisher {
	return &MultiPublisher{
		publishers: publishers,
	}
}

// Publish publishes to all configured publishers
func (p *MultiPublisher) Publish(ctx context.Context, event Event) error {
	for _, publisher := range p.publishers {
		if err := publisher.Publish(ctx, event); err != nil {
			log.Printf("Error publishing to backend: %v", err)
			// Continue with other publishers instead of failing fast
		}
	}
	return nil
}

// PublishBatch publishes batches to all configured publishers
func (p *MultiPublisher) PublishBatch(ctx context.Context, events []Event) error {
	for _, publisher := range p.publishers {
		if err := publisher.PublishBatch(ctx, events); err != nil {
			log.Printf("Error batch publishing to backend: %v", err)
		}
	}
	return nil
}

// Close closes all publishers
func (p *MultiPublisher) Close() error {
	for _, publisher := range p.publishers {
		if err := publisher.Close(); err != nil {
			log.Printf("Error closing publisher: %v", err)
		}
	}
	return nil
}

// EventType constants for common events
const (
	EventTypeStreamLive       = "stream.live"
	EventTypeStreamOffline    = "stream.offline"
	EventTypeNewFollower      = "user.new_follower"
	EventTypeChatMessage      = "chat.message"
	EventTypeRaidIncoming     = "raid.incoming"
	EventTypeRaidOutgoing     = "raid.outgoing"
	EventTypeSubscription     = "subscription.new"
	EventTypeGiftSubscription = "subscription.gift"
	EventTypeBitsCheered      = "bits.cheered"
	EventTypeStreamMilestone  = "stream.milestone"
)

// Helper functions to create common events

// NewStreamLiveEvent creates a stream live event
func NewStreamLiveEvent(streamID, streamerID string, data map[string]interface{}) Event {
	return Event{
		ID:        generateEventID(),
		Type:      EventTypeStreamLive,
		UserID:    streamerID,
		StreamID:  streamID,
		Data:      data,
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

// NewFollowerEvent creates a new follower event
func NewFollowerEvent(followerID, followedID string) Event {
	return Event{
		ID:     generateEventID(),
		Type:   EventTypeNewFollower,
		UserID: followedID,
		Data: map[string]interface{}{
			"follower_id": followerID,
			"followed_id": followedID,
		},
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

// NewChatMessageEvent creates a chat message event
func NewChatMessageEvent(streamID, userID, message string) Event {
	return Event{
		ID:       generateEventID(),
		Type:     EventTypeChatMessage,
		UserID:   userID,
		StreamID: streamID,
		Data: map[string]interface{}{
			"message": message,
		},
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}
