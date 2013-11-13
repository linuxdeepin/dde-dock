package main

import (
	"dlib"
	"dlib/dbus"
)

type DateTime struct {
	AutoSetTime     bool
	TimeShowFormat  bool
	CurrentTimeZone bool
}

func (date *DateTime) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.DateAndTime",
		"/com/deepin/daemon/DateAndTime",
		"com.deepin.daemon.DateAndTime",
	}
}

func GetTimeSettings() DateTime {
	dt := DateTime{}

	dtSettings := dlib.NewSettings("com.deepin.dde.datetime")
	dt.AutoSetTime = dtSettings.GetBoolean("is-auto-set")
	dt.TimeShowFormat = dtSettings.GetBoolean("is-24hour")

	return dt
}

func (date *DateTime) reset(propName string) {
}

func main() {
	date := GetTimeSettings()
	dbus.InstallOnSession(&date)
	select {}
}
