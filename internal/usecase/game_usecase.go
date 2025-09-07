package usecase

import (
	"context"
	"errors"
	"sync"

	"app/core-game/internal/domain/entity"
	"app/core-game/internal/domain/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
	ErrEventTypeMismatch  = errors.New("event type mismatch")
)

type GameUseCase interface {
	RegisterConnection(ctx context.Context, dto ConnectionDTO) (*entity.Connection, error)
	UnregisterConnection(ctx context.Context, id string) error
	HandleEvent(ctx context.Context, eventDTO EventDTO) error
	BroadcastState(ctx context.Context, ns entity.Namespace, message map[string]interface{}) error
	SendMessage(ctx context.Context, ns entity.Namespace, message map[string]interface{}) error
	BuyGood(ctx context.Context, ns entity.Namespace, details map[string]interface{}) error
}

type gameUseCase struct {
	connRepo  repository.ConnectionRepository
	eventRepo repository.EventRepository
	logger    *zap.Logger

	// In-memory connection map for fast concurrent access and message dispatch
	mu          sync.RWMutex
	connections map[string]*entity.Connection

	// Channel for async event processing
	eventChan chan *entity.Event
}

func NewGameUseCase(
	connRepo repository.ConnectionRepository,
	eventRepo repository.EventRepository,
	logger *zap.Logger,
) GameUseCase {
	uc := &gameUseCase{
		connRepo:    connRepo,
		eventRepo:   eventRepo,
		logger:      logger,
		connections: make(map[string]*entity.Connection),
		eventChan:   make(chan *entity.Event, 1000),
	}
	go uc.eventProcessor()
	return uc
}

func (g *gameUseCase) RegisterConnection(ctx context.Context, dto ConnectionDTO) (*entity.Connection, error) {
	if dto.ID == "" {
		dto.ID = uuid.NewString()
	}
	conn := entity.NewConnection(dto.ID, dto.Namespace)
	if err := g.connRepo.Save(ctx, conn); err != nil {
		return nil, err
	}
	g.mu.Lock()
	g.connections[conn.ID] = conn
	g.mu.Unlock()
	g.logger.Info("connection registered", zap.String("id", conn.ID), zap.String("namespace", string(conn.Namespace)))
	return conn, nil
}

func (g *gameUseCase) UnregisterConnection(ctx context.Context, id string) error {
	if err := g.connRepo.Delete(ctx, id); err != nil {
		return err
	}
	g.mu.Lock()
	delete(g.connections, id)
	g.mu.Unlock()
	g.logger.Info("connection unregistered", zap.String("id", id))
	return nil
}

func (g *gameUseCase) HandleEvent(ctx context.Context, eventDTO EventDTO) error {
	if eventDTO.ID == "" {
		eventDTO.ID = uuid.NewString()
	}

	event := entity.NewEvent(eventDTO.ID, eventDTO.Type, eventDTO.Namespace, eventDTO.Payload)

	// Save event to repository (etcd)
	if err := g.eventRepo.Save(ctx, event); err != nil {
		g.logger.Error("failed to save event", zap.Error(err))
		return err
	}

	// Push event to channel for asynchronous processing
	select {
	case g.eventChan <- event:
		// event sent for processing
	default:
		// channel full, drop event or log error
		g.logger.Warn("event channel full, dropping event", zap.String("event_id", event.ID))
		return errors.New("event processing overload")
	}
	return nil
}

func (g *gameUseCase) eventProcessor() {
	for event := range g.eventChan {
		// Process events concurrently but limit concurrency if needed
		go g.processEvent(event)
	}
}

func (g *gameUseCase) processEvent(event *entity.Event) {
	// Dispatch event based on type
	switch event.Type {
	case entity.EventTypeBroadcast:
		g.logger.Debug("processing broadcast event", zap.String("event_id", event.ID))
		g.broadcastToNamespace(event.Namespace, event.Payload)
	case entity.EventTypeSendMsg:
		g.logger.Debug("processing send_message event", zap.String("event_id", event.ID))
		g.sendMessageToNamespace(event.Namespace, event.Payload)
	case entity.EventTypeBuyGood:
		g.logger.Debug("processing buy_good event", zap.String("event_id", event.ID))
		g.processBuyGood(event.Namespace, event.Payload)
	default:
		g.logger.Warn("unknown event type", zap.String("event_id", event.ID), zap.String("type", string(event.Type)))
	}
}

func (g *gameUseCase) broadcastToNamespace(ns entity.Namespace, payload map[string]interface{}) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, conn := range g.connections {
		if conn.Namespace == ns {
			// Simulate sending message - in reality would push message to socket handler
			g.logger.Info("broadcast message to connection", zap.String("conn_id", conn.ID), zap.Any("payload", payload))
		}
	}
}

func (g *gameUseCase) sendMessageToNamespace(ns entity.Namespace, payload map[string]interface{}) {
	// For send_message, maybe only to specific users, check payload for target ids
	targetIDsRaw, ok := payload["target_ids"]
	if !ok {
		g.logger.Warn("send_message payload missing target_ids")
		return
	}
	targetIDs, ok := targetIDsRaw.([]interface{})
	if !ok {
		g.logger.Warn("send_message target_ids invalid format")
		return
	}

	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, tid := range targetIDs {
		if idStr, ok := tid.(string); ok {
			conn, ok := g.connections[idStr]
			if ok && conn.Namespace == ns {
				g.logger.Info("send message to connection", zap.String("conn_id", conn.ID), zap.Any("payload", payload))
			}
		}
	}
}

func (g *gameUseCase) processBuyGood(ns entity.Namespace, payload map[string]interface{}) {
	// Business logic for buying goods
	// For demo, just log and broadcast purchase event to namespace clients
	g.logger.Info("processing buy_good", zap.String("namespace", string(ns)), zap.Any("payload", payload))

	// Broadcast purchase update
	g.broadcastToNamespace(ns, map[string]interface{}{
		"type":    "purchase_update",
		"payload": payload,
	})
}

// Public method wrappers for use by delivery layer

func (g *gameUseCase) BroadcastState(ctx context.Context, ns entity.Namespace, message map[string]interface{}) error {
	event := entity.NewEvent(uuid.NewString(), entity.EventTypeBroadcast, ns, message)
	return g.HandleEvent(ctx, EventDTO{
		ID:        event.ID,
		Type:      event.Type,
		Namespace: event.Namespace,
		Payload:   event.Payload,
	})
}

func (g *gameUseCase) SendMessage(ctx context.Context, ns entity.Namespace, message map[string]interface{}) error {
	event := entity.NewEvent(uuid.NewString(), entity.EventTypeSendMsg, ns, message)
	return g.HandleEvent(ctx, EventDTO{
		ID:        event.ID,
		Type:      event.Type,
		Namespace: event.Namespace,
		Payload:   event.Payload,
	})
}

func (g *gameUseCase) BuyGood(ctx context.Context, ns entity.Namespace, details map[string]interface{}) error {
	event := entity.NewEvent(uuid.NewString(), entity.EventTypeBuyGood, ns, details)
	return g.HandleEvent(ctx, EventDTO{
		ID:        event.ID,
		Type:      event.Type,
		Namespace: event.Namespace,
		Payload:   event.Payload,
	})
}
