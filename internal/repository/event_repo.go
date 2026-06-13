package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/domain"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Save(ctx context.Context, event domain.Event) error {
	metadataJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO events (source, host, event_type, severity, timestamp, message, metadata)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err = r.db.Exec(ctx, query,
		event.Source,
		event.Host,
		event.EventType,
		event.Severity,
		event.Timestamp,
		event.Message,
		metadataJSON,
	)

	return err
}
