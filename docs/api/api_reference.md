# API Reference

## C2 Server API

### Authentication
- Endpoint: `/api/auth`
- Method: POST
- Description: Authenticate clients and issue tokens

### Command Execution
- Endpoint: `/api/commands`
- Method: POST
- Description: Send commands to connected clients

### Results Collection
- Endpoint: `/api/results`
- Method: GET
- Description: Retrieve command execution results

### WebSocket Connection
- Endpoint: `/ws`
- Description: Bidirectional communication channel

## Plugin API

### Plugin Interface
- Function: `initialize()`
- Description: Initialize the plugin

### Command Execution
- Function: `execute_command(command_data)`
- Description: Execute a specific command

### Resource Cleanup
- Function: `cleanup()`
- Description: Release resources and perform cleanup

## Core-Plugin Interface

### Plugin Loading
- Function: `load_plugin(plugin_path)`
- Description: Dynamically load a plugin

### Plugin Execution
- Function: `execute_plugin_command(plugin_id, command_data)`
- Description: Execute a command on a specific plugin

### Plugin Management
- Function: `list_plugins()`
- Description: List all available plugins