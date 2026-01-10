module github.com/yourusername/streaming-platform-api

go 1.21

require (
	github.com/99designs/gqlgen v0.17.49
	github.com/gorilla/websocket v1.5.1
	github.com/lib/pq v1.10.9
	github.com/redis/go-redis/v9 v9.5.1
	github.com/rabbitmq/amqp091-go v1.9.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/prometheus/client_golang v1.19.0
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/trace v1.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.24.0
	github.com/stretchr/testify v1.9.0
	github.com/golang/mock v1.6.0
	go.uber.org/zap v1.27.0
	github.com/jackc/pgx/v5 v5.5.5
	github.com/vektah/gqlparser/v2 v2.5.14
)

require (
	github.com/aws/aws-sdk-go-v2 v1.26.1
	github.com/aws/aws-sdk-go-v2/config v1.27.11
	github.com/aws/aws-sdk-go-v2/service/sqs v1.31.4
	github.com/aws/aws-sdk-go-v2/service/sns v1.29.4
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.30.4
)
