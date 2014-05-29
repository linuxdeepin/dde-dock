package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

type device struct {
	nmDevType     uint32
	nmDev         *nm.Device
	nmDevWireless *nm.DeviceWireless

	Path      dbus.ObjectPath
	State     uint32
	HwAddress string
	ActiveAp  dbus.ObjectPath // used for wireless device
}

func (m *Manager) newDevice(devPath dbus.ObjectPath) (dev *device) {
	nmDev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	dev = &device{
		Path:      nmDev.Path,
		State:     nmDev.State.Get(),
		nmDev:     nmDev,
		nmDevType: nmDev.DeviceType.Get(),
	}

	// dispatch for different device types
	switch dev.nmDevType {
	case NM_DEVICE_TYPE_ETHERNET:
		if nmDevWired, err := nmNewDeviceWired(dev.Path); err == nil {
			dev.HwAddress = nmDevWired.HwAddress.Get()
		}
	case NM_DEVICE_TYPE_WIFI:
		if nmDevWireless, err := nmNewDeviceWireless(dev.Path); err == nil {
			dev.HwAddress = nmDevWireless.HwAddress.Get()
			dev.nmDevWireless = nmDevWireless

			// connect property, about wireless active access point
			dev.ActiveAp = nmDevWireless.ActiveAccessPoint.Get()
			dev.nmDevWireless.ActiveAccessPoint.ConnectChanged(func() {
				dev.ActiveAp = nmDevWireless.ActiveAccessPoint.Get()
				m.updatePropDevices()
			})

			// connect signal AccessPointAdded() and AccessPointRemoved()
			for _, apPath := range nmGetAccessPoints(dev.Path) {
				m.addAccessPoint(dev.Path, apPath)
			}
			dev.nmDevWireless.ConnectAccessPointAdded(func(apPath dbus.ObjectPath) {
				m.addAccessPoint(dev.Path, apPath)
			})
			dev.nmDevWireless.ConnectAccessPointRemoved(func(apPath dbus.ObjectPath) {
				m.removeAccessPoint(dev.Path, apPath)
			})
		}
	}

	// connect signal, about device state
	dev.nmDev.ConnectStateChanged(func(newState uint32, oldState uint32, reason uint32) {
		dev.State = newState
		if m.DeviceStateChanged != nil {
			m.DeviceStateChanged(string(dev.Path), newState)
			m.updatePropDevices()
		}
	})

	return
}

func getDeviceAddress(devPath dbus.ObjectPath, devType uint32) (hwAddr string) {
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		dev, err := nmNewDeviceWired(devPath)
		if err != nil {
			return
		}
		defer func() { nmDestroyDeviceWired(dev) }()
		hwAddr = dev.HwAddress.Get()
	case NM_DEVICE_TYPE_WIFI:
		dev, err := nmNewDeviceWireless(devPath)
		if err != nil {
			return
		}
		defer func() { nmDestroyDeviceWireless(dev) }()
		hwAddr = dev.HwAddress.Get()
	}
	return
}

func getActiveAccessPoint(devPath dbus.ObjectPath, devType uint32) (activeAp dbus.ObjectPath) {
	if devType == NM_DEVICE_TYPE_WIFI {
		dev, err := nmNewDeviceWireless(devPath)
		if err != nil {
			return
		}
		activeAp = dev.ActiveAccessPoint.Get()
	}
	return
}

func (m *Manager) initDeviceManage() {
	m.devices = make(map[string][]*device)
	m.accessPoints = make(map[dbus.ObjectPath][]*accessPoint)
	nmManager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		m.handleDeviceChanged(opAdded, path)
	})
	nmManager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		m.handleDeviceChanged(opRemoved, path)
	})
	for _, p := range nmGetDevices() {
		m.handleDeviceChanged(opAdded, p)
	}
}

func (m *Manager) handleDeviceChanged(operation int32, devPath dbus.ObjectPath) {
	logger.Debugf("handleDeviceChanged: operation %d, devPath %s", operation, devPath)
	switch operation {
	case opAdded:
		m.addDevice(devPath)
	case opRemoved:
		m.removeDevice(devPath)
	default:
		logger.Error("didn't support operation")
	}
}

func (m *Manager) addDevice(devPath dbus.ObjectPath) {
	dev := m.newDevice(devPath)
	devName := getDeviceName(dev.nmDevType)
	m.devices[devName] = m.doAddDevice(m.devices[devName], dev)
	m.updatePropDevices()
}
func (m *Manager) doAddDevice(devs []*device, dev *device) []*device {
	if m.isDeviceExists(devs, dev.Path) {
		// device maybe repeat added
		return devs
	}
	devs = append(devs, dev)
	return devs
}

// func (m *Manager) addWiredDevice(nmDev *nm.Device) {
// 	wiredDevice := m.newDeviceOld(nmDev)
// 	if m.isDeviceExistsOld(m.WiredDevices, nmDev.Path) {
// 		// device maybe repeat added
// 		return
// 	}

// 	// connect signal DeviceStateChanged()
// 	nmDev.ConnectStateChanged(func(newState uint32, oldState uint32, reason uint32) {
// 		wiredDevice.State = newState
// 		if m.DeviceStateChanged != nil {
// 			m.DeviceStateChanged(string(nmDev.Path), newState)
// 		}

// 		// logger.Debug("dWirelessDevices:", wirelessDevice)

// 		// // connect signal diceStateChanged()
// 		// nmDev.ConnectStateChanged(func(newState uint32, oldState uint32, reason uint32) {
// 		// 	wirelessDevice.State = newState
// 		// 	dm.DeviceStateChanged != nil {
// 		// 		m.DeviceStateChanged(string(nmDev.Path), newState)
// 		// 	}
// 		// })

// 		// // connect signal AccessPointAdded() and AccessPointRemoved()
// 		// for _, apPath := range nmGetAccessPoints(nmDev.Path) {
// 	})
// 	m.WiredDevices = append(m.WiredDevices, wiredDevice)
// 	// m.updatePropWiredDevices()
// }
// func (m *Manager) addWirelessDevice(nmDev *nm.Device) {
// 	wirelessDevice := m.newDeviceOld(nmDev)
// 	if m.isDeviceExistsOld(m.WirelessDevices, nmDev.Path) {
// 		// device maybe repeat added
// 		return
// 	}
// 	logger.Debug("addWirelessDevices:", wirelessDevice)

// 	// connect signal DeviceStateChanged()
// 	nmDev.ConnectStateChanged(func(newState uint32, oldState uint32, reason uint32) {
// 		wirelessDevice.State = newState
// 		if m.DeviceStateChanged != nil {

// 			// 			m.addAccessPoint(nmDev.dh, apPath)
// 			// 		})
// 			// 		devWireless.ConnectAccessPointRemoved(dc(apPath dbus.ObjectPath) {
// 			// 			m.removeAccessPoint(nmDev.Path, apPath)
// 			// 		})
// 			// 	}

// 			// 	direlessDevices = append(m.WirelessDevices, wirelessDevice)
// 			// 	m.updatePropWirelessDevices()
// 			// }
// 			// 	func (m *Manager) removeDevice(path dbus.ObjectPath) {
// 			// 		for devType, devs := range m.devices {
// 			m.DeviceStateChanged(string(nmDev.Path), newState)
// 		}
// 	})

// 	// connect signal AccessPointAdded() and AccessPointRemoved()
// 	for _, apPath := range nmGetAccessPoints(nmDev.Path) {
// 		m.addAccessPoint(nmDev.Path, apPath)
// 	}
// 	if devWireless, err := nmNewDeviceWireless(nmDev.Path); err == nil {
// 		devWireless.ConnectAccessPointAdded(func(apPath dbus.ObjectPath) {
// 			m.addAccessPoint(nmDev.Path, apPath)
// 		})
// 		devWireless.ConnectAccessPointRemoved(func(apPath dbus.ObjectPath) {
// 			m.removeAccessPoint(nmDev.Path, apPath)
// 		})
// 	}

// 	m.WirelessDevices = append(m.WirelessDevices, wirelessDevice)
// 	m.updatePropWirelessDevices()
// }
func (m *Manager) removeDevice(path dbus.ObjectPath) {
	for devType, devs := range m.devices {
		if m.isDeviceExists(devs, path) {
			m.devices[devType] = m.doRemoveDevice(devs, path)
			break
		}
	}
	m.updatePropDevices()
}
func (m *Manager) doRemoveDevice(devs []*device, path dbus.ObjectPath) []*device {
	i := m.getDeviceIndex(devs, path)
	if i < 0 {
		return devs
	}

	// destroy object to reset all property connects
	dev := devs[i]
	if dev.nmDevWireless != nil {
		// nmDevWireless is optional, so check if is nil before
		// destroy it
		nmDestroyDeviceWireless(dev.nmDevWireless)
	}
	nmDestroyDevice(dev.nmDev)

	copy(devs[i:], devs[i+1:])
	devs[len(devs)-1] = nil
	devs = devs[:len(devs)-1]
	return devs
}
func (m *Manager) isDeviceExists(devs []*device, path dbus.ObjectPath) bool {
	if m.getDeviceIndex(devs, path) >= 0 {
		return true
	}
	return false
}
func (m *Manager) getDeviceIndex(devs []*device, path dbus.ObjectPath) int {
	for i, d := range devs {
		if d.Path == path {
			return i
		}
	}
	return -1
}

// TODO
func (m *Manager) enableDevice(devPath dbus.ObjectPath) (err error) {
	return
}
func (m *Manager) disableDevice(devPath dbus.ObjectPath) (err error) {
	return
}
