package model

import "context"

const (
	EventStatusPending   = "pending"
	EventStatusProcessed = "processed"
)

type Event struct {
	ID          string
	Channel     string
	Status      string
	Type        string
	LaunchSpeed int
	Speed       int
	Mission     string
	Time        string
	Number      int
	Reason      string
	Event       string
}

type EventStore interface {
	SaveEvent(ctx context.Context, event Event) error
	Pending(ctx context.Context) ([]Event, error)
	Processed(ctx context.Context, event Event) error
}
