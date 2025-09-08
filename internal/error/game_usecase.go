package error

import "errors"

var (
	ErrConnectionNotFound = errors.New("connection not found")
	ErrEventTypeMismatch  = errors.New("event type mismatch")
	ErrEventChannelFull   = errors.New("event processing overload")
	ErrInvalidTargetIDs   = errors.New("invalid target_ids format")
)
