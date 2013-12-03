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
	dtGSettings *gio.Settings
	busConn     *dbus.Conn

	setDT = setdatetime.GetSetDateTime("/com/deepin/daemon/setdatetime")
	gdate = datetimemechanism.GetDateTimeMechanism("/")
)

type DateTime struct {
	AutoSetTime      bool `access:"read"`
	Use24HourDisplay dbus.Property
	CurrentTimeZone  string `access:"read"`
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
	var err error
	busConn, err = dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	dt := DateTime{}
	dtGSettings = gio.NewSettings(_DATE_TIME_SCHEMA)

	dt.Use24HourDisplay = property.NewGSettingsPropertyFull(dtGSettings,
		"is-24hour", true, busConn, _DATE_TIME_PATH, _DATA_TIME_IFC,
		"Use24HourDisplay")
	dt.CurrentTimeZone = gdate.GetTimezone()

	dt.AutoSetTime = dtGSettings.GetBoolean("is-auto-set")
	dtGSettings.Connect("changed::is-auto-set", func(s *gio.Settings, name string) {
		fmt.Println("is-auto-set changed:", s.GetBoolean("is-auto-set"))
		dt.SetAutoSetTime(s.GetBoolean("is-auto-set"))
	})

	return &dt
}

func main() {
	date := NewDateAndTime()
	err := dbus.InstallOnAny(busConn, date)
	if err != nil {
		panic(err)
	}

	if date.AutoSetTime {
		go setDT.SetNtpUsing(true)
	}
	dlib.StartLoop()
}
