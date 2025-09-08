# Core-Game Real-Time Server

This project is a real-time game server implemented in Go. It leverages WebSocket connections for real-time communication between clients, uses etcd as a distributed storage backend, and supports event-driven game logic.

***

## Features

- WebSocket-based real-time communication
- Namespace-based connection management
- Event types including broadcast, send message, and buy good
- Connection and event persistence with etcd
- Graceful server shutdown
- Structured logging with Uber Zap
- Configurable via YAML

***

## Requirements

- Go 1.23+
- etcd cluster (local or remote)
- Dependencies managed via Go modules

***

## Project Structure

- **cmd/server**: Application entry point and HTTP server setup
- **config**: Configuration loading and defaults (`default.yaml` included)
- **internal/core-game**:
  - **domain/entity**: Core entities like `Connection` and `Event`
  - **domain/repository**: Repository interfaces for persistence
  - **infrastructure/etcd**: etcd client and etcd-based repository implementations
  - **infrastructure/logger**: Zap logger setup
  - **pkg/validator**: Validation helpers for namespaces, event types, and IDs
  - **usecase/game**: Business logic including connection management and event processing
  - **delivery/socket**: WebSocket handlers and routing
- **constants**: Shared constants for event types and namespaces
- **dto**: Data Transfer Objects for communication between layers

***

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd core-game
   ```

2. Build the server:
   ```bash
   make build
   ```

3. Run the server:
   ```bash
   make run
   ```

***

## Configuration

Server configuration is located in `config/default.yaml`. You may modify this YAML file to configure:

- Server port (default: 8080)
- etcd endpoints and dial timeout
- Log level (e.g., debug, info)

Example:

```yaml
server:
  port: 8080

etcd:
  endpoints:
    - "localhost:2379"
  dial_timeout_seconds: 5

log:
  level: "debug"
```

***

## Usage

- The server listens for WebSocket connections on the defined server port.
- Clients connect specifying a namespace (e.g., `/core/real-time/game/kiako/v1`) and optionally a connection ID.
- Supported event types:
  - `core.game.state.broadcast`
  - `core.game.state.send_message`
  - `core.game.state.buy_good`
- Events are validated, persisted, and processed asynchronously.
- The server supports graceful shutdown on system signals.

***

## Development

- Code is formatted using `go fmt`.
- Tests (if any) can be added under `internal` packages following Go testing conventions.
- Use `go run ./cmd/server` or `make run` to start locally.
- Debug configurations are available for VSCode (`.vscode/launch.json`).

***

## Dependencies

- `github.com/gorilla/websocket` - WebSocket implementation
- `go.etcd.io/etcd/client/v3` - etcd client for coordination and storage
- `go.uber.org/zap` - Structured logging
- `gopkg.in/yaml.v3` - YAML configuration parsing
- `github.com/google/uuid` - UUID generation

***

## License


***

If you have any questions or want to contribute, please open an issue or PR.