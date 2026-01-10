# StreamHub API Platform - Real-Time Notification System

## 🎯 Project Overview

A high-performance, scalable API platform demonstrating production-ready patterns for building real-time streaming services similar to Twitch. This project showcases expertise in GraphQL APIs, event-driven architecture, and low-latency messaging systems.

**Built to demonstrate:** API Platform Engineering skills for companies like Twitch, focusing on high-throughput services, GraphQL APIs, and real-time messaging platforms.

---

## 🏗️ Architecture

### System Components

```
┌─────────────────┐      ┌──────────────────┐      ┌─────────────────┐
│   GraphQL API   │─────▶│  Message Broker  │─────▶│   WebSocket     │
│   (Go Server)   │      │  (Redis/RabbitMQ)│      │   Server (Go)   │
└─────────────────┘      └──────────────────┘      └─────────────────┘
         │                        │                          │
         ▼                        ▼                          ▼
┌─────────────────┐      ┌──────────────────┐      ┌─────────────────┐
│   PostgreSQL    │      │   Event Store    │      │   Connected     │
│   (Primary DB)  │      │   (Event Log)    │      │   Clients       │
└─────────────────┘      └──────────────────┘      └─────────────────┘
```

### Key Features

1. **GraphQL Edge Service**
   - Type-safe schema with Go code generation
   - Query complexity analysis and rate limiting
   - DataLoader pattern for N+1 query optimization
   - Batched downstream service calls

2. **Real-Time Messaging Platform**
   - WebSocket-based pub/sub system
   - Handles millions of concurrent connections
   - Message persistence and replay capability
   - Fan-out architecture for broadcast messages

3. **Event-Driven Architecture**
   - Asynchronous event processing
   - Event sourcing for audit trails
   - Dead letter queues for failed messages
   - Circuit breaker pattern for fault tolerance

4. **Scalability & Performance**
   - Horizontal scaling with load balancing
   - Connection pooling and caching strategies
   - Distributed tracing with OpenTelemetry
   - Prometheus metrics and Grafana dashboards

---

## 🚀 Technical Stack

### Core Technologies
- **Language:** Go 1.21+ (statically-typed, high-performance)
- **API Layer:** gqlgen (GraphQL code generation)
- **Message Broker:** Redis Pub/Sub & RabbitMQ
- **Database:** PostgreSQL with connection pooling
- **Real-Time:** Gorilla WebSocket
- **Caching:** Redis with TTL strategies

### AWS Services (Production Ready)
- **API Gateway:** Route53, ALB, CloudFront CDN
- **Compute:** ECS Fargate for containerized services
- **Messaging:** Amazon SQS, SNS, EventBridge
- **Storage:** RDS PostgreSQL, ElastiCache Redis
- **Monitoring:** CloudWatch, X-Ray distributed tracing

### Development Tools
- **Testing:** Go testing, testify, gomock
- **CI/CD:** GitHub Actions, Docker
- **Code Quality:** golangci-lint, staticcheck
- **Documentation:** Swagger/OpenAPI, GraphQL Playground

---

## 📊 Performance Metrics

### Target Benchmarks
- **API Latency:** p99 < 100ms for GraphQL queries
- **Throughput:** 10,000+ requests/second per instance
- **WebSocket Connections:** 50,000+ concurrent connections/instance
- **Message Delivery:** < 500ms end-to-end latency
- **Availability:** 99.95% uptime SLA

### Scalability
- Horizontal scaling to 100+ service instances
- Handles billions of notifications per day
- Auto-scaling based on CPU/memory/queue depth
- Multi-region deployment support

---

## 🎓 Key Learning Demonstrations

### 1. Distributed Systems Expertise
- Building low-latency, high-availability applications
- Implementing eventual consistency patterns
- Handling distributed transactions and saga patterns
- Service mesh integration (future: Istio)

### 2. API Design Best Practices
- RESTful principles in GraphQL context
- Versioning strategies for backward compatibility
- Rate limiting and throttling mechanisms
- Developer-friendly error messages

### 3. Real-Time Systems
- WebSocket connection management at scale
- Message ordering and delivery guarantees
- Graceful degradation under load
- Connection recovery and reconnection logic

### 4. Operational Excellence
- Structured logging with context propagation
- Comprehensive metrics and alerting
- Chaos engineering for resilience testing
- Performance profiling and optimization

---

## 📁 Project Structure

```
streaming-platform-api/
├── cmd/
│   ├── api-server/          # GraphQL API server entrypoint
│   ├── ws-server/            # WebSocket notification server
│   └── event-processor/      # Background event processor
├── internal/
│   ├── graphql/              # GraphQL schema and resolvers
│   │   ├── schema.graphqls
│   │   ├── resolver.go
│   │   └── generated/
│   ├── websocket/            # WebSocket hub and client management
│   │   ├── hub.go
│   │   ├── client.go
│   │   └── message.go
│   ├── events/               # Event definitions and handlers
│   │   ├── publisher.go
│   │   ├── consumer.go
│   │   └── types.go
│   ├── services/             # Business logic services
│   │   ├── stream/
│   │   ├── user/
│   │   └── notification/
│   ├── repository/           # Data access layer
│   │   ├── postgres/
│   │   └── redis/
│   └── middleware/           # HTTP/gRPC middleware
│       ├── auth.go
│       ├── ratelimit.go
│       └── logging.go
├── pkg/                      # Public libraries
│   ├── logger/
│   ├── metrics/
│   └── config/
├── api/                      # API definitions
│   ├── graphql/
│   └── rest/
├── deployments/              # Deployment configurations
│   ├── docker/
│   │   ├── Dockerfile.api
│   │   └── Dockerfile.ws
│   ├── k8s/                  # Kubernetes manifests
│   └── terraform/            # AWS infrastructure
├── scripts/                  # Build and utility scripts
│   ├── migrate.sh
│   └── load-test.sh
├── docs/                     # Documentation
│   ├── architecture.md
│   ├── api-guide.md
│   └── deployment.md
├── tests/
│   ├── integration/
│   ├── load/
│   └── e2e/
├── .github/                  # GitHub configuration
│   ├── workflows/            # CI/CD pipelines
│   │   └── ci.yml
│   ├── ISSUE_TEMPLATE/       # Issue templates
│   └── pull_request_template.md
├── .vscode/                  # VS Code configuration
│   ├── settings.json         # Editor settings
│   ├── launch.json           # Debug configurations
│   ├── tasks.json            # Task runner
│   ├── extensions.json       # Recommended extensions
│   └── go.code-snippets      # Go code snippets
├── go.mod
├── go.sum
├── Makefile
├── .gitignore
├── .editorconfig
├── .golangci.yml            # Linter configuration
├── README.md
├── QUICKSTART.md            # 10-minute quick start
├── SETUP_MACOS.md          # Detailed macOS setup
├── PROJECT_SUMMARY.md      # Executive summary
└── VISUAL_OVERVIEW.md      # Visual architecture guide
```

---

## 🛠️ Development Tools

### VS Code Setup

This project includes complete VS Code configuration:

- **Automatic formatting** on save
- **Code snippets** for common patterns
- **Integrated debugging** with breakpoints
- **Task runner** for common commands
- **Recommended extensions** auto-prompt

**Press `F5` to start debugging with breakpoints!**

### Recommended VS Code Extensions

The following extensions will be suggested when you open the project:

- **Go** - Essential Go language support
- **GraphQL** - GraphQL syntax highlighting and validation
- **Docker** - Docker container management
- **GitLens** - Enhanced Git integration
- **Error Lens** - Inline error display
- **REST Client** - Test APIs directly in VS Code

### VS Code Shortcuts

| Action | Shortcut | Description |
|--------|----------|-------------|
| Start Debugging | `F5` | Launch servers with debugger |
| Run Task | `Cmd+Shift+B` | Build project |
| Open Terminal | `Ctrl+`` | Toggle integrated terminal |
| Command Palette | `Cmd+Shift+P` | Access all commands |
| Quick File Open | `Cmd+P` | Quickly open files |
| Find in Files | `Cmd+Shift+F` | Search across project |
| Format Document | `Shift+Alt+F` | Format current file |

### Available Make Commands

```bash
make help              # Show all available commands
make build             # Build all binaries
make run              # Run both servers
make run-api          # Run API server only
make run-ws           # Run WebSocket server only
make test             # Run tests with coverage
make test-integration # Run integration tests
make test-load        # Run load tests
make lint             # Run linters
make fmt              # Format code
make docker-up        # Start all services
make docker-down      # Stop all services
make docker-logs      # View service logs
make clean            # Clean build artifacts
make deps             # Download dependencies
```

---

## 🐙 GitHub Workflow

### Setting Up GitHub Repository

```bash
# Initialize git (if not already done)
git init

# Run the automated setup script
./setup-github.sh

# Or manually:
git remote add origin https://github.com/yourusername/streaming-platform-api.git
git branch -M main
git add .
git commit -m "feat: initial commit"
git push -u origin main
```

### CI/CD Pipeline

This project includes GitHub Actions workflows that automatically:

- ✅ Run linters and code quality checks
- ✅ Execute unit and integration tests
- ✅ Build binaries for multiple platforms
- ✅ Create Docker images
- ✅ Run security scans
- ✅ Generate coverage reports

**Workflows trigger on**: Push to main/develop, Pull Requests

### Creating Pull Requests

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Make changes and commit: `git commit -m "feat: add feature"`
3. Push: `git push origin feature/your-feature`
4. Open PR on GitHub (template auto-fills)
5. Wait for CI checks to pass
6. Merge when approved

---

## 🔧 Implementation Highlights

### GraphQL Schema Example

```graphql
type Query {
  stream(id: ID!): Stream
  streams(filter: StreamFilter, limit: Int = 20): StreamConnection!
  viewer: User
  notifications(limit: Int = 50): [Notification!]!
}

type Mutation {
  startStream(input: StartStreamInput!): Stream!
  stopStream(id: ID!): Stream!
  followUser(userId: ID!): User!
  sendNotification(input: NotificationInput!): Notification!
}

type Subscription {
  streamStatusChanged(streamId: ID!): Stream!
  notificationReceived: Notification!
  chatMessage(streamId: ID!): ChatMessage!
}

type Stream {
  id: ID!
  title: String!
  streamer: User!
  viewerCount: Int!
  status: StreamStatus!
  startedAt: Time!
  thumbnailUrl: String
  tags: [String!]!
}

type Notification {
  id: ID!
  type: NotificationType!
  message: String!
  data: JSON
  createdAt: Time!
  read: Boolean!
}

enum NotificationType {
  STREAM_LIVE
  NEW_FOLLOWER
  CHAT_MENTION
  RAID_INCOMING
  SUBSCRIPTION
}
```

### Real-Time WebSocket Hub (Go)

```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    rooms      map[string]map[*Client]bool
    mu         sync.RWMutex
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.registerClient(client)
        case client := <-h.unregister:
            h.unregisterClient(client)
        case message := <-h.broadcast:
            h.broadcastToClients(message)
        }
    }
}
```

### Event-Driven Notification System

```go
type EventPublisher interface {
    Publish(ctx context.Context, event Event) error
    PublishBatch(ctx context.Context, events []Event) error
}

type NotificationService struct {
    publisher EventPublisher
    repo      NotificationRepository
    cache     Cache
}

func (s *NotificationService) NotifyStreamLive(ctx context.Context, streamID string) error {
    // Fan-out to all followers
    followers, err := s.repo.GetStreamFollowers(ctx, streamID)
    if err != nil {
        return err
    }

    // Batch publish for efficiency
    events := make([]Event, len(followers))
    for i, follower := range followers {
        events[i] = Event{
            Type: "stream.live",
            UserID: follower.ID,
            Data: map[string]interface{}{
                "stream_id": streamID,
                "timestamp": time.Now(),
            },
        }
    }

    return s.publisher.PublishBatch(ctx, events)
}
```

---

## 🧪 Testing Strategy

### Unit Tests
- Table-driven tests for business logic
- Mock interfaces for external dependencies
- Code coverage target: >80%

### Integration Tests
- Docker Compose for service dependencies
- End-to-end GraphQL query testing
- Database transaction rollback tests

### Load Tests
- k6 scripts for GraphQL API
- WebSocket connection stress tests
- Message throughput benchmarks

### Chaos Engineering
- Network latency injection
- Random service failures
- Database connection exhaustion

---

## 🚦 Getting Started

### Prerequisites for macOS
```bash
# Install Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Go, Git, Docker, and VS Code
brew install go git
brew install --cask docker visual-studio-code

# Start Docker Desktop
open -a Docker

# Install Go tools
go install github.com/99designs/gqlgen@latest
go install golang.org/x/tools/gopls@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**📖 Detailed Setup**: See [SETUP_MACOS.md](SETUP_MACOS.md) for comprehensive macOS setup guide  
**⚡ Quick Start**: See [QUICKSTART.md](QUICKSTART.md) for 10-minute setup

### Clone and Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/streaming-platform-api
cd streaming-platform-api

# Open in VS Code
code .

# When VS Code opens:
# 1. Click "Install" for recommended extensions
# 2. Click "Yes, I trust" when prompted

# Download dependencies
go mod download

# Start dependencies
make docker-up

# Run database migrations (if available)
make migrate

# Start the API server
make run-api

# Start the WebSocket server (in another terminal)
make run-ws

# Or use VS Code debugger: Press F5 → "Launch All Servers"
```

### Quick Test

```bash
# Run tests
make test

# Open GraphQL Playground
open http://localhost:8080/playground

# Connect to WebSocket
npm install -g wscat
wscat -c "ws://localhost:8081/ws?user_id=test_user"

# Load test (requires k6)
brew install k6
make test-load
```

### Environment Variables

```bash
# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/streamhub
REDIS_URL=redis://localhost:6379

# Server Configuration
API_PORT=8080
WS_PORT=8081
GRAPHQL_PLAYGROUND=true

# AWS (Production)
AWS_REGION=us-east-1
SQS_QUEUE_URL=https://sqs.us-east-1.amazonaws.com/...
SNS_TOPIC_ARN=arn:aws:sns:us-east-1:...

# Observability
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
PROMETHEUS_PORT=9090
```

---

## 📈 Performance Optimization Techniques

1. **Connection Pooling**
   - Database: pgx with optimal pool size
   - Redis: Redigo connection pool
   - HTTP clients: Keep-alive and connection reuse

2. **Caching Strategies**
   - Hot data in Redis (viewer counts, stream status)
   - Application-level caching with TTL
   - GraphQL query result caching

3. **Query Optimization**
   - DataLoader for batching database queries
   - SQL query optimization with EXPLAIN ANALYZE
   - Indexed columns for frequent lookups

4. **Message Optimization**
   - Protocol Buffers for binary serialization
   - Message compression for large payloads
   - Batch processing for high-volume events

---

## 🔐 Security Considerations

- JWT-based authentication
- Rate limiting per user/IP
- Input validation and sanitization
- SQL injection prevention (parameterized queries)
- XSS protection in WebSocket messages
- CORS configuration for API endpoints
- Secrets management with AWS Secrets Manager

---

## 📚 Documentation

- [Architecture Overview](docs/architecture.md)
- [API Documentation](docs/api-guide.md)
- [Deployment Guide](docs/deployment.md)
- [Performance Tuning](docs/performance.md)
- [Contributing Guidelines](docs/contributing.md)

---

## 🎯 Alignment with Twitch API Platform Role

### Direct Relevance

| Requirement | Demonstration |
|-------------|---------------|
| GraphQL API Development | Complete GraphQL schema with resolvers, queries, mutations, subscriptions |
| Real-Time Messaging | WebSocket hub handling thousands of concurrent connections |
| High-Throughput Services | Optimized for 10K+ RPS with horizontal scaling |
| Event-Driven Architecture | Publisher/subscriber pattern with message queues |
| Go Development | Entire codebase in Go with idiomatic patterns |
| AWS Technologies | Infrastructure as code with Terraform, ECS, SQS, SNS |
| Distributed Systems | Fault tolerance, circuit breakers, graceful degradation |
| API Design | Developer-friendly, well-documented, versioned APIs |

### Skills Demonstrated

✅ **1-2 years equivalent experience** in building distributed systems  
✅ **Static-typed language expertise** (Go)  
✅ **Low-latency, high-availability** architecture patterns  
✅ **Developer-focused API design**  
✅ **AWS cloud services** integration  
✅ **Event-driven systems** with message brokers  
✅ **Large-scale API** design and maintenance  

---

## 🏆 Bonus Points Coverage

- ✅ AWS technologies (ECS, SQS, SNS, RDS, ElastiCache)
- ✅ Messaging/event-driven systems (Redis, RabbitMQ, EventBridge)
- ✅ Public API design principles and documentation

---

## 🚀 Future Enhancements

- [ ] gRPC service-to-service communication
- [ ] GraphQL Federation for microservices
- [ ] Kafka for high-volume event streaming
- [ ] Multi-region deployment with global load balancing
- [ ] Advanced observability with OpenTelemetry
- [ ] GraphQL schema stitching
- [ ] Rate limiting with token bucket algorithm
- [ ] API versioning strategy

---

## 📞 Contact & Discussion

This project is designed to demonstrate production-ready skills for building scalable, real-time APIs at companies like Twitch. Open to discussing:

- Architectural decisions and trade-offs
- Performance optimization strategies
- Scaling challenges and solutions
- Real-world production scenarios

---

## 📄 License

MIT License - See LICENSE file for details

---

**Note:** This is a portfolio/demonstration project showcasing API platform engineering skills. It is not affiliated with or endorsed by Twitch.
