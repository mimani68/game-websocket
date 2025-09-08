package game

import (
	"context"

	"app/core-game/constants"
	"app/core-game/internal/domain/entity"
	"app/core-game/internal/domain/repository"
	"app/core-game/internal/dto"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GameUseCase interface {
	RegisterConnection(ctx context.Context, dto dto.ConnectionDTO) (*entity.Connection, error)
	UnregisterConnection(ctx context.Context, id string) error
	HandleEvent(ctx context.Context, eventDTO dto.EventDTO) error
	BroadcastState(ctx context.Context, ns constants.Namespace, message map[string]interface{}) error
	SendMessage(ctx context.Context, ns constants.Namespace, message map[string]interface{}) error
	BuyGood(ctx context.Context, ns constants.Namespace, details map[string]interface{}) error
}

type gameUseCase struct {
	connectionManager *ConnectionManager
	eventProcessor    *EventProcessor
}

func NewGameUseCase(
	connRepo repository.ConnectionRepository,
	eventRepo repository.EventRepository,
	logger *zap.Logger,
) GameUseCase {
	connectionManager := NewConnectionManager(connRepo, logger)
	eventHandlers := NewEventHandlers(connectionManager, logger)
	eventProcessor := NewEventProcessor(eventRepo, eventHandlers, logger, 1000)

	eventProcessor.Start()

	return &gameUseCase{
		connectionManager: connectionManager,
		eventProcessor:    eventProcessor,
	}
}

func (g *gameUseCase) RegisterConnection(ctx context.Context, dto dto.ConnectionDTO) (*entity.Connection, error) {
	return g.connectionManager.Register(ctx, dto)
}

func (g *gameUseCase) UnregisterConnection(ctx context.Context, id string) error {
	return g.connectionManager.Unregister(ctx, id)
}

func (g *gameUseCase) HandleEvent(ctx context.Context, eventDTO dto.EventDTO) error {
	if eventDTO.ID == "" {
		eventDTO.ID = uuid.NewString()
	}

	event := entity.NewEvent(eventDTO.ID, eventDTO.Type, eventDTO.Namespace, eventDTO.Payload)
	return g.eventProcessor.SaveAndQueueEvent(ctx, event)
}

func (g *gameUseCase) createAndHandleEvent(ctx context.Context, eventType constants.EventType, ns constants.Namespace, payload map[string]interface{}) error {
	event := entity.NewEvent(uuid.NewString(), eventType, ns, payload)
	return g.HandleEvent(ctx, dto.EventDTO{
		ID:        event.ID,
		Type:      event.Type,
		Namespace: event.Namespace,
		Payload:   event.Payload,
	})
}

func (g *gameUseCase) BroadcastState(ctx context.Context, ns constants.Namespace, message map[string]interface{}) error {
	return g.createAndHandleEvent(ctx, constants.EventTypeBroadcast, ns, message)
}

func (g *gameUseCase) SendMessage(ctx context.Context, ns constants.Namespace, message map[string]interface{}) error {
	return g.createAndHandleEvent(ctx, constants.EventTypeSendMsg, ns, message)
}

func (g *gameUseCase) BuyGood(ctx context.Context, ns constants.Namespace, details map[string]interface{}) error {
	return g.createAndHandleEvent(ctx, constants.EventTypeBuyGood, ns, details)
}
