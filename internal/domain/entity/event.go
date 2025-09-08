package entity

import (
	"app/core-game/constants"
	"time"
)

type Event struct {
	ID        string
	Type      constants.EventType
	Namespace constants.Namespace
	Payload   map[string]interface{}
	CreatedAt time.Time
}

func NewEvent(id string, eventType constants.EventType, ns constants.Namespace, payload map[string]interface{}) *Event {
	return &Event{
		ID:        id,
		Type:      eventType,
		Namespace: ns,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
}
