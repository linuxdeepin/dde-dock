package main

import "dlib/dbus"

type Accounts struct {
	FaceRecog bool
}

func (info *Accounts) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Accounts",
		"/com/deepin/daemon/Accounts",
		"com.deepin.daemon.Accounts",
	}
}

func (info *Accounts) reset(propName string) {
}

func main() {
	info := Accounts{}
	dbus.InstallOnSession(&info)
	select {}
}
