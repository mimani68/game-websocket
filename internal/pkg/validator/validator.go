package validator

import (
	"app/core-game/constants"
	customError "app/core-game/internal/error"
	"errors"
	"regexp"
)

func ValidateNamespace(ns string) error {
	if ns != string(constants.KiakoV1) {
		return customError.ErrInvalidNamespace
	}
	return nil
}

func ValidateEventType(et string) error {
	switch et {
	case string(constants.EventTypeBroadcast),
		string(constants.EventTypeSendMsg),
		string(constants.EventTypeBuyGood):
		return nil
	default:
		return customError.ErrInvalidEventType
	}
}

func ValidateID(id string) error {
	if id == "" {
		return customError.ErrEmptyID
	}
	// Optionally validate UUID format
	// Simple regex UUID v4 check
	uuidRegex := regexp.MustCompile(`^[a-f0-9\-]{36}$`)
	if !uuidRegex.MatchString(id) {
		return errors.New("id is not a valid UUID")
	}
	return nil
}
