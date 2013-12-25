package main

import (
	"dbus/com/deepin/daemon/setdatetime"
	"dbus/org/gnome/settingsdaemon/datetimemechanism"
	"dlib"
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
	"fmt"
)

const (
	_DATE_TIME_DEST = "com.deepin.daemon.DateAndTime"
	_DATE_TIME_PATH = "/com/deepin/daemon/DateAndTime"
	_DATA_TIME_IFC  = "com.deepin.daemon.DateAndTime"

	_DATE_TIME_SCHEMA = "com.deepin.dde.datetime"
)

var (
	_busConn     *dbus.Conn
	_dtGSettings = gio.NewSettings(_DATE_TIME_SCHEMA)

	_setDT, _ = setdatetime.NewSetDateTime("/com/deepin/daemon/setdatetime")
	_gdate, _ = datetimemechanism.NewDateTimeMechanism("/")
)

type DateTime struct {
	AutoSetTime      bool
	Use24HourDisplay *property.GSettingsBoolProperty `access:"readwrite"`
	CurrentTimeZone  string
}

func (date *DateTime) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_DATE_TIME_DEST, _DATE_TIME_PATH, _DATA_TIME_IFC}
}

func (date *DateTime) SetDate(d string) bool {
	ret := _setDT.SetCurrentDate(d)
	return ret
}

func (date *DateTime) SetTime(t string) bool {
	ret := _setDT.SetCurrentTime(t)
	return ret
}

func (date *DateTime) SetTimeZone(zone string) {
	_gdate.SetTimezone(zone)
	date.CurrentTimeZone = zone
	dbus.NotifyChange(date, "CurrentTimeZone")
}

func (date *DateTime) SetAutoSetTime(auto bool) bool {
	var ret bool
	if auto {
		date.AutoSetTime = true
		ret = _setDT.SetNtpUsing(true)
	} else {
		date.AutoSetTime = false
		ret = _setDT.SetNtpUsing(false)
	}
	dbus.NotifyChange(date, "AutoSetTime")
	return ret
}

func (date *DateTime) SyncNtpTime() bool {
	ret := _setDT.SyncNtpTime()
	return ret
}

func NewDateAndTime() *DateTime {
	dt := &DateTime{}

	dt.Use24HourDisplay = property.NewGSettingsBoolProperty(dt, "Use24HourDisplay", _dtGSettings, "is-24hour")
	dt.CurrentTimeZone = _gdate.GetTimezone()

	dt.AutoSetTime = _dtGSettings.GetBoolean("is-auto-set")
	_dtGSettings.Connect("changed::is-auto-set", func(s *gio.Settings, name string) {
		fmt.Println("is-auto-set changed:", s.GetBoolean("is-auto-set"))
		dt.SetAutoSetTime(s.GetBoolean("is-auto-set"))
	})

	return dt
}

func main() {
	date := NewDateAndTime()
	err := dbus.InstallOnSession(date)
	if err != nil {
		panic(err)
	}

	if date.AutoSetTime {
		go _setDT.SetNtpUsing(true)
	}
	dlib.StartLoop()
}
