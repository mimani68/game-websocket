package validator

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidNamespace = errors.New("invalid namespace")
	ErrInvalidEventType = errors.New("invalid event type")
	ErrEmptyID          = errors.New("id cannot be empty")
)

func ValidateNamespace(ns string) error {
	if ns != "private" && ns != "guest" {
		return ErrInvalidNamespace
	}
	return nil
}

func ValidateEventType(et string) error {
	switch et {
	case "core.game.state.broadcast",
		"core.game.state.send_message",
		"core.game.state.buy_good":
		return nil
	default:
		return ErrInvalidEventType
	}
}

func ValidateID(id string) error {
	if id == "" {
		return ErrEmptyID
	}
	// Optionally validate UUID format
	// Simple regex UUID v4 check
	uuidRegex := regexp.MustCompile(`^[a-f0-9\-]{36}$`)
	if !uuidRegex.MatchString(id) {
		return errors.New("id is not a valid UUID")
	}
	return nil
}
