# HACKING guide for Dock

## Project layout

### Coding layout

Dock developed by Qt(QWidget).

It has five packages

- src/controller
This package control most of the non-visible data communication.
- src/dbus
This package include all of the DBus backend header that programs use.
- src/interfaces
This package include the header which export for plugins development. Those header should be packaged in a single software package in the future.
- src/panel
Panel was not the top level window. This package contains the most relevant widget of panel.
- src/widgets
The rest of the all other visible widgets.

### Others

- ~/.cache/deepin/dde-dock/
Store log files.
- ~/.config/deepin/dde-dock/style/ & /usr/share/dde-dock/style/
- Store style files.

## Core Design

### 布局管理

- 布局动画的实现

> - MovableLayout
>
>> 可拖拽排序的布局，是App布局（DockAppLayout）和插件布局（DockPluginLayout）的基类。
>
> - MovableSpacingItem
>
>> 此类继承于QFrame，不含有可视内容。在MovableLayout的某个位置插入此类的对象，这样就相当于插入一段空白项。并且此类在size改变时会使用动画处理，这样就能达到空白是有过渡动画的插入效果。当需要在新的位置插入空白项时，发出信号使之前的所有空白项都开始销毁，销毁过程也由动画控制以达到过度效果。当鼠标拖拽着快速移动时因为不断的插入和销毁，就能实现布局内控件的挤压移动动画。
>

- 布局鼠标事件处理

> - 鼠标移动时开始检测是否达到拖拽阈值（为防止轻微点击误操作），如果达到则开始拖拽处理。
> - 在dragEnter或者dragMove时根据鼠标的位置计算出当前悬停在哪个布局中的控件之上，根据需要插入空白项或者移除空白项；在dragLeave时通知所有空白项销毁；
> - 在drop事件发生时，根据drop之前布局的布局方式以及移动方向确定新项目的位置。并且drop事件传递的QDropEvent对象会通过信号向外发送，外部类或子类可根据该对象做进一步处理。

### App管理

- 把App固定到任务栏

> - 从Launcher拖拽项目到Dock的App布局中（drag的mime数据中有特殊标志值）
> - 从Launcher使用右键菜单发送到Dock
> - 从桌面或文件管理器中拖拽.desktop文件到Dock的App布局中（有特殊标志值的就调用后端提供的DBus方法请求固定到Dock上，如果没有特殊标志值则当成普通文件处理）

- 把App从任务栏移除

> - 使用右键菜单的“移除驻留”选项移除，调用后端移除函数进行移除
> - 拖拽到Dock显示区域外丢弃移除
当检测到鼠标位置离开布局的有效区域后，生成DropMask的一个对象跟随鼠标移动，该对象可以接受drop事件，并且大小正好就是拖拽时图像的大小，始终位于鼠标正下方。一旦鼠标在布局区域以外释放，则该对象接受相应数据作销毁动画后发出销毁信号完成项目的移除操作。

- App排序

>排序完成后调用后端的排序函数进行位置记录
- App的激活与窗口切换
> - 直接调用后端的Activate函数实现App的激活与窗口切换（如果有多个窗口）
> - 每个预览窗口都有对应的Xid，当点击某个预览窗口时，调用后端的ActiveWindow方法实现窗口切换。

- 文件拖拽

> - 有效的.desktop文件尝试固定到Dock的App布局中
> - 普通文件调用对应App的后端的HandleDragDrop方法尝试打开

### 插件系统

- 插件预检测

> - 检测插件文件是否有效
> - 检测插件版本是否有效
> - 检测插件是否为系统级别插件

- 插件通信

> - 插件主动向Dock传递数据
插件内部通过在init函数中保存DockPluginProxyInterface对象以随时调用该类提供的方法向Dock传递数据
> - Dock主动向插件传递数据
目前在插件接口类DockPluginInterface中实现的三个方法可以主动向插件传递数据，分别是：changeMode、setEnabled和invokeMenuItem

- 插件设置

>在每个插件初始化前（调用init函数前），插件设置窗口对象已经与传递到插件内的DockPluginProxyInterface子类DockPluginProxy的对应对象做好相应信号的连接。当插件初始化完成后，即可调用DockPluginProxy的infoChangedEvent方法（参数为InfoTypeConfigurable、InfoTypeEnable或者InfoTypeTitle），这样设置窗口就会读取插件的信息并显示到设置窗口中。

### 显示模式管理

- 不同显示模式的常量值（特别是size）

> - 常量值在dockconstants.h中定义
> - DockModeData提供的单例对象控制在不同显示模式下常量的使用

- 显示模式改变的通知

> 每个与显示模式相关的模块或者类都应该连接DockModeData类的dockModeChanged信号，在收到该信号后，重新调用该对象的方法（如：getDockHeight）即可获取不同显示模式下合适的常量值。

### 隐藏模式管理

- 移动鼠标切换Dock的隐藏与显示状态

>实际上Dock隐藏时只是把主窗口高度降低到了 1 像素，一个像素依然足够接受鼠标的移入事件。鼠标移入后即显示Dock主窗口，当鼠标移出到Dock的区域以外后使用后端提供的方法使主窗口隐藏。

- 移动窗口切换Dock的隐藏与显示状态

>由后端检测，只需要根据状态改变显示或者隐藏Dock主窗口即可

- 使用快捷键(Super+H)切换Dock的隐藏与显示状态

>由后端检测，只需要根据状态改变显示或者隐藏Dock主窗口即可。但是需要注意在按下快捷键时，鼠标是否停留在Dock主窗口区域的处理。

### 预览窗口

- 预览窗口的移动

> - 所有App的预览窗口内容与所有插件提供的Applet内容都是共用一个PreviewWindow对象的
> - 当鼠标从某个有预览内容的项移动到另一个同样有预览内容的项时，用当前鼠标所在的项的预览内容填充预览窗口（填充过程会加入延时以达到流畅过度的效果）即可，填充完成后请主动调用预览窗口的resizeWithContent函数和showPreview函数以更新预览窗的大小和位置。

- 每项的标题窗口

>每个DockItem都会有一个PreviewWindow对象用于显示标题。当没有预览内容时将使用该对象显示标题。

### 样式管理

- 样式表文件说明
目前样式表控制还有不完善的问题，部分控件的显示不是由样式表控制（如：预览窗口的外边框）。以后会针对样式表的内容整理出一份文档。

- 第三方样式表

> - 第三方样式安装目录
>
>> - /usr/share/dde-dock/style/
>> - ~/.config/deepin/dde-dock/style/
>
> - 第三方样式文件结构
>
>> - 必须在第三方样式安装目录下新建一个以样式名称命名的目录
>> - 以样式名称命名的目录下必须包含一个主要的样式表文件：style.qss
>> - 以样式名称命名的目录下可以拥有任意数量的其他资源
>

- 运行过程中切换样式表（如果未来系统统一了样式切换）
样式表的切换接口以DBus接口的方式给出

> - DBus接口

>> com.deepin.dde.dock
>
> - 相关方法
>
>> - currentStyleName(): 返回当前Dock正在使用的样式名称
>> - styleNameList(): 返回当前Dock默认样式和所有第三当样式名称的列表
>> - applyStyle(String styleName): 切换到styleNameList列表中某个有效的样式
>

### 特殊项

- Launcher启动项
使用QProcess启动dde-launcher进程。作为单独的项载入到DockPanel面板中。
- 未来可能加入的启动项
暂无

### 其他扩展

- 竖直方向布局
目前与布局相关的类是支持竖直方向布局的。只要根据需要改进DockPanel和MainWidget类即可。未来如果设计有这样的需求就可以这样做。

## UML & Diagram

### [UML File(Umbrello)](https://github.com/linuxdeepin/dde-dock/files/115295/dde-dock.xmi.zip)
### Diagram
- Use Case Diagram
![dde-dock](https://cloud.githubusercontent.com/assets/5242852/12776747/808aeee8-ca93-11e5-942a-c77ac08d5210.png)
- Class Diagram
1. mainwidget
![mainwidget](https://cloud.githubusercontent.com/assets/5242852/12776749/80a3cc1a-ca93-11e5-9431-1c06bacd828c.png)
2. layout
![layout](https://cloud.githubusercontent.com/assets/5242852/12776751/84e7b282-ca93-11e5-8838-fe9d1e5bca03.png)
3. appmanager
![appmanager](https://cloud.githubusercontent.com/assets/5242852/12776746/807ffd3a-ca93-11e5-8419-3fc2fb45371a.png)
4. pluginmanager
![pluginmanager](https://cloud.githubusercontent.com/assets/5242852/12776750/80b6d742-ca93-11e5-800f-a7c8e42a838e.png)
5. item
![item](https://cloud.githubusercontent.com/assets/5242852/12776748/8099af3c-ca93-11e5-8a17-0f036d18dd20.png)

## List of TODO (It’s the good way for contributing)

None.

## List of Workaround

None.

## Others

- [DDE Dock Plugin Development](https://github.com/linuxdeepin/developer-center/wiki/DDE-Dock-Plugin-Development)