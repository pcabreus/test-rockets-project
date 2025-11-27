package model

import (
	"context"
	"errors"
)

type RocketStatus string

const (
	StatusActive   RocketStatus = "ACTIVE"
	StatusExploded RocketStatus = "EXPLODED"
)

var (
	ErrAlreadyLaunched = errors.New("rocket already launched")
	ErrAlreadyExploded = errors.New("rocket already exploded")
	ErrRocketExploded  = errors.New("rocket has exploded")
	ErrRocketNotFound  = errors.New("rocket not found")
)

type Rocket struct {
	Channel string       `json:"channel"`
	Type    string       `json:"type"`
	Speed   int64        `json:"speed"`
	Mission string       `json:"mission"`
	Status  RocketStatus `json:"status"`
	Reason  string       `json:"reason,omitempty"`
}

// RocketLaunchedEvent Sent out once: when a rocket is launched for the first time.
type RocketLaunchedEvent struct {
	Type        string
	LaunchSpeed int64
	Mission     string
}

func (r *Rocket) ApplyLaunchEvent(e RocketLaunchedEvent) error {
	if r.Status == StatusActive || r.Status == StatusExploded {
		return ErrAlreadyLaunched
	}

	r.Type = e.Type
	r.Speed = e.LaunchSpeed
	r.Mission = e.Mission
	r.Status = StatusActive

	return nil
}

// RocketSpeedIncreasedEvent Continuously sent out: when the speed of a rocket is increased by a certain amount.
type RocketSpeedIncreasedEvent struct {
	By int64
}

func (r *Rocket) ApplySpeedIncreasedEvent(e RocketSpeedIncreasedEvent) error {
	if r.Status == StatusExploded {
		return ErrRocketExploded
	}

	r.Speed += e.By

	return nil
}

// RocketSpeedDecreasedEvent Continuously sent out: when the speed of a rocket is decreased by a certain amount.
type RocketSpeedDecreasedEvent struct {
	By int64
}

func (r *Rocket) ApplySpeedDecreasedEvent(e RocketSpeedDecreasedEvent) error {
	if r.Status == StatusExploded {
		return ErrRocketExploded
	}

	// TODO: Prevent speed from going negative?
	r.Speed -= e.By

	return nil
}

// RocketExploded Sent out once: if a rocket explodes due to an accident/malfunction.
type RocketExplodedEvent struct {
	Reason string
}

func (r *Rocket) ApplyExplodedEvent(e RocketExplodedEvent) error {
	if r.Status == StatusExploded {
		return ErrAlreadyExploded
	}

	r.Status = StatusExploded
	r.Reason = e.Reason

	return nil
}

// RocketMissionChangedEvent Continuously sent out: when the mission for a rocket is changed.
type RocketMissionChangedEvent struct {
	NewMission string
}

type ListRocketsFilter struct {
	FilterMission string
}

func (r *Rocket) ApplyMissionChangedEvent(e RocketMissionChangedEvent) error {
	r.Mission = e.NewMission
	return nil
}

type RocketStore interface {
	GetRocket(ctx context.Context, channel string) (*Rocket, error)
	SaveRocket(ctx context.Context, rocket *Rocket) error
	ListRockets(ctx context.Context, filter ListRocketsFilter) ([]*Rocket, error)
}
