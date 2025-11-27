package inmemory

import (
	"context"
	"log"
	"sync"

	"github.com/pcabreus/test-rockets-project/internal/model"
)

// InMemoryRocketStore is a simple in-memory implementation of RocketStore
type InMemoryRocketStore struct {
	rockets map[string]*model.Rocket // map of channel to Rocket
	// added mutex to protect rockets map
	mu sync.RWMutex
}

func NewRocketStore() *InMemoryRocketStore {
	return &InMemoryRocketStore{
		rockets: make(map[string]*model.Rocket),
	}
}

func (store *InMemoryRocketStore) GetRocket(ctx context.Context, channel string) (*model.Rocket, error) {
	store.mu.RLock()
	rocket, exists := store.rockets[channel]
	store.mu.RUnlock()

	if !exists {
		log.Println("Rocket not found for channel:", channel)
		return nil, model.ErrRocketNotFound
	}

	return rocket, nil
}

func (store *InMemoryRocketStore) SaveRocket(ctx context.Context, rocket *model.Rocket) error {
	// lock for write
	store.mu.Lock()
	store.rockets[rocket.Channel] = rocket
	store.mu.Unlock()

	log.Println("Rocket saved for channel:", rocket.Channel)

	return nil
}

func (store *InMemoryRocketStore) ListRockets(ctx context.Context, filter model.ListRocketsFilter) ([]*model.Rocket, error) {
	// read-lock while iterating the map
	store.mu.RLock()
	rockets := make([]*model.Rocket, 0, len(store.rockets))
	for _, rocket := range store.rockets {
		// Apply filtering logic here if needed
		rockets = append(rockets, rocket)
	}
	store.mu.RUnlock()

	return rockets, nil
}
