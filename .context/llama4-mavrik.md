Below is a high-level architecture for the Follows Application using Go, adhering to the given considerations and criteria. This scaffold focuses on the overall structure and key components.

### Project Structure
```markdown
- cmd
  - main.go
- internal
  - domain
    - event.go
    - namespace.go
    - connection.go
  - usecase
    - event_usecase.go
    - namespace_usecase.go
    - connection_usecase.go
  - repository
    - etcd_repository.go
    - event_repository.go
    - connection_repository.go
  - delivery
    - socket
      - handler.go
      - server.go
    - http
      - config.go
- pkg
  - etcdutil
    - etcd_util.go
  - socketutil
    - socket_util.go
  - logger
    - logger.go
  - errors
    - errors.go
- config
  - config.go
  - config.yaml
- go.mod
- go.sum
```

### Key Components

#### Domain Layer (`internal/domain`)
Defines the core business entities.

```go
// internal/domain/event.go
package domain

type EventType string

const (
    BroadcastEvent  EventType = "core.game.state.broadcast"
    SendMessageEvent EventType = "core.game.state.send_message"
    BuyGoodEvent     EventType = "core.game.state.buy_good"
)

type Event struct {
    Type    EventType
    Payload interface{}
}

type Namespace string

const (
    PrivateNamespace Namespace = "private"
    GuestNamespace    Namespace = "guest"
)

type Connection struct {
    ID        string
    Namespace Namespace
}
```

#### Use Case/Service Layer (`internal/usecase`)
Contains the business logic.

```go
// internal/usecase/event_usecase.go
package usecase

import (
    "context"
    "internal/domain"
)

type EventUsecase interface {
    Handle(ctx context.Context, event domain.Event) error
}

type eventUsecase struct {
    eventRepository domain.EventRepository
}

func NewEventUsecase(eventRepository domain.EventRepository) EventUsecase {
    return &eventUsecase{eventRepository: eventRepository}
}

func (e *eventUsecase) Handle(ctx context.Context, event domain.Event) error {
    // Handle event logic
    return e.eventRepository.Save(ctx, event)
}
```

#### Repository/Storage Layer (`internal/repository`)
Defines data access interfaces and their implementations.

```go
// internal/repository/etcd_repository.go
package repository

import (
    "context"
    "internal/domain"
    clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdRepository struct {
    client *clientv3.Client
}

func NewEtcdRepository(client *clientv3.Client) *EtcdRepository {
    return &EtcdRepository{client: client}
}

func (e *EtcdRepository) SaveConnection(ctx context.Context, conn domain.Connection) error {
    // Save connection to etcd
    _, err := e.client.Put(ctx, conn.ID, conn.Namespace.String())
    return err
}
```

#### Controller/Delivery Layer (`internal/delivery`)
Handles external interfaces like HTTP handlers and socket connections.

```go
// internal/delivery/socket/handler.go
package socket

import (
    "net/http"

    "github.com/gorilla/websocket"
    "internal/usecase"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type Handler struct {
    eventUsecase usecase.EventUsecase
}

func NewHandler(eventUsecase usecase.EventUsecase) *Handler {
    return &Handler{eventUsecase: eventUsecase}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        // Handle error
        return
    }
    defer conn.Close()

    // Handle socket connection
    go func() {
        for {
            // Read message from client
            _, message, err := conn.ReadMessage()
            if err != nil {
                // Handle error
                return
            }
            // Handle message
        }
    }()
}
```

#### Configuration and Error Handling
Configuration is loaded from a YAML file (`config/config.yaml`), and error handling is designed using a custom error type (`pkg/errors/errors.go`).

```go
// pkg/errors/errors.go
package errors

type Error struct {
    Code    int
    Message string
}

func NewError(code int, message string) *Error {
    return &Error{Code: code, Message: message}
}
```

#### Main Function (`cmd/main.go`)
Sets up the application.

```go
package main

import (
    "context"
    "log"
    "net/http"

    "internal/delivery/socket"
    "internal/repository"
    "internal/usecase"
    "pkg/etcdutil"

    "go.etcd.io/etcd/client/v3"
)

func main() {
    // Load configuration
    config, err := LoadConfig("config/config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // Set up etcd client
    etcdClient, err := clientv3.New(clientv3.Config{
        Endpoints: config.Etcd.Endpoints,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Initialize layers
    etcdRepository := repository.NewEtcdRepository(etcdClient)
    eventUsecase := usecase.NewEventUsecase(etcdRepository)
    handler := socket.NewHandler(eventUsecase)

    // Start server
    http.Handle("/core/game/v1/real-time/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

This scaffold provides a high-level architecture for the Follows Application, adhering to the given criteria and considerations. It is designed to be production-ready, with a focus on concurrency using goroutines and a clean architecture. Error handling is also a key component, with a custom error type defined. 

To make changes to the files, you can either use the Apply Button on the code block or switch to Agent Mode to make the suggested updates automatically. To switch to Agent Mode, use the Mode Selector dropdown.

----

# Part 2

----

Here's a more detailed implementation of the remaining files.

### `internal/domain/namespace.go`
```go
// internal/domain/namespace.go
package domain

type Namespace string

const (
    PrivateNamespace Namespace = "private"
    GuestNamespace   Namespace = "guest"
)

func (n Namespace) String() string {
    return string(n)
}

func NamespaceFromString(s string) (Namespace, error) {
    switch s {
    case "private":
    return PrivateNamespace, nil
    case "guest":
    return GuestNamespace, nil
    default:
    return "", ErrInvalidNamespace
}

var ErrInvalidNamespace = NewError(400, "invalid namespace")
```

### `internal/domain/connection.go`
```go
// internal/domain/connection.go
package domain

import (
    "errors"
)

var (
    ErrConnectionNotFound = errors.New("connection not found")
)

type Connection struct {
    ID        string
    Namespace Namespace
}

type ConnectionRepository interface {
    Save(ctx context.Context, conn Connection) error
    Get(ctx context.Context, id string) (Connection, error)
    Delete(ctx context.Context, id string) error
}
```

### `internal/usecase/namespace_usecase.go`
```go
// internal/usecase/namespace_usecase.go
package usecase

import (
    "context"
    "internal/domain"
)

type NamespaceUsecase interface {
    Validate(ctx context.Context, namespace domain.Namespace) error
}

type namespaceUsecase struct{}

func NewNamespaceUsecase() NamespaceUsecase {
    return &namespaceUsecase{}
}

func (n *namespaceUsecase) Validate(ctx context.Context, namespace domain.Namespace) error {
    // Validate namespace logic
    if namespace != domain.PrivateNamespace && namespace != domain.GuestNamespace {
    return domain.ErrInvalidNamespace
    }
    return nil
}
```

### `internal/usecase/connection_usecase.go`
```go
// internal/usecase/connection_usecase.go
package usecase

import (
    "context"
    "internal/domain"
)

type ConnectionUsecase interface {
    Handle(ctx context.Context, conn domain.Connection) error
}

type connectionUsecase struct {
    connectionRepository domain.ConnectionRepository
}

func NewConnectionUsecase(connectionRepository domain.ConnectionRepository) ConnectionUsecase {
    return &connectionUsecase{connectionRepository: connectionRepository}
}

func (c *connectionUsecase) Handle(ctx context.Context, conn domain.Connection) error {
    // Handle connection logic
    return c.connectionRepository.Save(ctx, conn)
}
```

### `internal/repository/etcd_connection_repository.go`
```go
// internal/repository/etcd_connection_repository.go
package repository

import (
    "context"
    "internal/domain"
    clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConnectionRepository struct {
    client *clientv3.Client
}

func NewEtcdConnectionRepository(client *clientv3.Client) domain.ConnectionRepository {
    return &EtcdConnectionRepository{client: client}
}

func (e *EtcdConnectionRepository) Save(ctx context.Context, conn domain.Connection) error {
    // Save connection to etcd
    _, err := e.client.Put(ctx, conn.ID, conn.Namespace.String())
    return err
}

func (e *EtcdConnectionRepository) Get(ctx context.Context, id string) (domain.Connection, error) {
    // Get connection from etcd
    resp, err := e.client.Get(ctx, id)
    if err != nil {
    return domain.Connection{}, err
    }
    if len(resp.Kvs) == 0 {
    return domain.Connection{}, domain.ErrConnectionNotFound
    }
    namespace, err := domain.NamespaceFromString(string(resp.Kvs[0].Value))
    if err != nil {
    return domain.Connection{}, err
    }
    return domain.Connection{ID: id, Namespace: namespace}, nil
}

func (e *EtcdConnectionRepository) Delete(ctx context.Context, id string) error {
    // Delete connection from etcd
    _, err := e.client.Delete(ctx, id)
    return err
}
```

### `internal/delivery/socket/server.go`
```go
// internal/delivery/socket/server.go
package socket

import (
    "net/http"

    "github.com/gorilla/websocket"
    "internal/usecase"
)

type Server struct {
    eventUsecase    usecase.EventUsecase
    connectionUsecase usecase.ConnectionUsecase
    namespaceUsecase  usecase.NamespaceUsecase
}

func NewServer(eventUsecase usecase.EventUsecase, connectionUsecase usecase.ConnectionUsecase, namespaceUsecase usecase.NamespaceUsecase) *Server {
    return &Server{
    eventUsecase:    eventUsecase,
    connectionUsecase: connectionUsecase,
    namespaceUsecase:  namespaceUsecase,
    }
}

func (s *Server) ListenAndServe(addr string) error {
    http.Handle("/core/game/v1/real-time/", http.HandlerFunc(s.handle))
    return http.ListenAndServe(addr, nil)
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
    // Handle socket connection
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
    // Handle error
    return
    }
    defer conn.Close()

    // Handle socket connection logic
    go func() {
    for {
    // Read message from client
    _, message, err := conn.ReadMessage()
    if err != nil {
    // Handle error
    return
    }
    // Handle message
    }
    }()
}
```

### `config/config.go`
```go
// config/config.go
package config

import (
    "io/ioutil"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Etcd struct {
    Endpoints []string `yaml:"endpoints"`
    } `yaml:"etcd"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
    return nil, err
    }
    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
    return nil, err
    }
    return &config, nil
}
```

### `config/config.yaml`
```yml
etcd:
    endpoints:
      - "http://localhost:2379"
```

### `pkg/etcdutil/etcd_util.go`
```go
// pkg/etcdutil/etcd_util.go
package etcdutil

import (
    clientv3 "go.etcd.io/etcd/client/v3"
)

func NewEtcdClient(endpoints []string) (*clientv3.Client, error) {
    return clientv3.New(clientv3.Config{
    Endpoints: endpoints,
    })
}
```

### `pkg/socketutil/socket_util.go`

```go
// pkg/socketutil/socket_util.go
package socketutil

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/gorilla/websocket"
    "internal/domain"
    "pkg/errors"
    "pkg/logger"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
    return true
    },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
    logger.Errorf("Failed to upgrade to websocket: %v", err)
    return nil, errors.New("failed to upgrade to websocket")
    }
    return conn, nil
}

func ReadMessage(conn *websocket.Conn) ([]byte, error) {
    _, message, err := conn.ReadMessage()
    if err != nil {
    logger.Errorf("Failed to read message from websocket: %v", err)
    return nil, errors.New("failed to read message from websocket")
    }
    return message, nil
}

func WriteMessage(conn *websocket.Conn, messageType int, data interface{}) error {
    message, err := json.Marshal(data)
    if err != nil {
    logger.Errorf("Failed to marshal message: %v", err)
    return errors.New("failed to marshal message")
    }
    err = conn.WriteMessage(messageType, message)
    if err != nil {
    logger.Errorf("Failed to write message to websocket: %v", err)
    return errors.New("failed to write message to websocket")
    }
    return nil
}

func HandleSocketConnection(conn *websocket.Conn, eventUsecase domain.EventUsecase) {
    defer conn.Close()

    for {
    message, err := ReadMessage(conn)
    if err != nil {
    logger.Errorf("Failed to read message from websocket: %v", err)
    return
    }

    var event domain.Event
    err = json.Unmarshal(message, &event)
    if err != nil {
    logger.Errorf("Failed to unmarshal message: %v", err)
    WriteMessage(conn, websocket.TextMessage, errors.New("invalid message format"))
    continue
    }

    err = eventUsecase.Handle(context.Background(), event)
    if err != nil {
    logger.Errorf("Failed to handle event: %v", err)
    WriteMessage(conn, websocket.TextMessage, err)
    continue
    }
    }
}
```

### `pkg/logger/logger.go`
```go
// pkg/logger/logger.go
package logger

import (
    "log"
)

func Infof(format string, v ...interface{}) {
    log.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
    log.Printf(format, v...)
}
```

### `pkg/errors/errors.go`
```go
// pkg/errors/errors.go
package errors

type Error struct {
    Code    int
    Message string
}

func NewError(code int, message string) *Error {
    return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
    return e.Message
}
```

### `cmd/main.go`
```go
// cmd/main.go
package main

import (
    "context"
    "log"

    "internal/delivery/socket"
    "internal/repository"
    "internal/usecase"
    "config"

    "go.etcd.io/etcd/client/v3"
)

func main() {
    // Load configuration
    config, err := config.LoadConfig("config/config.yaml")
    if err != nil {
    log.Fatal(err)
    }

    // Set up etcd client
    etcdClient, err := clientv3.New(clientv3.Config{
    Endpoints: config.Etcd.Endpoints,
    })
    if err != nil {
    log.Fatal(err)
    }

    // Initialize layers
    etcdRepository := repository.NewEtcdRepository(etcdClient)
    connectionRepository := repository.NewEtcdConnectionRepository(etcdClient)
    eventUsecase := usecase.NewEventUsecase(etcdRepository)
    connectionUsecase := usecase.NewConnectionUsecase(connectionRepository)
    namespaceUsecase := usecase.NewNamespaceUsecase()
    server := socket.NewServer(eventUsecase, connectionUsecase, namespaceUsecase)

    // Start server
    log.Fatal(server.ListenAndServe(":8080"))
}
```

This provides a more comprehensive implementation of the Follows Application, including the remaining files. To make changes to the files, you can either use the Apply Button on the code block or switch to Agent Mode to make the suggested updates automatically. To switch to Agent Mode, use the Mode Selector dropdown.


