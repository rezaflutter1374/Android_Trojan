# C2 Server Component

## Overview
This is the Command and Control (C2) server component of the Hermes: Spectre Edition project. It's built with Go and provides a WebSocket-based communication channel with infected clients, as well as a REST API for command issuance.

## Features
- WebSocket server for real-time bidirectional communication
- Client management system for tracking connected devices
- REST API for sending commands to specific or all clients
- Lightweight and efficient design using Go's concurrency model

## Requirements
- Go 1.16 or higher
- Gin web framework
- Gorilla WebSocket library

## Installation

1. Install Go from [golang.org](https://golang.org/)
2. Clone the repository
3. Install dependencies:
   ```
   cd c2_server
   go mod download
   ```

## Running the Server

```
go run cmd/server/main.go
```

The server will start on port 8080 by default.

## API Endpoints

### WebSocket Connection
- **URL**: `/ws`
- **Method**: GET
- **Query Parameters**: `id` (optional) - Client identifier

### Send Command
- **URL**: `/api/command`
- **Method**: POST
- **Body**:
  ```json
  {
    "target_id": "client_id_or_all",
    "action": "COMMAND_NAME",
    "payload": "command_parameters"
  }
  ```

## Security Considerations

This is a demonstration project. In a production environment, you would need to implement:

- TLS encryption for all communications
- Authentication for API access
- Origin restrictions for WebSocket connections
- Proper error handling and logging
- Input validation and sanitization