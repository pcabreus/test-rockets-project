package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pcabreus/test-rockets-project/internal/models"
)

func RocketEventHandler(eventStore models.EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simple logging for incoming requests
		log.Println("Received webhook request", r.Method, r.URL.Path)

		// Only accept POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Println("Invalid method:", r.Method)
			return
		}

		// TODO: Implement some mechanism to secure this endpoint.

		// Decode the incoming JSON payload
		var msg models.Event
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("Error decoding JSON:", err)
			return
		}

		// Process the message
		err = eventStore.SaveEvent(r.Context(), msg)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Error consuming message:", err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Webhook received"))
	}
}

func consume(_ context.Context, eventStore EventStore, message RocketEvent) error {
	log.Println("Consumed message", message)

	err := eventStore.SaveEvent(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}
