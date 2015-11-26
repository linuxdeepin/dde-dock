package power

import "pkg.deepin.io/lib/log"
import "pkg.deepin.io/lib/dbus/property"
import "pkg.deepin.io/lib/gio-2.0"
import ss "dbus/org/freedesktop/screensaver"
import "pkg.deepin.io/dde/api/soundutils"

const (
	settingKeyPowerButton = "button-power"
	settingKeyLidClose    = "lid-close"
	settingKeyLinePlan    = "ac-plan"
	settingKeyBatteryPlan = "battery-plan"
	settingKeyLockEnabled = "lock-enabled"
)

var logger = log.NewLogger("daemon/power")

type Power struct {
	coreSettings     *gio.Settings
	screensaver      *ss.ScreenSaver
	batGroup         *batteryGroup
	lidIsClosed      bool
	lowBatteryStatus uint32
	swQuit           chan struct{}

	// 按下电源键执行的操作
	PowerButtonAction *property.GSettingsEnumProperty `access:"readwrite"`
	// 合上笔记本盖时执行的操作
	LidClosedAction *property.GSettingsEnumProperty `access:"readwrite"`
	// 屏幕唤醒时是否锁屏
	LockWhenActive *property.GSettingsBoolProperty `access:"readwrite"`

	// 是否有显示器
	LidIsPresent bool

	// 接通电源时的电源计划
	LinePowerPlan *property.GSettingsEnumProperty `access:"readwrite"`
	// 接通电源时的挂起超时
	LinePowerSuspendDelay int32 `access:"readwrite"`
	// 接通电源时的空闲检测超时
	LinePowerIdleDelay int32 `access:"readwrite"`

	// 使用电池时的电源计划
	BatteryPlan *property.GSettingsEnumProperty `access:"readwrite"`
	// 使用电池时的挂起超时
	BatterySuspendDelay int32 `access:"readwrite"`
	// 使用电池时的空闲检测超时
	BatteryIdleDelay int32 `access:"readwrite"`

	// 剩余电量
	BatteryPercentage float64

	//Not in Charging, Charging, Full
	BatteryState uint32

	// 是否有电池设备
	BatteryIsPresent bool

	// 是否使用电池
	OnBattery bool

	// 电源计划列表
	PlanInfo string
}

func (p *Power) Reset() {
	_, str := p.coreSettings.GetDefaultValue(
		settingKeyPowerButton).GetString()
	action := queryPowerAction(str)
	if action != p.PowerButtonAction.Get() {
		p.PowerButtonAction.Set(action)
	}

	_, str = p.coreSettings.GetDefaultValue(
		settingKeyLidClose).GetString()
	action = queryPowerAction(str)
	if action != p.LidClosedAction.Get() {
		p.LidClosedAction.Set(action)
	}

	enabled := p.coreSettings.GetDefaultValue(
		settingKeyLockEnabled).GetBoolean()
	if enabled != p.LockWhenActive.Get() {
		p.LockWhenActive.Set(enabled)
	}

	_, str = p.coreSettings.GetDefaultValue(
		settingKeyLinePlan).GetString()
	plan := queryPowerPlan(str)
	if plan != p.LinePowerPlan.Get() {
		p.LinePowerPlan.Set(plan)
	}

	_, str = p.coreSettings.GetDefaultValue(
		settingKeyBatteryPlan).GetString()
	plan = queryPowerPlan(str)
	if plan != p.BatteryPlan.Get() {
		p.BatteryPlan.Set(plan)
	}
}

func NewPower() *Power {
	p := &Power{}
	p.coreSettings = gio.NewSettings("com.deepin.daemon.power")
	p.PowerButtonAction = property.NewGSettingsEnumProperty(p, "PowerButtonAction", p.coreSettings, settingKeyPowerButton)
	p.LidClosedAction = property.NewGSettingsEnumProperty(p, "LidClosedAction", p.coreSettings, settingKeyLidClose)
	p.LockWhenActive = property.NewGSettingsBoolProperty(p, "LockWhenActive", p.coreSettings, settingKeyLockEnabled)

	var err error
	if p.screensaver, err = ss.NewScreenSaver("org.freedesktop.ScreenSaver", "/org/freedesktop/ScreenSaver"); err != nil {
		logger.Warning("Can't build org.freedesktop.ScreenSaver:", err)
	}

	p.initPlan()
	p.initUpower()
	p.initEventHandle()

	p.LinePowerPlan = property.NewGSettingsEnumProperty(p, "LinePowerPlan", p.coreSettings, settingKeyLinePlan)
	p.LinePowerPlan.ConnectChanged(func() {
		p.setLinePowerPlan(p.LinePowerPlan.Get())
	})
	p.setLinePowerPlan(p.LinePowerPlan.Get())

	p.BatteryPlan = property.NewGSettingsEnumProperty(p, "BatteryPlan", p.coreSettings, settingKeyBatteryPlan)
	p.BatteryPlan.ConnectChanged(func() {
		p.setBatteryPlan(p.BatteryPlan.Get())
	})
	p.setBatteryPlan(p.BatteryPlan.Get())

	p.swQuit = make(chan struct{})

	return p
}

func sendNotify(icon, summary, body string) {
	//TODO: close previous notification
	if notifier != nil {
		notifier.Notify(dbusDest, 0, icon, summary, body, nil, nil, 0)
	} else {
		logger.Warning("failed to show notify message:", summary, body)
	}
}

func playSound(name string) {
	soundutils.PlaySystemSound(name, "", false)
}
