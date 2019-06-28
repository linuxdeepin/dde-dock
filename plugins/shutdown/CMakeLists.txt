
set(PLUGIN_NAME "shutdown")

project(${PLUGIN_NAME})

# Sources files
file(GLOB_RECURSE SRCS "*.h" "*.cpp")

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)

if (${CMAKE_SYSTEM_PROCESSOR}  STREQUAL "aarch64")
    add_definitions("-DDISABLE_POWER_OPTIONS")
endif()

add_library(${PLUGIN_NAME} STATIC ${SRCS} shutdown.qrc)
set(LIBRARY_OUTPUT_PATH ../)
target_include_directories(${PLUGIN_NAME} PUBLIC ${DtkWidget_INCLUDE_DIRS}
                                                 ${Qt5DBus_INCLUDE_DIRS}
                                                 ../../interfaces)
target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${DtkWidget_LIBRARIES}
    ${Qt5Widgets_LIBRARIES}
    ${Qt5Svg_LIBRARIES}
)
