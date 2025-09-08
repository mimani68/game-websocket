package error

import "errors"

var (
	ErrInvalidNamespace = errors.New("invalid namespace")
	ErrInvalidEventType = errors.New("invalid event type")
	ErrEmptyID          = errors.New("id cannot be empty")
)
