package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/domain"
	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/repository"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Consumer struct {
	js   jetstream.JetStream
	repo *repository.EventRepository
}

func NewConsumer(natsURL string, repo *repository.EventRepository) (*Consumer, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("create jetstream context: %w", err)
	}

	return &Consumer{js: js, repo: repo}, nil
}

// Start subscribes to the EVENTS stream and processes messages until ctx is cancelled.
func (c *Consumer) Start(ctx context.Context) error {
	stream, err := c.js.Stream(ctx, "EVENTS")
	if err != nil {
		return fmt.Errorf("get stream: %w", err)
	}

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "ingestion-worker",
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: "events.>",
	})
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	log.Println("worker: subscribed to events.>, waiting for messages...")

	_, err = cons.Consume(func(msg jetstream.Msg) {
		var event domain.Event
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			log.Printf("worker: failed to unmarshal event: %v", err)
			msg.Nak() // tell NATS this message wasn't processed; it will be redelivered
			return
		}

		if err := c.repo.Save(ctx, event); err != nil {
			log.Printf("worker: failed to save event: %v", err)
			msg.Nak()
			return
		}

		log.Printf("worker: saved event from host=%s type=%s", event.Host, event.EventType)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	<-ctx.Done()
	return nil
}
