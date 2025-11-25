package inmemory

import (
	"context"
	"log"

	"github.com/pcabreus/test-rockets-project/internal/model"
)

type InMemoryEventStore struct {
	// Simple in-memory store for demonstration purposes
	events  map[string]model.Event // map of event ID to Event
	pending []string               // simulate an index to query for pending events
	// TODO: add mutex for concurrent access in real implementation
}

func NewEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[string]model.Event),
	}
}

func (store *InMemoryEventStore) SaveEvent(ctx context.Context, event model.Event) error {
	if _, exists := store.events[event.ID]; exists {
		log.Println("Duplicate event detected:", event.ID)
		return nil // Ignore duplicate
	}

	// TODO: Mutex for concurrent access in real implementation
	store.events[event.ID] = event
	store.pending = append(store.pending, event.ID)

	log.Println("Event saved:", event.ID)
	return nil
}

// Pending returns all events with status "pending"
func (store *InMemoryEventStore) Pending(ctx context.Context) ([]model.Event, error) {
	var pendingEvents []model.Event
	for _, idx := range store.pending {
		if event, exists := store.events[idx]; exists {
			pendingEvents = append(pendingEvents, event)
		}
	}

	return pendingEvents, nil
}

// Processed marks an event as processed and removes it from pending list
func (store *InMemoryEventStore) Processed(ctx context.Context, event model.Event) error {
	if _, exists := store.events[event.ID]; !exists {
		return nil // Event does not exist- Question: what to do if event does not exist?
	}

	event.Status = model.EventStatusProcessed

	store.events[event.ID] = event

	// Remove from pending list
	for i, id := range store.pending {
		if id == event.ID {
			store.pending = append(store.pending[:i], store.pending[i+1:]...)
			break
		}
	}

	log.Println("Event marked as processed:", event.ID)
	return nil
}
