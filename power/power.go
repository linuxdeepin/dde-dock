// +build ignore
package main

import (
	"dlib/dbus"
	"dlib"
	"os"
)

type upower struct {
	bus_name             string
	object_path          string

	//method names
	m_EnumerateDevices  string
	
	//property names
	p_CanSuspend        string
	p_LidIsPresent       string
}

type battery struct {
	bus_name                 string
	object_path              string

	m_Refresh                string

	p_IsPresent              string
	p_PowerSupply            bool
	p_Percentage             float64
	p_Voltage                float64
	p_TimeToEmpty            int64
	p_TImeToFull             int64

}

const (
	opsuspend = "suspend"
	oppoweroff= "poweroff"
	ophibernate= "hibernate"
)

const (
	sustime_0=0           //don't suspend
	sustime_5=5
	sustime_10=10
	sustime_30=30
	sustime_60=60         //after one hour
)



type Power struct {
	BatteryIsPresent        bool    `access:"read"`  //battery present
	BatteryPercentage		float64 `access:"read"`  //batter valtage
	PluginedIn              int32   `access:"read"` //power pluged in
	SuspendTime             []int32 `access:"read/write"` //with or without battery
	HandleLowPower          []string`access:"read/write"`
	HandleClosedLid         []string
}

var BAT0=battery {
	"org.freedesktop.UPower",    //bus name
	"org/freedesktop/UPower/devices/battery_BAT0",
}


func (p *Power)GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo {
		"com.deepin.daemon.Power",
		"/com/deepin/daemon/Power",
		"com.deepin.daemon.Power",
	}

}

func (p *power) Refresh() int32 {
	conn,err := dbus.SystemBus()
}


func main() {
	dbus.InstallOnSession(&Power{})
	select {}
}

