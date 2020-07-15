cmake_minimum_required(VERSION 3.7)

set(BIN_NAME dde_dock_unit_test)

# 自动生成moc文件
set(CMAKE_AUTOMOC ON)

# 源文件
file(GLOB_RECURSE SRCS "*.h" "*.cpp")

# 查找依赖库
find_package(PkgConfig REQUIRED)
find_package(Qt5 COMPONENTS Test DBus REQUIRED)

pkg_check_modules(DFrameworkDBus REQUIRED dframeworkdbus)

# 添加执行文件信息
add_executable(${BIN_NAME} ${SRCS} ${INTERFACES})

# 包含路径
target_include_directories(${BIN_NAME} PUBLIC ${DFrameworkDBus_INCLUDE_DIRS})

# 链接库
target_link_libraries(${BIN_NAME} PRIVATE
    ${Qt5Test_LIBRARIES}
    ${Qt5DBus_LIBRARIES}
    ${DFrameworkDBus_LIBRARIES}
)
