package main

import (
	"log"
	"net/http"
)

func main() {

	defer producer.Close()
	http.HandleFunc("/order", orderHandler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}

}
