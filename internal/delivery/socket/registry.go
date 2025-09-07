package socket

import (
	"net/http"

	"app/core-game/internal/usecase"

	"go.uber.org/zap"
)

type SocketRegistry struct {
	handler *SocketHandler
}

func NewSocketRegistry(gameUC usecase.GameUseCase, logger *zap.Logger) *SocketRegistry {
	return &SocketRegistry{
		handler: NewSocketHandler(gameUC, logger),
	}
}

func (r *SocketRegistry) RegisterRoutes(mux *http.ServeMux) {
	// Register websocket handlers for both namespaces
	mux.Handle("/core/game/v1/real-time/private", r.handler)
	mux.Handle("/core/game/v1/real-time/guest", r.handler)
}
