package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

type device struct {
	Path  dbus.ObjectPath
	State uint32
}

func newDevice(core *nm.Device) *device {
	return &device{core.Path, core.State.Get()}
}

// DisconnectDevice will disconnect all connection in target device.
func (m *Manager) DisconnectDevice(devPath dbus.ObjectPath) (err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	err = dev.Disconnect()
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

// TODO remove
// func (m *Manager) DisconnectDevice(path dbus.ObjectPath) error {
// 	if dev, err := nmNewDevice(path); err != nil {
// 		return err
// 	} else {
// 		dev.Disconnect()
// 		nm.DestroyDevice(dev)
// 		switch dev.DeviceType.Get() {
// 		case NM_DEVICE_TYPE_WIFI:
// 			dbus.NotifyChange(m, "WirelessConnections")
// 		case NM_DEVICE_TYPE_ETHERNET:
// 			logger.Debug("DisconnectDevice...", path)
// 			dbus.NotifyChange(m, "WiredConnections")
// 		}
// 		return nil
// 	}
// }

func (m *Manager) initDeviceManage() {
	nmManager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		m.handleDeviceChanged(opAdded, path)
	})
	nmManager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		m.handleDeviceChanged(opRemoved, path)
	})
	devs, err := nmGetDevices()
	if err != nil {
		panic(err)
	}
	for _, p := range devs {
		m.handleDeviceChanged(opAdded, p)
	}
}

func (m *Manager) addWiredDevice(dev *nm.Device) {
	wiredDevice := newDevice(dev)
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
		m.updatePropWiredDevices()
	})
	m.WiredDevices = append(m.WiredDevices, wiredDevice)
	m.updatePropWiredDevices()
}
func (m *Manager) addWirelessDevice(dev *nm.Device) {
	wirelessDevice := newDevice(dev)
	if isDeviceExists(m.WirelessDevices, wirelessDevice) {
		// device maybe repeat added
		return
	}
	logger.Debug("addWirelessDevices:", wirelessDevice)

	// connect signal DeviceStateChanged()
	dev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		wirelessDevice.State = newState
		if m.DeviceStateChanged != nil {
			m.DeviceStateChanged(string(dev.Path), newState)
		}
		// TODO remove
		m.updatePropWirelessDevices()
	})

	// connect signal AccessPointAdded() and AccessPointRemoved()
	if devWireless, err := nmNewDeviceWireless(dev.Path); err == nil {
		devWireless.ConnectAccessPointAdded(func(apPath dbus.ObjectPath) {
			if m.AccessPointAdded != nil {
				if ap, err := NewAccessPoint(apPath); err == nil {
					// logger.Debug("AccessPointAdded:", ap.Ssid, apPath) // TODO test
					m.AccessPointAdded(string(dev.Path), string(ap.Path))
				}
			}
		})
		devWireless.ConnectAccessPointRemoved(func(apPath dbus.ObjectPath) {
			if m.AccessPointRemoved != nil {
				// logger.Debug("AccessPointRemoved:", apPath) // TODO test
				m.AccessPointRemoved(string(dev.Path), string(apPath))
			}
		})
	}

	m.WirelessDevices = append(m.WirelessDevices, wirelessDevice)
	m.updatePropWirelessDevices()
}
func (m *Manager) addOtherDevice(dev *nm.Device) {
	m.OtherDevices = append(m.OtherDevices, newDevice(dev))

	otherDevice := newDevice(dev)
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
		m.updatePropOtherDevices()
	})
	m.OtherDevices = append(m.OtherDevices, otherDevice)
	m.updatePropOtherDevices()
}
func isDeviceExists(devs []*device, dev *device) bool {
	for _, d := range devs {
		if d.Path == dev.Path {
			return true
		}
	}
	return false
}

func (m *Manager) handleDeviceChanged(operation int32, path dbus.ObjectPath) {
	logger.Debugf("handleDeviceChanged: operation %d, path %s", operation, path)
	switch operation {
	case opAdded:
		dev, err := nmNewDevice(path)
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
	case opRemoved:
		var removed bool
		if m.WirelessDevices, removed = tryRemoveDevice(path, m.WirelessDevices); removed {
			m.updatePropWirelessDevices()
			logger.Debug("WirelessRemoved..")
		} else if m.WiredDevices, removed = tryRemoveDevice(path, m.WiredDevices); removed {
			m.updatePropWiredDevices()
		}
	default:
		panic("Didn't support operation")
	}
}

// GetAccessPoints return all access point's dbus path of target device.
func (m *Manager) GetAccessPoints(path dbus.ObjectPath) (aps []dbus.ObjectPath, err error) {
	dev, err := nmNewDeviceWireless(path)
	if err != nil {
		return
	}
	aps, err = dev.GetAccessPoints()
	return
}

// GetAccessPointProperty return access point's detail information.
func (m *Manager) GetAccessPointProperty(apPath dbus.ObjectPath) (ap accessPoint, err error) {
	ap, err = NewAccessPoint(apPath)
	return
}

// TODO
func (m *Manager) getDeviceAddress(devPath dbus.ObjectPath, devType uint32) string {
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		dev, err := nmNewDeviceWired(devPath)
		if err != nil {
			panic(err)
		}
		defer func() { nm.DestroyDeviceWired(dev) }()
		return dev.HwAddress.Get()
	case NM_DEVICE_TYPE_WIFI:
		dev, err := nmNewDeviceWireless(devPath)
		if err != nil {
			panic(err)
		}
		defer func() { nm.DestroyDeviceWireless(dev) }()
		return dev.HwAddress.Get()
	}
	return ""
}

func tryRemoveDevice(path dbus.ObjectPath, devices []*device) ([]*device, bool) {
	var newDevices []*device
	found := false
	for _, dev := range devices {
		if dev.Path != path {
			newDevices = append(newDevices, dev)
		} else {
			found = true
		}
	}
	return newDevices, found
}
