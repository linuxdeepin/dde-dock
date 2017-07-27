## 描述

deepin 个性化后端, 提供了GTK主题，Icon 主题，Cursor 主题, 背景和字体的管理功能。


## 功能实现简介

首先Gtk, Icon, Cursor和字体都需要先设置对应的 xsettings 属性(xsettings 的实现见 startdde)，部分还需要设置 xresources , 使用的接口见 dde-api

Cursor 还需要单独监听 `gtk-cursor-theme-name` 的改变，来让 Gtk 程序实时生效，对于 Qt 程序则需要在设置时更改下一些光标的映射关系，具体见 dde-api

背景的绘制和设置接口都是由 deepin-wm 提供的，只需要调用接口就可以了。

字体使用了 fontconfig 来获取字体列表，并更改了它的配置问题来实现实时生效的。

而主题列表的获取都是遍历相关的安装目录而的到的(字体除外)，然后会监听这些目录的变化来刷新列表。


## 目录结构

+ *subthemes, fsnotify.go*: 管理 GTK, Icon, Cursor 主题, 并监听列表的改变，发出 `Refreshed` 信号
+ *background, bg_wrapper.go*: 管理桌面背景
+ *fonts, default_font_config.go*: 管理字体，包括标准字体，等款字体及字体大小
+ *listener.go, cursor.c, cursor.h*: 处理 `gtk cursor` 的改变事件，让改变实时生效
+ *handle_gsetting.go*: 监听 `gsettings` 的改变，并应用
+ *manager.go, stup.go, ifc.go*: 个性化后端的接口
+ *appearance.go*: 个性化模块的入口


## DBus 接口简介

*Dest*: com.deepin.daemon.Appearance 
*Path*: /com/deepin/daemon/Appearance 
*Interface*: com.deepin.daemon.Appearance 


### 支持的主题类型

+ TypeGtkTheme          ("gtk")
+ TypeIconTheme         ("icon")
+ TypeCursorTheme       ("cursor")
+ TypeBackground        ("background")
+ TypeGreeterBackground ("greeterbackground")
+ TypeStandardFont      ("standardfont")
+ TypeMonospaceFont     ("monospacefont")
+ TypeFontSize          ("fontsize")


### Methods

+ List(type string) (string, error)
    获取指定类型的主题列表，返回的是json格式的字符串。如果类型错误将返回错误。
+ Show(type, name string) (string, error)
    获取指定类型主题的详细信息，包含主题名称，路径，是否可删除。如果类型错误或者主题不存在将返回错误。
+ Set(type, name string) error
    这是指定类型的主题，如果类型错误或者主题不存在将返回错误。
+ Delete(type, name string) error
    删除指定类型主题，注意只可删除用户目录下的。如果类型错误或者主题不存在将返回错误。
+ Thumbnail(type, name string) (string, error)
    获取指定类型主题的缩略图，返回的是缩略图的路径。如果类型错误或者主题不存在将返回错误。
+ Reset()
    重置所有的设置为默认值


### Properties

+ GtkTheme
    显示当前窗口主题
+ IconTheme
    显示当前图标主题
+ CursorTheme
    显示当前光标主题
+ Background
    显示当前桌面背景
+ StandardFont
    显示当前标准字体
+ MonospaceFont
    显示当前等宽字体
+ FontSize
    显示当前字体大小


### Signals

+ Changed(type, name string)
    当上面的属性改变时，会发送此信号，包含改变的属性类型及改变后的值
+ Refreshed(type string)
    当 gtk, icon, cursor, background 的安装目录改变后，有主题或壁纸被添加或删除后，就会发出此信号
