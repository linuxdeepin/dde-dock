# 从零构建 dde-dock 的插件
本教程将展示一个简单的 dde-dock 插件的开发过程，插件开发者可跟随此步骤为 dde-dock 创造出更多具有丰富功能的插件。

在本教程中，将创建一个可以实时显示用户家目录(`~/`)使用情况的小工具。

## 插件的工作原理
dde-dock 插件本质是一个按 Qt 插件标准所开发的共享库文件(`so`)。通过 dde-dock 预定的规范与提供的接口，共同完成 dde-dock 的功能扩展。

## 准备环境
插件的开发环境可以是任意的，只要是符合 Qt 插件规范及 dde-dock 插件规范的共享库文件，都可以被当作 dde-dock 插件载入。
下面以 Qt + cmake 为例进行说明，以 Deepin 15.9 环境为基础，先安装如下的包：

- dde-dock-dev
- cmake
- qtbase5-dev-tools
- libqt5core5a
- libqt5widgets5
- pkg-config

## 基本的项目结构

### 创建必需的项目目录与文件
插件名称叫做`home_monitor`，所以创建以下的目录结构：
```
home_monitor
├── home_monitor.json
├── homemonitorplugin.cpp
├── homemonitorplugin.h
└── CMakeLists.txt
```

`CMakeLists.txt` 是 cmake 命令要读取的配置文件，其内容如下：
``` cmake
# 学习 cmake 时建议直接从命令列表作为入口，遇到不清楚意思的命令都可以在此处查阅：
# https://cmake.org/cmake/help/latest/manual/cmake-commands.7.html
# 另外下面时完整的文档入口：
# https://cmake.org/cmake/help/latest/

# 设置运行被配置所需的 cmake 最低版本
cmake_minimum_required(VERSION 3.11)

# 使用 set 命令设置一个变量
set(PLUGIN_NAME "home_monitor")

# 设置项目名称
project(${PLUGIN_NAME})

# 启用 qt moc 的支持
set(CMAKE_AUTOMOC ON)
# 启用 qrc 资源文件的支持
set(CMAKE_AUTORCC ON)

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
# 根据 cmake 文档描述此命令已经有另外几个功能划分的更为细致的命令所取代，具体请查阅文档
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
```

`home_monitor.json`文件指明了当前插件所使用的 dde-dock 的接口版本，dde-dock 在加载此插件时，会检测自己的接口版本是否与插件的接口版本一致，当双方的接口版本不一致或者不兼容时，dde-dock 为了安全将阻止加载对应的插件。

在 dde-dock 内建的插件代码中，可以找到当前具体的接口版本，目前最新的版本是 `1.2` 。

``` json
{
    "api": "1.2"
}
```

另外(可选的)还支持指定一个 dbus 服务，dock 在加载插件时会检查此插件所依赖的 dbus 服务，如果服务没有启动则不会初始化这个插件，直到服务启动，
如下表示依赖地址为 "com.deepin.daemon.Network" 的 dbus 服务。

``` json
{
    "api": "1.2",
    "depends-daemon-dbus-service": "com.deepin.daemon.Network"
}
```

`homemonitorplugin.h` 包含了类 `HomeMonitorPlugin`，它继承自 `PluginItemInterface`，这代表了它是一个实现了 dde-dock 接口的插件。

`PluginItemInterface` 中包含众多的功能接口以丰富插件的功能，具体的接口功能与用法可以查看对应文件中的文档。大多数接口在没有特定需求的时候都是无需处理的，需要所有插件显式处理的接口只有 `pluginName`、`init`、`itemWidget` 三个接口。

`PluginItemInterface` 中的接口都是被动的，即插件本身无法确定这些接口什么时刻会被调用，为此 dde-dock 的插件机制还提供了 `PluginProxyInterface`， `PluginProxyInterface` 的具体实例将会通过 `PluginItemInterface` 中的 `init` 接口传递给插件，因此在 `init` 接口中总是会先把这个传入的对象保存起来以供后续使用。`PluginProxyInterface` 中提供了一些可供插件随时调用的接口，以便让 dock 相应插件的请求，开发者们可以将 `PluginProxyInterface` 视为 dde-dock 中所有插件的管理者。

`PluginItemInterface` 和 `PluginProxyInterface` 的文件可以打开其头文件查看：

```
/usr/include/dde-dock/pluginproxyinterface.h
/usr/include/dde-dock/pluginsiteminterface.h
```

下面是最小化实现了一个 dock 插件的源码，请注意，本文的代码只是为了简述开发一个插件的主要过程，代码可能已经过时或不完整，详细的示例代码应该查看 `home-monitor` 目录下的内容。

`homemonitorplugin.h` 中包含对应接口的声明：

``` c++
#ifndef HOMEMONITORPLUGIN_H
#define HOMEMONITORPLUGIN_H

#include <dde-dock/pluginsiteminterface.h>

#include <QObject>

class HomeMonitorPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "home_monitor.json")

public:
    explicit HomeMonitorPlugin(QObject *parent = nullptr);

    // 返回插件的名称，必须是唯一值，不可以和其它插件冲突
    const QString pluginName() const override;
    // 插件初始化函数
    void init(PluginProxyInterface *proxyInter) override;

    // 返回插件的 widget
    QWidget *itemWidget(const QString &itemKey) override;
};

#endif // HOMEMONITORPLUGIN_H
```

`homemonitorplugin.cpp` 中包含对应接口的实现

``` c++
#include "homemonitorplugin.h"

HomeMonitorPlugin::HomeMonitorPlugin(QObject *parent)
    : QObject(parent)
{

}

const QString HomeMonitorPlugin::pluginName() const
{
    return QStringLiteral("home_monitor");
}

void HomeMonitorPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;
}

QWidget *HomeMonitorPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    // 这里暂时返回空指针，这意味着不会有任何东西被添加到 dock 上
    return nullptr;
}
```

## 测试插件加载
### 从源码构建
为了不污染源码目录，推荐在源码目录中创建`build`进行构建：

``` sh
cd home_monitor
mkdir build
cd build
cmake ..
make -j4
```
### 安装
```
sudo make install
```

可以看到有`home_monitor.so`文件被安装在了 dde-dock 的插件目录。

``` sh
install -m 755 -p ./home_monitor/libhome_monitor.so /usr/lib/dde-dock/plugins/libhome_monitor.so
```

### 测试加载
执行 `pkill dde-dock; dde-dock` 来重新运行 dde-dock，在终端输出中如果出现以下的输出，说明插件的加载已经正常。
```
init plugin:  "home_monitor"
init plugin finished:  "home_monitor"
```

## 创建自己的 widget
按照一般的业务逻辑处理，这部分不是本教程的重点，可以参考完整代码及其它插件进行实现，这里只简单使用一个 QLable。

在 `homemonitorplugin.h` 中相应位置添加成员声明：

``` c++
class HomeMonitorPlugin : public QObject, PluginsItemInterface
{
private:
    QWidget *m_pluginWidget;
};
```

然后在 `homemonitorplugin.cpp` 中将添加的成员初始化，比如在 `init` 接口中初始化：

``` c++
void HomeMonitorPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    m_pluginWidget = new QLable("Hello Dock");
}
```

## 添加 widget 到 dde-dock 面板上
在 `init` 方法中获取到了 `PluginProxyInterface` 对象，调用此对象的 `itemAdded` 即可实现向 dde-dock 面板上添加项目。
第二个 `QString` 类型的参数代表了本插件所提供的 item 的 id，当一个插件提供多个 item 时，不同 item 之间的 id 要保证唯一。

``` c++
void HomeMonitorPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    m_pluginWidget = new QLable("Hello Dock");

    m_proxyInter->itemAdded(this, QString());
}
```

在调用 `itemAdded` 之后，dde-dock 会在合适的时机调用插件的`itemWidget`接口以获取需要显示的 widget。在 itemWidget 接口中分析 itemKey，返回与之对应的 widget 对象，当插件只有一个可显示项目时，itemKey 可以忽略 (但不建议忽略)。

``` c++
QWidget *HomeMonitorPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_pluginWidget;
}
```
