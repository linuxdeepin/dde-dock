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
		"com.deepin.dss.systeminfo",
		"/com/deepin/dss/systeminfo",
		"com.deepin.dss.systeminfo",
	}
}

func (sys *SystemInfo) GetSystemInfo() map[string]string {
	return nil
}

func main () {
        sys := SystemInfo {}
        dbus.InstallOnSession (&sys)
        select {}
}
