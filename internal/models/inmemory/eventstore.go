package inmemory

import (
	"context"
	"log"

	"github.com/pcabreus/test-rockets-project/internal/models"
)

type InMemoryEventStore struct {
	// Simple in-memory store for demonstration purposes
	// First key is channel/ID, second key is message number/order
	// Key could represent different indexes in a real DB
	events map[string]map[int]models.Event
}

func NewEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[string]map[int]models.Event),
	}
}

func (store *InMemoryEventStore) SaveEvent(ctx context.Context, event models.Event) error {
	if _, exists := store.events[event.Metadata.Channel]; !exists {
		store.events[event.Metadata.Channel] = make(map[int]models.Event)
	}

	if _, exists := store.events[event.Metadata.Channel][event.Metadata.MessageNumber]; exists {
		log.Println("Duplicate event detected:", event.Metadata.MessageNumber)
		return nil // Ignore duplicate
	}

	store.events[event.Metadata.Channel][event.Metadata.MessageNumber] = event
	log.Println("Event saved:", event.Metadata.MessageNumber)
	return nil
}
