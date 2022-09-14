# DDE Dock

DDE Dock是深度桌面环境的任务栏组件。

构建 dde-dock 插件的教程：[plugins-developer-guide](plugins/plugin-guide/plugins-developer-guide.md)

### 依赖

请查看“debian/control”文件中提供的“Depends”。

### 编译依赖

请查看“debian/control”文件中提供的“Build-Depends”。

## 安装

### 构建过程

1. 确保已经安装了所有的编译依赖

2. 构建
```
$ cd dde-dock
$ mkdir Build
$ cd Build
$ cmake ..
$ make
```

3. 安装
```
$ sudo make install
```

## 帮助

* [Matrix](https://matrix.to/#/#deepin-community:matrix.org)
* [WiKi](https://wiki.deepin.org)
* [官方论坛](https://bbs.deepin.org)
* [开发者中心](https://github.com/linuxdeepin/developer-center/issues) 

## 贡献指南

我们鼓励您报告问题并做出更改

* [Contribution guide for developers](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers-en). (English)
* [开发者代码贡献指南](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers) (中文)

## 开源许可证

dde-dock在 [LGPL-3.0-or-later](LICENSE)下发布。