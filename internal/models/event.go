package models

import "context"

type Event struct {
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
	Mission     string `json:"mission"`
}

type EventStore interface {
	SaveEvent(ctx context.Context, event Event) error
}
