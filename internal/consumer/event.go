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
	eventStore model.EventStore
	next       map[string]int // track next expected message number per channel.
}

func NewRocketEventConsumer(eventStore model.EventStore) *RocketEventConsumer {
	return &RocketEventConsumer{
		eventStore: eventStore,
		next:       make(map[string]int),
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
	log.Printf("Processing event ID=%s Channel=%s Number=%d Mission=%s\n",
		event.ID, event.Channel, event.Number, event.Mission)
	time.Sleep(100 * time.Millisecond) // simulate work
	// Mark event as processed
	if err := c.eventStore.Processed(ctx, event); err != nil {
		return err
	}

	log.Printf("Event ID=%s processed successfully\n", event.ID)
	return nil
}
