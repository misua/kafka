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
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: kafka.String(kafkaTopic), Partition: kafka.PartitionAny},
		Value:          orderBytes,
		Key:            []byte(order.OrderID),
	}, nil)

	if err != nil {
		log.Printf("Failed to produce message: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Order received: %s\n", order.OrderID)
}
