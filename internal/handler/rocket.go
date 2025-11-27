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
	w.Header().Set("Content-Type", "application/json")

	rockets, err := h.rocketStore.ListRockets(r.Context(), model.ListRocketsFilter{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Use a custom response struct if pagination is implemented
	json.NewEncoder(w).Encode(rockets)
}

// GetRocket handles retrieving a specific rocket by channel
func (h *RocketStatusHandler) GetRocket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	channel := r.URL.Path[len("/rockets/"):]
	rocket, err := h.rocketStore.GetRocket(r.Context(), channel)
	if err != nil {
		if err == model.ErrRocketNotFound {
			w.WriteHeader(http.StatusNotFound)
			// TODO: return a JSON error message
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rocket)
}
