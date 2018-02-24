# 从零构建 dde-dock 的插件
本教程将展示一个简单的 dde-dock 插件的开发过程，插件开发者可跟随此步骤为 dde-dock 创造出更多具有丰富功能的插件。

在本教程中，将创建一个可以实时显示用户家目录(`~/`)使用情况的小工具。

### 插件的工作原理
dde-dock 插件本质是一个按 Qt 插件标准所开发的共享库文件(`so`)。通过 dde-dock 预定的规范与提供的接口，共同完成 dde-dock 的功能扩展。

## 准备环境
插件的开发环境可以是任意的，只要是符合 Qt 插件规范及 dde-dock 插件规范的共享库文件，都可以被当作 dde-dock 插件载入。下面以 Qt + qmake 为例进行说明：

### 安装依赖
以 Deepin 15.5 环境为基础，至少先安装如下的包：

- dde-dock-dev
- qt5-qmake
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
└── home_monitor.pro
```

`home_monitor.pro`文件内容如下：
``` qmake
# 添加所需的 Qt 模块
QT              += widgets
# 指定生成目标为共享库
TEMPLATE         = lib
# 指定生成目标为 Qt 插件
CONFIG          += plugin c++11

# 指定生成目标的名称
TARGET          = $$qtLibraryTarget(home_monitor)
# 指定生成目标的目录
DESTDIR          = $$_PRO_FILE_PWD_
# 添加必要的文件到插件中
DISTFILES       += home_monitor.json

HEADERS += \
    homemonitorplugin.h

SOURCES += \
    homemonitorplugin.cpp

# 以下是安装相关的设定
isEmpty(PREFIX) {
    PREFIX = /usr
}

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target
```

`home_monitor.json`文件指明了当前插件所使用的 dde-dock 的接口版本，dde-dock 在加载此插件时，会检测自己的接口版本是否与插件的接口版本一致，当双方的接口版本不一致时，dde-dock 为了安全将阻止加载对应的插件。

在 dde-dock 内建的插件代码中，可以找到当前具体的接口版本，目前只有 1.0 版本。
``` json
{
	"api": "1.0"
}
```

`homemonitorplugin.h`包含了类`HomeMonitorPlugin`，它继承自`PluginItemInterface`，这代表了它是一个实现了 dde-dock 接口的插件。

`PluginItemInterface`中包含众多的功能接口以丰富插件的功能，具体的接口功能与用法可以查看对应文件中的文档。大多数接口在没有特定需求的时候都是无需处理的，需要所有插件显式处理的接口只有`pluginName`、`init`、`itemWidget`三个接口。
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

`homemonitorplugin.cpp`中包含对应接口的实现
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
qmake ..
make -j4
```
### 安装
```
sudo make install
```

可以看到有`home_monitor.so`文件被安装在了 dde-dock 的插件目录。

``` sh
install -m 755 -p ../../home_monitor/libhome_monitor.so /usr/lib/dde-dock/plugins/libhome_monitor.so
```

### 测试加载
执行`pkill dde-dock; dde-dock`来重新运行 dde-dock，在终端输出中如果出现以下的输出，说明插件的加载已经正常。
```
init plugin:  "home_monitor"
init plugin finished:  "home_monitor"
```

## 创建自己的 widget
按照一般的业务逻辑处理，这部分不是本教程的重点，可以参考完整代码及其它插件进行实现。

## 添加 widget 到 dde-dock 面板上
在`init`方法中获取到了`PluginProxyInterface`对象，调用此对象的`itemAdded`即可实现向 dde-dock 面板上添加项目。
第二个`QString`类型的参数代表了本插件所提供的 item 的 id，当一个插件提供多个 item 时，不同 item 之间的 id 要保证唯一。
``` c++
proxyInter->itemAdded(this, QString());
```

在调用`itemAdded`之后，dde-dock 会在合适的时机调用插件的`itemWidget`接口以获取需要显示的 widget。在 itemWidget 接口中分析 itemKey，返回与之对应的 widget 对象，当插件只有一个可显示项目时，itemKey 可以忽略。
``` c++
QWidget *HomeMonitorPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_pluginWidget;
}
```