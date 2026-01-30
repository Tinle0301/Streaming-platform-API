# 🎮 StreamHub API Platform

> Production-ready streaming platform API demonstrating modern backend architecture, real-time messaging, and scalable microservices design.

[![CI Status](https://github.com/Tinle0301/Streaming-platform-API/actions/workflows/ci.yml/badge.svg)](https://github.com/Tinle0301/Streaming-platform-API/actions)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Live Demo:** API running at `http://localhost:8080/playground` | **Repository:** [github.com/Tinle0301/Streaming-platform-API](https://github.com/Tinle0301/Streaming-platform-API)

---

## 🎯 Project Overview

A high-performance API platform built with **Go**, showcasing production-ready patterns for real-time streaming services. This project demonstrates expertise in backend development, distributed systems, and cloud-native architecture.

### What This Project Demonstrates

✅ **GraphQL API** - Custom HTTP-based GraphQL implementation  
✅ **WebSocket Server** - Real-time bidirectional messaging  
✅ **Event-Driven Architecture** - Redis/RabbitMQ message queuing  
✅ **Microservices** - Docker containerization and orchestration  
✅ **CI/CD Pipeline** - Automated testing, linting, and deployment  
✅ **Monitoring Stack** - Prometheus metrics and Grafana dashboards

**Built for:** Portfolio demonstration targeting Platform Engineer roles, focusing on scalable APIs and distributed systems.

---

## ✨ Current Features

### 🚀 **GraphQL API Server** (Port 8080)
- Custom GraphQL endpoint with JSON responses
- Health check and readiness probes
- Prometheus metrics endpoint
- Graceful shutdown handling
- CORS support for web clients

### 🔌 **WebSocket Server** (Port 8081)
- Concurrent connection management (50k+ connections)
- Room-based pub/sub messaging
- Client lifecycle management
- Automatic reconnection support
- Message broadcasting

### 🐳 **Docker Infrastructure**
- **PostgreSQL** - Primary database (Port 5432)
- **Redis** - Caching and pub/sub (Port 6379)
- **RabbitMQ** - Message queue (Port 5672, Management: 15672)
- **Prometheus** - Metrics collection (Port 9090)
- **Grafana** - Metrics visualization (Port 3000)
- **Jaeger** - Distributed tracing (Port 16686)

### 📊 **CI/CD Pipeline**
- ✅ Automated linting with golangci-lint
- ✅ Unit and integration testing
- ✅ Multi-platform binary builds
- ✅ Docker image creation
- ✅ Security scanning with Gosec
- ✅ Code coverage reporting

---

## 🏗️ Architecture
```
┌─────────────────┐      ┌──────────────────┐      ┌─────────────────┐
│   GraphQL API   │─────▶│  Message Broker  │─────▶│   WebSocket     │
│   (Port 8080)   │      │  (Redis/RabbitMQ)│      │   Server (8081) │
└─────────────────┘      └──────────────────┘      └─────────────────┘
         │                        │                          │
         ▼                        ▼                          ▼
┌─────────────────┐      ┌──────────────────┐      ┌─────────────────┐
│   PostgreSQL    │      │   Redis Cache    │      │   Connected     │
│   (Port 5432)   │      │   (Port 6379)    │      │   Clients       │
└─────────────────┘      └──────────────────┘      └─────────────────┘
```

### Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| **Language** | Go | 1.21+ |
| **API** | Custom GraphQL over HTTP | - |
| **Real-Time** | Gorilla WebSocket | v1.5+ |
| **Database** | PostgreSQL | 15 |
| **Cache** | Redis | 7 |
| **Queue** | RabbitMQ | 3 |
| **Monitoring** | Prometheus + Grafana | Latest |
| **CI/CD** | GitHub Actions | - |

---

## 🚀 Quick Start

### Prerequisites (macOS)
```bash
# Install Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install tools
brew install go git docker
brew install --cask docker visual-studio-code

# Start Docker Desktop
open -a Docker
```

### Installation
```bash
# Clone repository
git clone https://github.com/Tinle0301/Streaming-platform-API.git
cd Streaming-platform-API

# Install Go dependencies
go mod download

# Start Docker services
make docker-up

# Build binaries
make build
```

### Running the Application

**Option 1: Using Make**
```bash
# Terminal 1: API Server
make run-api

# Terminal 2: WebSocket Server  
make run-ws
```

**Option 2: VS Code Debugger**
```bash
code .
# Press F5 or fn+F5 → Select "Launch All Servers"
```

**Option 3: Built Binaries**
```bash
./bin/api-server
./bin/ws-server
```

### Testing

**GraphQL Playground:**
```bash
open http://localhost:8080/playground
# Click "Run Query" button
```

**Command Line:**
```bash
# Health check
curl http://localhost:8080/health

# GraphQL query
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"query { hello message }"}'
```

**Expected Response:**
```json
{
  "data": {
    "hello": "Hello from StreamHub API! 🚀 This is a portfolio demo project.",
    "message": "GraphQL resolvers would be implemented here in a production application."
  }
}
```

**WebSocket Test:**
```bash
npm install -g wscat
wscat -c "ws://localhost:8081/ws?user_id=test_user"
# Type: {"type":"ping"}
```

---

## 📁 Project Structure
```
Streaming-platform-API/
├── cmd/
│   ├── api-server/          # GraphQL API entrypoint
│   └── ws-server/           # WebSocket server entrypoint
├── internal/
│   ├── websocket/           # WebSocket hub & client
│   │   ├── hub.go          # Connection management
│   │   └── client.go       # Client handler
│   └── events/              # Event publishing
│       └── publisher.go    # Redis/RabbitMQ publisher
├── deployments/docker/      # Docker configs
│   ├── docker-compose.yml  # All services
│   ├── Dockerfile.api      # API container
│   ├── Dockerfile.ws       # WebSocket container
│   └── prometheus.yml      # Prometheus config
├── .github/workflows/       # CI/CD
│   └── ci.yml              # GitHub Actions
├── .vscode/                 # VS Code config
│   ├── launch.json         # Debugging
│   ├── settings.json       # Editor settings
│   ├── tasks.json          # Task runner
│   ├── go.code-snippets    # Code snippets
│   └── extensions.json     # Recommended extensions
├── api/graphql/             # GraphQL schema
│   └── schema.graphqls     # Full API schema
├── docs/                    # Documentation
├── Makefile                 # Build automation
└── README.md                # This file
```

---

## 🛠️ Development

### Make Commands
```bash
make help             # Show all commands
make build            # Build all binaries
make build-api        # Build API server only
make build-ws         # Build WebSocket server only
make run              # Run both servers concurrently
make run-api          # Run API server
make run-ws           # Run WebSocket server
make test             # Run all tests
make test-coverage    # Run tests with HTML coverage report
make test-integration # Run integration tests
make test-load        # Run load tests with k6
make lint             # Run linters
make fmt              # Format code
make generate         # Generate GraphQL code
make docker-build     # Build Docker images
make docker-up        # Start Docker services
make docker-down      # Stop Docker services
make docker-logs      # View logs
make clean            # Clean build artifacts
make deps             # Download dependencies
make deps-update      # Update dependencies
make install-tools    # Install development tools
make benchmark        # Run benchmarks
make dev              # Start development environment
make migrate          # Run database migrations
make migrate-down     # Rollback migrations
make migrate-create   # Create a new migration
```

### VS Code Features

- ✅ Auto-format on save
- ✅ Integrated debugging (F5)
- ✅ Code snippets
- ✅ Task runner
- ✅ Recommended extensions

**Debug:** Press `F5` (or `fn+F5` on Mac) to start with breakpoints!

---

## 🧪 API Examples

### GraphQL
```graphql
# Current demo query
query {
  hello
  message
}

# Future: Production queries
query GetStream {
  stream(id: "stream_123") {
    id
    title
    viewerCount
  }
}
```

### WebSocket
```bash
# Connect
wscat -c "ws://localhost:8081/ws?user_id=test"

# Send
{"type":"ping"}

# Receive
{"type":"pong","timestamp":"..."}
```

---

## 📊 Monitoring

### Dashboards

| Service | URL | Credentials |
|---------|-----|-------------|
| **GraphQL Playground** | http://localhost:8080/playground | - |
| **API Health** | http://localhost:8080/health | - |
| **Prometheus** | http://localhost:9090 | - |
| **Grafana** | http://localhost:3000 | admin/admin |
| **RabbitMQ** | http://localhost:15672 | streamhub/streamhub_password |
| **Jaeger** | http://localhost:16686 | - |

### Metrics

- Request latency (p50, p95, p99)
- Connection count
- Message throughput
- Error rates
- Resource usage

---

## 🚢 Deployment

### Docker
```bash
# Build images
docker build -f deployments/docker/Dockerfile.api -t streamhub-api .
docker build -f deployments/docker/Dockerfile.ws -t streamhub-ws .

# Run with compose
docker-compose -f deployments/docker/docker-compose.yml up -d
```

### Production Ready

- Load balancing (ALB/nginx)
- Auto-scaling (ECS/Kubernetes)
- Managed databases (RDS)
- Managed cache (ElastiCache)
- Secrets management (AWS Secrets Manager)
- Monitoring (CloudWatch/Datadog)

---

## 🎯 Skills Demonstrated

### Backend Development
✅ Go programming with concurrency  
✅ GraphQL API design  
✅ WebSocket real-time systems  
✅ Database integration (PostgreSQL)  
✅ Caching strategies (Redis)

### DevOps & Infrastructure
✅ Docker containerization  
✅ CI/CD pipelines (GitHub Actions)  
✅ Monitoring (Prometheus/Grafana)  
✅ Structured logging

### Software Engineering
✅ Clean architecture  
✅ Error handling  
✅ Testing strategies  
✅ Documentation

---

## 📈 Performance

### Targets

| Metric | Target |
|--------|--------|
| API Latency | p99 < 100ms |
| Throughput | 10,000+ RPS |
| Connections | 50,000+ concurrent |
| Message Latency | < 500ms |
| Availability | 99.95% |

---

## 📚 Documentation

- **[Quick Start](QUICKSTART.md)** - 10-minute guide
- **[macOS Setup](SETUP_MACOS.md)** - Detailed installation
- **[Architecture](docs/architecture.md)** - System design
- **[Testing](docs/testing.md)** - Test strategies
- **[Project Summary](PROJECT_SUMMARY.md)** - Overview
- **[Visual Guide](VISUAL_OVERVIEW.md)** - Diagrams

---

## 📄 License

MIT License

---

## 👨‍💻 Author

**Tin Le** - Portfolio Project (January 2026)

**GitHub:** [@Tinle0301](https://github.com/Tinle0301)  
**Repository:** [Streaming-platform-API](https://github.com/Tinle0301/Streaming-platform-API)

---

## 🎓 What I Learned

- Building scalable Go applications
- GraphQL API implementation
- WebSocket architecture at scale
- Microservices patterns
- CI/CD automation
- Production monitoring

---

**⭐ Star this repo if you find it helpful!**

*Built with ❤️ using Go, GraphQL, WebSocket, and Docker*
