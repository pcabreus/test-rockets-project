package model

import "context"

const (
	EventStatusPending   = "pending"
	EventStatusProcessed = "processed"

	EventTypeRocketLaunched       = "RocketLaunched"
	EventTypeRocketSpeedIncreased = "RocketSpeedIncreased"
	EventTypeRocketSpeedDecreased = "RocketSpeedDecreased"
	EventTypeRocketExploded       = "RocketExploded"
	EventTypeRocketMissionChanged = "RocketMissionChanged"
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
	EventType   string
	By          int
}

type EventStore interface {
	SaveEvent(ctx context.Context, event Event) error
	Pending(ctx context.Context) ([]Event, error)
	Processed(ctx context.Context, event Event) error
}
