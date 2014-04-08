package main

import "dlib"
import "dlib/logger"
import "dlib/dbus"
import "dlib/dbus/property"
import "dlib/gio-2.0"
import "dbus/org/freedesktop/notifications"
import "fmt"

var LOGGER = logger.NewLogger("com.deepin.daemon.Power").SetLogLevel(logger.LEVEL_INFO)

type Power struct {
	coreSettings     *gio.Settings
	notifier         *notifications.Notifier
	lidIsClosed      bool
	lowBatteryStatus uint32

	PowerButtonAction *property.GSettingsEnumProperty `access:"readwrite"`
	LidClosedAction   *property.GSettingsEnumProperty `access:"readwrite"`
	LockWhenActive    *property.GSettingsBoolProperty `access:"readwrite"`

	LidIsPresent bool

	LinePowerPlan         int32 `access:"readwrite"`
	LinePowerSuspendDelay int32 `access:"readwrite"`
	LinePowerIdleDelay    int32 `access:"readwrite"`

	BatteryPlan         int32 `access:"readwrite"`
	BatterySuspendDelay int32 `access:"readwrite"`
	BatteryIdleDelay    int32 `access:"readwrite"`

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

	p.setBatteryPlan(int32(p.coreSettings.GetEnum("battery-plan")))
	p.setLinePowerPlan(int32(p.coreSettings.GetEnum("ac-plan")))
	fmt.Println("LidClosedAction:", p.LidClosedAction.Get())

	p.initUpower()
	p.initEventHandle()

	var err error
	if p.notifier, err = notifications.NewNotifier("org.freedesktop.Notifications", "/org/freedesktop/Notifications"); err != nil {
		LOGGER.Warning("Can't build org.freedesktop.Notficaations:", err)
	}

	return p
}

func (p *Power) sendNotify(icon, summary, body string) {
	//TODO: close previous notification
	if p.notifier != nil {
		p.notifier.Notify("com.deepin.daemon.power", 0, icon, summary, body, nil, nil, 0)
	} else {
		LOGGER.Warning("failed to show notify message:", summary, body)
	}
}

func main() {
	if !dlib.UniqueOnSession("com.deepin.daemon.Power") {
		LOGGER.Warning("There already has an Power daemon running.")
		return
	}
	p := NewPower()

	dbus.InstallOnSession(p)
	fmt.Println("GetBattery:", getBattery().Vendor.Get())
	dbus.Wait()
}
