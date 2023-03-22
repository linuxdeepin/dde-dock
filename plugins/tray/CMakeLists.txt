
set(PLUGIN_NAME "tray")

project(${PLUGIN_NAME})

generation_dbus_interface(${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/xml ${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/generation_dbus_interface)

# Sources files
file(GLOB_RECURSE SRCS "*.h"
    "*.cpp"
    "../../widgets/*.h"
    "../../widgets/*.cpp"
    "../../frame/util/themeappicon.h"
    "../../frame/util/themeappicon.cpp"
    "../../frame/util/dockpopupwindow.h"
    "../../frame/util/dockpopupwindow.cpp"
    "../../frame/util/abstractpluginscontroller.h"
    "../../frame/util/abstractpluginscontroller.cpp"
    "../../frame/util/pluginloader.h"
    "../../frame/util/pluginloader.cpp"
    "../../frame/dbus/sni/*.h"
    "../../frame/dbus/sni/*.cpp"
    "../../frame/dbus/dbusmenu.h"
    "../../frame/dbus/dbusmenu.cpp"
    "../../frame/dbus/dbusmenumanager.h"
    "../../frame/dbus/dbusmenumanager.cpp"
    "../../frame/dbus/dockinterface.h"
    "../../frame/dbus/dockinterface.cpp"
    "../../widgets/*.h"
    "../../widgets/*.cpp"
    "../../frame/util/imageutil.h"
    "../../frame/util/imageutil.cpp"
    "../../frame/util/menudialog.h"
    "../../frame/util/menudialog.cpp"
    "../../frame/util/touchsignalmanager.h"
    "../../frame/util/touchsignalmanager.cpp"
    "../../frame/controller/proxyplugincontroller.h"
    "../../frame/controller/proxyplugincontroller.cpp"
    "../../frame/qtdbusextended/*.h"
    "../../frame/qtdbusextended/*.cpp"
)

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(Qt5X11Extras REQUIRED)
find_package(Qt5Concurrent REQUIRED)
find_package(DtkWidget REQUIRED)
find_package(dbusmenu-qt5 REQUIRED)

pkg_check_modules(XCB_LIBS REQUIRED xcb-ewmh xcb xcb-image xcb-composite xtst x11 xext xcb-icccm dbusmenu-qt5 xcursor)
pkg_check_modules(QGSettings REQUIRED gsettings-qt)

add_definitions("${QT_DEFINITIONS} -DQT_PLUGIN")
add_library(${PLUGIN_NAME} SHARED ${SRCS} tray.qrc)
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ../)
target_include_directories(${PLUGIN_NAME} PUBLIC ${DtkWidget_INCLUDE_DIRS}
                                                 ../../interfaces
                                                 ../../frame
                                                 ../../frame/qtdbusextended
                                                 ./dbusinterface
                                                 ./dbusinterface/generation_dbus_interface)

target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${DtkWidget_LIBRARIES}
    PkgConfig::QGSettings
    PkgConfig::XCB_LIBS
    Qt5::Concurrent
    Qt5::X11Extras
    Qt5::Widgets
    Qt5::DBus
    Qt5::Svg
    pthread)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-dock/plugins)
