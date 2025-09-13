cmake_minimum_required(VERSION 3.10.2)

project(takephoto_plugin)

# Add the plugin library
add_library(takephoto SHARED
            takephoto.cpp)

# Include directories (plugin public API)
target_include_directories(takephoto PRIVATE
                           ${CMAKE_CURRENT_SOURCE_DIR}/../../include)

# Link required Android NDK libs when building for Android
if(ANDROID)
    target_link_libraries(takephoto PRIVATE log android camera2ndk)
endif()

# Set output name to match the expected plugin naming convention
set_target_properties(takephoto PROPERTIES
                      OUTPUT_NAME "takephoto"
                      PREFIX "lib"
                      SUFFIX ".so")