package power

import "pkg.linuxdeepin.com/lib/gio-2.0"
import "time"
import . "pkg.linuxdeepin.com/lib/gettext"
import libupower "dbus/org/freedesktop/upower"

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

func (p *Power) refreshUpower() {

	if upower != nil && upower.OnBattery.Get() != p.OnBattery {
		p.setPropOnBattery(upower.OnBattery.Get())

		if p.OnBattery {
			playSound("power-unplug")
		} else {
			playSound("power-plug")
		}
		//OnBattery will effect current PowerPlan idle value
		p.updateIdletimer()
	}

	p.setPropLidIsPresent(upower.LidIsPresent.Get())

	present, state, percentage, err := getBatteryInfo()
	if err != nil {
		p.setPropBatteryIsPresent(present)
		p.setPropBatteryState(state)
		p.setPropBatteryPercentage(percentage)
		p.handleBatteryPercentage()
	} else {
		p.setPropBatteryIsPresent(false)
	}
	//TODO: handle lowe battery
}

func (p *Power) handleBatteryPercentage() {
	if !p.OnBattery {
		if p.lowBatteryStatus == lowBatteryStatusAction {
			p.lowBatteryStatus = lowBatteryStatusNormal
			doCloseLowpower()
			if p.LockWhenActive.Get() {
				doLock()
			}
		}
		return
	}
	switch {
	case p.BatteryPercentage < abnormalBatteryPercentage:
		//Battery state abnormal
		if p.OnBattery && p.lowBatteryStatus != lowBatteryStatusAbnormal {
			p.lowBatteryStatus = lowBatteryStatusAbnormal
			sendNotify("battery-0", Tr("Abnormal battery power"), Tr("Battery power can not be predicted, please save important documents properly and  not do important operations."))
		}
	case p.BatteryPercentage < float64(p.coreSettings.GetInt("percentage-action")):
		if p.lowBatteryStatus != lowBatteryStatusAction {
			p.lowBatteryStatus = lowBatteryStatusAction
			sendNotify("battery-0", Tr("Battery Critical Low"), Tr("Computer has been in suspend mode, please plug in."))
			doSuspend()
			go func() {
				for p.lowBatteryStatus == lowBatteryStatusAction {
					<-time.After(time.Second * 30)
					//TODO: suspend when there hasn't user input event
					if p.lowBatteryStatus == lowBatteryStatusAction {
						doSuspend()
					}
				}
			}()
		}
	case p.BatteryPercentage < float64(p.coreSettings.GetInt("percentage-critical")):
		if p.lowBatteryStatus != lowBatteryStatusCritcal {
			p.lowBatteryStatus = lowBatteryStatusCritcal
			sendNotify("battery-10", Tr("Battery Critical Low"), Tr("Please plug in, or computer will be in suspend mode."))

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
			sendNotify("battery-25", Tr("Battery Low"), Tr("Computer will be in suspend mode, please plug in now."))
			playSound("power-low")
		}
	default:
		if p.lowBatteryStatus == lowBatteryStatusAction {
			p.lowBatteryStatus = lowBatteryStatusNormal
			doCloseLowpower()
			if p.LockWhenActive.Get() {
				doLock()
			}
		}
	}
}

func (p *Power) initUpower() {
	if upower != nil {
		upower.ConnectChanged(func() {
			p.refreshUpower()
		})
		upower.ConnectDeviceChanged(func(path string) {
			p.refreshUpower()
		})
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
}

func getBatteryInfo() (bool, uint32, float64, error) {
	devs, err := upower.EnumerateDevices()
	if err != nil {
		logger.Error("Can't EnumerateDevices", err)
		return false, 0, 0, nil
	}
	for _, path := range devs {
		dev, err := libupower.NewDevice(UPOWER_BUS_NAME, path)
		if err == nil && dev.Type.Get() == DeviceTypeBattery {
			return dev.IsPresent.Get(),
				dev.State.Get(),
				dev.Percentage.Get(),
				nil
		}
	}
	return false, 0, 0, nil
}
