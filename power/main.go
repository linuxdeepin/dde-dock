package main

import "dlib/logger"
import "dlib/dbus"
import "dlib/dbus/property"
import "dlib/gio-2.0"
import "fmt"

var LOGGER = logger.NewLogger("com.deepin.daemon.Power").SetLogLevel(logger.LEVEL_INFO)

type Power struct {
	coreSettings *gio.Settings
	lidIsClosed  bool

	PowerButtonAction *property.GSettingsEnumProperty
	LidClosedAction   *property.GSettingsEnumProperty
	LockWhenActive    *property.GSettingsBoolProperty

	LidIsPresent bool

	CurrentPlan  int32 `dbus:"readwrite"`
	SuspendDelay int32 `dbus:"readwrite"`
	IdleDelay    int32 `dbus:"readwrite"`

	BatteryPercentage float64

	//Not in Charging, Charging, Full
	BatteryState uint32

	BatteryIsPresent bool

	OnBattery bool
}

func (*Power) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Power",
		"/com/deepin/daemon/Power",
		"com.deepin.daemon.Power",
	}
}

func NewPower() *Power {
	p := &Power{}
	p.coreSettings = gio.NewSettings("com.deepin.daemon.power")
	p.PowerButtonAction = property.NewGSettingsEnumProperty(p, "PowerButtonAction", p.coreSettings, "button-power")
	p.LidClosedAction = property.NewGSettingsEnumProperty(p, "LidClosedAction", p.coreSettings, "lid-close")
	p.LockWhenActive = property.NewGSettingsBoolProperty(p, "LockWhenActive", p.coreSettings, "lock-enabled")

	p.setPlan(int32(p.coreSettings.GetEnum("current-power-plan")))
	fmt.Println("LidClosedAction:", p.LidClosedAction.Get())

	p.initUpower()
	p.initEventHandle()

	return p
}

func main() {
	p := NewPower()
	dbus.InstallOnSession(p)
	fmt.Println("GetBattery:", getBattery().Vendor.Get())
	dbus.Wait()
}
