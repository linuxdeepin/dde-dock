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
	busConn     *dbus.Conn
	dtGSettings = gio.NewSettings(_DATE_TIME_SCHEMA)

	setDT = setdatetime.GetSetDateTime("/com/deepin/daemon/setdatetime")
	gdate = datetimemechanism.GetDateTimeMechanism("/")
)

type DateTime struct {
	AutoSetTime      bool
	Use24HourDisplay dbus.Property `access:"readwrite"`
	CurrentTimeZone  string
}

func (date *DateTime) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_DATE_TIME_DEST, _DATE_TIME_PATH, _DATA_TIME_IFC}
}

func (date *DateTime) SetDate(d string) bool {
	ret := setDT.SetCurrentDate(d)
	return ret
}

func (date *DateTime) SetTime(t string) bool {
	ret := setDT.SetCurrentTime(t)
	return ret
}

func (date *DateTime) SetTimeZone(zone string) {
	gdate.SetTimezone(zone)
	date.CurrentTimeZone = zone
	dbus.NotifyChange(date, "CurrentTimeZone")
}

func (date *DateTime) SetAutoSetTime(auto bool) bool {
	var ret bool
	if auto {
		date.AutoSetTime = true
		ret = setDT.SetNtpUsing(true)
	} else {
		date.AutoSetTime = false
		ret = setDT.SetNtpUsing(false)
	}
	dbus.NotifyChange(date, "AutoSetTime")
	return ret
}

func (date *DateTime) SyncNtpTime() bool {
	ret := setDT.SyncNtpTime()
	return ret
}

func NewDateAndTime() *DateTime {
	dt := &DateTime{}

	dt.Use24HourDisplay = property.NewGSettingsBoolProperty(dt, "Use24HourDisplay", dtGSettings, "is-24hour")
	dt.CurrentTimeZone = gdate.GetTimezone()

	dt.AutoSetTime = dtGSettings.GetBoolean("is-auto-set")
	dtGSettings.Connect("changed::is-auto-set", func(s *gio.Settings, name string) {
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
		go setDT.SetNtpUsing(true)
	}
	dlib.StartLoop()
}
