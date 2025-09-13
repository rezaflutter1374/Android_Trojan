cmake_minimum_required(VERSION 3.10.2)

project(core_payload_jni)

# Export compile commands to aid IDEs (e.g., clangd) with correct include paths/flags
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Map ANDROID_ABI to Rust target triple
if(${ANDROID_ABI} STREQUAL "arm64-v8a")
    set(RUST_TARGET "aarch64-linux-android")
elseif(${ANDROID_ABI} STREQUAL "armeabi-v7a")
    set(RUST_TARGET "armv7-linux-androideabi")
elseif(${ANDROID_ABI} STREQUAL "x86_64")
    set(RUST_TARGET "x86_64-linux-android")
else()
    message(FATAL_ERROR "Unsupported ANDROID_ABI: ${ANDROID_ABI}")
endif()

# Set the path to the Rust library for this ABI
set(RUST_LIB_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../jniLibs/${ANDROID_ABI}")

# Find the Rust library
find_library(RUST_LIB core_payload
             PATHS ${RUST_LIB_DIR}
             NO_CMAKE_FIND_ROOT_PATH
             REQUIRED)

# Add the JNI wrapper library
add_library(core_payload_jni SHARED
            core_payload_jni.cpp
            # jni_wrapper.cpp is intentionally excluded to avoid duplicate symbol errors
            )

# Include directories
# Always include our current directory
target_include_directories(core_payload_jni PRIVATE
                           ${CMAKE_CURRENT_SOURCE_DIR})

# Help IDEs (like clangd) locate Android NDK JNI headers to prevent false 'jni.h not found' diagnostics.
# This does not affect real builds, as the Android toolchain already provides correct include paths.
if(DEFINED ANDROID_NDK)
    file(GLOB NDK_SYSROOT_INCLUDE_DIR "${ANDROID_NDK}/toolchains/llvm/prebuilt/*/sysroot/usr/include")
    if(NDK_SYSROOT_INCLUDE_DIR)
        target_include_directories(core_payload_jni PRIVATE ${NDK_SYSROOT_INCLUDE_DIR})
    endif()
endif()

# Link against the Rust library and Android log library
target_link_libraries(core_payload_jni
                      ${RUST_LIB}
                      android
                      log)