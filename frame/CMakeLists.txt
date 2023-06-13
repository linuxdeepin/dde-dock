cmake_minimum_required(VERSION 3.16)

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
find_package(DtkGui REQUIRED)
find_package(DtkWidget REQUIRED)
find_package(DtkCMake REQUIRED)
find_package(dbusmenu-qt5 REQUIRED)
find_package(ECM REQUIRED NO_MODULE)
set(CMAKE_MODULE_PATH ${ECM_MODULE_PATH})
find_package(DWayland REQUIRED)

set(THREADS_PREFER_PTHREAD_FLAG ON)
find_package(Threads REQUIRED)

pkg_check_modules(XCB_EWMH REQUIRED IMPORTED_TARGET x11 xcb xcb-icccm xcb-image xcb-ewmh xcb-composite xtst dbusmenu-qt5 xext xcursor xkbcommon xres)
pkg_check_modules(QGSettings REQUIRED IMPORTED_TARGET gsettings-qt)
pkg_check_modules(WAYLAND REQUIRED IMPORTED_TARGET wayland-client wayland-cursor wayland-egl)

# driver-manager
add_executable(${BIN_NAME}
    ${SRCS}
    ${INTERFACES}
    ${SRC_PATH}
    item/item.qrc)

target_include_directories(${BIN_NAME} PUBLIC
    ${DtkWidget_INCLUDE_DIRS}
    ${PROJECT_BINARY_DIR}
    ${Qt5WaylandClient_PRIVATE_INCLUDE_DIRS}
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
    ${DtkWidget_LIBRARIES}
    PkgConfig::QGSettings
    PkgConfig::XCB_EWMH
    PkgConfig::WAYLAND
    Dtk::Gui
    Qt5::Widgets
    Qt5::Gui
    Qt5::Concurrent
    Qt5::X11Extras
    Qt5::DBus
    Qt5::Svg
    Qt5::WaylandClient
    Qt5::XkbCommonSupport
    DWaylandClient
    Threads::Threads
    -lm
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
install(TARGETS ${BIN_NAME} DESTINATION ${CMAKE_INSTALL_BINDIR})
# window_patterns
install(FILES taskmanager/window_patterns.json DESTINATION ${CMAKE_INSTALL_DATADIR}/dde-dock/)
