
set(PLUGIN_NAME "overlay-warning")

project(${PLUGIN_NAME})

# Sources files
file(GLOB_RECURSE SRCS "*.h" "*.cpp" "../../widgets/tipswidget.h" "../../widgets/tipswidget.cpp")

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)

add_definitions("${QT_DEFINITIONS} -DQT_PLUGIN")
add_library(${PLUGIN_NAME} SHARED ${SRCS} overlay-warning.qrc)
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ../)
target_include_directories(${PLUGIN_NAME} PUBLIC ${DtkWidget_INCLUDE_DIRS}
                                                 ../../interfaces)
target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${DtkWidget_LIBRARIES}
    Qt5::Widgets
    Qt5::DBus
    Qt5::Svg)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-dock/plugins)

# policy kit
install(FILES org.deepin.dde.dock.overlay.policy
        DESTINATION share/polkit-1/actions/)
