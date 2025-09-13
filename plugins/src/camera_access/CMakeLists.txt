cmake_minimum_required(VERSION 3.10.2)

project(camera_access_plugin)

# Add the plugin library
add_library(cameraaccess SHARED
            camera_access.cpp)

# Include directories (plugin public API)
target_include_directories(cameraaccess PRIVATE
                           ${CMAKE_CURRENT_SOURCE_DIR}/../../include)

# Set output name to match the expected plugin naming convention
set_target_properties(cameraaccess PROPERTIES
                      OUTPUT_NAME "cameraaccess"
                      PREFIX "lib"
                      SUFFIX ".so")