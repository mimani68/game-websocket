package entity

import (
	"app/core-game/constants"
	"sync"
	"time"
)

// Connection represents a client socket connection
type Connection struct {
	ID        string
	Namespace constants.Namespace
	CreatedAt time.Time
	// Add more metadata if needed
	mu sync.RWMutex
	// Store connection state or metadata if needed
	State map[string]interface{}
}

func NewConnection(id string, namespace constants.Namespace) *Connection {
	return &Connection{
		ID:        id,
		Namespace: namespace,
		CreatedAt: time.Now(),
		State:     make(map[string]interface{}),
	}
}

// UpdateState safely updates the state
func (c *Connection) UpdateState(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.State[key] = value
}

// GetState safely retrieves state
func (c *Connection) GetState(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.State[key]
	return v, ok
}
