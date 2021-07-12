
set(PLUGIN_NAME "network")

project(${PLUGIN_NAME})

# Sources files
file(GLOB_RECURSE SRCS "*.h" "*.cpp" "../../widgets/*.h" "../../widgets/*.cpp"
    "../../frame/util/imageutil.h" "../../frame/util/imageutil.cpp"
    "../../frame/util/statebutton.h" "../../frame/util/statebutton.cpp")

find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Network REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)
find_package(KF5NetworkManagerQt REQUIRED)

pkg_check_modules(DDE-Network-Utils REQUIRED dde-network-utils)
pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)
pkg_check_modules(QGSettings REQUIRED gsettings-qt)

get_target_property(KF5_QT_INCLUDE_DIRS KF5::NetworkManagerQt INTERFACE_INCLUDE_DIRECTORIES)

set(TRY_KF5_LIBRARIES
IMPORTED_LOCATION_DEBIAN
IMPORTED_LOCATION_NOCONFIG
)

if (NOT KF5_QT_LIBRARIES)
    find_library(KF5_QT_LIBRARIES libKF5NetworkManagerQt.so)
endif()

add_definitions("${QT_DEFINITIONS} -DQT_PLUGIN")
add_library(${PLUGIN_NAME} SHARED ${SRCS} network.qrc)
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ../system-trays)
target_include_directories(${PLUGIN_NAME} PUBLIC ${DtkWidget_INCLUDE_DIRS}
                                                 ${Qt5DBus_INCLUDE_DIRS}
                                                 ${DFrameworkDBus_INCLUDE_DIRS}
                                                 ${QGSettings_INCLUDE_DIRS}
                                                 ${DDE-Network-Utils_INCLUDE_DIRS}
                                                 ${KF5_QT_INCLUDE_DIRS}
                                                 ${Qt5Network_INCLUDE_DIRS}
                                                 ../../interfaces
                                                 ../../frame)
target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${DtkWidget_LIBRARIES}
    ${Qt5Widgets_LIBRARIES}
    ${Qt5Svg_LIBRARIES}
    ${Qt5DBus_LIBRARIES}
    ${QGSettings_LIBRARIES}
    ${DDE-Network-Utils_LIBRARIES}
    ${DFrameworkDBus_LIBRARIES}
    ${KF5_QT_LIBRARIES}
    ${Qt5Network_LIBRARIES}
)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-dock/plugins/system-trays)
