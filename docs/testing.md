# Testing Strategy & Implementation

## Overview

This document outlines the comprehensive testing strategy for the StreamHub API Platform, covering unit tests, integration tests, load tests, and chaos engineering practices.

## Testing Pyramid

```
           /\
          /  \
         / E2E\     E2E Tests (5%)
        /______\
       /        \
      /Integration\ Integration Tests (15%)
     /____________\
    /              \
   /   Unit Tests   \  Unit Tests (80%)
  /__________________\
```

## 1. Unit Tests

### Purpose
Test individual functions and methods in isolation.

### Coverage Target
- **Minimum**: 80% code coverage
- **Goal**: 90% for critical paths

### Example: WebSocket Hub Test

```go
package websocket_test

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/yourusername/streaming-platform-api/internal/websocket"
)

func TestHub_RegisterClient(t *testing.T) {
    // Arrange
    hub := websocket.NewHub()
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    go hub.Run(ctx)
    
    client := &websocket.Client{
        UserID: "test_user_1",
    }
    
    // Act
    hub.Register <- client
    time.Sleep(100 * time.Millisecond) // Allow processing
    
    // Assert
    assert.Equal(t, 1, hub.GetTotalClients())
    
    metrics := hub.GetMetrics()
    assert.Equal(t, int32(1), metrics.ActiveConnections)
}

func TestHub_BroadcastToRoom(t *testing.T) {
    tests := []struct {
        name          string
        room          string
        messageType   string
        clientsInRoom int
        expectedMsgs  int
    }{
        {
            name:          "broadcast to existing room",
            room:          "stream_123",
            messageType:   "viewer_count_update",
            clientsInRoom: 3,
            expectedMsgs:  3,
        },
        {
            name:          "broadcast to empty room",
            room:          "stream_999",
            messageType:   "test_message",
            clientsInRoom: 0,
            expectedMsgs:  0,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Example: Event Publisher Test

```go
package events_test

func TestRedisPublisher_Publish(t *testing.T) {
    // Setup Redis test container
    redis, err := miniredis.Run()
    require.NoError(t, err)
    defer redis.Close()
    
    publisher, err := events.NewRedisPublisher("redis://" + redis.Addr())
    require.NoError(t, err)
    defer publisher.Close()
    
    // Test event
    event := events.Event{
        ID:   "test_event_1",
        Type: "stream.live",
        Data: map[string]interface{}{
            "stream_id": "stream_123",
        },
    }
    
    // Act
    err = publisher.Publish(context.Background(), event)
    
    // Assert
    assert.NoError(t, err)
}

func TestRedisPublisher_PublishBatch(t *testing.T) {
    // Benchmark batch publishing
    events := make([]events.Event, 1000)
    for i := 0; i < 1000; i++ {
        events[i] = events.NewStreamLiveEvent(
            fmt.Sprintf("stream_%d", i),
            "user_1",
            nil,
        )
    }
    
    start := time.Now()
    err := publisher.PublishBatch(context.Background(), events)
    duration := time.Since(start)
    
    assert.NoError(t, err)
    assert.Less(t, duration, 1*time.Second, "Batch should complete within 1s")
}
```

### Running Unit Tests

```bash
# Run all unit tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test -v ./internal/websocket/...

# Run with race detector
go test -race ./...

# Benchmark tests
go test -bench=. -benchmem ./...
```

## 2. Integration Tests

### Purpose
Test interaction between components with real dependencies.

### Setup
Uses Docker Compose to spin up test dependencies.

### Example: GraphQL Integration Test

```go
package integration_test

import (
    "testing"
    "context"
    
    "github.com/stretchr/testify/suite"
)

type GraphQLTestSuite struct {
    suite.Suite
    client     *graphql.Client
    db         *sql.DB
    redis      *redis.Client
    cleanup    func()
}

func (suite *GraphQLTestSuite) SetupSuite() {
    // Start test containers
    compose := testcontainers.NewLocalDockerCompose(
        []string{"../../deployments/docker/docker-compose.test.yml"},
        "integration-test",
    )
    
    compose.WithCommand([]string{"up", "-d"}).Invoke()
    suite.cleanup = func() {
        compose.Down()
    }
    
    // Initialize clients
    suite.client = graphql.NewClient("http://localhost:8080/graphql")
    suite.db = initTestDB()
    suite.redis = initTestRedis()
}

func (suite *GraphQLTestSuite) TearDownSuite() {
    suite.cleanup()
}

func (suite *GraphQLTestSuite) TestStartStream() {
    // Arrange
    mutation := `
        mutation StartStream($input: StartStreamInput!) {
            startStream(input: $input) {
                id
                title
                status
                streamer {
                    username
                }
            }
        }
    `
    
    variables := map[string]interface{}{
        "input": map[string]interface{}{
            "title":      "Test Stream",
            "categoryId": "category_1",
            "tags":       []string{"test", "gaming"},
        },
    }
    
    // Act
    var response struct {
        StartStream struct {
            ID     string
            Title  string
            Status string
        }
    }
    
    err := suite.client.Query(context.Background(), mutation, &response, variables)
    
    // Assert
    suite.NoError(err)
    suite.Equal("Test Stream", response.StartStream.Title)
    suite.Equal("LIVE", response.StartStream.Status)
    
    // Verify event was published
    // Verify database record created
    // Verify cache updated
}

func TestGraphQLSuite(t *testing.T) {
    suite.Run(t, new(GraphQLTestSuite))
}
```

### Example: WebSocket Integration Test

```go
func TestWebSocketIntegration(t *testing.T) {
    // Start WebSocket server
    hub := websocket.NewHub()
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        websocket.ServeWs(hub, w, r)
    }))
    defer server.Close()
    
    // Connect client
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
    conn, _, err := websocket.DefaultDialer.Dial(wsURL+"?user_id=test_user", nil)
    require.NoError(t, err)
    defer conn.Close()
    
    // Subscribe to room
    subscribeMsg := map[string]interface{}{
        "type": "subscribe",
        "data": map[string]interface{}{
            "room": "stream_123",
        },
    }
    
    err = conn.WriteJSON(subscribeMsg)
    require.NoError(t, err)
    
    // Wait for acknowledgment
    var ackMsg map[string]interface{}
    err = conn.ReadJSON(&ackMsg)
    require.NoError(t, err)
    assert.Equal(t, "ack", ackMsg["type"])
    
    // Broadcast message to room
    hub.BroadcastToRoom("stream_123", "test_notification", map[string]interface{}{
        "message": "Hello, World!",
    })
    
    // Verify client receives message
    var receivedMsg map[string]interface{}
    err = conn.ReadJSON(&receivedMsg)
    require.NoError(t, err)
    assert.Equal(t, "test_notification", receivedMsg["type"])
}
```

### Running Integration Tests

```bash
# Start test dependencies
docker-compose -f deployments/docker/docker-compose.test.yml up -d

# Run integration tests
make test-integration

# Or with Go tags
go test -tags=integration ./tests/integration/...

# Stop test dependencies
docker-compose -f deployments/docker/docker-compose.test.yml down
```

## 3. Load Testing

### Purpose
Validate system performance under expected and peak loads.

### Tools
- **k6**: Modern load testing tool
- **vegeta**: HTTP load testing
- **wrk**: HTTP benchmarking

### GraphQL API Load Test (k6)

```javascript
// tests/load/api-load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '2m', target: 100 },  // Ramp up to 100 users
        { duration: '5m', target: 100 },  // Stay at 100 users
        { duration: '2m', target: 200 },  // Ramp up to 200 users
        { duration: '5m', target: 200 },  // Stay at 200 users
        { duration: '2m', target: 0 },    // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(99)<100'],  // 99% of requests under 100ms
        http_req_failed: ['rate<0.01'],    // Error rate under 1%
    },
};

const BASE_URL = 'http://localhost:8080';

export default function () {
    // Test: Get streams
    let streamQuery = `
        query GetStreams {
            streams(limit: 20) {
                edges {
                    node {
                        id
                        title
                        viewerCount
                        status
                    }
                }
            }
        }
    `;
    
    let response = http.post(
        `${BASE_URL}/graphql`,
        JSON.stringify({ query: streamQuery }),
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer test_token',
            },
        }
    );
    
    check(response, {
        'status is 200': (r) => r.status === 200,
        'response time OK': (r) => r.timings.duration < 100,
        'has data': (r) => JSON.parse(r.body).data !== undefined,
    });
    
    sleep(1);
}
```

### WebSocket Load Test (k6)

```javascript
// tests/load/websocket-load-test.js
import ws from 'k6/ws';
import { check } from 'k6';

export let options = {
    vus: 1000,  // 1000 virtual users
    duration: '5m',
};

export default function () {
    const url = 'ws://localhost:8081/ws?user_id=user_' + __VU;
    
    ws.connect(url, {}, function (socket) {
        socket.on('open', () => {
            console.log('WebSocket connected');
            
            // Subscribe to a room
            socket.send(JSON.stringify({
                type: 'subscribe',
                data: { room: 'stream_123' }
            }));
        });
        
        socket.on('message', (data) => {
            let message = JSON.parse(data);
            check(message, {
                'message has type': (m) => m.type !== undefined,
            });
        });
        
        socket.on('close', () => {
            console.log('WebSocket closed');
        });
        
        socket.setTimeout(() => {
            socket.close();
        }, 300000); // Keep alive for 5 minutes
    });
}
```

### Running Load Tests

```bash
# Install k6
brew install k6  # macOS
# or
wget -qO- https://github.com/grafana/k6/releases/download/v0.47.0/k6-v0.47.0-linux-amd64.tar.gz | tar xvz

# Run API load test
make test-load
# or
k6 run tests/load/api-load-test.js

# Run with custom configuration
k6 run --vus 500 --duration 10m tests/load/api-load-test.js

# Run WebSocket load test
k6 run tests/load/websocket-load-test.js

# Generate HTML report
k6 run --out json=results.json tests/load/api-load-test.js
k6 report results.json --export results.html
```

## 4. End-to-End (E2E) Tests

### Purpose
Test complete user workflows from client to database.

### Example: Stream Lifecycle E2E Test

```go
func TestStreamLifecycle_E2E(t *testing.T) {
    // 1. User starts stream
    streamID := startStream(t, "Test Stream")
    
    // 2. Verify stream appears in listings
    streams := getStreams(t)
    assert.Contains(t, streams, streamID)
    
    // 3. Follower receives notification
    notification := waitForNotification(t, followerUserID, 5*time.Second)
    assert.Equal(t, "stream.live", notification.Type)
    
    // 4. Viewer joins stream via WebSocket
    wsConn := connectWebSocket(t, viewerUserID)
    subscribeToStream(t, wsConn, streamID)
    
    // 5. Send chat message
    sendChatMessage(t, wsConn, streamID, "Hello!")
    
    // 6. Verify message received by other viewers
    msg := receiveMessage(t, wsConn, 2*time.Second)
    assert.Equal(t, "Hello!", msg.Content)
    
    // 7. Update stream title
    updateStream(t, streamID, "Updated Title")
    
    // 8. End stream
    stopStream(t, streamID)
    
    // 9. Verify analytics updated
    analytics := getStreamAnalytics(t, streamID)
    assert.Greater(t, analytics.TotalViews, 0)
}
```

## 5. Chaos Engineering

### Purpose
Test system resilience under failure conditions.

### Chaos Scenarios

#### 1. Database Connection Loss
```go
func TestChaos_DatabaseFailure(t *testing.T) {
    // Stop database
    stopPostgreSQL()
    defer startPostgreSQL()
    
    // Verify circuit breaker opens
    _, err := makeGraphQLRequest(t, streamQuery)
    assert.Error(t, err)
    
    // Verify graceful degradation
    // - Cached data still served
    // - Error messages user-friendly
    
    // Restart database
    startPostgreSQL()
    
    // Verify recovery
    waitForHealthy(t, 30*time.Second)
    response, err := makeGraphQLRequest(t, streamQuery)
    assert.NoError(t, err)
}
```

#### 2. Redis Cache Failure
```go
func TestChaos_CacheFailure(t *testing.T) {
    // Stop Redis
    stopRedis()
    defer startRedis()
    
    // Verify system continues operating
    // (slower, but functional)
    response, err := makeGraphQLRequest(t, streamQuery)
    assert.NoError(t, err)
    
    // Verify latency increased but within acceptable range
    assert.Less(t, response.Latency, 500*time.Millisecond)
}
```

#### 3. Message Queue Failure
```go
func TestChaos_MessageQueueFailure(t *testing.T) {
    // Stop RabbitMQ
    stopRabbitMQ()
    defer startRabbitMQ()
    
    // Verify events buffered in memory
    startStream(t, "Test Stream")
    
    // Restart RabbitMQ
    startRabbitMQ()
    
    // Verify buffered events published
    waitForEvents(t, 10*time.Second)
}
```

#### 4. Network Latency Injection
```bash
# Inject 500ms latency
tc qdisc add dev eth0 root netem delay 500ms

# Run tests
make test-integration

# Remove latency
tc qdisc del dev eth0 root
```

## 6. Performance Benchmarks

### GraphQL Query Benchmarks

```go
func BenchmarkGraphQL_GetStreams(b *testing.B) {
    client := setupGraphQLClient()
    query := `
        query {
            streams(limit: 20) {
                edges { node { id title viewerCount } }
            }
        }
    `
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := client.Query(context.Background(), query)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkWebSocket_BroadcastMessage(b *testing.B) {
    hub := websocket.NewHub()
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    go hub.Run(ctx)
    
    // Create 1000 clients
    for i := 0; i < 1000; i++ {
        client := createMockClient(hub, fmt.Sprintf("user_%d", i))
        hub.JoinRoom("stream_123", client)
    }
    
    message := map[string]interface{}{
        "type": "test",
        "data": map[string]interface{}{"message": "benchmark"},
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        hub.BroadcastToRoom("stream_123", "test", message)
    }
}
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkGraphQL_GetStreams -benchmem ./...

# With CPU profiling
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof

# With memory profiling
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof
```

## 7. Test Coverage

### Generating Coverage Reports

```bash
# Generate coverage
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# View in browser
open coverage.html
```

### Coverage Goals

| Component | Target Coverage |
|-----------|----------------|
| Business Logic | 90%+ |
| API Handlers | 85%+ |
| WebSocket Hub | 90%+ |
| Event Publishers | 85%+ |
| Utilities | 80%+ |
| Overall | 85%+ |

## 8. CI/CD Pipeline Tests

### GitHub Actions Workflow

```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
      
      redis:
        image: redis:7
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run linters
        run: make lint
      
      - name: Run unit tests
        run: make test
      
      - name: Run integration tests
        run: make test-integration
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

## Summary

This comprehensive testing strategy ensures:
- ✅ **High code quality** with extensive unit test coverage
- ✅ **Integration validation** with real dependencies
- ✅ **Performance guarantees** through load testing
- ✅ **Production readiness** via E2E and chaos testing
- ✅ **Continuous quality** through CI/CD automation

All tests are automated and run on every commit to maintain system reliability and performance standards.
