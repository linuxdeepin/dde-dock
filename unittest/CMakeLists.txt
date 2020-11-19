cmake_minimum_required(VERSION 3.7)

set(BIN_NAME dde_dock_unit_test)

# 自动生成moc文件
set(CMAKE_AUTOMOC ON)

# 源文件
file(GLOB_RECURSE SRCS "*.h" "*.cpp")

# 用于测试覆盖率的编译条件
#set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fprofile-arcs -ftest-coverage -lgcov")

# 查找依赖库
find_package(PkgConfig REQUIRED)
find_package(Qt5 COMPONENTS Widgets Test DBus REQUIRED)
find_package(GTest REQUIRED)

pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)
pkg_check_modules(QGSettings REQUIRED gsettings-qt)

# 添加执行文件信息
add_executable(${BIN_NAME} ${SRCS} ${INTERFACES})

# 包含路径
target_include_directories(${BIN_NAME} PUBLIC ${DFrameworkDBus_INCLUDE_DIRS}  ../../interfaces)

# 链接库
target_link_libraries(${BIN_NAME} PRIVATE
    ${Qt5Test_LIBRARIES}
    ${Qt5DBus_LIBRARIES}
    ${DFrameworkDBus_LIBRARIES}
    ${QGSettings_LIBRARIES}
    ${Qt5Widgets_LIBRARIES}
    ${GTEST_LIBRARIES}
    -lpthread
    -lm
)
