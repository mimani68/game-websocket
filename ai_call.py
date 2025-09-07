import requests

url = "https://api.avalapis.ir/v1/chat/completions"
content = """
Generate the full codebase for follows application

Specification:
- Application goal: Real-time socket game server

Objects:
- Transmit client-server connection 
- Connecting player via low-latency connection protocol

Considerations
1- Answer must be High-level architecture 
2- No pseudocode, placeholder or mock response
3- Don't return fake and invalid response
4- It should be production ready
5- High level of concurrency using goroutine
6- Answer should not be structure or scaffold only
7- All logics have to be implemented

Application architecture
1- Three event have to be fire 
  * core.game.state.broadcast
  * core.game.state.send_message
  * core.game.state.buy_good
2- Two namespaces have to be define 
  * private
  * guest
3- socket general path prefix is `/core/game/v1/real-time/{private|public}`
  
Criteria
1- Stack must be golang
2- Using generic is mandate
3- Object oriented programming over procedural
4- Use socket for connection to client
5- Store the state of event and connection inside etcd key value store
6- No hard coded configuration
7- Clean architecture must be implement (the Domain layer (core business entities), the Use Case/Service layer (business logic), the Repository/Storage layer (data access interfaces), and the Controller/Delivery layer (external interfaces like HTTP handlers))
8- Interface-based design for repositories and use cases

Error handling
1- Strong error handling must be design


Directory tree first, then every file listed in full.

```
.
├── cmd
│   └── server
│       └── main.go
├── config
│   ├── config.go
│   └── default.yaml
├── internal
│   ├── delivery
│   │   └── socket
│   │       ├── handler.go
│   │       ├── mapper.go
│   │       └── registry.go
│   ├── domain
│   │   ├── entity
│   │   │   ├── connection.go
│   │   │   └── event.go
│   │   └── repository
│   │       ├── connection_repository.go
│   │       └── event_repository.go
│   ├── infrastructure
│   │   ├── etcd
│   │   │   ├── client.go
│   │   │   ├── connection_repo_impl.go
│   │   │   └── event_repo_impl.go
│   │   └── logger
│   │       └── zap.go
│   ├── usecase
│   │   ├── abstract_usecase.go
│   │   ├── game_usecase.go
│   │   └── dto.go
│   └── pkg
│       ├── generic
│       │   └── optional.go
│       │   └── mandate.go
│       └── validator
│           └── validator.go
├── go.mod
└── Makefile
```
"""

headers = {
    "Content-Type": "application/json", 
    "Authorization": "Bearer aa-J2xvjBVso3kyot6X14lcMMebjrCdPr6mzS1nzEosUWu6u0ph"
}

data = {
    "model": "deepseek-coder",
    "messages": [{
        "role": "user",
        "content": content
    }]
}

response = requests.post(url, headers=headers, json=data)

print(response.text)
