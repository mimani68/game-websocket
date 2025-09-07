package etcd

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"go.uber.org/zap"
)

type EtcdClient struct {
	client *clientv3.Client
	logger *zap.Logger
}

func NewEtcdClient(endpoints []string, dialTimeout time.Duration, logger *zap.Logger) (*EtcdClient, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdClient{
		client: cli,
		logger: logger,
	}, nil
}

func (e *EtcdClient) Close() error {
	return e.client.Close()
}

func (e *EtcdClient) Put(ctx context.Context, key, val string) error {
	_, err := e.client.Put(ctx, key, val)
	if err != nil {
		e.logger.Error("etcd put error", zap.String("key", key), zap.Error(err))
	}
	return err
}

func (e *EtcdClient) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := e.client.Get(ctx, key)
	if err != nil {
		e.logger.Error("etcd get error", zap.String("key", key), zap.Error(err))
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	return resp.Kvs[0].Value, nil
}

func (e *EtcdClient) Delete(ctx context.Context, key string) error {
	_, err := e.client.Delete(ctx, key)
	if err != nil {
		e.logger.Error("etcd delete error", zap.String("key", key), zap.Error(err))
	}
	return err
}
