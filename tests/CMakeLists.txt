cmake_minimum_required(VERSION 3.7)

set(BIN_NAME dde_dock_unit_test)

# 自动生成moc文件
set(CMAKE_AUTOMOC ON)

# 源文件
file(GLOB_RECURSE SRCS "*.h" "*.cpp" "../widgets/*.h" "../widgets/*.cpp")

# 用于测试覆盖率的编译条件
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fprofile-arcs -ftest-coverage -lgcov")

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

pkg_check_modules(QGSettings REQUIRED gsettings-qt)
pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)
pkg_check_modules(XCB_EWMH REQUIRED xcb-ewmh x11)

# 添加执行文件信息
add_executable(${BIN_NAME} ${SRCS} ${INTERFACES} ${SRC_PATH} ../frame/item/item.qrc ../frame/frame.qrc)

# 包含路径
target_include_directories(${BIN_NAME} PUBLIC
    ${DtkWidget_INCLUDE_DIRS}
    ${XCB_EWMH_INCLUDE_DIRS}
    ${DFrameworkDBus_INCLUDE_DIRS}
    ${Qt5Gui_PRIVATE_INCLUDE_DIRS}
    ${QGSettings_INCLUDE_DIRS}
    ../interfaces
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
    -lpthread
    -lm
)

add_custom_target(check)

add_custom_command(TARGET check
    COMMAND ./${BIN_NAME}
    )

add_dependencies(check ${BIN_NAME})
