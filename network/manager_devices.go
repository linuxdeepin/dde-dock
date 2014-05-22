package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

type deviceOld struct {
	Path  dbus.ObjectPath
	State uint32
}
type device struct {
	nmDevType     uint32
	nmDev         *nm.Device
	nmDevWireless *nm.DeviceWireless

	Path      dbus.ObjectPath
	State     uint32
	HwAddress string
	ActiveAp  dbus.ObjectPath // used for wireless device
}

func (m *Manager) newDeviceOld(nmDev *nm.Device) (dev *deviceOld) {
	return &deviceOld{
		Path:  nmDev.Path,
		State: nmDev.State.Get(),
	}
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
		}
	}

	// connect signal, about device state
	dev.nmDev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		dev.State = newState
		if m.DeviceStateChanged != nil {
			m.DeviceStateChanged(string(nmDev.Path), newState)
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
		defer func() { nm.DestroyDeviceWired(dev) }()
		hwAddr = dev.HwAddress.Get()
	case NM_DEVICE_TYPE_WIFI:
		dev, err := nmNewDeviceWireless(devPath)
		if err != nil {
			return
		}
		defer func() { nm.DestroyDeviceWireless(dev) }()
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
		// TODO remove
		nmDev, err := nmNewDevice(devPath)
		if err != nil {
			return
		}
		switch nmDev.DeviceType.Get() {
		case NM_DEVICE_TYPE_ETHERNET:
			m.addWiredDevice(nmDev)
		case NM_DEVICE_TYPE_WIFI:
			m.addWirelessDevice(nmDev)
		}
		m.addDevice(devPath)
	case opRemoved:
		if m.isDeviceExistsOld(m.WiredDevices, devPath) {
			m.WiredDevices = m.doRemoveDeviceOld(m.WiredDevices, devPath)
		} else if m.isDeviceExistsOld(m.WirelessDevices, devPath) {
			m.WirelessDevices = m.doRemoveDeviceOld(m.WirelessDevices, devPath)
			logger.Debug("WirelessRemoved..")
		}
		m.removeDevice(devPath)
	default:
		logger.Error("didn't support operation")
	}
}

func (m *Manager) addDevice(devPath dbus.ObjectPath) {
	dev := m.newDevice(devPath)

	// // connect signal about device state
	// nmDev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
	// 	dev.State = newState
	// 	if m.DeviceStateChanged != nil {
	// 		m.DeviceStateChanged(string(nmDev.Path), newState)
	// 		m.updatePropDevices()
	// 	}
	// })

	devName := getDeviceTypeName(dev.nmDevType)
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

// TODO remove
// func (m *Manager) addDevice(devs []*device, nmDev *nm.Device) []*device {
// 	dev := newDevice(nmDev)
// 	if m.isDeviceExists(devs, nmDev.Path) {
// 		// device maybe repeat added
// 		return devs
// 	}

// 	// TODO
// 	// connect signals, state changed and wireless active access point
// 	nmDev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
// 		dev.State = newState
// 		if m.DeviceStateChanged != nil {
// 			m.DeviceStateChanged(string(nmDev.Path), newState)
// 			m.updatePropDevices()
// 		}
// 	})
// 	devs = append(devs, dev)
// 	return devs
// }
func (m *Manager) addWiredDevice(nmDev *nm.Device) {
	wiredDevice := m.newDeviceOld(nmDev)
	if m.isDeviceExistsOld(m.WiredDevices, nmDev.Path) {
		// device maybe repeat added
		return
	}

	// connect signal DeviceStateChanged()
	nmDev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		wiredDevice.State = newState
		if m.DeviceStateChanged != nil {
			m.DeviceStateChanged(string(nmDev.Path), newState)
		}
	})
	m.WiredDevices = append(m.WiredDevices, wiredDevice)
	// m.updatePropWiredDevices()
}
func (m *Manager) addWirelessDevice(nmDev *nm.Device) {
	wirelessDevice := m.newDeviceOld(nmDev)
	if m.isDeviceExistsOld(m.WirelessDevices, nmDev.Path) {
		// device maybe repeat added
		return
	}
	logger.Debug("addWirelessDevices:", wirelessDevice)

	// connect signal DeviceStateChanged()
	nmDev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		wirelessDevice.State = newState
		if m.DeviceStateChanged != nil {
			m.DeviceStateChanged(string(nmDev.Path), newState)
		}
	})

	// connect signal AccessPointAdded() and AccessPointRemoved()
	for _, apPath := range nmGetAccessPoints(nmDev.Path) {
		m.addAccessPoint(nmDev.Path, apPath)
	}
	if devWireless, err := nmNewDeviceWireless(nmDev.Path); err == nil {
		devWireless.ConnectAccessPointAdded(func(apPath dbus.ObjectPath) {
			m.addAccessPoint(nmDev.Path, apPath)
		})
		devWireless.ConnectAccessPointRemoved(func(apPath dbus.ObjectPath) {
			m.removeAccessPoint(nmDev.Path, apPath)
		})
	}

	m.WirelessDevices = append(m.WirelessDevices, wirelessDevice)
	m.updatePropWirelessDevices()
}
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

func (m *Manager) doRemoveDeviceOld(devs []*deviceOld, path dbus.ObjectPath) []*deviceOld {
	i := m.getDeviceIndexOld(devs, path)
	if i < 0 {
		return devs
	}
	copy(devs[i:], devs[i+1:])
	devs[len(devs)-1] = nil
	devs = devs[:len(devs)-1]
	return devs
}
func (m *Manager) isDeviceExistsOld(devs []*deviceOld, path dbus.ObjectPath) bool {
	if m.getDeviceIndexOld(devs, path) >= 0 {
		return true
	}
	return false
}
func (m *Manager) getDeviceIndexOld(devs []*deviceOld, path dbus.ObjectPath) int {
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
