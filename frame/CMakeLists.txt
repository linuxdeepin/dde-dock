cmake_minimum_required(VERSION 3.7)

set(BIN_NAME dde-dock)

configure_file(environments.h.in environments.h @ONLY)

# Sources files
file(GLOB_RECURSE SRCS "*.h" "*.cpp" "../widgets/*.h" "../widgets/*.cpp")

set(DATETIME_SRCS
    plugins/datetime/datetimeplugin.cpp
    plugins/datetime/datetimeplugin.h
    plugins/datetime/datetimewidget.h
    plugins/datetime/datetimewidget.cpp
)

set(NETWORK_SRCS
    plugins/network/item/applet/accesspoint.cpp
    plugins/network/item/applet/accesspoint.h
    plugins/network/item/applet/accesspointwidget.h
    plugins/network/item/applet/devicecontrolwidget.cpp
    plugins/network/item/applet/devicecontrolwidget.h
    plugins/network/item/applet/horizontalseperator.cpp
    plugins/network/item/applet/horizontalseperator.h
    plugins/network/item/applet/wirelesslist.h
    plugins/network/item/applet/accesspointwidget.cpp
    plugins/network/item/applet/wirelesslist.cpp
    plugins/network/item/deviceitem.cpp
    plugins/network/item/deviceitem.h
    plugins/network/item/wireditem.h
    plugins/network/item/wirelessitem.h
    plugins/network/item/wireditem.cpp
    plugins/network/item/wirelessitem.cpp
    plugins/network/networkplugin.h
    plugins/network/util/imageutil.cpp
    plugins/network/util/imageutil.h
    plugins/network/networkplugin.cpp
)

set(SHUTDOWN_SRCS
    plugins/shutdown/shutdownplugin.cpp
    plugins/shutdown/shutdownplugin.h
    )

set(POWER_SRCS
    plugins/power/powerstatuswidget.h
    plugins/power/powerplugin.cpp
    plugins/power/powerplugin.h
    plugins/power/powerstatuswidget.cpp
)

set(SOUND_SRCS
    plugins/sound/componments/horizontalseparator.cpp
    plugins/sound/componments/horizontalseparator.h
    plugins/sound/componments/volumeslider.cpp
    plugins/sound/componments/volumeslider.h
    plugins/sound/dbus/dbusaudio.cpp
    plugins/sound/dbus/dbusaudio.h
    plugins/sound/dbus/dbussink.cpp
    plugins/sound/dbus/dbussink.h
    plugins/sound/dbus/dbussinkinput.cpp
    plugins/sound/dbus/dbussinkinput.h
    plugins/sound/sinkinputwidget.h
    plugins/sound/soundapplet.h
    plugins/sound/sounditem.h
    plugins/sound/soundplugin.cpp
    plugins/sound/soundplugin.h
    plugins/sound/sinkinputwidget.cpp
    plugins/sound/soundapplet.cpp
    plugins/sound/sounditem.cpp
)

set(TRAY_SRCS
    plugins/tray/abstracttraywidget.cpp
    plugins/tray/abstracttraywidget.h
    plugins/tray/dbus/dbustraymanager.cpp
    plugins/tray/dbus/dbustraymanager.h
    plugins/tray/fashiontray/containers/abstractcontainer.cpp
    plugins/tray/fashiontray/containers/abstractcontainer.h
    plugins/tray/fashiontray/containers/attentioncontainer.cpp
    plugins/tray/fashiontray/containers/attentioncontainer.h
    plugins/tray/fashiontray/containers/holdcontainer.cpp
    plugins/tray/fashiontray/containers/holdcontainer.h
    plugins/tray/fashiontray/containers/normalcontainer.cpp
    plugins/tray/fashiontray/containers/normalcontainer.h
    plugins/tray/fashiontray/containers/spliteranimated.cpp
    plugins/tray/fashiontray/containers/spliteranimated.h
    plugins/tray/fashiontray/fashiontrayconstants.h
    plugins/tray/fashiontray/fashiontraycontrolwidget.cpp
    plugins/tray/fashiontray/fashiontraycontrolwidget.h
    plugins/tray/fashiontray/fashiontrayitem.cpp
    plugins/tray/fashiontray/fashiontrayitem.h
    plugins/tray/fashiontray/fashiontraywidgetwrapper.h
    plugins/tray/fashiontray/fashiontraywidgetwrapper.cpp
    plugins/tray/indicatortray.h
    plugins/tray/indicatortraywidget.cpp
    plugins/tray/indicatortraywidget.h
    plugins/tray/snitraywidget.h
    plugins/tray/system-trays/systemtrayitem.cpp
    plugins/tray/system-trays/systemtrayitem.h
    plugins/tray/system-trays/systemtrayscontroller.cpp
    plugins/tray/system-trays/systemtrayscontroller.h
    plugins/tray/trayplugin.h
    plugins/tray/xembedtraywidget.cpp
    plugins/tray/xembedtraywidget.h
    plugins/tray/indicatortray.cpp
    plugins/tray/snitraywidget.cpp
    plugins/tray/trayplugin.cpp
    )

set(TRASH_SRCS
    plugins/trash/popupcontrolwidget.h
    plugins/trash/trashplugin.cpp
    plugins/trash/trashplugin.h
    plugins/trash/trashwidget.h
    plugins/trash/popupcontrolwidget.cpp
    plugins/trash/trashwidget.cpp
    )

set(KEYBOARD_LAYOUT_SRCS
    plugins/keyboard-layout/dbusadaptors.cpp
    plugins/keyboard-layout/dbusadaptors.h
    plugins/keyboard-layout/keyboardplugin.cpp
    plugins/keyboard-layout/keyboardplugin.h
    )

set(OVERLAY_SRCS
    plugins/overlay-warning/overlay-warning-plugin.h
    plugins/overlay-warning/overlay-warning-plugin.cpp
    )

set(SRCS
    ${SRCS}
    ${DATETIME_SRCS}
    ${NETWORK_SRCS}
    ${SHUTDOWN_SRCS}
    ${POWER_SRCS}
    ${SOUND_SRCS}
    ${TRAY_SRCS}
    ${TRASH_SRCS}
    ${KEYBOARD_LAYOUT_SRCS}
    ${OVERLAY_SRCS}
)

set(QRCS
    plugins/datetime/datetime.qrc
    plugins/network/network.qrc
    plugins/overlay-warning/overlay-warning.qrc
    plugins/power/power.qrc
    plugins/shutdown/shutdown.qrc
    plugins/sound/resources/sound.qrc
    plugins/tray/tray.qrc
    item/item.qrc
    frame.qrc
    )

# Find the library
find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Concurrent REQUIRED)
find_package(Qt5X11Extras REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)
find_package(Qt5Svg REQUIRED)

pkg_check_modules(XCB_EWMH REQUIRED xcb-ewmh x11 xcb xcb-image xcb-composite xtst xext xcb-icccm dbusmenu-qt5)
pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)
pkg_check_modules(QGSettings REQUIRED gsettings-qt)
pkg_check_modules(DDE-Network-Utils REQUIRED dde-network-utils)
pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)

if (${CMAKE_SYSTEM_PROCESSOR}  STREQUAL "aarch64")
    add_definitions("-DDISABLE_POWER_OPTIONS")
endif()

# driver-manager
add_executable(${BIN_NAME} ${SRCS} ${INTERFACES} ${QRCS})
target_include_directories(${BIN_NAME} PUBLIC ${DtkWidget_INCLUDE_DIRS}
                                              ${XCB_EWMH_INCLUDE_DIRS}
                                              ${DFrameworkDBus_INCLUDE_DIRS}
                                              ${Qt5Gui_PRIVATE_INCLUDE_DIRS}
                                              ${PROJECT_BINARY_DIR}
                                              ${QGSettings_INCLUDE_DIRS}
                                              ${DDE-Network-Utils_INCLUDE_DIRS}
                                              ../interfaces)
target_link_libraries(${BIN_NAME} PRIVATE
    ${XCB_EWMH_LIBRARIES}
    ${DFrameworkDBus_LIBRARIES}
    ${DtkWidget_LIBRARIES}
    ${Qt5Widgets_LIBRARIES}
    ${Qt5Concurrent_LIBRARIES}
    ${Qt5X11Extras_LIBRARIES}
    ${Qt5DBus_LIBRARIES}
    ${QGSettings_LIBRARIES}
    ${Qt5Svg_LIBRARIES}
    ${DDE-Network-Utils_LIBRARIES}
    pthread
)

# bin
install(TARGETS ${BIN_NAME} DESTINATION bin)
