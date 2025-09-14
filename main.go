package main

import ("net/http")
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {

	defer producer.Close()
	http.HandleFunc("/order", orderHandler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}

}
