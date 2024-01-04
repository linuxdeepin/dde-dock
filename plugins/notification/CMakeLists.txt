set(PLUGIN_NAME "notification")
find_package(Qt5 REQUIRED COMPONENTS DBus)
find_package(Dtk REQUIRED COMPONENTS Widget)

add_library(${PLUGIN_NAME} SHARED
    notification.h
    notification.cpp
    notificationplugin.h
    notificationplugin.cpp
    ${CMAKE_SOURCE_DIR}/widgets/tipswidget.h
    ${CMAKE_SOURCE_DIR}/widgets/tipswidget.cpp
    notification.qrc
)
target_compile_definitions(${PLUGIN_NAME} PRIVATE QT_PLUGIN)
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/plugins)

target_include_directories(${PLUGIN_NAME} PRIVATE
    $<BUILD_INTERFACE:${CMAKE_SOURCE_DIR}/interfaces>
    $<BUILD_INTERFACE:${CMAKE_SOURCE_DIR}/widgets>
    $<BUILD_INTERFACE:${CMAKE_SOURCE_DIR}/frame/qtdbusextended>
)

target_link_libraries(${PLUGIN_NAME} PRIVATE
    Dtk::Widget
    Qt5::DBus
)

install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-dock/plugins)
