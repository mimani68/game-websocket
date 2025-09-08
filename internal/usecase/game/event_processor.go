package game

import (
	"context"

	"app/core-game/internal/domain/entity"
	"app/core-game/internal/domain/repository"
	customError "app/core-game/internal/error"

	"go.uber.org/zap"
)

type EventProcessor struct {
	eventRepo repository.EventRepository
	eventChan chan *entity.Event
	logger    *zap.Logger
	handlers  *EventHandlers
}

func NewEventProcessor(eventRepo repository.EventRepository, handlers *EventHandlers, logger *zap.Logger, bufferSize int) *EventProcessor {
	return &EventProcessor{
		eventRepo: eventRepo,
		eventChan: make(chan *entity.Event, bufferSize),
		logger:    logger,
		handlers:  handlers,
	}
}

func (ep *EventProcessor) Start() {
	go ep.processEvents()
}

func (ep *EventProcessor) processEvents() {
	for event := range ep.eventChan {
		go ep.handlers.ProcessEvent(event)
	}
}

func (ep *EventProcessor) SaveAndQueueEvent(ctx context.Context, event *entity.Event) error {
	if err := ep.eventRepo.Save(ctx, event); err != nil {
		ep.logger.Error("failed to save event", zap.Error(err))
		return err
	}

	select {
	case ep.eventChan <- event:
		return nil
	default:
		ep.logger.Warn("event channel full, dropping event", zap.String("event_id", event.ID))
		return customError.ErrEventChannelFull
	}
}

func (ep *EventProcessor) Stop() {
	close(ep.eventChan)
}
