package main

import (
	"dbus-gen/SetDateTime"
	"dbus-gen/gdatetime"
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
	setDT       *SetDateTime.SetDateTime
)

type DateTime struct {
	AutoSetTime     bool `access:"read"`
	TimeShowFormat  dbus.Property
	CurrentTimeZone string `access:"read"`
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
	gdate := gdatetime.GetDateTimeMechanism("/")
	gdate.SetTimezone(zone)
	date.CurrentTimeZone = zone
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
	return ret
}

func (date *DateTime) SyncNtpTime() bool {
	ret := setDT.SyncNtpTime()
	return ret
}

func NewDateAndTime() *DateTime {
	dt := DateTime{}
	dtGSettings = gio.NewSettings(_DATE_TIME_SCHEMA)

	dt.TimeShowFormat = property.NewGSettingsPropertyFull(dtGSettings,
		"is-24hour", true, busConn, _DATE_TIME_DEST, _DATA_TIME_IFC,
		"TimeShowFormat")
	d := gdatetime.GetDateTimeMechanism("/")
	dt.CurrentTimeZone = d.GetTimezone()

	dt.AutoSetTime = dtGSettings.GetBoolean("is-auto-set")
	dtGSettings.Connect("changed::is-auto-set", func(s *gio.Settings, name string) {
		fmt.Println("is-auto-set changed:", s.GetBoolean("is-auto-set"))
		dt.SetAutoSetTime(s.GetBoolean("is-auto-set"))
	})

	return &dt
}

func main() {
	var err error
	busConn, err = dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	setDT = SetDateTime.GetSetDateTime("/com/deepin/daemon/SetDateTime")
	date := NewDateAndTime()
	err = dbus.InstallOnAny(busConn, date)
	if err != nil {
		panic(err)
	}

	if date.AutoSetTime {
		go setDT.SetNtpUsing(true)
	}
	dlib.StartLoop()
}
