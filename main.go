package main

import (
	"log"
	"net/http"

	"github.com/pcabreus/test-rockets-project/internal/handler"
	"github.com/pcabreus/test-rockets-project/internal/models/inmemory"
)

// main function to start the webhook server and handle incoming Rocket messages.
func main() {
	log.Println("Starting webhook server on :8080")

	eventStore := inmemory.NewInMemoryEventStore()

	http.HandleFunc("/messages", handler.RocketEventHandler(eventStore))

	http.HandleFunc("/rockets")

	http.ListenAndServe(":8080", nil)
}
