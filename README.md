# DDE Dock

DDE Dock is the dock of Deepin Desktop Environment.

A tutorial for build dde-dock plugin: [plugins-developer-guide](plugins/plugin-guide/plugins-developer-guide.md)

### Dependencies

You can also check the "Depends" provided in the `debian/control` file.

### Build dependencies

You can also check the "Build-Depends" provided in the `debian/control` file.

## Installation

### Build from source code

1. Make sure you have installed all dependencies.

2. Build:
```
$ cd dde-dock
$ mkdir Build
$ cd Build
$ cmake ..
$ make
```

3. Install:
```
$ sudo make install
```

## Getting help

- [Official Forum](https://bbs.deepin.org/) for generic discussion and help.
- [Developer Center](https://github.com/linuxdeepin/developer-center) for BUG report and suggestions.
- [Wiki](https://wiki.deepin.org/)
- [Developer Center](https://github.com/linuxdeepin/dde-dock) 

## Getting involved

We encourage you to report issues and contribute changes

* [Contribution guide for developers](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers-en). (English)
* [开发者代码贡献指南](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers) (中文)

## License

dde-dock is licensed under [LGPL-3.0-or-later](LICENSE).