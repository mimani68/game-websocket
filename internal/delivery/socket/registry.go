package socket

import (
	"net/http"

	"app/core-game/constants"
	"app/core-game/internal/usecase/game"

	"go.uber.org/zap"
)

type SocketRegistry struct {
	handler *SocketHandler
}

func NewSocketRegistry(gameUC game.GameUseCase, logger *zap.Logger) *SocketRegistry {
	return &SocketRegistry{
		handler: NewSocketHandler(gameUC, logger),
	}
}

func (r *SocketRegistry) RegisterRoutes(mux *http.ServeMux) {
	// Register websocket handlers for both namespaces
	mux.Handle(string(constants.KiakoV1), r.handler)
}
