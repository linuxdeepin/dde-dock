package main

import "dlib/dbus"

type DateTime struct {
	CurrentDate	string
	CurrentTime	string
	AutoSetTime	bool
	TimeShowFormat	string
	CurrentTimeZone	string

	CurrentDateChanged	func (curDate string)
	CurrentTimeChanged	func (curTime string)
	AutoSetTimeChanged	func (autoSet bool)
	TimeShowFormatChanged	func (format string)
	CurrentTimeZoneChanged	func (curZone string)
}

func (date *DateTime) GetDBusInfo () dbus.DBusInfo {
	return dbus.DBusInfo {
		"com.deepin.daemon.DateAndTime",
		"/com/deepin/daemon/DateAndTime",
		"com.deepin.daemon.DateAndTime",
	}
}

func (date *DateTime) SetCurrentDate (curDate string) bool {
	return true
}

func (date *DateTime) SetCurrentTime (curTime string) bool {
	return true
}

func (date *DateTime) SetAutoSetTime (autoSet bool) bool {
	return true
}

func (date *DateTime) SetAutoShowFormat (format string) bool {
	return true
}

func (date *DateTime) SetCurrentTimeZone (curZone string) bool {
	return true
}

func main () {
	date := DateTime {}
	dbus.InstallOnSession (&date)
	select {}
}
