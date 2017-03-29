## DDE Daemon

DDE Daemon is a daemon for handling  the deepin session settings

## Dependencies


### Build dependencies

* [dde-api](https://github.com/linuxdeepin/dde-api)
* [startdde](https://github.com/linuxdeepin/startdde)
* libudev
* fontconfig
* libbamf3

### Runtime dependencies

* upower
* udisks2
* acpid
* bluez
* systemd
* pulseaudio
* network-manager
* policykit-1-gnome
* grub-themes-deepin
* gnome-keyring
* deepin-notifications
* xserver-xorg-input-wacom
* libinput
* xdotool

### Optional Dependencies

* network-manager-vpnc-gnome
* network-manager-pptp-gnome
* network-manager-l2tp-gnome
* network-manager-strongswan-gnome
* network-manager-openvpn-gnome
* network-manager-openconnect-gnome
* iso-codes
* iw (check if wireless device support hotspot mode)
* mobile-broadband-provider-info
* xserver-xorg-input-synaptics (provide mode features, such as disable touchpad when typing ...)

## Installation

Build:
```
$ make GOPATH=/usr/share/gocode
```

Or, build through gccgo
```
$ make GOPATH=/usr/share/gocode USE_GCCGO=1
```

Install:
```
sudo make install
```

## Usage

### dde-system-daemon

`dde-system-daemon` primarily provide account services, need to run as root.

### dde-session-daemon

#### Flags:

```
memprof      : Write memory profile to specific file
cpuprof      : Write cpu profile to specific file, can not use memprof and
               cpuprof together
-i --Ignore  : Ignore missing modules, --no-ignore to revert it, default is true
-v --verbose : Show much more message, the shorthand for --loglevel debug,
               if specificed, loglevel is ignored
-l --loglevel: Set log level, possible value is error/warn/info/debug/no
```

#### Commands:

```
list   : List all the modules or the dependencies of one module.
auto   : Automatically get enabled and disabled modules from settings.
enable : Enable modules and their dependencies, ignore settings.
disable: Disable modules, ignore settings.
```

## Getting help

Any usage issues can ask for help via

* [Gitter](https://gitter.im/orgs/linuxdeepin/rooms)
* [IRC channel](https://webchat.freenode.net/?channels=deepin)
* [Forum](https://bbs.deepin.org/)
* [WiKi](http://wiki.deepin.org/)

## Getting involved

We encourage you to report issues and contribute changes.

* [Contribution guide for users](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Users)
* [Contribution guide for developers](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Developers)

## License

DDE Daemon is licensed under [GPLv3](LICENSE).
