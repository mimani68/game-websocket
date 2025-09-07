package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"app/core-game/internal/domain/entity"
	"app/core-game/internal/domain/repository"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const connectionPrefix = "core/game/connection/"

type ConnectionRepoImpl struct {
	etcd *EtcdClient
}

func NewConnectionRepo(etcd *EtcdClient) repository.ConnectionRepository {
	return &ConnectionRepoImpl{etcd: etcd}
}

func (r *ConnectionRepoImpl) Save(ctx context.Context, conn *entity.Connection) error {
	key := connectionPrefix + conn.ID
	data, err := json.Marshal(conn)
	if err != nil {
		return err
	}
	return r.etcd.Put(ctx, key, string(data))
}

func (r *ConnectionRepoImpl) Delete(ctx context.Context, id string) error {
	key := connectionPrefix + id
	return r.etcd.Delete(ctx, key)
}

func (r *ConnectionRepoImpl) GetByID(ctx context.Context, id string) (*entity.Connection, error) {
	key := connectionPrefix + id
	val, err := r.etcd.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("connection not found")
	}
	var conn entity.Connection
	err = json.Unmarshal(val, &conn)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (r *ConnectionRepoImpl) ListByNamespace(ctx context.Context, ns entity.Namespace) ([]*entity.Connection, error) {
	resp, err := r.etcd.client.Get(ctx, connectionPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	conns := make([]*entity.Connection, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var conn entity.Connection
		if err := json.Unmarshal(kv.Value, &conn); err != nil {
			continue // skip bad data
		}
		if conn.Namespace == ns {
			conns = append(conns, &conn)
		}
	}
	return conns, nil
}
