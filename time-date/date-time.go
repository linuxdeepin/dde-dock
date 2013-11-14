package main

import (
	"dlib"
	"dlib/dbus"
)

const (
	_DATE_TIME_DEST = "com.deepin.daemon.DataAndTime"
	_DATE_TIME_PATH = "/com/deepin/daemon/DateAndTime"
	_DATA_TIME_IFC  = "com.deepin.daemon.DataAndTime"

	_DATE_TIME_SCHEMA = "com.deepin.dde.datetime"
)

type DateTime struct {
	AutoSetTime     dbus.Property
	TimeShowFormat  dbus.Property
	CurrentTimeZone bool
}

func (date *DateTime) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_DATE_TIME_DEST, _DATE_TIME_PATH, _DATA_TIME_IFC}
}

func GetDateAndTime() (dt DateTime) {
	busType, _ := dbus.SystemBus()
	dtSettings := dlib.NewSettings(_DATE_TIME_SCHEMA)

	dt.AutoSetTime = property.NewGSettingsPropertyFull(dtSettings,
		"is-auto-set", true, busType, _DATE_TIME_PATH, _DATA_TIME_IFC,
		"AutoSetTime")
	dt.TimeShowFormat = property.NewGSettingsPropertyFull(dtSettings,
		"is-24hour", true, busType, _DATE_TIME_DEST, _DATA_TIME_IFC,
		"TimeShowFormat")

	return dt
}

func main() {
	date := GetDateAndTime()
	dbus.InstallOnSession(&date)
	select {}
}
