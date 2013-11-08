package main

import "dlib/dbus"

type ExtDevManager struct {
	ExtDevs []*ExtDevEntry
}

type ExtDevEntry struct {
	UseHabit       string
	MoveSpeed      float64
	MoveAccuracy   float64
	ClickFrequency float64
	DragDelay      float64
	DevType        string
	DbusID         string

	UseHabitChanged       func(string, string)
	MoveSpeedChanged      func(string, float64)
	MoveAccuracyChanged   func(string, float64)
	ClickFrequencyChanged func(string, float64)
	DragDelayChanged      func(string, float64)
}

func (ext *ExtDevManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.ExtDevManager",
		"/com/deepin/daemon/ExtDevManager",
		"com.deepin.daemon.ExtDevManager",
	}
}

func (ext *ExtDevManager) GetDevEntryById(id string) *ExtDevEntry {
	return NewExtDevEntry(id)
}

func (entry *ExtDevEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.ExtDevManager",
		"/com/deepin/daemon/ExtDevManager/Dev" + entry.DbusID,
		"com.deepin.daemon.ExtDevManager.DevEntry",
	}
}

func NewExtDevEntry(id string) *ExtDevEntry {
	return &ExtDevEntry{
		DbusID: id,
	}
}

func main() {
	m := ExtDevManager{}
	dbus.InstallOnSession(&m)
	e := m.GetDevEntryById("0")
	dbus.InstallOnSession(e)
	select {}
	/*timer := time.NewTimer (time.Second * 20)*/
	/*<-timer.C*/
}
