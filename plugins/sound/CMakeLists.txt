
set(PLUGIN_NAME "sound")

project(${PLUGIN_NAME})

# Sources files
file(GLOB_RECURSE SRCS "*.h" "*.cpp")

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)

pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)
pkg_check_modules(QGSettings REQUIRED gsettings-qt)

add_library(${PLUGIN_NAME} STATIC ${SRCS} resources/sound.qrc)
set(LIBRARY_OUTPUT_PATH ../system-trays)
target_include_directories(${PLUGIN_NAME} PUBLIC ${DtkWidget_INCLUDE_DIRS}
  ${DFrameworkDBus_INCLUDE_DIRS}
  ${QGSettings_INCLUDE_DIRS}
  ../../interfaces
  ../../frame)
target_link_libraries(${PLUGIN_NAME} PRIVATE
  ${DtkWidget_LIBRARIES}
  ${DFrameworkDBus_LIBRARIES}
  ${QGSettings_LIBRARIES}
  ${Qt5DBus_LIBRARIES}
  ${Qt5Widgets_LIBRARIES}
  ${Qt5Svg_LIBRARIES}
)
