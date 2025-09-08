package game

import (
	"app/core-game/constants"
	"app/core-game/internal/domain/entity"
	customError "app/core-game/internal/error"

	"go.uber.org/zap"
)

type EventHandlers struct {
	connectionManager *ConnectionManager
	logger            *zap.Logger
}

func NewEventHandlers(connectionManager *ConnectionManager, logger *zap.Logger) *EventHandlers {
	return &EventHandlers{
		connectionManager: connectionManager,
		logger:            logger,
	}
}

func (eh *EventHandlers) ProcessEvent(event *entity.Event) {
	switch event.Type {
	case constants.EventTypeBroadcast:
		eh.handleBroadcast(event)
	case constants.EventTypeSendMsg:
		eh.handleSendMessage(event)
	case constants.EventTypeBuyGood:
		eh.handleBuyGood(event)
	default:
		eh.logger.Warn("unknown event type",
			zap.String("event_id", event.ID),
			zap.String("type", string(event.Type)),
		)
	}
}

func (eh *EventHandlers) handleBroadcast(event *entity.Event) {
	eh.logger.Debug("processing broadcast event", zap.String("event_id", event.ID))
	connections := eh.connectionManager.GetConnectionsByNamespace(event.Namespace)

	for _, conn := range connections {
		eh.logger.Info("broadcast message to connection",
			zap.String("conn_id", conn.ID),
			zap.Any("payload", event.Payload),
		)
	}
}

func (eh *EventHandlers) handleSendMessage(event *entity.Event) {
	eh.logger.Debug("processing send_message event", zap.String("event_id", event.ID))

	targetIDs, err := extractTargetIDs(event.Payload)
	if err != nil {
		eh.logger.Warn("invalid send_message payload", zap.Error(err))
		return
	}

	for _, targetID := range targetIDs {
		if conn, exists := eh.connectionManager.GetConnectionByID(targetID); exists {
			if conn.Namespace == event.Namespace {
				eh.logger.Info("send message to connection",
					zap.String("conn_id", conn.ID),
					zap.Any("payload", event.Payload),
				)
			}
		}
	}
}

func (eh *EventHandlers) handleBuyGood(event *entity.Event) {
	eh.logger.Info("processing buy_good",
		zap.String("namespace", string(event.Namespace)),
		zap.Any("payload", event.Payload),
	)

	// Business logic for buying goods
	purchaseUpdate := map[string]interface{}{
		"type":    "purchase_update",
		"payload": event.Payload,
	}

	// Simulate broadcast to namespace
	connections := eh.connectionManager.GetConnectionsByNamespace(event.Namespace)
	for _, conn := range connections {
		eh.logger.Info("broadcast purchase update to connection",
			zap.String("conn_id", conn.ID),
			zap.Any("payload", purchaseUpdate),
		)
	}
}

// extractTargetIDs is kept here since it's specific to event handling
func extractTargetIDs(payload map[string]interface{}) ([]string, error) {
	targetIDsRaw, ok := payload["target_ids"]
	if !ok {
		return nil, customError.ErrInvalidTargetIDs
	}

	targetIDsInterface, ok := targetIDsRaw.([]interface{})
	if !ok {
		return nil, customError.ErrInvalidTargetIDs
	}

	var targetIDs []string
	for _, tid := range targetIDsInterface {
		if idStr, ok := tid.(string); ok {
			targetIDs = append(targetIDs, idStr)
		}
	}

	return targetIDs, nil
}
