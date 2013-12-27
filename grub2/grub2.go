package main

import (
	"dlib/dbus"
	"fmt"
)

const (
	_GRUB2_SETTINGS_SERV = "com.deepin.daemon.Grub2"
	_GRUB2_SETTINGS_PATH = "/com/deepin/daemon/Grub2"
	_GRUB2_SETTINGS_IFC  = "com.deepin.daemon.Grub2"
)

type Grub2Settings struct {
	simple_parser *SimpleParser
}

func NewGrub2Settings () *Grub2Settings {
	settings := &Grub2Settings{}
	settings.simple_parser = &SimpleParser{}
	return settings
}

func (g *Grub2Settings) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_GRUB2_SETTINGS_SERV,
		_GRUB2_SETTINGS_PATH,
		_GRUB2_SETTINGS_IFC,
	}
}

func main() {
	grub := NewGrub2Settings()
	grub.simple_parser.Parse()
	fmt.Println(grub.simple_parser.timeout)
	err := dbus.InstallOnSession(grub)
	if err != nil {
		panic(err)
	}
	select {}
}
