package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// main function to start the webhook server and handle incoming Rocket messages.
func main() {
	log.Println("Starting webhook server on :8080")

	http.HandleFunc("/messages", handler)

	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Simple logging for incoming requests
	log.Println("Received webhook request", r.Method, r.URL.Path)

	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Invalid method:", r.Method)
		return
	}

	// TODO: Implement some mechanism to secure this endpoint with token

	// Decode the incoming JSON payload
	var msg RocketMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	// Process the message
	err = consume(r.Context(), msg)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error consuming message:", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Webhook received"))
}

type RocketMessage struct {
	Metadata Metadata `json:"metadata"`
	Message  Message  `json:"message"`
}

type Metadata struct {
	Channel       string `json:"channel"`
	MessageNumber int    `json:"messageNumber"`
	MessageTime   string `json:"messageTime"`
	MessageType   string `json:"messageType"`
}

type Message struct {
	Type        string `json:"type"`
	LaunchSpeed int    `json:"launchSpeed"`
	Mission     string `json:"mission"`
}

func consume(_ context.Context, message RocketMessage) error {
	log.Println("Consumed message", message)
	return nil
}
