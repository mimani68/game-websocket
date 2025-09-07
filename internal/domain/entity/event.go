package entity

import (
	"time"
)

type EventType string

const (
	EventTypeBroadcast EventType = "core.game.state.broadcast"
	EventTypeSendMsg   EventType = "core.game.state.send_message"
	EventTypeBuyGood   EventType = "core.game.state.buy_good"
)

type Event struct {
	ID        string
	Type      EventType
	Namespace Namespace
	Payload   map[string]interface{}
	CreatedAt time.Time
}

func NewEvent(id string, eventType EventType, ns Namespace, payload map[string]interface{}) *Event {
	return &Event{
		ID:        id,
		Type:      eventType,
		Namespace: ns,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
}
