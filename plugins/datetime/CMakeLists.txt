
set(PLUGIN_NAME "datetime")

project(${PLUGIN_NAME})

# Sources files
file(GLOB SRCS "*.h" "*.cpp" "../../widgets/*.h" "../../widgets/*.cpp")

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)
pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)

pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)
pkg_check_modules(QGSettings REQUIRED gsettings-qt)

add_definitions("${QT_DEFINITIONS} -DQT_PLUGIN")
add_library(${PLUGIN_NAME} SHARED ${SRCS} datetime.qrc)
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ../)
target_include_directories(${PLUGIN_NAME} PUBLIC
  ${DtkWidget_INCLUDE_DIRS}
  ${DFrameworkDBus_INCLUDE_DIRS}
  ${QGSettings_INCLUDE_DIRS}
  ../../interfaces)
target_link_libraries(${PLUGIN_NAME} PRIVATE
  ${Qt5DBus_LIBRARIES}
  ${DtkWidget_LIBRARIES}
  ${Qt5Widgets_LIBRARIES}
  ${Qt5Svg_LIBRARIES}
  ${DFrameworkDBus_LIBRARIES}
  ${QGSettings_LIBRARIES}
)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION ${PLUGIN_INSTALL_DIR}/dde-dock/plugins)
