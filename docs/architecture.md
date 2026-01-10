# StreamHub API Platform - Architecture Overview

## Table of Contents
1. [System Architecture](#system-architecture)
2. [Component Details](#component-details)
3. [Data Flow](#data-flow)
4. [Scalability Strategy](#scalability-strategy)
5. [Performance Optimization](#performance-optimization)
6. [Security Architecture](#security-architecture)

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Load Balancer (ALB)                      │
│                   Route53 + CloudFront CDN                       │
└──────────────────────────┬──────────────────────────────────────┘
                           │
        ┌──────────────────┴──────────────────┐
        │                                     │
        ▼                                     ▼
┌────────────────┐                  ┌────────────────┐
│   GraphQL API  │                  │   WebSocket    │
│   Servers      │                  │   Servers      │
│   (ECS Tasks)  │                  │   (ECS Tasks)  │
└────────┬───────┘                  └────────┬───────┘
         │                                   │
         │                                   │
         ▼                                   ▼
┌─────────────────────────────────────────────────────────┐
│              Event-Driven Message Layer                 │
│  ┌──────────┐  ┌──────────┐  ┌───────────────────┐    │
│  │  Redis   │  │ RabbitMQ │  │  Amazon SQS/SNS   │    │
│  │ Pub/Sub  │  │  Queues  │  │  EventBridge      │    │
│  └──────────┘  └──────────┘  └───────────────────┘    │
└─────────────────────────────────────────────────────────┘
         │                                   │
         │                                   │
         ▼                                   ▼
┌─────────────────┐              ┌─────────────────────┐
│   PostgreSQL    │              │   ElastiCache       │
│   (RDS)         │              │   Redis Cluster     │
│   - Primary     │              │   - Session Store   │
│   - Read Replica│              │   - Cache Layer     │
└─────────────────┘              └─────────────────────┘
```

### Microservices Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   GraphQL    │  │     REST     │  │   WebSocket  │     │
│  │   Endpoint   │  │   Endpoint   │  │   Endpoint   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   Business Logic Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Stream     │  │     User     │  │ Notification │     │
│  │   Service    │  │   Service    │  │   Service    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │    Chat      │  │   Analytics  │  │    Moderation│     │
│  │   Service    │  │   Service    │  │   Service    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      Data Access Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ Repository   │  │    Cache     │  │    Event     │     │
│  │   Pattern    │  │   Manager    │  │  Publisher   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. GraphQL API Server

**Purpose**: Provides a unified, type-safe API for clients

**Key Features**:
- Schema-first development with gqlgen
- Query complexity analysis
- Automatic batching and caching (DataLoader pattern)
- Rate limiting per user/IP
- Request validation and sanitization

**Technology Stack**:
- Go 1.21+
- gqlgen for code generation
- pgx for PostgreSQL connections
- Redis for caching

**Scaling Strategy**:
- Horizontal scaling with ECS Auto Scaling
- Read replicas for database queries
- Redis cluster for distributed caching
- Connection pooling (max 100 connections per instance)

### 2. WebSocket Server

**Purpose**: Handles real-time bidirectional communication

**Key Features**:
- Concurrent connection handling (50K+ per instance)
- Room-based pub/sub system
- Automatic reconnection handling
- Message persistence and replay
- Connection health monitoring

**Technology Stack**:
- Go with Gorilla WebSocket
- Redis Pub/Sub for message distribution
- In-memory client registry

**Scaling Strategy**:
- Sticky sessions via ALB
- Horizontal scaling based on connection count
- Message fan-out through Redis Pub/Sub
- Graceful connection migration on scale-down

### 3. Event-Driven Architecture

**Purpose**: Decouples services and enables async processing

**Components**:

#### Redis Pub/Sub
- Low-latency message delivery (<10ms)
- Best for real-time notifications
- Fire-and-forget semantics

#### RabbitMQ
- Guaranteed message delivery
- Dead letter queues for failed messages
- Message persistence and replay
- Complex routing patterns

#### AWS SQS/SNS/EventBridge
- Production-grade managed services
- Automatic scaling and durability
- Integration with other AWS services

**Event Types**:
- `stream.live` - Stream went live
- `stream.offline` - Stream ended
- `user.new_follower` - New follower
- `chat.message` - Chat message sent
- `raid.incoming` - Stream raid
- `subscription.new` - New subscription

### 4. Database Layer

**Primary Database: PostgreSQL**
- ACID transactions
- Complex queries and joins
- Full-text search
- JSON support for flexible schemas

**Schema Design**:
```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    follower_count INT DEFAULT 0,
    INDEX idx_username (username)
);

-- Streams table
CREATE TABLE streams (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    viewer_count INT DEFAULT 0,
    started_at TIMESTAMP,
    INDEX idx_status_viewers (status, viewer_count DESC)
);

-- Notifications table (partitioned by created_at)
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    data JSONB,
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
) PARTITION BY RANGE (created_at);
```

**Caching Strategy**:
- Hot data in Redis (viewer counts, online status)
- Application-level caching with TTL
- GraphQL query result caching
- Cache invalidation on updates

## Data Flow

### Stream Goes Live Flow

```
1. Streamer starts stream
   ↓
2. GraphQL Mutation: startStream()
   ↓
3. Create stream record in PostgreSQL
   ↓
4. Publish "stream.live" event to Redis/RabbitMQ
   ↓
5. Event consumers:
   - NotificationService: Fan out to all followers
   - AnalyticsService: Start tracking metrics
   - SearchService: Update stream index
   ↓
6. WebSocket Server broadcasts to subscribed clients
   ↓
7. Clients receive real-time notification
```

### Chat Message Flow

```
1. Viewer sends message via WebSocket
   ↓
2. WebSocket server validates message
   ↓
3. Publish to "chat.{stream_id}" channel
   ↓
4. All WebSocket servers subscribed to channel
   ↓
5. Broadcast to viewers in same stream room
   ↓
6. Async: Store message in database for history
   ↓
7. Async: Update chat analytics
```

### Notification Flow

```
1. Event triggers notification (e.g., new follower)
   ↓
2. NotificationService creates notification
   ↓
3. Store in PostgreSQL (for persistence)
   ↓
4. Cache in Redis (for quick access)
   ↓
5. Publish to user's notification channel
   ↓
6. WebSocket delivers to connected clients
   ↓
7. Push notification sent if offline (FCM/APNs)
```

## Scalability Strategy

### Horizontal Scaling

**API Servers**:
- Auto-scale based on CPU/Memory (target: 70% utilization)
- Min: 2 instances, Max: 50 instances
- Scale-up trigger: CPU > 70% for 2 minutes
- Scale-down trigger: CPU < 30% for 5 minutes

**WebSocket Servers**:
- Auto-scale based on connection count
- Target: 30,000 connections per instance
- Sticky sessions via ALB for connection affinity
- Graceful shutdown with connection migration

**Database**:
- Read replicas for read-heavy workloads
- Connection pooling (pgBouncer)
- Query optimization and indexing
- Partitioning for large tables (notifications)

### Caching Strategy

**Multi-Level Caching**:
1. Application-level (in-memory)
2. Redis cluster (distributed)
3. CDN (CloudFront for static assets)

**Cache Patterns**:
- **Cache-aside**: Application checks cache, then DB
- **Write-through**: Update cache on every write
- **TTL-based**: Expire after fixed duration

### Message Queue Scaling

**RabbitMQ**:
- Cluster with multiple nodes
- Queue sharding by stream ID
- Message TTL and size limits
- Dead letter queues for failures

**SQS/SNS**:
- Automatic scaling by AWS
- FIFO queues for ordering
- Message deduplication

## Performance Optimization

### Query Optimization

**GraphQL DataLoader**:
```go
// Batch multiple user queries into single DB call
userLoader := dataloader.NewBatchedLoader(
    func(keys []string) []*User {
        // SELECT * FROM users WHERE id IN (?, ?, ?)
        return repo.GetUsersByIDs(keys)
    },
)
```

**Database Indexing**:
- Composite indexes for common queries
- Partial indexes for filtered queries
- GIN indexes for JSONB columns

**Connection Pooling**:
```go
// PostgreSQL connection pool
pool, err := pgxpool.New(ctx, connString)
pool.Config().MaxConns = 100
pool.Config().MinConns = 10
pool.Config().MaxConnIdleTime = 5 * time.Minute
```

### Latency Reduction

**Target Latencies**:
- GraphQL queries: p99 < 100ms
- WebSocket messages: < 500ms end-to-end
- Database queries: < 50ms
- Cache hits: < 5ms

**Techniques**:
- Read replicas geographically distributed
- CDN for static content
- Compression for large payloads
- Protocol Buffers for binary serialization

## Security Architecture

### Authentication & Authorization

**JWT-Based Auth**:
```go
type Claims struct {
    UserID   string   `json:"user_id"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}
```

**GraphQL Directives**:
```graphql
type Mutation {
    startStream(...): Stream! @auth @rateLimit(limit: 10, window: 60)
    sendMessage(...): Message! @auth @rateLimit(limit: 100, window: 60)
}
```

### Rate Limiting

**Strategy**: Token bucket algorithm
- Default: 100 requests/minute per user
- GraphQL: 1000 queries/minute per user
- WebSocket: 10 connections per user

### Input Validation

- GraphQL schema validation
- SQL injection prevention (parameterized queries)
- XSS protection in chat messages
- File upload validation

### Network Security

- TLS 1.3 for all connections
- VPC isolation for database
- Security groups for service communication
- WAF rules for common attacks

---

## Monitoring & Observability

### Metrics (Prometheus)

**Key Metrics**:
- Request rate and latency (p50, p95, p99)
- Error rate and types
- WebSocket connection count
- Database connection pool utilization
- Cache hit/miss ratio
- Message queue depth

### Logging (Structured)

```go
log.Info("stream_started",
    zap.String("stream_id", streamID),
    zap.String("user_id", userID),
    zap.Int("viewer_count", 0),
    zap.Duration("startup_time", time.Since(start)),
)
```

### Distributed Tracing (OpenTelemetry)

- End-to-end request tracking
- Service dependency mapping
- Performance bottleneck identification
- Error propagation analysis

---

## Future Enhancements

1. **GraphQL Federation**: Microservices with independent schemas
2. **gRPC**: Service-to-service communication
3. **Kafka**: High-throughput event streaming
4. **Kubernetes**: Container orchestration
5. **Service Mesh**: Istio for traffic management
6. **Multi-Region**: Global deployment with data replication
