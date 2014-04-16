package main

import "dbus/org/freedesktop/upower"
import "dlib/gio-2.0"
import "time"
import "dlib"

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
	lowBatteryStatusLow
	lowBatteryStatusCritcal
	lowBatteryStatusAction
)

func (p *Power) refreshUpower(up *upower.Upower) {

	if up.OnBattery.Get() != p.OnBattery {
		p.setPropOnBattery(up.OnBattery.Get())

		if p.OnBattery {
			p.player.PlaySystemSound("power-unplug")
		} else {
			p.player.PlaySystemSound("power-plug")
		}
		//OnBattery will effect current PowerPlan idle value
		p.updateIdletimer()
	}

	p.setPropLidIsPresent(up.LidIsPresent.Get())

	if dev := getBattery(); dev != nil {
		p.setPropBatteryIsPresent(dev.IsPresent.Get())
		p.setPropBatteryState(dev.State.Get())
		p.setPropBatteryPercentage(dev.Percentage.Get())
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
	case p.BatteryPercentage < float64(p.coreSettings.GetInt("percentage-action")):
		if p.lowBatteryStatus != lowBatteryStatusAction {
			p.lowBatteryStatus = lowBatteryStatusAction
			p.sendNotify("battery-action", dlib.Tr("Battery cirtical low"), dlib.Tr("Computer will suspend very soon unless it is plugged in."))
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
			p.sendNotify("battery-critical", dlib.Tr("Battery cirtical low"), dlib.Tr("Computer will suspend very soon unless it is plugged in."))

			p.player.PlaySystemSound("power-caution")
			go func() {
				for p.lowBatteryStatus == lowBatteryStatusCritcal {
					<-time.After(time.Second * 10)
					p.player.PlaySystemSound("power-caution")
				}
			}()
		}
	case p.BatteryPercentage < float64(p.coreSettings.GetInt("percentage-low")):
		if p.lowBatteryStatus != lowBatteryStatusLow {
			p.lowBatteryStatus = lowBatteryStatusLow
			p.sendNotify("battery-low", dlib.Tr("Battery low"), dlib.Tr("Computer will suspend very soon unless it is plugged in(TODO:Calucate remaining)."))
			p.player.PlaySystemSound("power-low")
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
	up, err := upower.NewUpower(UPOWER_BUS_NAME, "/org/freedesktop/UPower")
	if err != nil {
		LOGGER.Error("Can't build org.freedesktop.UPower:", err)
	} else {
		up.ConnectChanged(func() {
			p.refreshUpower(up)
		})
		up.ConnectDeviceChanged(func(path string) {
			p.refreshUpower(up)
		})
	}
	p.setPropOnBattery(up.OnBattery.Get())
	p.refreshUpower(up)

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
func getBattery() *upower.Device {
	up, err := upower.NewUpower(UPOWER_BUS_NAME, "/org/freedesktop/UPower")
	if err != nil {
		LOGGER.Error("Can't build org.freedesktop.UPower:", err)
	}
	devs, err := up.EnumerateDevices()
	if err != nil {
		LOGGER.Error("Can't EnumerateDevices", err)
	}
	for _, path := range devs {
		dev, err := upower.NewDevice(UPOWER_BUS_NAME, path)
		if err == nil && dev.Type.Get() == DeviceTypeBattery {
			return dev
		}
	}
	return nil
}
