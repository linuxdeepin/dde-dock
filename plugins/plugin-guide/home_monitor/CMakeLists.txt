# 学习 cmake 时建议直接从命令列表作为入口，遇到不清楚意思的命令都可以在此处查阅：
# https://cmake.org/cmake/help/latest/manual/cmake-commands.7.html
# 另外下面是完整的文档入口：
# https://cmake.org/cmake/help/latest/

# 设置运行该脚本所需的 cmake 最低版本
cmake_minimum_required(VERSION 3.11)

# 使用 set 命令设置一个变量
set(PLUGIN_NAME "home_monitor")

# 设置项目名称
project(${PLUGIN_NAME})

# 启用 qt moc 的支持
set(CMAKE_AUTOMOC ON)
# 启用 qrc 资源文件的支持
set(CMAKE_AUTORCC ON)

# generate a compile commands file as complete database for vim-YouCompleteMe or some other similar tools
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# 指定所有源码文件
# 使用了 cmake 的 file 命令，递归查找项目目录下所有头文件和源码文件，
# 并将结果放入 SRCS 变量中，SRCS 变量可用于后续使用
file(GLOB_RECURSE SRCS "*.h" "*.cpp")

# 指定要用到的库
# 使用了 cmake 的 find_package 命令，查找库 Qt5Widgets 等，
# REQUIRED 参数表示如果没有找到则报错
# find_package 命令在找到并加载指定的库之后会设置一些变量，
# 常用的有：
# <库名>_FOUND          是否找到（Qt5Widgets_FOUND）
# <库名>_DIR            在哪个目录下找到的（Qt5Widgets_DIR）
# <库名>_INCLUDE_DIRS   有哪些头文件目录（Qt5Widgets_INCLUDE_DIRS）
# <库名>_LIBRARIES      有哪些库文件（Qt5Widgets_LIBRARIES）
find_package(Qt5Widgets REQUIRED)
find_package(DtkWidget REQUIRED)

# find_package 命令还可以用来加载 cmake 的功能模块
# 并不是所有的库都直接支持 cmake 查找的，但大部分都支持了 pkg-config 这个标准，
# 因此 cmake 提供了间接加载库的模块：FindPkgConfig， 下面这行命令表示加载 FindPkgConfig 模块，
# 这个 cmake 模块提供了额外的基于 “pkg-config” 加载库的能力
# 执行下面的命令后后会设置如下变量，不过一般用不到：
# PKG_CONFIG_FOUND            pkg-config 可执行文件是否找到了
# PKG_CONFIG_EXECUTABLE       pkg-config 可执行文件的路径
# PKG_CONFIG_VERSION_STRING   pkg-config 的版本信息
find_package(PkgConfig REQUIRED)

# 加载 FindPkgConfig 模块后就可以使用 pkg_check_modules 命令加载需要的库
# pkg_check_modules 命令是由 FindPkgConfig 模块提供的，因此要使用这个命令必须先加载 FindPkgConfig 模块。
# 执行 pkg_check_modules 命令加载库也会设置一些类似执行 find_package 加载库后设置的变量：
# DdeDockInterface_FOUND
# DdeDockInterface_INCLUDE_DIRS
# DdeDockInterface_LIBRARIES
# 还有有另外的一些变量以及更灵活的用法，比如一次性查找多个库，这些请自行查找 cmake 文档学习。
pkg_check_modules(DdeDockInterface REQUIRED dde-dock)

# add_definitions 命令用于声明/定义一些编译/预处理参数
# 根据 cmake 文档描述此命令已经由另外几个功能划分的更为细致的命令所取代，具体请查阅文档
# 在我们这里的例子应该使用较新的 add_compile_definitions 命令，不过为了保持与 dock 已有插件一致，
# 暂时仍然使用 add_definitions，add_definitions 的语法很简单就是直接写要定义的 flag 并在前面加上 "-D" 即可
# 括号中的 ${QT_DEFINITIONS} 变量会在执行 cmake 时展开为它的值，这个变量属于历史遗留，应该是在 qt3/qt4 时有用，
# 基于 qt5 或更高版本的新插件不必使用此变量。要查看 qt5 的库定义了哪些变量应该查看变量：${Qt5Widgets_DEFINITIONS}
add_definitions("${QT_DEFINITIONS} -DQT_PLUGIN")

# 新增一个编译目标
# 这里使用命令 add_library 来表示本项目要生成一个库文件目标，
# 类似的还有命令 add_executable 添加一个可执行二进制目标，甚至 add_custom_target(使用较少) 添加自定义目标
# SHARED 表示生成的库应该是动态库，
# 变量 ${PLUGIN_NAME} 和 ${SRCS} 都是前面处理好的，
# 另外 qrc 资源文件也应该追加在后面以编译进目标中。
add_library(${PLUGIN_NAME} SHARED ${SRCS} home_monitor.qrc)

# 设置目标的生成位置，这里表示生成在执行 make 的目录,
# 另外还有很多可用于设置的属性，可查阅 cmake 文档。
set_target_properties(${PLUGIN_NAME} PROPERTIES LIBRARY_OUTPUT_DIRECTORY ./)

# 设置目标要使用的 include 目录，即头文件目录
# 变量 ${DtkWidget_INCLUDE_DIRS} 是在前面执行 find_package 命令时引入的
# 当出现编译失败提示找不到某些库的头文件时应该检查此处是否将所有需要的头文件都包含了
target_include_directories(${PLUGIN_NAME} PUBLIC
    ${Qt5Widgets_INCLUDE_DIRS}
    ${DtkWidget_INCLUDE_DIRS}
    ${DdeDockInterface_INCLUDE_DIRS}
)

# 设置目标要使用的链接库
# 变量 ${DtkWidget_LIBRARIES} 和 ${Qt5Widgets_LIBRARIES} 是在前面执行 find_package 命令时引入的
# 当出现运行时错误提示某些符号没有定义时应该检查此处是否将所有用的库都写在了这里
target_link_libraries(${PLUGIN_NAME} PRIVATE
    ${Qt5Widgets_LIBRARIES}
    ${DtkWidget_LIBRARIES}
    ${DdeDockInterface_LIBRARIES}
)

# 设置安装路径的前缀(默认为"/usr/local")
set(CMAKE_INSTALL_PREFIX "/usr")

# 设置执行 make install 时哪个目标应该被 install 到哪个位置
install(TARGETS ${PLUGIN_NAME} LIBRARY DESTINATION lib/dde-dock/plugins)
