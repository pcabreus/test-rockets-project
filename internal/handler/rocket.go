package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pcabreus/test-rockets-project/internal/model"
)

type RocketStatusHandler struct {
	rocketStore model.RocketStore
}

func NewRocketHandlers(rocketStore model.RocketStore) RocketStatusHandler {
	return RocketStatusHandler{
		rocketStore: rocketStore,
	}
}

// ListRockets handles listing all rockets
func (h *RocketStatusHandler) ListRockets(w http.ResponseWriter, r *http.Request) {
	rockets, err := h.rocketStore.ListRockets(r.Context(), model.ListRocketsFilter{})
	if err != nil {
		// handle error
		return
	}

	// respond with rockets
	// TODO: Use a custom response struct if pagination is implemented
	json.NewEncoder(w).Encode(rockets)
}

// GetRocket handles retrieving a specific rocket by channel
func (h *RocketStatusHandler) GetRocket(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Path[len("/rockets/"):]
	rocket, err := h.rocketStore.GetRocket(r.Context(), channel)
	if err != nil {
		// handle error
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// respond with rocket
	json.NewEncoder(w).Encode(rocket)
}
