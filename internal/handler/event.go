package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pcabreus/test-rockets-project/internal/model"
)

type Request struct {
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
	Speed       int    `json:"speed"`
	Mission     string `json:"mission"`
	NewMission  string `json:"newMission"`
	Reason      string `json:"reason"`
}

func RocketEventHandler(eventStore model.EventStore) http.HandlerFunc {
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
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("Error decoding JSON:", err)
			return
		}

		// Decision: used compound ID for simplicity and uniqueness
		// We are relying on Channel + MessageNumber to detect duplicates
		// In real implementation, consider UUIDs or database-generated IDs
		id := fmt.Sprintf("%s-%d", req.Metadata.Channel, req.Metadata.MessageNumber)

		mission := req.Message.Mission
		if req.Message.NewMission != "" {
			mission = req.Message.NewMission
		}

		status := model.EventStatusPending

		msg := model.Event{
			ID:          id,
			Status:      status,
			Mission:     mission,
			Type:        req.Message.Type,
			LaunchSpeed: req.Message.LaunchSpeed,
			Speed:       req.Message.Speed,
			Reason:      req.Message.Reason,
			Time:        req.Metadata.MessageTime,
			Number:      req.Metadata.MessageNumber,
			Channel:     req.Metadata.Channel,
			Event:       req.Metadata.MessageType,
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
