# CMake build configuration for C++ plugins
cmake_minimum_required(VERSION 3.10.2)

project(hermes_plugins)

# Set Android specific options if building for Android
if(ANDROID)
    set(CMAKE_POSITION_INDEPENDENT_CODE ON)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fPIC")
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -fPIC")
endif()

# Note: Plugins use simple string manipulation for JSON-like functionality
# No external JSON library dependency needed

# Add subdirectories for each plugin
add_subdirectory(src/camera_access)
add_subdirectory(src/takephoto)

# You can add more plugins here as they are developed
# add_subdirectory(src/location_access)
# add_subdirectory(src/contacts_access)
# etc.