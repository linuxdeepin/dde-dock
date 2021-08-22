
set(PLUGIN_NAME "onboard")

project(${PLUGIN_NAME})

# Sources files
file(GLOB_RECURSE SRCS "*.h" "*.cpp" "../../widgets/*.h" "../../widgets/*.cpp")

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)

add_definitions("${QT_DEFINITIONS} -DQT_PLUGIN")
add_library(${PLUGIN_NAME} SHARED ${SRCS} onboard.qrc)
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ../)
target_include_directories(${PLUGIN_NAME} PUBLIC ${Qt5DBus_INCLUDE_DIRS}
                                                 ${DFrameworkDBus_INCLUDE_DIRS}
                                                 ../../interfaces)
target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${Qt5Widgets_LIBRARIES}
    ${DFrameworkDBus_LIBRARIES}
)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION ${PLUGIN_INSTALL_DIR}/dde-dock/plugins)
