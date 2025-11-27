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
	log.Println("Starting webhook server on :8080")

	// TODO: capture signals to stop the consumer gracefully

	eventStore := inmemory.NewEventStore()
	rocketStore := inmemory.NewRocketStore()

	http.HandleFunc("/messages", handler.RocketEventHandler(eventStore))

	api := handler.NewRocketHandlers(rocketStore)
	http.HandleFunc("/rockets", api.ListRockets)
	http.HandleFunc("/rockets/", api.GetRocket)

	// This can be done in a separate service/process in real implementation
	// For simplicity and inmemory store, we run it here
	consumer := consumer.NewRocketEventConsumer(eventStore, rocketStore)

	if err := consumer.Start(context.Background()); err != nil {
		log.Println("Error starting consumer:", err)
	}

	http.ListenAndServe(":8080", nil)
}
