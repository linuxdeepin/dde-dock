package power

import "gir/gio-2.0"
import "time"
import . "pkg.deepin.io/lib/gettext"
import "pkg.deepin.io/lib/dbus"
import "pkg.deepin.io/dde/api/soundutils"

const (
	UPOWER_BUS_NAME = "org.freedesktop.UPower"
)

const (
	//defined at http://upower.freedesktop.org/docs/Device.html#Device:Type
	DeviceTypeUnknow    = 0
	DeviceTypeLinePower = 1
	DeviceTypeBattery   = 2
	DeviceTypeUps       = 3
	DeviceTypeMonitor   = 4
	DeviceTypeMouse     = 5
	DeviceTypeKeyboard  = 6
	DeviceTypePda       = 7
	DeviceTypePhone     = 8
)

const (
	//defined at http://upower.freedesktop.org/docs/Device.html#Device:State
	BatteryStateUnknown          = 0
	BatteryStateCharging         = 1
	BatteryStateDischarging      = 2
	BatteryStateEmpty            = 3
	BatteryStateFullyCharged     = 4
	BatteryStatePendingCharge    = 5
	BatteryStatePendingDischarge = 6
)

const (
	//internal used
	lowBatteryStatusNormal = iota
	lowBatteryStatusAbnormal
	lowBatteryStatusLow
	lowBatteryStatusCritcal
	lowBatteryStatusAction
)

const abnormalBatteryPercentage = float64(1.0)

var hasSleepInLowPower bool

func (p *Power) refreshUpower() {
	if upower != nil && upower.OnBattery.Get() != p.OnBattery {
		p.setPropOnBattery(upower.OnBattery.Get())

		if p.OnBattery {
			playSound(soundutils.EventPowerUnplug)
		} else {
			playSound(soundutils.EventPowerPlug)
		}
		//OnBattery will effect current PowerPlan idle value
		p.updateIdletimer()
	}

	p.setPropLidIsPresent(upower.LidIsPresent.Get())
	p.updateBatteryInfo()
}

func (p *Power) handleBatteryPercentage() {
	if !p.OnBattery {
		// Close low power saver if power plug
		if p.lowBatteryStatus == lowBatteryStatusAction {
			doCloseLowpower()
			p.lowBatteryStatus = lowBatteryStatusNormal
			if hasSleepInLowPower && p.LockWhenActive.Get() {
				hasSleepInLowPower = false
				doLock()
			}
		}

		// Reset lowBatteryStatus status if power plug
		if p.lowBatteryStatus != lowBatteryStatusNormal {
			p.lowBatteryStatus = lowBatteryStatusNormal
		}
		return
	}
	switch {
	case p.BatteryPercentage < abnormalBatteryPercentage:
		//Battery state abnormal
		if p.OnBattery && p.lowBatteryStatus != lowBatteryStatusAbnormal {
			p.lowBatteryStatus = lowBatteryStatusAbnormal
			sendNotify("battery_empty", Tr("Abnormal battery power"), Tr("Battery power can not be predicted, please save important documents properly and  not do important operations."))
		}
	case p.BatteryPercentage < float64(p.coreSettings.GetInt("percentage-action")):
		if p.lowBatteryStatus != lowBatteryStatusAction {
			p.lowBatteryStatus = lowBatteryStatusAction
			playSound(soundutils.EventBatteryLow)
			sendNotify("battery_empty", Tr("Battery Critical Low"), Tr("Computer has been in suspend mode, please plug in."))
			go func() {
				for p.lowBatteryStatus == lowBatteryStatusAction {
					<-time.After(time.Second * 30)
					if !p.OnBattery {
						break
					}
					hasSleepInLowPower = false
					//TODO: suspend when there hasn't user input event
					if p.lowBatteryStatus == lowBatteryStatusAction {
						doSuspend()
					}
				}
			}()
		}
		doShowLowpower()
	case p.BatteryPercentage < float64(p.coreSettings.GetInt("percentage-critical")):
		if p.lowBatteryStatus != lowBatteryStatusCritcal {
			p.lowBatteryStatus = lowBatteryStatusCritcal
			playSound(soundutils.EventBatteryLow)
			sendNotify("battery_low", Tr("Battery Critical Low"), Tr("Please plug in, or computer will be in suspend mode."))

			playSound("power-caution")
			go func() {
				for p.lowBatteryStatus == lowBatteryStatusCritcal {
					<-time.After(time.Second * 10)
					playSound("power-caution")
				}
			}()
		}
	case p.BatteryPercentage < float64(p.coreSettings.GetInt("percentage-low")):
		if p.lowBatteryStatus != lowBatteryStatusLow {
			p.lowBatteryStatus = lowBatteryStatusLow
			playSound(soundutils.EventBatteryLow)
			sendNotify("battery_caution", Tr("Battery Low"), Tr("Computer will be in suspend mode, please plug in now."))
			playSound("power-low")
		}
	default:
		if p.lowBatteryStatus == lowBatteryStatusAction {
			p.lowBatteryStatus = lowBatteryStatusNormal
			doCloseLowpower()
			if hasSleepInLowPower && p.LockWhenActive.Get() {
				hasSleepInLowPower = false
				doLock()
			}
		}
	}
}

func (p *Power) initUpower() {
	if upower != nil {
		upower.OnBattery.ConnectChanged(func() {
			p.refreshUpower()
		})
		upower.ConnectDeviceAdded(func(path dbus.ObjectPath) {
			if p.batGroup != nil {
				p.batGroup.AddBatteryDevice(path)
			}
			p.refreshUpower()
		})
		upower.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
			if p.batGroup != nil {
				p.batGroup.RemoveBatteryDevice(path)
			}
			p.refreshUpower()
		})
		p.batGroup = NewBatteryGroup(p.updateBatteryInfo)
		p.setPropOnBattery(upower.OnBattery.Get())
		p.refreshUpower()
	}

	p.coreSettings.Connect("changed::percentage-action", func(s *gio.Settings, name string) {
		p.handleBatteryPercentage()
	})
	p.coreSettings.Connect("changed::percentage-low", func(s *gio.Settings, name string) {
		p.handleBatteryPercentage()
	})
	p.coreSettings.Connect("changed::percentage-critical", func(s *gio.Settings, name string) {
		p.handleBatteryPercentage()
	})
	p.coreSettings.GetInt("percentage-action")
}

func (p *Power) updateBatteryInfo() {
	if p.batGroup == nil {
		logger.Debug("No battery device")
		return
	}

	present, state, percentage, err := p.batGroup.GetBatteryInfo()
	if err != nil {
		logger.Warning(err)
		return
	}

	if present && (state == BatteryStateDischarging) {
		p.setPropOnBattery(true)
	} else {
		p.setPropOnBattery(false)
	}
	p.setPropBatteryIsPresent(present)
	p.setPropBatteryState(state)
	p.setPropBatteryPercentage(percentage)
	p.handleBatteryPercentage()
	//TODO: handle lower battery
}
