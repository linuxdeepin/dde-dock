
set(PLUGIN_NAME "media")

project(${PLUGIN_NAME})

generation_dbus_interface(${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/xml ${CMAKE_CURRENT_SOURCE_DIR}/dbusinterface/generation_dbus_interface)

# Sources files
file(GLOB_RECURSE SRCS "*.h" "*.cpp"
    "../../widgets/*.h"
    "../../widgets/*.cpp"
    "../../frame/util/imageutil.h"
    "../../frame/util/imageutil.cpp"
    "../../frame/util/statebutton.h"
    "../../frame/util/statebutton.cpp"
    "../../frame/util/horizontalseperator.h"
    "../../frame/util/horizontalseperator.cpp"
    "../../frame/qtdbusextended/*.h"
    "../../frame/qtdbusextended/*.cpp")

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)

pkg_check_modules(XCB_LIBS REQUIRED IMPORTED_TARGET xcursor)
pkg_check_modules(QGSettings REQUIRED IMPORTED_TARGET gsettings-qt)

add_definitions("${QT_DEFINITIONS} -DQT_PLUGIN")
add_library(${PLUGIN_NAME} SHARED ${SRCS})
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ../quick-trays)
target_include_directories(${PLUGIN_NAME} PUBLIC ${DtkWidget_INCLUDE_DIRS}
  ../../interfaces
  ../../widgets
  ../../frame
  ../../frame/qtdbusextended
  ./dbusinterface/generation_dbus_interface
  componments)

target_link_libraries(${PLUGIN_NAME} PRIVATE
  ${DtkWidget_LIBRARIES}
  PkgConfig::QGSettings
  PkgConfig::XCB_LIBS
  Qt5::Widgets
  Qt5::DBus
  Qt5::Svg)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-dock/plugins/quick-trays)
