package main

import (
	"dbus-gen/SetDateTime"
	"dbus-gen/datetime"
	"dlib"
	"dlib/dbus"
	"dlib/dbus/property"
	"fmt"
)

const (
	_DATE_TIME_DEST = "com.deepin.daemon.DateAndTime"
	_DATE_TIME_PATH = "/com/deepin/daemon/DateAndTime"
	_DATA_TIME_IFC  = "com.deepin.daemon.DateAndTime"

	_DATE_TIME_SCHEMA = "com.deepin.dde.datetime"
)

var (
	dtGSettings *dlib.Settings
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

func (date *DateTime) SetDate(d string) {
	setDT.SetCurrentDate(d)
}

func (date *DateTime) SetTime(t string) {
	setDT.SetCurrentTime(t)
}

func (date *DateTime) SetTimeZone(zone string) {
	gdate := datetime.GetDateTimeMechanism("/")
	gdate.SetTimezone(zone)
	date.CurrentTimeZone = zone
}

func (date *DateTime) SetAutoSetTime(auto bool) {
	if auto {
		date.AutoSetTime = true
		setDT.SetNtpUsing(true)
	} else {
		date.AutoSetTime = false
		setDT.SetNtpUsing(false)
	}
}

func (date *DateTime) SyncNtpTime() {
	setDT.SyncNtpTime()
}

func NewDateAndTime() *DateTime {
	dt := DateTime{}
	dtGSettings = dlib.NewSettings(_DATE_TIME_SCHEMA)

	dt.TimeShowFormat = property.NewGSettingsPropertyFull(dtGSettings,
		"is-24hour", true, busConn, _DATE_TIME_DEST, _DATA_TIME_IFC,
		"TimeShowFormat")
	d := datetime.GetDateTimeMechanism("/")
	dt.CurrentTimeZone = d.GetTimezone()

	dt.AutoSetTime = dtGSettings.GetBoolean("is-auto-set")
	dtGSettings.Connect("changed::is-auto-set", func(s *dlib.Settings, name string) {
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
		setDT.SetNtpUsing(true)
	}
	fmt.Println("Start Loop ...")
	dlib.StartLoop()
}
