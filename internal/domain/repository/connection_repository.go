package repository

import (
	"context"

	"app/core-game/constants"
	"app/core-game/internal/domain/entity"
)

// ConnectionRepository defines the interface for connection persistence
type ConnectionRepository interface {
	Save(ctx context.Context, conn *entity.Connection) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.Connection, error)
	ListByNamespace(ctx context.Context, ns constants.Namespace) ([]*entity.Connection, error)
}
