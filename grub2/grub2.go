package main

import (
	"dlib/dbus"
)

const (
	_GRUB2_SETTINGS_SERV = "com.deepin.daemon.Grub2"
	_GRUB2_SETTINGS_PATH = "/com/deepin/daemon/Grub2"
	_GRUB2_SETTINGS_IFC  = "com.deepin.daemon.Grub2"
)

type Grub2Settings struct {
}

func (g *Grub2Settings) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_GRUB2_SETTINGS_SERV,
		_GRUB2_SETTINGS_PATH,
		_GRUB2_SETTINGS_IFC,
	}
}

func main() {
	grub := &Grub2Settings{}
	err := dbus.InstallOnSession(grub)
	if err != nil {
		panic(err)
	}
	select {}
}
