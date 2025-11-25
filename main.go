package main

import (
	"context"
	"log"
	"net/http"

	"github.com/pcabreus/test-rockets-project/internal/consumer"
	"github.com/pcabreus/test-rockets-project/internal/handler"
	"github.com/pcabreus/test-rockets-project/internal/model/inmemory"
)

// main function to start the webhook server and handle incoming Rocket messages.
func main() {
	log.Println("Starting webhook server on :8080")

	eventStore := inmemory.NewEventStore()
	rocketStore := inmemory.NewRocketStore()

	http.HandleFunc("/messages", handler.RocketEventHandler(eventStore))

	consumer := consumer.NewRocketEventConsumer(eventStore, rocketStore)

	if err := consumer.Start(context.Background()); err != nil {
		log.Println("Error starting consumer:", err)
	}

	http.ListenAndServe(":8080", nil)
}
