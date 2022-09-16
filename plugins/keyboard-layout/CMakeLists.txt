
set(PLUGIN_NAME "keyboard-layout")

project(${PLUGIN_NAME})

generation_dbus_interface(${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/xml ${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/generation_dbus_interface)

# Sources files
file(GLOB_RECURSE SRCS "*.h"
    "*.cpp"
    "./dbusinterface/generation_dbus_interface/*.h"
    "./dbusinterface/generation_dbus_interface/*.cpp"
    "../../frame/qtdbusextended/*.h"
    "../../frame/qtdbusextended/*.cpp"
    "./dbusinterface/types/*.h"
    "./dbusinterface/types/*.cpp")

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(DtkWidget REQUIRED)
find_package(Qt5DBus REQUIRED)

pkg_check_modules(QGSettings REQUIRED gsettings-qt)

add_definitions("${QT_DEFINITIONS} -DQT_PLUGIN")
add_library(${PLUGIN_NAME} SHARED ${SRCS})
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ../)
target_include_directories(${PLUGIN_NAME} PUBLIC ${DtkWidget_INCLUDE_DIRS}
                                                 ${QGSettings_INCLUDE_DIRS}
                                                 ../../interfaces
                                                 ./dbusinterface
                                                 ./dbusinterface/generation_dbus_interface
                                                 ../../frame/qtdbusextended)

target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${DtkWidget_LIBRARIES}
    ${QGSettings_LIBRARIES}
    ${Qt5Widgets_LIBRARIES}
    ${Qt5Svg_LIBRARIES}
    ${Qt5DBus_LIBRARIES}
)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-dock/plugins/system-trays/)
install(FILES ./keybord_layout.json DESTINATION /etc/dde-dock/indicator)
