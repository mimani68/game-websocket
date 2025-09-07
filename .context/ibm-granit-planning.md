### AI-Generated Code Planning Document

Below is a structured plan to develop a Golang application with high concurrency utilizing WebSockets, adhering to the provided requirements:

#### Project Structure

    app/
    │
    ├── core/
    │   ├── domain/
    │   │   └── (Business entities interfaces)
    │   ├── use_cases/
    │   │   ├── broadcast_game_state.go
    │   │   ├── send_game_message.go
    │   │   └── buy_good.go
    │   └── repositories/
    │       ├── game_state_repository.go
    │       └── etcd_game_state_repository.go
    │
    ├── controller/
    │   ├── websocket/
    │   │   ├── websocket_handler.go
    │   │   └── connection_manager.go
    │   └── http/
    │       └── handler.go
    │
    ├── infra/
    │   ├── etcd/
    │   │   └── etcd_client.go
    │   └── config/
    │       └── config_loader.go
    │
    ├── utils/
    │   └── error_handling.go
    │   └── logging.go
    │
    ├── main.go
    └── configs/
        └── app. tomlfile 
    

#### File Map

*   **`core/`**: Business logic.
    
    *   **`domain/`**: Domain entities (interfaces) and value objects.
    *   **`use_cases/`**: Use cases encapsulating use case logic.
    *   **`repositories/`**: Repository interfaces and etcd implementations.
*   **`controller/`**: External interface handling.
    
    *   **`websocket/`**: WebSocket handler and connection manager.
    *   **`http/`**: HTTP handlers (if applicable).
*   **`infra/`**: Infrastructure implementations.
    
    *   **`etcd/`**: etcd client interactions.
    *   **`config/`**: Configuration loading.
*   **`utils/`**: Utilities, error handling, logging.
    

#### Architecture

*   **Clean Architecture**:
    *   **Domain Layer**: Business logic (entities, values).
    *   **Use Case Layer**: Realizes business rules (use cases).
    *   **Repository Layer**: Abstractions for data accessibility (etcd).
    *   **Controller Layer**: External interfaces (WebSocket for real-time).

#### Suitable Protocols

*   **WebSockets**: For real-time, bidirectional communication, matching the given requirement.

#### Modules and Components

1.  **Domain Module**:
    
    *   Classes: `GameState`, `Player`, `Good`, etc.
    *   Interfaces: Define the behaviour of entities like `GoodService` interface.
2.  **Use Case Module**:
    
    *   Use Cases: `BroadcastGameStateUseCase`, `SendMessageUseCase`, `BuyGoodUseCase`.
    *   Coordinate between repositories and apply necessary business logic.
3.  **Repository Module**:
    
    *   Interfaces: `GameStateRepository` with methods `Get`, `Update`, `Delete`.
    *   Implementations: `EtcdGameStateRepository` handling etcd interactions.
4.  **Controller Module**:
    
    *   WebSocket Controller: Manages connection lifecycle and dispatches messages to use cases.
    *   Connection Manager manages goroutines for concurrent WebSocket connections ensuring high concurrency.
5.  **Infra Module**:
    
    *   etcd Components: `EtcdClient`, wraps etcd interactions.
    *   Configuration Loader: `ConfigLoader` to load and manage configurations from environment variables/files.

#### Configuration List

*   Etcd as the primary key-value store for:
    *   Global game state persistence.
    *   Player session management.
    *   Active WebSocket connection states.
*   Configuration handled via `ConfigLoader` in `infra/config/`.

#### Goroutine Usage for Concurrency

*   WebSocket connections will each run in their goroutine within the `controller/websocket/` module leveraging standard Golang's goroutine management (`sync`/`waitgroup` for managing them efficiently).

#### Error Handling

*   Define domain-specific error types for comprehensive error handling.
*   Centralized logging with `logrus` or `zap` for detailed logging during development and production.
*   Error handling routines to centralize logging, stack traces, and appropriate HTTP response codes for external interfaces.

### Implementation Notes

1.  Utilize dependency injection throughout to adhere to the Dependency Inversion Principle (DIP) and facilitate testing.
2.  Implement a pub-sub pattern via etcd for state management and event broadcasting.
3.  Employ context for managing timeouts and cancellations to support resilient communication.
4.  Consider using standard Golang struct tags for efficient serialization between domain objects and etcd keys.
