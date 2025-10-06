# Kafka Order Processing Service

A Go-based HTTP service that receives order data via REST API and publishes messages to Apache Kafka for asynchronous processing.

## Architecture Overview

This application implements a producer pattern where:
- HTTP server accepts POST requests with order data
- Orders are serialized to JSON and published to Kafka topic
- Kafka producer runs asynchronously with event monitoring
- Messages are keyed by OrderID for partition consistency

## Application Components

### Order Structure
```go
type Order struct {
    OrderID   string    // Unique order identifier
    UserID    string    // User who placed the order
    Total     float64   // Order total amount
    Status    string    // Order status (e.g., "pending")
    Timestamp time.Time // Order creation time (RFC3339 format)
}
```

### Configuration
- **Kafka Broker**: localhost:9092
- **Kafka Topic**: orders
- **HTTP Server Port**: 8080
- **API Endpoint**: POST /order

### API Usage

**Endpoint**: `POST http://localhost:8080/order`

**Request Body Example**:
```json
{
  "id": "123",
  "user_id": "u456",
  "total": 99.99,
  "status": "pending",
  "timestamp": "2024-10-01T12:34:56Z"
}
```

**Response**:
- 202 Accepted: Order received successfully
- 400 Bad Request: Invalid JSON payload
- 405 Method Not Allowed: Non-POST request
- 500 Internal Server Error: Kafka or serialization failure

**Example cURL**:
```bash
curl -X POST http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -d '{
    "id": "order-001",
    "user_id": "user-123",
    "total": 149.99,
    "status": "pending",
    "timestamp": "2025-10-06T17:00:00Z"
  }'
```

## Implementation Details

### Kafka Producer Initialization
- Producer is initialized in `init()` function before main execution
- Runs a goroutine to monitor delivery events asynchronously
- Logs successful deliveries and failures to console
- Producer is closed via `defer` in main function

### Message Publishing
- Messages are keyed by OrderID for consistent partition routing
- Partition assignment is automatic (kafka.PartitionAny)
- Synchronous produce call with asynchronous delivery confirmation
- Error handling at both produce and delivery stages

### Dependencies
- `github.com/confluentinc/confluent-kafka-go/v2/kafka` - Confluent Kafka Go client

## Setup Sequence

### 1. Start Kafka Infrastructure
```bash
docker-compose up -d
```

### 2. Access Kafka Container
```bash
docker exec -it kafka bash
```

### 3. Create Kafka Topic
```bash
kafka-topics.sh --create \
  --topic orders \
  --partitions 3 \
  --replication-factor 1 \
  --bootstrap-server localhost:9092
```

**Topic Configuration Parameters**:
- `--topic`: Topic name (orders)
- `--partitions`: Number of partitions (dictates max parallel consumers)
- `--replication-factor`: Number of broker replicas (for fault tolerance)
- `--bootstrap-server`: Kafka broker address

**Partition Considerations**:
- More partitions = higher parallelism for consumers
- Each partition maintains message ordering
- Messages with same key route to same partition

### 4. Monitor Kafka Messages in Real-Time
```bash
kafka-console-consumer.sh \
  --topic orders \
  --from-beginning \
  --bootstrap-server localhost:9092
```

This will display all messages published to the orders topic.

### 5. Start the Go Application
```bash
go run .
```

The server will start on port 8080 and begin accepting order requests.

## Testing the Service

### Send a Test Order
```bash
curl -X POST http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test-order-001",
    "user_id": "test-user-456",
    "total": 299.99,
    "status": "pending",
    "timestamp": "2025-10-06T17:03:20Z"
  }'
```

### Verify Message Delivery
Check the console consumer output to see the published message, and check the Go application logs for delivery confirmation.

## TODO
- How to write to database from Kafka consumer?
- Create CloudNativePG cluster - if in GCP, distribute nodes across regions
- Kubernetes cluster should be up and running with Cilium CNI
- Implement consumer service to process orders from Kafka
- Add message schema validation
- Implement dead letter queue for failed messages
- Add metrics and monitoring (Prometheus/Grafana)
- Implement graceful shutdown with producer flush