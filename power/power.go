package main

// #cgo CFLAGS: -DLIBEXECDIR=""
// #cgo amd64 386 CFLAGS: -g -Wall
// #cgo pkg-config:glib-2.0 gtk+-3.0 x11 xext xtst xi gnome-desktop-3.0 upower-glib libnotify libcanberra-gtk3 gudev-1.0
// #cgo LDFLAGS: -lm
// #define GNOME_DESKTOP_USE_UNSTABLE_API
// #include "gnome-idle-monitor.h"
// #include "gsd-power-manager.h"
// int deepin_power_manager_start()
// {
//      GsdPowerManager *manager = gsd_power_manager_new();
//      GError *error = NULL;
//      gtk_init(NULL,NULL);
//      g_setenv("G_MESSAGES_DEBUG","all",FALSE);
//      notify_init("deepin-power-manager");
//      XInitThreads();
//      gsd_power_manager_start(manager,&error);
//      gtk_main();
//      return 0;
// }
import "C"

import (
	"dbus/org/freedesktop/upower"
	"dlib"
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
	//"os"
	"regexp"
	//"unsafe"
)

type dbusBattery struct {
	bus_name    string
	object_path string

	//device upower.device
}

const (
	power_bus_name    = "com.deepin.daemon.Power"
	power_object_path = "/com/deepin/daemon/Power"
	power_interface   = "com.deepin.daemon.Power"

	schema_gsettings_power               = "com.deepin.daemon.power"
	schema_gsettings_power_settings_id   = "com.deepin.daemon.power.settings"
	schema_gsettings_power_settings_path = "/com/deepin/daemon/power/profiles/"
	schema_gsettings_screensaver         = "org.gnome.desktop.screensaver"
)

type Power struct {
	//plugins.power keys
	powerProfile    *gio.Settings
	powerSettings   *gio.Settings
	ButtonHibernate *property.GSettingsStringProperty `access:"readwrite"`
	ButtonPower     *property.GSettingsStringProperty `access:"readwrite"`
	ButtonSleep     *property.GSettingsStringProperty `access:"readwrite"`
	ButtonSuspend   *property.GSettingsStringProperty `access:"readwrite"`

	CriticalBatteryAction *property.GSettingsStringProperty `access:"read"`
	LidCloseAcAction      *property.GSettingsStringProperty `access:"readwrite"`
	LidCloseBatteryAction *property.GSettingsStringProperty `access:"readwrite"`

	ShowTray *property.GSettingsBoolProperty `access:"readwrite"`

	SleepDisplayAc      *property.GSettingsIntProperty `access:"readwrite"`
	SleepDisplayBattery *property.GSettingsIntProperty `access:"readwrite"`

	SleepInactiveAcTimeout      *property.GSettingsIntProperty `access:"readwrite"`
	SleepInactiveBatteryTimeout *property.GSettingsIntProperty `access:"readwrite"`

	SleepInactiveAcType      *property.GSettingsStringProperty `access:"readwrite"`
	SleepInactiveBatteryType *property.GSettingsStringProperty `access:"readwrite"`

	CurrentProfile *property.GSettingsStringProperty `access:"readwrite"`

	//upower interface
	upower *upower.Upower

	//upower battery interface
	upowerBattery     *upower.Device
	BatteryIsPresent  dbus.Property `access:"read` //battery present
	IsRechargable     dbus.Property `access:"read`
	BatteryPercentage dbus.Property `access:"read` //
	Model             dbus.Property `access:"read`
	Vendor            dbus.Property `access:"read`
	TimeToEmpty       dbus.Property `access:"read` //
	TimeToFull        dbus.Property `access:"read` //time to fully charged
	State             dbus.Property `access:"read` //1 for in,2 for out
	Type              dbus.Property `access:"read` //type,2

	//gnome.desktop.screensaver keys
	screensaverSettings *gio.Settings
	LockEnabled         *property.GSettingsBoolProperty `access:"readwrite"`
}

func NewPower() (*Power, error) {
	power := &Power{}

	power.powerProfile = gio.NewSettings(schema_gsettings_power)
	power.CurrentProfile = property.NewGSettingsStringProperty(
		power, "CurrentProfile", power.powerProfile,
		"current-profile")

	power.powerSettings = gio.NewSettingsWithPath(
		schema_gsettings_power_settings_id,
		string(schema_gsettings_power_settings_path)+
			power.CurrentProfile.Get()+"/")
	power.screensaverSettings = gio.NewSettings(schema_gsettings_screensaver)
	power.getGsettingsProperty()

	power.upower, _ = upower.NewUpower("/org/freedesktop/UPower")
	if power.upower == nil {
		println("WARNING:UPower not provided by dbus\n")
	} else {
		println("enumerating devices\n")
		devices, _ := power.upower.EnumerateDevices()
		paths := getUpowerDeviceObjectPath(devices)
		println(paths)
		if len(paths) >= 1 {
			power.upowerBattery, _ = upower.NewDevice(dbus.ObjectPath(paths[0]))
			if power.upowerBattery != nil {
				power.getUPowerProperty()
			}
		} else {
			println("upower battery interface not found\n")
		}
	}
	return power, nil
}

func (p *Power) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Power",  //bus name
		"/com/deepin/daemon/Power", //object path
		"com.deepin.daemon.Power",
	}
}

func (power *Power) getGsettingsProperty() int32 {
	power.CurrentProfile = property.NewGSettingsStringProperty(
		power, "CurrentProfile", power.powerProfile, "current-profile")
	power.ButtonHibernate = property.NewGSettingsStringProperty(
		power, "ButtonHibernate", power.powerSettings, "button-hibernate")
	power.ButtonPower = property.NewGSettingsStringProperty(
		power, "ButtonPower", power.powerSettings, "button-power")
	power.ButtonSleep = property.NewGSettingsStringProperty(
		power, "ButtonSleep", power.powerSettings, "button-sleep")
	power.ButtonSuspend = property.NewGSettingsStringProperty(
		power, "ButtonSuspend", power.powerSettings, "button-suspend")

	power.CriticalBatteryAction = property.NewGSettingsStringProperty(
		power, "CriticalBatteryAction", power.powerSettings, "critical-battery-action")
	power.LidCloseAcAction = property.NewGSettingsStringProperty(
		power, "LidCloseAction", power.powerSettings, "lid-close-ac-action")
	power.LidCloseBatteryAction = property.NewGSettingsStringProperty(
		power, "LidCloseBatteryAction", power.powerSettings, "lid-close-battery-action")
	//power.ShowTray = property.NewGSettingsBoolProperty(
	//power, "ShowTray", power.powerSettings, "show-tray")
	power.SleepInactiveAcTimeout = property.NewGSettingsIntProperty(
		power, "SleepInactiveAcTimeout", power.powerSettings, "sleep-inactive-ac-timeout")
	power.SleepInactiveBatteryTimeout = property.NewGSettingsIntProperty(
		power, "SleepInactiveBatteryTimeout", power.powerSettings, "sleep-inactive-battery-timeout")
	//power.SleepDisplayAc = property.NewGSettingsIntProperty(
	//power, "SleepDisplayAc", power.powerSettings, "sleep-display-ac")
	//power.SleepDisplayBattery = property.NewGSettingsIntProperty(
	//power, "SleepDisplayBattery", power.powerSettings, "sleep-display-battery")

	power.SleepInactiveAcType = property.NewGSettingsStringProperty(
		power, "SleepInactiveAcType", power.powerSettings,
		"sleep-inactive-ac-type")
	power.SleepInactiveBatteryType = property.NewGSettingsStringProperty(
		power, "SleepInactiveBatteryType", power.powerSettings, "sleep-inactive-battery-type")

	power.LockEnabled = property.NewGSettingsBoolProperty(
		power, "LockEnabled", power.screensaverSettings, "lock-enabled")

	return 0
}

func (p *Power) getUPowerProperty() int32 {
	if p.upowerBattery == nil {
		return -1
	}
	p.BatteryIsPresent = property.NewWrapProperty(p, "IsPresent", p.upowerBattery.IsPresent)
	p.IsRechargable = property.NewWrapProperty(p, "IsRechargable", p.upowerBattery.IsRechargeable)
	p.BatteryPercentage = property.NewWrapProperty(p, "BatteryPercentage", p.upowerBattery.Percentage)
	p.TimeToEmpty = property.NewWrapProperty(p, "TimeToEmpty", p.upowerBattery.TimeToEmpty)
	p.TimeToFull = property.NewWrapProperty(p, "TimeToFull", p.upowerBattery.TimeToFull)
	p.Model = property.NewWrapProperty(p, "Model", p.upowerBattery.Model)
	p.Vendor = property.NewWrapProperty(p, "Vendor", p.upowerBattery.Vendor)
	p.State = property.NewWrapProperty(p, "State", p.upowerBattery.State)
	p.Type = property.NewWrapProperty(p, "Type", p.upowerBattery.Type)
	return 1
}

func (power *Power) EnumerateDevices() []dbus.ObjectPath {
	if power.upower == nil {
		println("WARNING:Upower object it nil\n")
	}
	devices, _ := power.upower.EnumerateDevices()
	for _, v := range devices {
		println(v)
	}
	return devices
}

func getUpowerDeviceObjectPath(devices []dbus.ObjectPath) []dbus.ObjectPath {
	paths := make([]dbus.ObjectPath, len(devices))
	batPattern, err := regexp.Compile(
		"/org/freedesktop/UPower/devices/battery_BAT[[:digit:]]+")
	if err != nil {
		panic(err)
	}
	linePattern, err := regexp.Compile(
		"/org/freedesktop/UPower/devices/line_power_ADP[[:digit:]]+")
	if err != nil {
		panic(err)
	}

	i := 0
	for _, path := range devices {
		ret := batPattern.FindString(string(path))
		println("findString " + ret)
		if ret == "" {
			ret = linePattern.FindString(string(path))
			if ret == "" {
				continue
			} else {
				println("findString " + ret)
				paths[1] = path
				i = i + 1
			}
		} else {
			paths[0] = path
			i = i + 1
		}
	}
	return paths[0:i]
}

func main() {
	power, err := NewPower()
	if err != nil {
		return
	}
	C.deepin_power_manager_start()
	dbus.InstallOnSession(power)
	dbus.DealWithUnhandledMessage()
	dlib.StartLoop()
	//select {}
}
