
set(PLUGIN_NAME "onboard")

project(${PLUGIN_NAME})

generation_dbus_interface(${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/xml ${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/generation_dbus_interface)

# Sources files
file(GLOB_RECURSE SRCS "*.h"
    "*.cpp"
    "../../widgets/*.h"
    "../../widgets/*.cpp"
    "./dbusinterface/generation_dbus_interface/*.h"
    "./dbusinterface/generation_dbus_interface/*.cpp"
    "../../frame/qtdbusextended/*.h"
    "../../frame/qtdbusextended/*.cpp")

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkGui REQUIRED)
find_package(DtkWidget REQUIRED)

add_definitions("${QT_DEFINITIONS} -DQT_PLUGIN")
add_library(${PLUGIN_NAME} SHARED ${SRCS} onboard.qrc)
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ../)

target_include_directories(${PLUGIN_NAME} PUBLIC ${Qt5DBus_INCLUDE_DIRS}
    ../../interfaces
    ../../frame/qtdbusextended
    ./dbusinterface/generation_dbus_interface)

target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${Qt5Widgets_LIBRARIES}
    ${Qt5DBus_LIBRARIES}
    ${DtkGui_LIBRARIES}
    ${DtkWidget_LIBRARIES}
)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-dock/plugins)
