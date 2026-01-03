package events

import (
	"context"

	"github.com/kuldeepstechwork/gocart-api/internal/config"
)

// EventPublisher is a placeholder for an event publisher.
type EventPublisher struct{}

// NewEventPublisher creates a new event publisher.
func NewEventPublisher(ctx context.Context, cfg *config.AWSConfig) (*EventPublisher, error) {
	return &EventPublisher{}, nil
}
