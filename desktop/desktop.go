package main

import "dlib/dbus"

type DesktopManager struct {
	ComputerShow bool
	HomeDIrShow  bool
	TrashShow    bool
	SoftCenter   bool

	DockShowMode string

	LeftTop     string
	RightBottom string
}

func (desk *DesktopManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Desktop",
		"/com/deepin/daemon/Desktop",
		"com.deepin.daemon.Desktop",
	}
}

func (desk *DesktopManager) reset (propName string) {
}

func main() {
	desk := DesktopManager{}
	dbus.InstallOnSession(&desk)
	select {}
}
