ADD_COMPILE_OPTIONS(-fno-access-control)

cmake_minimum_required(VERSION 3.7)

set(BIN_NAME dde_dock_unit_test)

# 自动生成moc文件
set(CMAKE_AUTOMOC ON)

# 源文件
file(GLOB_RECURSE SRCS
    "*.h"
    "*.cpp"
    "../widgets/*.h"
    "../widgets/*.cpp")

# Sources files
file(GLOB_RECURSE PLUGIN_SRCS
    "../plugins/bluetooth/*.h"
    "../plugins/bluetooth/*.cpp"
    "../plugins/bluetooth/componments/*.h"
    "../plugins/bluetooth/componments/*.cpp"
    "../frame/util/horizontalseperator.h"
    "../frame/util/horizontalseperator.cpp")

# 用于测试覆盖率的编译条件
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fprofile-arcs -ftest-coverage -lgcov")

if (CMAKE_BUILD_TYPE STREQUAL "Debug")
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -g -fsanitize=address -O2")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -g -fsanitize=address -O2")
endif()

# 查找gmock的cmake文件
set(CMAKE_MODULE_PATH ${CMAKE_SOURCE_DIR}/tests/cmake/modules)

# 查找依赖库
find_package(PkgConfig REQUIRED)
find_package(Qt5Widgets REQUIRED)
find_package(Qt5Concurrent REQUIRED)
find_package(Qt5X11Extras REQUIRED)
find_package(Qt5DBus REQUIRED)
find_package(DtkWidget REQUIRED)
find_package(Qt5Svg REQUIRED)
find_package(Qt5 COMPONENTS Test REQUIRED)
find_package(GTest REQUIRED)
find_package(GMock REQUIRED)

pkg_check_modules(QGSettings REQUIRED gsettings-qt)
pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)
pkg_check_modules(XCB_EWMH REQUIRED xcb-ewmh x11)

# 添加执行文件信息
add_executable(${BIN_NAME}
    ${SRCS}
    ${INTERFACES}
    ${SRC_PATH}
    ${PLUGIN_SRCS}
    ../frame/item/item.qrc
    ut_res.qrc)

# 包含路径
target_include_directories(${BIN_NAME} PUBLIC
    ${DtkWidget_INCLUDE_DIRS}
    ${XCB_EWMH_INCLUDE_DIRS}
    ${DFrameworkDBus_INCLUDE_DIRS}
    ${Qt5Gui_PRIVATE_INCLUDE_DIRS}
    ${QGSettings_INCLUDE_DIRS}
    ../interfaces
    fakedbus
    ../plugins/bluetooth
    ../plugins/bluetooth/componments
)

# 链接库
target_link_libraries(${BIN_NAME} PRIVATE
    ${Qt5Test_LIBRARIES}
    ${XCB_EWMH_LIBRARIES}
    ${DFrameworkDBus_LIBRARIES}
    ${DtkWidget_LIBRARIES}
    ${Qt5Widgets_LIBRARIES}
    ${Qt5Concurrent_LIBRARIES}
    ${Qt5X11Extras_LIBRARIES}
    ${Qt5DBus_LIBRARIES}
    ${QGSettings_LIBRARIES}
    ${Qt5Svg_LIBRARIES}
    ${GTEST_LIBRARIES}
    ${GMOCK_LIBRARIES}
    -lpthread
    -lm
)

add_custom_target(check)

add_custom_command(TARGET check
    COMMAND ./${BIN_NAME}
    )

add_dependencies(check ${BIN_NAME})
