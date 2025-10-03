// order.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Order struct {
	OrderID   string    `json:"id"`
	UserID    string    `json:"user_id"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	kafkaTopic  = "orders"
	kafkaBroker = "localhost:9092"
)

var producer *kafka.Producer

func init() {
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBroker})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %s", err)
	}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Message delivered to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// process the order from json payload in request body

	// 1. json.NewDecoder(r.Body) is the entire raw json payload
	// .Decode knows how to map json fields to struct fields
	// &order is a pointer to the order struct we want to populate
	// if the json payload is malformed, Decode will return an error
	// and we respond with a 400 Bad Request status code
	// otherwise, order struct is populated with data from json payload
	// e.g. {"id":"123","user_id":"u456","total":99.99,"status":"pending","timestamp":"2024-10-01T12:34:56Z"}
	// will populate order.OrderID = "123", order.UserID = "u456", etc.
	// Note: Timestamp should be in RFC3339 format for proper parsing
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// 2. Serialize the order to JSON
	orderBytes, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// 3. Publish the order message to Kafka
	topicString := kafkaTopic
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topicString,
			Partition: kafka.PartitionAny},
		Value: orderBytes,
		Key:   []byte(order.OrderID),
	}, nil)

	if err != nil {
		log.Printf("Failed to produce message: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Order received: %s\n", order.OrderID)
}
