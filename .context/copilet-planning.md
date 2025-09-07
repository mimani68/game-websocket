# Planning Document for Real-Time Game Service in Go

This document outlines a production-ready, highly concurrent Go service using clean architecture, WebSocket communication, etcd for state storage, and generic, interface-based design. It covers project structure, file mappings, architecture layers, protocols, modules/components, configuration, and error-handling strategies.

---

## Project Structure

```
/cmd/server
    └── main.go
/internal
    /config
        └── config.go
    /domain
        ├── event.go
        ├── game_state.go
        └── errors.go
    /usecase
        ├── game_broadcast.go
        ├── send_message.go
        ├── buy_good.go
        └── interfaces.go
    /repository
        ├── etcd_client.go
        ├── game_state_repo.go
        └── interfaces.go
    /controller
        ├── websocket_handler.go
        └── router.go
    /transport
        ├── websocket.go
        └── protocols.go
    /socket
        ├── connection_manager.go
        └── namespaces.go
    /event
        ├── dispatcher.go
        └── subscriber.go
    /store
        └── etcd_store.go
    /logger
        └── logger.go
/internal/utils
    └── generics.go
```

---

## File Map

- `/cmd/server/main.go`: Application bootstrap, configuration loading, dependency injection.
- `/internal/config/config.go`: Reads environment or file-based settings (no hardcodes).
- `/internal/domain/`: Core entities and error definitions.
- `/internal/usecase/`: Business logic implementations and generic interfaces.
- `/internal/repository/`: etcd-backed storage implementations and repository interfaces.
- `/internal/controller/`: WebSocket controllers mapping incoming frames to use cases.
- `/internal/transport/`: Low-level WebSocket protocol handling and message framing.
- `/internal/socket/`: Connection pooling, namespace segregation, concurrent write broadcasting.
- `/internal/event/`: Event dispatcher, registering subscribers, firing business events.
- `/internal/store/etcd_store.go`: Wrapper over etcd client for key-value state operations.
- `/internal/logger/logger.go`: Structured logging with levels and context.
- `/internal/utils/generics.go`: Generic helpers (e.g., concurrent-safe maps, channels).

---

## Architecture Layers

### Domain Layer
Defines core types and errors:
- GameStateEvent (type parameterized over payload)
- Message, Purchase
- Custom error types for business conditions  

### Use Case Layer
Implements business logic:
- BroadcastStateUseCase (fires `core.game.state.broadcast`)
- SendMessageUseCase (fires `core.game.state.send_message`)
- BuyGoodUseCase (fires `core.game.state.buy_good`)
- All use cases expose generic interfaces:  
  ```go
  type UseCase[T any] interface {
      Execute(ctx context.Context, input T) error
  }
  ```  

### Repository Layer
Interfaces for data access:
- GameStateRepository interface (Get, Set, Watch)
- ConnectionRepository interface (Register, Deregister, List)
Concrete implementations leverage `clientv3` etcd with leases and watchers.

### Controller/Delivery Layer
Handles external I/O:
- WebSocketHandler registers routes under `/core/game/v1/real-time/{private|public}`
- Upgrades HTTP to WebSocket, maintains namespaces via ConnectionManager
- Maps inbound messages to UseCase.Execute calls  

---

## Protocols

- WebSocket for client connections (RFC 6455)
- HTTP/1.1 upgrade endpoint at `/core/game/v1/real-time/{namespace}`
- gRPC or etcd clientv3 API for key-value operations
- JSON over WebSocket frames with envelope containing `eventType` and `payload`

---

## Modules & Components

- Configuration Module: centralized environment and file-based settings
- Logging Module: structured JSON logger with context propagation
- Event Module: dispatch and subscribe to typed channels via Go generics
- Socket Module: concurrent-safe connection registry, namespace segmentation
- UseCase Module: business logic orchestrating repositories and events
- Repository Module: etcd-backed storage with retry and backoff
- Controller Module: WebSocket HTTP handlers, input validation, error mapping
- Utils Module: generic concurrency primitives (e.g., `Map[K,V]`, `Pool[T]`)

---

## Configuration List

Load via environment variables or configuration file:
- ETCD_ENDPOINTS (comma-separated host:port)
- ETCD_DIAL_TIMEOUT (duration)
- SOCKET_READ_TIMEOUT (duration)
- SOCKET_WRITE_TIMEOUT (duration)
- LOG_LEVEL (debug/info/warn/error)
- NAMESPACE_PRIVATE (e.g., `private`)
- NAMESPACE_GUEST (e.g., `guest`)
- API_PATH_PREFIX (`/core/game/v1/real-time`)
- MAX_CONNECTIONS_PER_NAMESPACE (int)

---

## Concurrency & Goroutine Strategy

- ConnectionManager uses sync.Map for thread-safe registry.
- Event dispatcher loops fire events; each subscriber runs in its own goroutine.
- UseCase.Execute spawns background goroutines for long-running tasks (e.g., etcd watches) and returns immediately if appropriate.
- Worker pools for processing high-throughput buy_good requests.
- Context and cancellation propagated across goroutines to ensure graceful shutdown.

---

## Error Handling

- Domain errors defined as typed errors for use-case validation.
- Wrap repository errors with `fmt.Errorf("%w: context", err)` for root-cause inspection.
- Controller maps errors to WebSocket close codes or JSON error frames.
- Retries with exponential backoff in etcd operations; fail-fast on configuration issues.
- Panic recovery middleware around goroutine entry points to prevent crashes.

---

Further considerations include setting up CI/CD pipelines with linters, security scans, and integration tests to validate WebSocket workflows and etcd state persistence. Continuous profiling and monitoring (e.g., Prometheus, Jaeger) should be integrated to track goroutine usage, latency, and error rates.