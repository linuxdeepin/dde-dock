package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

type Device struct {
	Path  dbus.ObjectPath
	State uint32
}

func NewDevice(core *nm.Device) *Device {
	return &Device{core.Path, core.State.Get()}
}

// DisconnectDevice will disconnect all connection in target device.
func (m *Manager) DisconnectDevice(devPath dbus.ObjectPath) (err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	err = dev.Disconnect()
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

// TODO remove
// func (m *Manager) DisconnectDevice(path dbus.ObjectPath) error {
// 	if dev, err := nm.NewDevice(NMDest, path); err != nil {
// 		return err
// 	} else {
// 		dev.Disconnect()
// 		nm.DestroyDevice(dev)
// 		switch dev.DeviceType.Get() {
// 		case NM_DEVICE_TYPE_WIFI:
// 			dbus.NotifyChange(m, "WirelessConnections")
// 		case NM_DEVICE_TYPE_ETHERNET:
// 			Logger.Debug("DisconnectDevice...", path)
// 			dbus.NotifyChange(m, "WiredConnections")
// 		}
// 		return nil
// 	}
// }

func (m *Manager) initDeviceManage() {
	NMManager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		m.handleDeviceChanged(OpAdded, path)
	})
	NMManager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		m.handleDeviceChanged(OpRemoved, path)
	})
	devs, err := NMManager.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, p := range devs {
		m.handleDeviceChanged(OpAdded, p)
	}
}

func (m *Manager) addWirelessDevice(dev *nm.Device) {
	wirelessDevice := NewDevice(dev)
	if isDeviceExists(m.WirelessDevices, wirelessDevice) {
		// device maybe repeat added
		return
	}
	Logger.Debug("addWirelessDevices:", wirelessDevice)

	// connect signal DeviceStateChanged()
	dev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		wirelessDevice.State = newState
		if m.DeviceStateChanged != nil {
			m.DeviceStateChanged(string(dev.Path), newState)
		}
		// TODO remove
		dbus.NotifyChange(m, "WirelessDevices")
	})

	// connect signal AccessPointAdded() and AccessPointRemoved()
	if devWireless, err := nm.NewDeviceWireless(NMDest, dev.Path); err == nil {
		devWireless.ConnectAccessPointAdded(func(apPath dbus.ObjectPath) {
			if m.AccessPointAdded != nil {
				if ap, err := NewAccessPoint(apPath); err == nil {
					// Logger.Debug("AccessPointAdded:", ap.Ssid, apPath) // TODO test
					m.AccessPointAdded(string(dev.Path), string(ap.Path))
				}
			}
		})
		devWireless.ConnectAccessPointRemoved(func(apPath dbus.ObjectPath) {
			if m.AccessPointRemoved != nil {
				// Logger.Debug("AccessPointRemoved:", apPath) // TODO test
				m.AccessPointRemoved(string(dev.Path), string(apPath))
			}
		})
	}

	m.WirelessDevices = append(m.WirelessDevices, wirelessDevice)
	dbus.NotifyChange(m, "WirelessDevices")
}
func (m *Manager) addWiredDevice(dev *nm.Device) {
	wiredDevice := NewDevice(dev)
	if isDeviceExists(m.WiredDevices, wiredDevice) {
		// device maybe repeat added
		return
	}

	// connect signal DeviceStateChanged()
	dev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		wiredDevice.State = newState
		if m.DeviceStateChanged != nil {
			m.DeviceStateChanged(string(dev.Path), newState)
		}
		// TODO remove
		dbus.NotifyChange(m, "WirelessDevices")
	})
	m.WiredDevices = append(m.WiredDevices, wiredDevice)
	dbus.NotifyChange(m, "WiredDevices")
}
func (m *Manager) addOtherDevice(dev *nm.Device) {
	m.OtherDevices = append(m.OtherDevices, NewDevice(dev))

	otherDevice := NewDevice(dev)
	if isDeviceExists(m.OtherDevices, otherDevice) {
		// device maybe repeat added
		return
	}

	// connect signal DeviceStateChanged()
	dev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		otherDevice.State = newState
		if m.DeviceStateChanged != nil {
			m.DeviceStateChanged(string(dev.Path), newState)
		}
		// TODO remove
		dbus.NotifyChange(m, "WirelessDevices")
	})
	m.OtherDevices = append(m.OtherDevices, otherDevice)
	dbus.NotifyChange(m, "OtherDevices")
}
func isDeviceExists(devs []*Device, dev *Device) bool {
	for _, d := range devs {
		if d.Path == dev.Path {
			return true
		}
	}
	return false
}

func (m *Manager) handleDeviceChanged(operation int32, path dbus.ObjectPath) {
	Logger.Debugf("handleDeviceChanged: operation %d, path %s", operation, path)
	switch operation {
	case OpAdded:
		dev, err := nm.NewDevice(NMDest, path)
		if err != nil {
			panic(err)
		}
		switch dev.DeviceType.Get() {
		case NM_DEVICE_TYPE_WIFI:
			m.addWirelessDevice(dev)
		case NM_DEVICE_TYPE_ETHERNET:
			m.addWiredDevice(dev)
		default:
			m.addOtherDevice(dev)
		}
	case OpRemoved:
		var removed bool
		if m.WirelessDevices, removed = tryRemoveDevice(path, m.WirelessDevices); removed {
			dbus.NotifyChange(m, "WirelessDevices")
			Logger.Debug("WirelessRemoved..")
		} else if m.WiredDevices, removed = tryRemoveDevice(path, m.WiredDevices); removed {
			dbus.NotifyChange(m, "WiredDevices")
		}
	default:
		panic("Didn't support operation")
	}
}

// GetAccessPoints return all access point's dbus path of target device.
func (m *Manager) GetAccessPoints(path dbus.ObjectPath) (aps []dbus.ObjectPath, err error) {
	dev, err := nm.NewDeviceWireless(NMDest, path)
	if err != nil {
		return
	}
	aps, err = dev.GetAccessPoints()
	return
}

// GetAccessPointProperty return access point's detail information.
func (m *Manager) GetAccessPointProperty(apPath dbus.ObjectPath) (ap AccessPoint, err error) {
	ap, err = NewAccessPoint(apPath)
	return
}

func (m *Manager) getDeviceAddress(devPath dbus.ObjectPath, devType uint32) string {
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		dev, err := nm.NewDeviceWired(NMDest, devPath)
		if err != nil {
			panic(err)
		}
		defer func() { nm.DestroyDeviceWired(dev) }()
		return dev.HwAddress.Get()
	case NM_DEVICE_TYPE_WIFI:
		dev, err := nm.NewDeviceWireless(NMDest, devPath)
		if err != nil {
			panic(err)
		}
		defer func() { nm.DestroyDeviceWireless(dev) }()
		return dev.HwAddress.Get()
	}
	return ""
}
