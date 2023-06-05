
set(PLUGIN_NAME "pluginmanager")

project(${PLUGIN_NAME})

# Sources files
file(GLOB_RECURSE SRCS "*.h" "*.cpp" "*.qrc" "../../frame/drag/quickdragcore.h" "../../frame/drag/quickdragcore.cpp"
"../../frame/util/docksettings.h" "../../frame/util/docksettings.cpp"
"../../frame/util/settings.h" "../../frame/util/settings.cpp"
"../../frame/util/pluginloader.h" "../../frame/util/pluginloader.cpp"
"../../frame/dbus/dockinterface.h" "../../frame/dbus/dockinterface.cpp"
"../../frame/dbusinterface/generation_dbus_interface/org_deepin_dde_daemon_dock1.h"
"../../frame/dbusinterface/generation_dbus_interface/org_deepin_dde_daemon_dock1.cpp"
"../../frame/dbusinterface/types/dockrect.h"
"../../frame/dbusinterface/types/dockrect.cpp"
"../../interfaces/pluginmanagerinterface.h"
)

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)

pkg_check_modules(QGSettings REQUIRED IMPORTED_TARGET gsettings-qt)

add_library(${PLUGIN_NAME} SHARED ${SRCS})
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ../loader/)
target_include_directories(${PLUGIN_NAME} PUBLIC ${DtkWidget_INCLUDE_DIRS}
                                                 ./pluginadapter
                                                 ../../frame/drag
                                                 ../../frame/dbusinterface
                                                 ../../frame/dbusinterface/generation_dbus_interface/
                                                 ../../interfaces
                                             )

target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${DtkWidget_LIBRARIES}
    PkgConfig::QGSettings
    Qt5::Widgets
    Qt5::DBus
    Qt5::Svg)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-dock/plugins/loader)
