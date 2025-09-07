package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/core-game/config"
	"app/core-game/internal/delivery/socket"
	"app/core-game/internal/infrastructure/etcd"
	"app/core-game/internal/infrastructure/logger"
	"app/core-game/internal/usecase"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig("config/default.yaml")
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("starting server", zap.Any("config", cfg))

	etcdCli, err := etcd.NewEtcdClient(cfg.Etcd.Endpoints, cfg.EtcdDialTimeout(), log)
	if err != nil {
		log.Fatal("failed to create etcd client", zap.Error(err))
	}
	defer etcdCli.Close()

	connRepo := etcd.NewConnectionRepo(etcdCli)
	eventRepo := etcd.NewEventRepo(etcdCli)

	gameUC := usecase.NewGameUseCase(connRepo, eventRepo, log)

	socketRegistry := socket.NewSocketRegistry(gameUC, log)

	mux := http.NewServeMux()
	socketRegistry.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: mux,
	}

	// Graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		log.Info("shutdown signal received")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error("server shutdown error", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	log.Info("server listening", zap.String("address", server.Addr))
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("server failed", zap.Error(err))
	}

	<-idleConnsClosed
	log.Info("server stopped")
}
