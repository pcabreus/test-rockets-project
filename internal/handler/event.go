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
	By          int    `json:"by"`
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

		// NOTE: This endpoint is unauthenticated for the PoC. In production
		// add authentication method and rate limiting.

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("Error decoding JSON:", err)
			return
		}

		// Minimal validation is expected here (channel, messageNumber).
		// TODO: More validation and timestamp parsing can be added; omitted for time.
		if req.Metadata.Channel == "" || req.Metadata.MessageNumber < 1 {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			log.Println("Invalid payload:", req)
			return
		}

		// Decision: use composite ID for deduplication
		// Rely on `channel + messageNumber` to detect duplicates. This is
		// simple and effective given the contract that messageNumber is
		// unique per channel.
		id := fmt.Sprintf("%s-%d", req.Metadata.Channel, req.Metadata.MessageNumber)

		mission := req.Message.Mission
		if req.Message.NewMission != "" {
			// Reuse mission field for mission updates
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
			EventType:   req.Metadata.MessageType,
			By:          req.Message.By,
		}

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
