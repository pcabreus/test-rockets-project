package consumer

import (
	"context"
	"log"
	"time"

	"github.com/pcabreus/test-rockets-project/internal/model"
)

const firstExpectedEventNumber = 1

type EventConsumer interface {
	Start(ctx context.Context) error
	Consume(ctx context.Context, event model.Event) error
}

// RocketEventConsumer processes rocket events in order per channel
// It was designed for only one instance running.
type RocketEventConsumer struct {
	eventStore  model.EventStore
	rocketStore model.RocketStore
	next        map[string]int // track next expected message number per channel.
}

func NewRocketEventConsumer(eventStore model.EventStore, rocketStore model.RocketStore) *RocketEventConsumer {
	return &RocketEventConsumer{
		eventStore:  eventStore,
		rocketStore: rocketStore,
		next:        make(map[string]int),
	}
}

// Start the consumer
// Open a goroutine pending events
func (c *RocketEventConsumer) Start(ctx context.Context) error {
	log.Println("RocketEventConsumer started")

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("RocketEventConsumer stopped")
				return
			case <-ticker.C:
				log.Println("Checking for pending events...")
				pending, err := c.eventStore.Pending(ctx)
				if err != nil {
					log.Println("error listing pending events:", err)
					continue
				}

				for _, ev := range pending {
					ch := ev.Channel
					num := ev.Number

					// default first expected number is 1
					// we are to using mutex here assuming single instance
					if _, ok := c.next[ch]; !ok {
						c.next[ch] = firstExpectedEventNumber
					}

					// only process if it matches the next expected number
					if num == c.next[ch] {
						if err := c.Consume(ctx, ev); err != nil {
							log.Println("consume error:", err)
							// leave next unchanged so it will be retried later
							// go to a dead-letter queue after some retries
							continue
						}
						c.next[ch]++
					}
					// otherwise leave the event pending; it will be picked up again
				}
			}
		}
	}()

	return nil
}

// Consume processes a single event
func (c *RocketEventConsumer) Consume(ctx context.Context, event model.Event) error {
	// Simulate processing time
	log.Printf("Processing event ID=%s Channel=%s Number=%d Mission=%s\n", event.ID, event.Channel, event.Number, event.Mission)

	rocket, err := c.rocketStore.GetRocket(ctx, event.Channel)
	if err != nil {
		if err != model.ErrRocketNotFound {
			return err
		}

		// Create new rocket if not found
		rocket = &model.Rocket{
			Channel: event.Channel,
		}
	}

	switch event.Event {
	case "RocketLaunched":
		launchEvent := model.RocketLaunchedEvent{
			Type:        event.Type,
			LaunchSpeed: int64(event.LaunchSpeed),
			Mission:     event.Mission,
		}
		if err := rocket.ApplyLaunchEvent(launchEvent); err != nil {
			return err
		}
	case "RocketSpeedIncreased":
		speedEvent := model.RocketSpeedIncreasedEvent{
			By: int64(event.Speed),
		}
		if err := rocket.ApplySpeedIncreasedEvent(speedEvent); err != nil {
			return err
		}
	case "RocketSpeedDecreased":
		speedEvent := model.RocketSpeedDecreasedEvent{
			By: int64(event.Speed),
		}
		if err := rocket.ApplySpeedDecreasedEvent(speedEvent); err != nil {
			return err
		}
	case "RocketExploded":
		explodedEvent := model.RocketExplodedEvent{}
		if err := rocket.ApplyExplodedEvent(explodedEvent); err != nil {
			return err
		}

	default:
		log.Printf("Unknown event type: %s\n", event.Type)
		return nil // Ignore unknown event types
	}

	// Save updated rocket state
	if err := c.rocketStore.SaveRocket(ctx, rocket); err != nil {
		return err
	}

	// Mark event as processed
	if err := c.eventStore.Processed(ctx, event); err != nil {
		return err
	}

	// Use a single transaction in real implementation to avoid inconsistencies between rocket and event state

	log.Printf("Event ID=%s processed successfully\n", event.ID)
	return nil
}
