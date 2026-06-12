package queue

import (
	"encoding/json"
	"fmt"

	"context"

	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/domain"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Publisher struct {
	js jetstream.JetStream
}

// NewPublisher connects to NATS and ensures the stream exists.
func NewPublisher(natsURL string) (*Publisher, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}
	fmt.Printf("nats connected: %v, status: %v", nc != nil, nc.Status())

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("create jetstream context: %w", err)
	}
	fmt.Printf("jetstream created: %v", js != nil)

	ctx := context.Background()

	// Create the stream if it doesn't exist. "events.>" captures
	// any subject starting with "events." (e.g. events.security)
	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "EVENTS",
		Subjects: []string{"events.>"},
	})
	if err != nil {
		return nil, fmt.Errorf("create stream: %w", err)
	}
	fmt.Printf("stream created: %+v", stream)

	return &Publisher{js: js}, nil
}

// Publish marshals the event and sends it to the EVENTS stream.
func (p *Publisher) Publish(ctx context.Context, event domain.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	_, err = p.js.Publish(ctx, "events.security", data)
	if err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	return nil
}
