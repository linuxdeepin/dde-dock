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
		"com.deepin.daemon.systeminfo",
		"/com/deepin/daemon/systeminfo",
		"com.deepin.daemon.systeminfo",
	}
}

func (sys *SystemInfo) GetSystemInfo() map[string]string {
	m := make (map[string]string)

	sys.Version = "Version"
	sys.Processor = "CPU"
	sys.MemorySize = "Memory"
	sys.SystemType = "System Type"
	sys.DiskCap = "Disk Capacity"

	m[sys.Version]		= "2013"
	m[sys.Processor]	= "i3 310M"
	m[sys.MemorySize]	= "4G"
	m[sys.SystemType]	= "64 bit"
	m[sys.DiskCap]		= "500G"

	return m
}

func main () {
        sys := SystemInfo {}
        dbus.InstallOnSession (&sys)
        select {}
}
