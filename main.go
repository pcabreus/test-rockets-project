package main

import (
	"context"
	"log"
	"net/http"

	"github.com/pcabreus/test-rockets-project/internal/consumer"
	"github.com/pcabreus/test-rockets-project/internal/handler"
	"github.com/pcabreus/test-rockets-project/internal/model/inmemory"
)

func main() {
	// Start server and in-memory components.
	// Design note: for this PoC we run the consumer and HTTP server in the
	// same process and use in-memory stores. This keeps the submission
	// self-contained and focuses on the required guarantees: per-channel
	// ordering and idempotency. In production these components should be
	// separated and the stores replaced with durable implementations.

	log.Println("Starting webhook server on :8080")

	// TODO: capture OS signals and implement graceful shutdown (server + consumer)

	eventStore := inmemory.NewEventStore()
	rocketStore := inmemory.NewRocketStore()

	// HTTP handlers
	http.HandleFunc("/messages", handler.RocketEventHandler(eventStore))

	api := handler.NewRocketHandlers(rocketStore)
	http.HandleFunc("/rockets", api.ListRockets)
	http.HandleFunc("/rockets/", api.GetRocket)

	// Start the consumer in background. Single-instance assumption applies
	// (ordering state is kept in memory). To scale, replace EventStore with a
	// durable store and use partitioning or DB transactions to coordinate workers.
	consumer := consumer.NewRocketEventConsumer(eventStore, rocketStore)

	if err := consumer.Start(context.Background()); err != nil {
		log.Println("Error starting consumer:", err)
	}

	http.ListenAndServe(":8080", nil)
}
