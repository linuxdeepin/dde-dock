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

	power.powerSettings = gio.NewSettings(schema_gsettings_power)
	power.screensaverSettings = gio.NewSettings(schema_gsettings_screensaver)

	power.CurrentPlan = property.NewGSettingsProperty(
		&power, "CurrentPlan", power.powerSettings, "current-plan")
	power.ButtonHibernate = property.NewGSettingsProperty(
		&power, "ButtonHibernate", power.powerSettings, "button-hibernate")
	power.ButtonPower = property.NewGSettingsProperty(
		&power, "ButtonPower", power.powerSettings, "button-power")
	power.ButtonSleep = property.NewGSettingsProperty(
		&power, "ButtonSleep", power.powerSettings, "button-sleep")
	power.ButtonSuspend = property.NewGSettingsProperty(
		&power, "ButtonSuspend", power.powerSettings, "button-suspend")

	power.CriticalBatteryAction = property.NewGSettingsProperty(
		&power, "CriticalBatteryAction", power.powerSettings, "critical-battery-action")
	power.LidCloseAcAction = property.NewGSettingsProperty(
		&power, "LidCloseAction", power.powerSettings, "lid-close-ac-action")
	power.LidCloseBatteryAction = property.NewGSettingsProperty(
		&power, "LidCloseBatteryAction", power.powerSettings, "lid-close-battery-action")
	power.ShowTray = property.NewGSettingsProperty(
		&power, "ShowTray", power.powerSettings, "show-tray")
	power.SleepInactiveAcTimeout = property.NewGSettingsProperty(
		&power, "SleepInactiveAcTimeout", power.powerSettings, "sleep-inactive-ac-timeout")
	power.SleepInactiveBatteryTimeout = property.NewGSettingsProperty(
		&power, "SleepInactiveBatteryTimeout", power.powerSettings, "sleep-inactive-battery-timeout")
	power.SleepDisplayAc = property.NewGSettingsProperty(
		&power, "SleepDisplayAc", power.powerSettings, "sleep-display-ac")
	power.SleepDisplayBattery = property.NewGSettingsProperty(
		&power, "SleepDisplayBattery", power.powerSettings, "sleep-display-battery")

	power.SleepInactiveAcType = property.NewGSettingsProperty(
		&power, "SleepInactiveAcType", power.powerSettings, "sleep-inactive-ac-type")
	power.SleepInactiveBatteryType = property.NewGSettingsProperty(
		&power, "SleepInactiveBatteryType", power.powerSettings, "sleep-inactive-battery-type")

	power.LockEnabled = property.NewGSettingsProperty(
		&power, "LockEnabled", power.screensaverSettings, "lock-enabled")

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
