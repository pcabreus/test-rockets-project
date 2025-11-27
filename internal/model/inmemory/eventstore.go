package inmemory

import (
	"context"
	"log"
	"slices"
	"sync"

	"github.com/pcabreus/test-rockets-project/internal/model"
)

// InMemoryEventStore is a simple in-memory implementation of EventStore
type InMemoryEventStore struct {
	events  map[string]model.Event // map of event ID to Event
	pending []string               // simulate an index to query for pending events
	// added mutex to protect events and pending
	mu sync.RWMutex
}

func NewEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events:  make(map[string]model.Event),
		pending: nil,
	}
}

// SaveEvent saves a new event to the store
func (store *InMemoryEventStore) SaveEvent(ctx context.Context, event model.Event) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.events[event.ID]; exists {
		log.Println("Duplicate event detected:", event.ID)
		return nil // Ignore duplicate
	}

	store.events[event.ID] = event
	store.pending = append(store.pending, event.ID)

	log.Println("Event saved:", event.ID)
	return nil
}

// Pending returns all events with status "pending"
func (store *InMemoryEventStore) Pending(ctx context.Context) ([]model.Event, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

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
	store.mu.Lock()
	defer store.mu.Unlock()

	storedEvent, exists := store.events[event.ID]
	if !exists {
		return nil // Event does not exist
	}

	storedEvent.Status = model.EventStatusProcessed
	store.events[event.ID] = storedEvent

	// Remove from pending list
	store.pending = slices.DeleteFunc(store.pending, func(id string) bool {
		return id == event.ID
	})

	log.Println("Event marked as processed:", event.ID)
	return nil
}
