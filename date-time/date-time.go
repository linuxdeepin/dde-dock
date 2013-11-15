package main

import (
	"datetime"
	"dlib"
	"dlib/dbus"
	"dlib/dbus/property"
)

const (
	_DATE_TIME_DEST = "com.deepin.daemon.DateAndTime"
	_DATE_TIME_PATH = "/com/deepin/daemon/DateAndTime"
	_DATA_TIME_IFC  = "com.deepin.daemon.DateAndTime"

	_DATE_TIME_SCHEMA = "com.deepin.dde.datetime"
)

type DateTime struct {
	AutoSetTime     dbus.Property
	TimeShowFormat  dbus.Property
	CurrentTimeZone string
}

func (date *DateTime) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_DATE_TIME_DEST, _DATE_TIME_PATH, _DATA_TIME_IFC}
}

func (date *DateTime) SetDate(day, month, year uint32) {
	d := datetime.GetDateTimeMechanism("/")
	d.SetDate(day, month, year)
}

func (date *DateTime) SyncTime() {
	d := datetime.GetDateTimeMechanism("/")
	d.SyncTime()
	d.SetUsingNtp(true)
}

func (date *DateTime) SetTime(secondSinceEpoch int64) bool {
	d := datetime.GetDateTimeMechanism("/")
	b := d.CanSetTime()
	if b != 2 {
		return false
	}
	d.SetTime(secondSinceEpoch)
	return true
}

func (date *DateTime) SetTimeZone(tz string) bool {
	d := datetime.GetDateTimeMechanism("/")
	b := d.CanSetTimezone()
	if b != 2 {
		return false
	}
	d.SetTimezone(tz)
	return true
}

func NewDateAndTime() *DateTime {
	dt := DateTime{}
	busType, _ := dbus.SystemBus()
	dtSettings := dlib.NewSettings(_DATE_TIME_SCHEMA)

	dt.AutoSetTime = property.NewGSettingsPropertyFull(dtSettings,
		"is-auto-set", true, busType, _DATE_TIME_PATH, _DATA_TIME_IFC,
		"AutoSetTime")
	dt.TimeShowFormat = property.NewGSettingsPropertyFull(dtSettings,
		"is-24hour", true, busType, _DATE_TIME_DEST, _DATA_TIME_IFC,
		"TimeShowFormat")
	d := datetime.GetDateTimeMechanism("/")
	dt.CurrentTimeZone = d.GetTimezone()

	return &dt
}

func main() {
	date := NewDateAndTime()
	dbus.InstallOnSession(date)
	select {}
}
