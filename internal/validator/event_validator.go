package validator

import (
	"errors"

	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/domain"
)

var validSeverities = map[string]bool{
	"low":      true,
	"medium":   true,
	"high":     true,
	"critical": true,
}

func ValidateEvent(e domain.Event) error {
	if !validSeverities[e.Severity] {
		return errors.New("invalid severity: must be low, medium, high, or critical")
	}
	return nil
}
