package main

import (
	"dlib"
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
	"upower"
)

type dbusBattery struct {
	bus_name    string
	object_path string

	//device upower.device
}

const (
	power_bus_name         = "com.deepin.daemon.Power"
	power_object_path      = "/com/deepin/daemon/Power"
	power_interface        = "com.deepin.daemon.Power"
	schema_gsettings_power = "org.gnome.settings-daemon.plugins.power"
)

const (
	operation_suspend   = "suspend"
	operation_poweroff  = "poweroff"
	operation_hibernate = "hibernate"
)

type Power struct {
	ButtonHibernate			    dbus.Property
	ButtonPower				    dbus.Property
	ButtonSleep				    dbus.Property
	ButtonSuspend			    dbus.Property

	CriticalBatteryAction	    dbus.Property
	LidCloseAcAcAction		    dbus.Property
	LidCloseBatteryAction	    dbus.Property

	SleepDisplayAc              dbus.Property
	SleepDisplayBattery         dbus.Property

	SleepInactiveAcTimeout      dbus.Property
	SleepInactiveBatteryTimeout dbus.Property

	SleepInactiveAcType         dbus.Property
	SleepInactiveBatteryType    dbus.Property

	CurrentPlan					dbus.Property

	BatteryIsPresent  bool     `access:"read"`       //battery present
	BatteryPercentage float64  `access:"read"`       //batter voltage
	charging          int32    `access:"read"`       //charging or discharging
	PlugedIn          int32    `access:"read"`       //1 for in,2 for out
	TimeToEmpty       int64    `access:"read"`       //
	TimeToFull        int64    `access:"read"`       //time to fully charged
	SuspendTime       []int32  `access:"read/write"` //with or without battery
}

var device *upower.Device = nil

func NewPower() (*Power, error) {
	power := Power{}
	busConn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	powerSettings := gio.NewSettings(schema_gsettings_power)
	power.CurrentPlan = property.NewGSettingsPropertyFull(
		powerSettings, "current-plan", "", busConn,
		power_object_path, power_interface, "CurrentPlan")
	power.ButtonHibernate = property.NewGSettingsPropertyFull(
		powerSettings, "button-hibernate", "", busConn,
		power_object_path, power_interface, "ButtonHibernate")
	power.ButtonPower = property.NewGSettingsPropertyFull(
		powerSettings, "button-power", "", busConn,
		power_object_path, power_interface, "ButtonPower")
	power.ButtonSleep = property.NewGSettingsPropertyFull(
		powerSettings, "button-sleep", "", busConn,
		power_object_path, power_interface, "ButtonSleep")
	power.ButtonSuspend = property.NewGSettingsPropertyFull(
		powerSettings, "button-suspend", "", busConn,
		power_object_path, power_interface, "ButtonSuspend")

	power.CriticalBatteryAction = property.NewGSettingsPropertyFull(
		powerSettings, "critical-battery-action", "", busConn,
		power_object_path, power_interface, "CriticalBatteryAction")
	power.LidCloseAcAcAction = property.NewGSettingsPropertyFull(
		powerSettings, "lid-close-ac-action", "", busConn,
		power_object_path, power_interface, "LidCloseAcAction")
	power.LidCloseBatteryAction = property.NewGSettingsPropertyFull(
		powerSettings, "lid-close-battery-action", "", busConn,
		power_object_path, power_interface, "LidCloseBatteryAction")
	power.SleepInactiveAcTimeout=property.NewGSettingsPropertyFull(
		powerSettings,"sleep-inactive-ac-timeout",int32(0),busConn,
		power_object_path,power_interface,"SleepInactiveAcTimeout")
	power.SleepInactiveBatteryTimeout=property.NewGSettingsPropertyFull(
		powerSettings,"sleep-inactive-battery-timeout",int32(0),busConn,
		power_object_path,power_interface,"SleepInactiveBatteryTimeout")
	power.SleepDisplayAc=property.NewGSettingsPropertyFull(
		powerSettings,"sleep-display-ac",int32(0),busConn,
		power_object_path,power_interface,"SleepDisplayAc")
	power.SleepDisplayBattery=property.NewGSettingsPropertyFull(
		powerSettings,"sleep-display-battery",int32(0),busConn,
		power_object_path,power_interface,"SleepDisplayBattery")

	power.SleepInactiveAcType=property.NewGSettingsPropertyFull(
		powerSettings,"sleep-inactive-ac-type","",busConn,
		power_object_path,power_interface,"SleepInactiveAcType")
	power.SleepInactiveBatteryType=property.NewGSettingsPropertyFull(
		powerSettings,"sleep-inactive-battery-type","",busConn,
		power_object_path,power_interface,"SleepInactiveBatteryType")

	device = upower.GetDevice("/org/freedesktop/UPower/devices/battery_BAT0")
	return &power, nil
}

func (p *Power) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Power",  //bus name
		"/com/deepin/daemon/Power", //object path
		"com.deepin.daemon.Power",
	}
}

func (p *Power) Refresh() int32 {
	if device == nil {
		return -1
	}
	p.BatteryPercentage = device.GetPercentage()
	//p.charging=
	p.PlugedIn = int32(device.GetState())
	p.TimeToEmpty = device.GetTimeToEmpty()
	p.TimeToFull = device.GetTimeToFull()

	return 1
}

func main() {
	power, err := NewPower()
	if err != nil {
		return
	}
	dbus.InstallOnSession(power)
	dlib.StartLoop()
	//select {}
}
