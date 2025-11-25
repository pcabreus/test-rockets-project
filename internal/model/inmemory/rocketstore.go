package inmemory

import (
	"context"
	"log"

	"github.com/pcabreus/test-rockets-project/internal/model"
)

type InMemoryRocketStore struct {
	// Simple in-memory store for demonstration purposes
	rockets map[string]*model.Rocket // map of channel to Rocket
	// TODO: add mutex for concurrent access in real implementation
}

func NewRocketStore() *InMemoryRocketStore {
	return &InMemoryRocketStore{
		rockets: make(map[string]*model.Rocket),
	}
}

func (store *InMemoryRocketStore) GetRocket(ctx context.Context, channel string) (*model.Rocket, error) {
	rocket, exists := store.rockets[channel]
	if !exists {
		return nil, model.ErrRocketNotFound
	}
	log.Println("Rocket not found for channel:", channel)
	return rocket, nil
}

func (store *InMemoryRocketStore) SaveRocket(ctx context.Context, rocket *model.Rocket) error {
	store.rockets[rocket.Channel] = rocket
	log.Println("Rocket saved for channel:", rocket.Channel)
	return nil
}
