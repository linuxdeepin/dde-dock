package main

import "dlib/dbus"

type SystemInfo struct {
	Version    string
	Processor  string
	MemorySize string
	SystemType string
	DiskCap    string
}

func (sys *SystemInfo) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.SystemInfo",
		"/com/deepin/daemon/SystemInfo",
		"com.deepin.daemon.SystemInfo",
	}
}

func main() {
	sys := SystemInfo{}

	sys.Version =		"2013"
	sys.Processor =		"i3 310M"
	sys.MemorySize =	"4G"
	sys.SystemType =	"64 bit"
	sys.DiskCap =		"500G"

	dbus.InstallOnSession(&sys)
	select {}
}
