# Hermes: Spectre Edition - Architecture Overview

## Components

### 1. C2 Server (Go)
- Command and control server that manages communications
- Sends commands and receives results
- Built with Go for lightweight, fast network operations

### 2. Dropper App (Android - Kotlin/Java)
- Deceptive application installed by the user
- Executes the core payload
- Manages permissions

### 3. Core Payload (Rust)
- ELF binary executed by the Dropper
- Communicates with the C2 server
- Dynamically loads plugins
- Executes commands

### 4. Plugins (C++ with NDK)
- Shared libraries (.so) providing specific capabilities
- Examples: audio recording, photography, location tracking
- Loaded by the Rust core

## Data Flow

```
[User] <–> [C2 Server UI] <–> [C2 Go Backend] <–> [Internet] <–> [Core Rust Payload] <–> [C++ Plugins]
```

## Communication Interfaces

### C2 Protocol
- Defined using Protocol Buffers
- Secure, encrypted communication channel
- Command structure and response format

### Plugin Protocol
- C++ interface for plugin development
- Standard API for all plugins
- Dynamic loading mechanism

### Dropper-Core Bridge
- JNI interface between Kotlin/Java and Rust
- Handles initialization and communication

## Clean Architecture

Each component follows clean architecture principles:

### C2 Server
- Domain layer: Core business logic
- Use cases: Application-specific business rules
- Repositories: Data access abstractions
- Services: External service integrations
- API: External interface

### Dropper App
- Presentation: UI components
- Domain: Business logic
- Data: Data sources and repositories
- DI: Dependency injection
- Native: JNI bridges

### Core Payload
- Communication: C2 client implementation
- Plugin: Plugin management
- Utils: Utility functions
- FFI: Foreign function interfaces

### Plugins
- Include: Public API headers
- Src: Implementation code
- Bindings: JNI and Rust bindings
- Tests: Unit and integration tests