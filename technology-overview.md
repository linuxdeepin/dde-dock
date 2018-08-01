
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

其中，AppItem 与 PluginsItem 是比较复杂的类型，详细说明：

#### AppItem

AppItem 是在 dock 上最经常与用户交互的类型，它关联着每个应用程序的窗口列表。所有的数据都是从后端(dde-initializer)的 DBus 服务中所获取的，具体的数据格式可以参考后端相关文档。

这里需要注意的是，后端数据分为两大块。一个是整体上的，即获取系统中有哪些需要显示在 dock 上的应用程序；另一个是每个应用程序，它有哪些窗口，应用程序的属性及它的各个窗口的属性等数据。

由于 dock 的管理单元是每个 Item，所以有几个应用程序，即总共应该创建几个 Item 这种控制策略应该由 ItemController 负责，而每个应用程序的窗口数据、属性数据，则由对应的 AppItem 自己去负责。由于都是读取同一个 DBus 服务，在这部分的数据处理一定要分清楚模块，否则会使整体上的数据流比较混乱。

##### Window Preview

窗口预览是应用程序类型特有的一个功能，由 `item/components` 下面的几个类提供。这部分的代码被封装在了 AppItem 内部，并利用 `DockItem` 标准的显示 Popup 的接口来显示预览窗口。这部分的代码比较独立，只与 AppItem 自己的实现有关。

#### PluginsItem

PluginsItem 是与插件所注册的某个具体 Item 相关联，__并不是与某个插件进行直接关联__。因为一个插件可能注册多个 Item，也可能一个 Item 也不注册。

PluginsItem 是一个对外来控件的包装类，所以在这里面大多工作都是将 DockItem 的一些事件或者行为转发或者加入到外来控件上，实现对外来控件的一个控制效果。

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

### Popup Window

`Popup Window` 是一个特殊的控件。它是所有 Item 中用来显示提示信息，或是显示弹出式控件、列表的一个容器。但是考虑到这种控件使用频率非常高，如果在每个 Item 中创建多个的话很浪费内存，所以将 `PopupWindow` 设计为一个全局的控件，所有的 Item 共用这个控件。

## 插件

插件是符合标准的 Qt Plugins。插件的开发不必熟悉 dock 的所有代码，只需要熟悉一般的 Qt 插件开发过程，并了解 dock 所提供的接口。dock 的接口安装 `dde-dock-dev` 包即可。这也是方便插件开发者在无需配置完整的 dock 开发环境的情况下，更方便的进行 dock 插件的开发。

### 插件的一般组织形式

一般来说，一个插件由一个主控制类和至少一个控件类组成。控制类通过 dock 的插件接口与主程序通信，并获知当前 dock 的一些状态。通过插件自己的业务需求和 dock 的状态，可以调用接口添加新的 Items 到 dock 面板上，或是从面板上删除之前自己添加的 Items。

对于插件请求创建的每个 Item，主程序都会调用插件的主控制类获取一个 Widget 作为显示内容，并创建一个 PluginsItem 对此 Widget 进行包装。包装后的 PluginsItem 将会作为标准的 DockItem 注册到 MainPanel 上显示出来。

### 插件开发中的调试方法

在加载插件失败时，主程序会打印相关信息，仔细参考相关日志即可发现大部分问题。一般就是对应插件的 so 中有某些符号没有解析成功，或是插件版本与主程序的版本不相同。

如果插件可以成功加载，即可使用 gdb 等程序进行调试。

# 接口设计

## 插件接口

插件接口定义在 `interfaces/*.h` 中，参考具体类或函数的注释。

## DBus 接口

dock 主程序提供了一个 DBus 服务，可供外部访问到 dock 的 geometryRect 信息。这个信息也可以从后端的接口中读取，创建这个接口的最初目的是用作调试。当 dock 位置不正确时，可以比对此接口的信息与后端的信息。

其中后端的信息代表了 dock 主程序计算的结果，如果这个数据错误，就说明 dock 在位置计算的部分有 bug。

如果后端信息正确，而本接口中的数据错误，那就是计算正确，但是在向 X Server 发送对应的控制请求时出错。

通过检查两个接口的数据是否相同，也是一个进行自我检查的方法。目前在 `MainWindow::positionCheck` 中就进行了这样的操作，当发现两个数据不相同时，就重新向 X Server 发送请求，以此来 workaround 某些情况下 dock 位置不正确的问题。

# 优化

## MainWindow

目前所有的动画都被放在了 `MainWindow` 中进行处理。经过多次的改动，现在 `MainWindow` 中已经有很多动画相关的代码，这使得在进行窗口管理时不得不考虑动画的很多事情。

以后可以尝试将动画部分剔除出来，`MainWindow` 只进行窗口位置、大小等操作，尤其是应该把涉及到 `MainPanel` 动画的部分移动到 `MainPanel` 类中去，以此来减少在 `MainWindow` 中控制其它控件所带来的混乱。

## Popup Window

共用 `PopupWindow` 带来了很多好处，但是如果在使用时没有好好处理 data race、或是在 Tips Window 与 Model Window 切换中没有处理好顺序，就会造成很难调试，也很难处理的问题。可以尝试在这方面做一些优化。
