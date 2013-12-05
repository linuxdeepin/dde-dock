package main

import (
	"dbus/org/freedesktop/upower"
	"dlib"
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
)

type dbusBattery struct {
	bus_name    string
	object_path string

	//device upower.device
}

const (
	power_bus_name               = "com.deepin.daemon.Power"
	power_object_path            = "/com/deepin/daemon/Power"
	power_interface              = "com.deepin.daemon.Power"
	schema_gsettings_power       = "org.gnome.settings-daemon.plugins.power"
	schema_gsettings_screensaver = "org.gnome.desktop.screensaver"
)

const (
	operation_suspend   = "suspend"
	operation_poweroff  = "poweroff"
	operation_hibernate = "hibernate"
)

type Power struct {
	//plugins.power keys
	ButtonHibernate dbus.Property
	ButtonPower     dbus.Property
	ButtonSleep     dbus.Property
	ButtonSuspend   dbus.Property

	CriticalBatteryAction dbus.Property
	LidCloseAcAction      dbus.Property
	LidCloseBatteryAction dbus.Property

	ShowTray dbus.Property

	SleepDisplayAc      dbus.Property
	SleepDisplayBattery dbus.Property

	SleepInactiveAcTimeout      dbus.Property
	SleepInactiveBatteryTimeout dbus.Property

	SleepInactiveAcType      dbus.Property
	SleepInactiveBatteryType dbus.Property

	CurrentPlan dbus.Property

	//upower interface
	BatteryIsPresent  bool    `access:"read"` //battery present
	BatteryPercentage float64 `access:"read"` //batter voltage
	charging          int32   `access:"read"` //charging or discharging
	PlugedIn          int32   `access:"read"` //1 for in,2 for out
	TimeToEmpty       int64   `access:"read"` //
	TimeToFull        int64   `access:"read"` //time to fully charged

	//gnome.desktop.screensaver keys
	LockEnabled dbus.Property

	powerSettings       *gio.Settings
	screensaverSettings *gio.Settings

	upowerDevice *upower.Device
}

func NewPower() (*Power, error) {
	power := Power{}
	busConn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	power.powerSettings = gio.NewSettings(schema_gsettings_power)
	power.screensaverSettings = gio.NewSettings(schema_gsettings_screensaver)
	power.CurrentPlan = property.NewGSettingsPropertyFull(
		power.powerSettings, "current-plan", "", busConn,
		power_object_path, power_interface, "CurrentPlan")
	power.ButtonHibernate = property.NewGSettingsPropertyFull(
		power.powerSettings, "button-hibernate", "", busConn,
		power_object_path, power_interface, "ButtonHibernate")
	power.ButtonPower = property.NewGSettingsPropertyFull(
		power.powerSettings, "button-power", "", busConn,
		power_object_path, power_interface, "ButtonPower")
	power.ButtonSleep = property.NewGSettingsPropertyFull(
		power.powerSettings, "button-sleep", "", busConn,
		power_object_path, power_interface, "ButtonSleep")
	power.ButtonSuspend = property.NewGSettingsPropertyFull(
		power.powerSettings, "button-suspend", "", busConn,
		power_object_path, power_interface, "ButtonSuspend")

	power.CriticalBatteryAction = property.NewGSettingsPropertyFull(
		power.powerSettings, "critical-battery-action", "", busConn,
		power_object_path, power_interface, "CriticalBatteryAction")
	power.LidCloseAcAction = property.NewGSettingsPropertyFull(
		power.powerSettings, "lid-close-ac-action", "", busConn,
		power_object_path, power_interface, "LidCloseAction")
	power.LidCloseBatteryAction = property.NewGSettingsPropertyFull(
		power.powerSettings, "lid-close-battery-action", "", busConn,
		power_object_path, power_interface, "LidCloseBatteryAction")
	power.ShowTray = property.NewGSettingsPropertyFull(
		power.powerSettings, "show-tray", true, busConn,
		power_object_path, power_interface, "ShowTray")
	power.SleepInactiveAcTimeout = property.NewGSettingsPropertyFull(
		power.powerSettings, "sleep-inactive-ac-timeout", int32(0), busConn,
		power_object_path, power_interface, "SleepInactiveAcTimeout")
	power.SleepInactiveBatteryTimeout = property.NewGSettingsPropertyFull(
		power.powerSettings, "sleep-inactive-battery-timeout", int32(0), busConn,
		power_object_path, power_interface, "SleepInactiveBatteryTimeout")
	power.SleepDisplayAc = property.NewGSettingsPropertyFull(
		power.powerSettings, "sleep-display-ac", int32(0), busConn,
		power_object_path, power_interface, "SleepDisplayAc")
	power.SleepDisplayBattery = property.NewGSettingsPropertyFull(
		power.powerSettings, "sleep-display-battery", int32(0), busConn,
		power_object_path, power_interface, "SleepDisplayBattery")

	power.SleepInactiveAcType = property.NewGSettingsPropertyFull(
		power.powerSettings, "sleep-inactive-ac-type", "", busConn,
		power_object_path, power_interface, "SleepInactiveAcType")
	power.SleepInactiveBatteryType = property.NewGSettingsPropertyFull(
		power.powerSettings, "sleep-inactive-battery-type", "", busConn,
		power_object_path, power_interface, "SleepInactiveBatteryType")

	power.LockEnabled = property.NewGSettingsPropertyFull(
		power.screensaverSettings, "lock-enabled", true, busConn,
		power_object_path, power_interface, "WakePassword")

	power.upowerDevice = upower.GetDevice("/org/freedesktop/UPower/devices/battery_BAT0")

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
	if p.upowerDevice == nil {
		return -1
	}
	p.BatteryPercentage = p.upowerDevice.GetPercentage()
	//p.charging=
	p.PlugedIn = int32(p.upowerDevice.GetState())
	p.TimeToEmpty = p.upowerDevice.GetTimeToEmpty()
	p.TimeToFull = p.upowerDevice.GetTimeToFull()

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
