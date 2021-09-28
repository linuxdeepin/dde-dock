cmake_minimum_required(VERSION 3.7)

set(PLUGIN_NAME "dcc-dock-settings-plugin")

project(${PLUGIN_NAME})

set(CMAKE_AUTOMOC ON)

file(GLOB_RECURSE SRCS
    "*.h"
    "*.cpp")

find_package(Qt5 COMPONENTS Core Widgets DBus REQUIRED)
find_package(DdeControlCenter REQUIRED)
find_package(PkgConfig REQUIRED)
find_package(DtkWidget REQUIRED)

pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)
pkg_check_modules(QGSettings REQUIRED gsettings-qt)

add_library(${PLUGIN_NAME} SHARED ${SRCS})
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ./)
target_include_directories(${PLUGIN_NAME} PUBLIC
    ../src
    ${Qt5Widgets_INCLUDE_DIRS}
    ${DtkWidget_INCLUDE_DIRS}
    ${DdeControlCenter_INCLUDE_DIR}
    ${DFrameworkDBus_INCLUDE_DIRS}
    ${QGSettings_INCLUDE_DIRS}
    ${Qt5DBus_INCLUDE_DIRS}
    )
target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${Qt5Widgets_LIBRARIES}
    ${DtkWidget_LIBRARIES}
    ${DdeControlCenter_LIBRARIES}
    ${DFrameworkDBus_LIBRARIES}
    ${QGSettings_LIBRARIES}
    ${Qt5DBus_LIBRARIES}
    )

set(CMAKE_INSTALL_PREFIX "/usr")

install(FILES com.deepin.dde.control-center.dock-settings.gschema.xml
    DESTINATION share/glib-2.0/schemas)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-control-center/modules)
