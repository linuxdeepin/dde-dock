package main

import (
	"dlib"
	"dlib/dbus"
)

type DesktopManager struct {
	ShowComputerIcon bool
	ShowHomeIcon     bool
	ShowTrashIcon    bool
	ShowDSCIcon      bool

	DockShowMode int32
	LeftTop      int32
	RightBottom  int32
}

var defaults DesktopManager

func (desk *DesktopManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Desktop",
		"/com/deepin/daemon/Desktop",
		"com.deepin.daemon.Desktop",
	}
}

func (desk *DesktopManager) OnPropertiesChanged(name string, old interface{}) {
	switch name {
	case "ShowComputerIcon":
		SetShowIconBoolean(name, desk.ShowComputerIcon)
	case "ShowHomeIcon":
		SetShowIconBoolean(name, desk.ShowHomeIcon)
	case "ShowTrashIcon":
		SetShowIconBoolean(name, desk.ShowTrashIcon)
	case "ShowDSCIcon":
		SetShowIconBoolean(name, desk.ShowDSCIcon)
	case "DockShowMode":
		SetDockInt(name, desk.DockShowMode)
	case "LeftTop":
		SetPlaceAction(name, desk.LeftTop)
	case "RightBottom":
		SetPlaceAction(name, desk.RightBottom)
	}
}

func (desk *DesktopManager) reset(propName string) {
}

func InitDefaults() {
	defaults.ShowComputerIcon = true
	defaults.ShowHomeIcon = true
	defaults.ShowTrashIcon = true
	defaults.ShowDSCIcon = true
	defaults.DockShowMode = 0
	defaults.LeftTop = 2
	defaults.RightBottom = 1
}

func InitListenSignal(desk *DesktopManager) {
	deskSettings := dlib.NewSettings("com.deepin.dde.desktop")
	deskSettings.Connect("changed", func(s *dlib.Settings, name string) {
		value := deskSettings.GetBoolean(name)
		switch name {
		case "show-computer-icon":
			desk.ShowComputerIcon = value
		case "show-home-icon":
			desk.ShowHomeIcon = value
		case "show-trash-icon":
			desk.ShowTrashIcon = value
		case "show-dsc-icon":
			desk.ShowDSCIcon = value
		}
	})

	dockSettings := dlib.NewSettings("com.deepin.dde.dock")
	dockSettings.Connect("changed", func(s *dlib.Settings, name string) {
		if name == "hide-mode" {
			value := dockSettings.GetString(name)
			switch value {
			case "default":
				desk.DockShowMode = 0
			case "autohide":
				desk.DockShowMode = 1
			case "keephidden":
				desk.DockShowMode = 2
			}
		}
	})
}

func GetDesktopSettings() DesktopManager {
	desk := DesktopManager{}

	gs := dlib.NewSettings("com.deepin.dde.desktop")
	desk.ShowComputerIcon = gs.GetBoolean("show-computer-icon")
	desk.ShowHomeIcon = gs.GetBoolean("show-home-icon")
	desk.ShowTrashIcon = gs.GetBoolean("show-trash-icon")
	desk.ShowDSCIcon = gs.GetBoolean("show-dsc-icon")

	dock := dlib.NewSettings("com.deepin.dde.dock")
	mode := dock.GetString("hide-mode")
	switch mode {
	case "default":
		desk.DockShowMode = 0
	case "autohide":
		desk.DockShowMode = 1
	case "keephidden":
		desk.DockShowMode = 2
	}

	return desk
}

func SetShowIconBoolean(propName string, propValue bool) {
	show := dlib.NewSettings("com.deepin.dde.desktop")
	show.SetBoolean(propName, propValue)
}

func SetDockInt(propName string, propValue int32) {
	dock := dlib.NewSettings("com.deepin.dde.dock")

	switch propValue {
	case 0:
		dock.SetString(propName, "default")
	case 1:
		dock.SetString(propName, "autohide")
	case 2:
		dock.SetString(propName, "keephidden")
	}
}

func SetPlaceAction(propName string, propValue int32) {
}

func main() {
	desk := GetDesktopSettings()
	InitListenSignal(&desk)
	dbus.InstallOnSession(&desk)
	select {}
}
