package game

import (
	"context"
	"sync"

	"app/core-game/constants"
	"app/core-game/internal/domain/entity"
	"app/core-game/internal/domain/repository"
	"app/core-game/internal/dto"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConnectionManager struct {
	connRepo    repository.ConnectionRepository
	connections map[string]*entity.Connection
	mu          sync.RWMutex
	logger      *zap.Logger
}

func NewConnectionManager(connRepo repository.ConnectionRepository, logger *zap.Logger) *ConnectionManager {
	return &ConnectionManager{
		connRepo:    connRepo,
		connections: make(map[string]*entity.Connection),
		logger:      logger,
	}
}

func (cm *ConnectionManager) Register(ctx context.Context, dto dto.ConnectionDTO) (*entity.Connection, error) {
	if dto.ID == "" {
		dto.ID = uuid.NewString()
	}

	conn := entity.NewConnection(dto.ID, dto.Namespace)
	if err := cm.connRepo.Save(ctx, conn); err != nil {
		return nil, err
	}

	cm.mu.Lock()
	cm.connections[conn.ID] = conn
	cm.mu.Unlock()

	cm.logger.Info("connection registered",
		zap.String("id", conn.ID),
		zap.String("namespace", string(conn.Namespace)),
	)
	return conn, nil
}

func (cm *ConnectionManager) Unregister(ctx context.Context, id string) error {
	if err := cm.connRepo.Delete(ctx, id); err != nil {
		return err
	}

	cm.mu.Lock()
	delete(cm.connections, id)
	cm.mu.Unlock()

	cm.logger.Info("connection unregistered", zap.String("id", id))
	return nil
}

func (cm *ConnectionManager) GetConnectionsByNamespace(ns constants.Namespace) []*entity.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var connections []*entity.Connection
	for _, conn := range cm.connections {
		if conn.Namespace == ns {
			connections = append(connections, conn)
		}
	}
	return connections
}

func (cm *ConnectionManager) GetConnectionByID(id string) (*entity.Connection, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	conn, exists := cm.connections[id]
	return conn, exists
}

func (cm *ConnectionManager) GetAllConnections() []*entity.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	connections := make([]*entity.Connection, 0, len(cm.connections))
	for _, conn := range cm.connections {
		connections = append(connections, conn)
	}
	return connections
}
