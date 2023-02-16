cmake_minimum_required(VERSION 3.7)

set(BIN_NAME dde-dock)

configure_file(environments.h.in environments.h @ONLY)

if (CMAKE_BUILD_TYPE STREQUAL "Debug")
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -g -fsanitize=address -O0")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -g -fsanitize=address -O0")
endif()

generation_dbus_interface(${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/xml ${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/generation_dbus_interface)

# Sources files
file(GLOB_RECURSE SRCS "*.h" "*.cpp" "../widgets/*.h" "../widgets/*.cpp" "../interfaces/*.h")

# Find the library
find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Concurrent REQUIRED)
find_package(Qt5X11Extras REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5WaylandClient REQUIRED)
find_package(Qt5XkbCommonSupport REQUIRED)
find_package(DtkWidget REQUIRED)
find_package(DtkCMake REQUIRED)
find_package(dbusmenu-qt5 REQUIRED)

pkg_check_modules(XCB_EWMH REQUIRED xcb-image xcb-ewmh xcb-composite xtst x11 dbusmenu-qt5 xext xcursor)
pkg_check_modules(QGSettings REQUIRED gsettings-qt)
pkg_check_modules(DtkGUI REQUIRED dtkgui)

set(Wayland_INCLUDE_DIRS /usr/include/DWayland/Client)
set(Wayland_LIBRARIES /usr/lib/${CMAKE_LIBRARY_ARCHITECTURE}/libDWaylandClient.so)

# driver-manager
add_executable(${BIN_NAME}
    ${SRCS}
    ${INTERFACES}
    ${SRC_PATH}
    item/item.qrc)

target_include_directories(${BIN_NAME} PUBLIC
    ${DtkWidget_INCLUDE_DIRS}
    ${XCB_EWMH_INCLUDE_DIRS}
    ${Qt5Gui_PRIVATE_INCLUDE_DIRS}
    ${PROJECT_BINARY_DIR}
    ${QGSettings_INCLUDE_DIRS}
    ${DtkGUI_INCLUDE_DIRS}
    ${Qt5Svg_INCLUDE_DIRS}
    ${dbusmenu-qt5_INCLUDE_DIRS}
    ${Wayland_INCLUDE_DIRS}
    ${Qt5WaylandClient_INCLUDE_DIRS}
    ${Qt5WaylandClient_PRIVATE_INCLUDE_DIRS}
    ${Qt5XkbCommonSupport_PRIVATE_INCLUDE_DIRS}
    ../interfaces
    ../widgets
    ./dbusinterface/generation_dbus_interface
    ./qtdbusextended/
    ./dbusinterface
    accessible
    controller
    dbus
    display
    item
    item/components
    model
    pluginadapter
    screenspliter
    util
    window
    window/components
    window/tray
    window/tray/widgets
    drag
    xcb
    ../plugins/tray
    ../plugins/show-desktop
    ../plugins/datetime
    ../plugins/onboard
    ../plugins/trash
    ../plugins/shutdown
    ../plugins/multitasking
    ../plugins/overlay-warning
)

target_link_libraries(${BIN_NAME} PRIVATE
    ${XCB_EWMH_LIBRARIES}
    ${DtkWidget_LIBRARIES}
    ${Qt5Widgets_LIBRARIES}
    ${Qt5Gui_LIBRARIES}
    ${Qt5Concurrent_LIBRARIES}
    ${Qt5X11Extras_LIBRARIES}
    ${Qt5DBus_LIBRARIES}
    ${QGSettings_LIBRARIES}
    ${DtkGUI_LIBRARIES}
    ${Qt5Svg_LIBRARIES}
    ${Wayland_LIBRARIES}
    ${Qt5Wayland_LIBRARIES}
    ${Qt5WaylandClient_LIBRARIES}
    ${Qt5XkbCommonSupport_LIBRARIES}
    -lpthread -lm
)

if (${CMAKE_SYSTEM_PROCESSOR} STREQUAL "sw_64")
    target_compile_definitions(${BIN_NAME} PUBLIC DISABLE_SHOW_ANIMATION)
endif()

if (${CMAKE_SYSTEM_PROCESSOR} STREQUAL "mips64")
    target_compile_definitions(${BIN_NAME} PUBLIC DISABLE_SHOW_ANIMATION)
endif()

if (${CMAKE_SYSTEM_PROCESSOR} STREQUAL "aarch64")
    target_compile_definitions(${BIN_NAME} PUBLIC DISABLE_SHOW_ANIMATION)
endif()

# bin
install(TARGETS ${BIN_NAME} DESTINATION ${CMAKE_INSTALL_INCLUDEDIR})
