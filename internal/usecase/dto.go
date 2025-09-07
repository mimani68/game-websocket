package usecase

import (
	"app/core-game/internal/domain/entity"
)

// EventDTO represents data transfer object for event creation
type EventDTO struct {
	ID        string                 `json:"id"`
	Type      entity.EventType       `json:"type"`
	Namespace entity.Namespace       `json:"namespace"`
	Payload   map[string]interface{} `json:"payload"`
}

// ConnectionDTO represents data transfer object for connection creation
type ConnectionDTO struct {
	ID        string           `json:"id"`
	Namespace entity.Namespace `json:"namespace"`
}
