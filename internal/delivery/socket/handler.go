package socket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"app/core-game/constants"
	"app/core-game/internal/dto"
	"app/core-game/internal/pkg/validator"
	"app/core-game/internal/usecase/game"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type SocketHandler struct {
	gameUseCase game.GameUseCase
	logger      *zap.Logger

	// Map connection ID to websocket connection
	connections sync.Map // map[string]*websocket.Conn
}

func NewSocketHandler(gameUC game.GameUseCase, logger *zap.Logger) *SocketHandler {
	return &SocketHandler{
		gameUseCase: gameUC,
		logger:      logger,
	}
}

func (h *SocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract namespace from URL path
	ns, err := h.extractNamespace(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid namespace", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("failed to upgrade websocket", zap.Error(err))
		http.Error(w, "failed to upgrade websocket", http.StatusInternalServerError)
		return
	}

	// Register connection
	connID := r.URL.Query().Get("id")
	if connID == "" {
		connID = generateConnectionID()
	}
	connEntity, err := h.gameUseCase.RegisterConnection(r.Context(), dto.ConnectionDTO{
		ID:        connID,
		Namespace: ns,
	})
	fmt.Println(connEntity)
	if err != nil {
		h.logger.Error("failed to register connection", zap.Error(err))
		ws.Close()
		return
	}

	h.connections.Store(connID, ws)
	defer func() {
		h.connections.Delete(connID)
		h.gameUseCase.UnregisterConnection(r.Context(), connID)
		ws.Close()
		h.logger.Info("connection closed", zap.String("conn_id", connID))
	}()

	h.logger.Info("new connection established", zap.String("conn_id", connID), zap.String("namespace", string(ns)))

	// Read loop
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			h.logger.Info("read message error or connection closed", zap.Error(err), zap.String("conn_id", connID))
			break
		}

		// Handle incoming message as event
		var evt dto.EventDTO
		if err := json.Unmarshal(msg, &evt); err != nil {
			h.logger.Warn("invalid message format", zap.Error(err), zap.String("conn_id", connID))
			continue
		}

		// Validate event
		if err := validator.ValidateEventType(string(evt.Type)); err != nil {
			h.logger.Warn("invalid event type", zap.Error(err), zap.String("conn_id", connID))
			continue
		}
		if err := validator.ValidateNamespace(string(evt.Namespace)); err != nil {
			h.logger.Warn("invalid event namespace", zap.Error(err), zap.String("conn_id", connID))
			continue
		}

		// Set event namespace to connection namespace forcibly for security
		evt.Namespace = ns

		// Process event async
		go func(e dto.EventDTO) {
			if err := h.gameUseCase.HandleEvent(context.Background(), e); err != nil {
				h.logger.Error("failed to handle event", zap.Error(err), zap.String("conn_id", connID))
			}
		}(evt)
	}
}

func (h *SocketHandler) extractNamespace(path string) (constants.Namespace, error) {
	// Expected path:  /core/real-time/game/kiako/v1
	parts := strings.Split(path, "/")
	if len(parts) < 6 {
		return "", errors.New("invalid path")
	}
	ns := parts[4]
	if ns != string(constants.NamespacePrivate) && ns != string(constants.NamespaceGuest) {
		return "", errors.New("invalid namespace")
	}
	return constants.Namespace(ns), nil
}

func generateConnectionID() string {
	return uuidNew()
}

// uuidNew helper to generate UUID string without external import here
func uuidNew() string {
	// Using random from crypto/rand for UUID v4 generation
	// For brevity use google/uuid package:
	return uuidNewString()
}

func uuidNewString() string {
	// Minimal wrapper, add import if necessary:
	// github.com/google/uuid
	// Import is already in usecase so safe to use here
	return "conn-" + entityNamespaceRandomID()
}

func entityNamespaceRandomID() string {
	// fallback unique string generator
	return uuid.NewString()
}
