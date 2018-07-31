
# 框架设计

`dde-dock` 主要分为两部分，即主界面 `frame` 部分与插件 `plugins` 部分。

## 主界面

主界面是指 dock 栏的主窗口。主窗口负责 dock 显示的位置、形状及相关的动画与特效处理。

主面板(MainPanel)是主窗口(MainWindow)的唯一子控件，它是一个 BoxLayout，负责容纳 dock 上存放的所有项目。根据设置的不同，它的排版方向有水平和竖直两种。

### Item
dock 上所存放的所有条目都继承自 `Item` 类。这样设计的原因是简化 dock 布局，使得 dock 主面板的布局上所有条目的管理都能统一起来。

目前 dock 上的 Items 有以下几类：

- DockItem： 所有 Item 的基类，抽象类。
- LauncherItem： 启动器类型的 Item。
- AppItem： 关联所有应用程序的 Item。
- PlaceholderItem： 占位空间，用于在交换、拖拽 Item 时，显示一个临时的、空白的 Item。
- StretchItem： 在时尚模式下，充当中间的可伸缩空白区域。
- ContainerItem： 容器空间，用于存放被收纳隐藏的插件 Item。
- PluginsItem： 插件条目，每个插件可以注册 0 个或多个 PluginsItem 用于显示数据。

Dock 上从左至右（或从上至下）有不同的 Items 区域，在不同的显示模式下，不同工作区的显示状态或者调整策略都不一样，将它们统一为 Items 进行管理，可以极大的减小在这方面的工作。

### Item Controller

`DockItemController` 类是控制与管理所有 Items 的地方。任何 Item 的创建、销毁操作，移动、交换、刷新等信号的起始点都从这里开始。

其中，AppItem 的相关数据是从后端获取的。这些与后端通信的操作被封装在了 AppItem 中。Item Controller 并不处理这些具体某个 Item 的事情。

### DockItem Controller

由于插件的复杂性与特殊性，专门为插件管理加了一层包装。DockItem Controller 是 ItemController 的一部分，专门负责插件类型的 Items 的创建、排序等相关操作。同时，也是作为 dock 主程序到插件之间的一个 proxy 的作用。

### MainPanel

`MainPanel` 是主界面上的唯一控件，是容纳所有 Items 的地方。这个类接受来自 ItemController 的控制消息，来更新界面上的 Items 列表。

它主动进行的操作只有两种：

- d&d 操作的处理。它接受 drag & drop 事件，对事件进行处理并显示动画。中间过程全部是临时数据。当用户操作完毕后，它将最终的控制信号发送给 ItemsController，并由接收它发送的信号来更新界面顺序。
- 布局调整处理。在 dock 位置、大小、Items 数量等发生变化后，MainPanel 负责调整每个 Item 的大小并刷新布局。

需要注意的是，主面板类并不直接去控制 Items 列表的顺序，更不会去添加与销毁某个 Item。为了保证解耦，功能上不能与 Controller 混淆，所以对 Item 的控制操作应该 __全部__ 来自于 ItemController 的控制信号。

## 插件

插件是符合标准的 Qt Plugins。插件的开发不必熟悉 dock 的所有代码，只需要熟悉一般的 Qt 插件开发过程，并了解 dock 所提供的接口。dock 的接口安装 `dde-dock-dev` 包即可。这也是方便插件开发者在无需配置完整的 dock 开发环境的情况下，更方便的进行 dock 插件的开发。

### 插件开发中的调试方法

在加载插件失败时，主程序会打印相关信息，仔细参考相关日志即可发现大部分问题。一般就是对应插件的 so 中有某些符号没有解析成功，或是插件版本与主程序的版本不相同。

如果插件可以成功加载，即可使用 gdb 等程序进行调试。

# 接口设计

## 插件接口

插件接口定义在 `interfaces/*.h` 中，参考具体类或函数的注释。


