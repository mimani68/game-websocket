package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"

	"app/core-game/internal/domain/entity"
	"app/core-game/internal/domain/repository"
)

const eventPrefix = "core/game/event/"

type EventRepoImpl struct {
	etcd *EtcdClient
}

func NewEventRepo(etcd *EtcdClient) repository.EventRepository {
	return &EventRepoImpl{etcd: etcd}
}

func (r *EventRepoImpl) Save(ctx context.Context, event *entity.Event) error {
	key := eventPrefix + event.ID
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return r.etcd.Put(ctx, key, string(data))
}

func (r *EventRepoImpl) GetByID(ctx context.Context, id string) (*entity.Event, error) {
	key := eventPrefix + id
	val, err := r.etcd.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("event not found")
	}
	var event entity.Event
	if err := json.Unmarshal(val, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepoImpl) ListByNamespace(ctx context.Context, ns entity.Namespace) ([]*entity.Event, error) {
	resp, err := r.etcd.client.Get(ctx, eventPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	events := make([]*entity.Event, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var event entity.Event
		if err := json.Unmarshal(kv.Value, &event); err != nil {
			continue // skip invalid data
		}
		if event.Namespace == ns {
			events = append(events, &event)
		}
	}
	return events, nil
}
