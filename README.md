# Hermes: Spectre Edition

A multi-language project demonstrating clean architecture principles across Go, Kotlin/Java, Rust, and C++.

## Project Structure

```
Hermes_Spectre_Edition/
├── c2_server/           # Go backend server
├── dropper_app/         # Android Kotlin/Java app
├── core_payload/        # Rust core library
├── plugins/             # C++ plugin modules
├── interfaces/          # Communication interfaces
└── docs/                # Documentation
```

## Components

### C2 Server (Go)

Command and control server built with Go, providing a backend for managing communications with deployed payloads.

### Dropper App (Android - Kotlin/Java)

Android application that serves as the delivery mechanism for the core payload. It extracts and executes the Rust core payload.

### Core Payload (Rust)

Core functionality implemented in Rust, providing a secure and efficient execution environment. It communicates with the C2 server via WebSockets and dynamically loads plugins based on commands received.

### Plugins (C++ with NDK)

Extendable plugin system built with C++ and the Android NDK, allowing for modular functionality. Each plugin must export a C function with the signature `char* execute_action(const char* payload)`.

## Communication Protocol

### C2 to Core Payload

```json
{
  "action": "PLUGIN_NAME",
  "payload": "json_encoded_parameters"
}
```

### Plugin Result Format

Plugins should return results as JSON-encoded strings that will be forwarded to the C2 server.

## Build Instructions

Use the provided build script to compile all components:

```bash
./build_android.sh
```

This will:
1. Build the Rust core payload for ARM32 and ARM64
2. Copy the libraries to the Android project
3. Build the Android dropper app

## Documentation

Comprehensive documentation is available in the `docs/` directory:

- Architecture overview
- API reference
- Deployment guide
- Development guide

## Getting Started

See the [Development Guide](docs/development/development_guide.md) for setup instructions and the [Deployment Guide](docs/deployment/deployment_guide.md) for deployment instructions.