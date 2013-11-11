package main

import (
	"dlib/dbus"
)

const (
	opsuspend = 0
	oppoweroff=1
	ophibernate=2
)

const (
	sustime_0=0           //don't suspend
	sustime_5=5
	sustime_10=10
	sustime_30=30
	sustime_60=60         //after one hour
)



type Power struct {
	BatteryPre              int32
	BatteryVoltageNow       float64
	PluginedIn              int32
	SuspendTime             []int32//with or without battery
	HandleLowPower          []int32
	HandleClosedLid         []int32
}


func (p *Power)GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo {
		"com.deepin.daemon.Power",
		"/com/deepin/daemon/Power",
		"com.deepin.daemon.Power",
	}

}

func main() {
	dbus.InstallOnSession(&Power{})
	select {}
}

