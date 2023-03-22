
set(PLUGIN_NAME "power")

project(${PLUGIN_NAME})

generation_dbus_interface(${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/xml ${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/generation_dbus_interface)

# Sources files
file(GLOB_RECURSE SRCS "*.h"
    "*.cpp"
    "../../widgets/tipswidget.h"
    "../../widgets/tipswidget.cpp"
    "../../frame/qtdbusextended/*.h"
    "../../frame/qtdbusextended/*.cpp")

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)

pkg_check_modules(QGSettings REQUIRED IMPORTED_TARGET gsettings-qt)

add_definitions("${QT_DEFINITIONS} -DQT_PLUGIN")
add_library(${PLUGIN_NAME} SHARED ${SRCS} power.qrc)
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ../quick-trays)
target_include_directories(${PLUGIN_NAME} PUBLIC ${DtkWidget_INCLUDE_DIRS}
                                                 ../../interfaces
                                                 ../../frame/qtdbusextended
                                                 ./dbusinterface/generation_dbus_interface)

target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${DtkWidget_LIBRARIES}
    PkgConfig::QGSettings
    Qt5::Widgets
    Qt5::DBus
    Qt5::Svg)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-dock/plugins/quick-trays)
