package repository

import (
	"context"

	"app/core-game/internal/domain/entity"
)

// EventRepository defines interface for event persistence
type EventRepository interface {
	Save(ctx context.Context, event *entity.Event) error
	GetByID(ctx context.Context, id string) (*entity.Event, error)
	ListByNamespace(ctx context.Context, ns entity.Namespace) ([]*entity.Event, error)
}
