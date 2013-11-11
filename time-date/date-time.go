package main

import "dlib/dbus"

type DateTime struct {
	CurrentDate     string
	CurrentTime     string
	AutoSetTime     bool
	TimeShowFormat  string
	CurrentTimeZone string
}

func (date *DateTime) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.DateAndTime",
		"/com/deepin/daemon/DateAndTime",
		"com.deepin.daemon.DateAndTime",
	}
}

func (date *DateTime) reset(propName string) {
}

func main() {
	date := DateTime{}
	dbus.InstallOnSession(&date)
	select {}
}
