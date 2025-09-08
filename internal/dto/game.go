package dto

import (
	"app/core-game/constants"
)

// EventDTO represents data transfer object for event creation
type EventDTO struct {
	ID        string                 `json:"id"`
	Type      constants.EventType    `json:"type"`
	Namespace constants.Namespace    `json:"namespace"`
	Payload   map[string]interface{} `json:"payload"`
}

// ConnectionDTO represents data transfer object for connection creation
type ConnectionDTO struct {
	ID        string              `json:"id"`
	Namespace constants.Namespace `json:"namespace"`
}
