package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

type device struct {
	Path   dbus.ObjectPath
	State  uint32
	HwAddr string
}

func newDevice(nmDev *nm.Device) *device {
	return &device{
		Path:   nmDev.Path,
		State:  nmDev.State.Get(),
		HwAddr: getDeviceAddress(nmDev.Path, nmDev.DeviceType.Get()),
	}
}

func getDeviceAddress(devPath dbus.ObjectPath, devType uint32) (hwAddr string) {
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		dev, err := nmNewDeviceWired(devPath)
		if err != nil {
			return
		}
		// defer func() { nm.DestroyDeviceWired(dev) }() // TODO remove
		hwAddr = dev.HwAddress.Get()
	case NM_DEVICE_TYPE_WIFI:
		dev, err := nmNewDeviceWireless(devPath)
		if err != nil {
			return
		}
		// defer func() { nm.DestroyDeviceWireless(dev) }() // TODO remove
		hwAddr = dev.HwAddress.Get()
	}
	return
}

func (m *Manager) initDeviceManage() {
	m.devices = make(map[string][]*device)
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

func (m *Manager) handleDeviceChanged(operation int32, devPath dbus.ObjectPath) {
	logger.Debugf("handleDeviceChanged: operation %d, devPath %s", operation, devPath)
	switch operation {
	case opAdded:
		nmDev, err := nmNewDevice(devPath)
		if err != nil {
			return
		}
		switch nmDev.DeviceType.Get() {
		case NM_DEVICE_TYPE_ETHERNET:
			m.addWiredDevice(nmDev)
			m.devices[deviceTypeEthernet] = m.addDevice(m.devices[deviceTypeEthernet], nmDev)
		case NM_DEVICE_TYPE_WIFI:
			m.addWirelessDevice(nmDev)
			m.devices[deviceTypeWifi] = m.addDevice(m.devices[deviceTypeWifi], nmDev)
		case NM_DEVICE_TYPE_UNUSED1:
			m.devices[deviceTypeUnused1] = m.addDevice(m.devices[deviceTypeUnused1], nmDev)
		case NM_DEVICE_TYPE_UNUSED2:
			m.devices[deviceTypeUnused2] = m.addDevice(m.devices[deviceTypeUnused2], nmDev)
		case NM_DEVICE_TYPE_BT:
			m.devices[deviceTypeBt] = m.addDevice(m.devices[deviceTypeBt], nmDev)
		case NM_DEVICE_TYPE_OLPC_MESH:
			m.devices[deviceTypeOlpcMesh] = m.addDevice(m.devices[deviceTypeOlpcMesh], nmDev)
		case NM_DEVICE_TYPE_WIMAX:
			m.devices[deviceTypeWimax] = m.addDevice(m.devices[deviceTypeWimax], nmDev)
		case NM_DEVICE_TYPE_MODEM:
			m.devices[deviceTypeModem] = m.addDevice(m.devices[deviceTypeModem], nmDev)
		case NM_DEVICE_TYPE_INFINIBAND:
			m.devices[deviceTypeInfiniband] = m.addDevice(m.devices[deviceTypeInfiniband], nmDev)
		case NM_DEVICE_TYPE_BOND:
			m.devices[deviceTypeBond] = m.addDevice(m.devices[deviceTypeBond], nmDev)
		case NM_DEVICE_TYPE_VLAN:
			m.devices[deviceTypeVlan] = m.addDevice(m.devices[deviceTypeVlan], nmDev)
		case NM_DEVICE_TYPE_ADSL:
			m.devices[deviceTypeAdsl] = m.addDevice(m.devices[deviceTypeAdsl], nmDev)
		case NM_DEVICE_TYPE_BRIDGE:
			m.devices[deviceTypeBridge] = m.addDevice(m.devices[deviceTypeBridge], nmDev)
		default:
			logger.Error("unknown device type", nmDev.DeviceType.Get())
		}
		m.updatePropDevices()
	case opRemoved:
		if m.isDeviceExists(m.WiredDevices, devPath) {
			m.WiredDevices = m.removeDevice(m.WiredDevices, devPath)
		} else if m.isDeviceExists(m.WirelessDevices, devPath) {
			m.WirelessDevices = m.removeDevice(m.WirelessDevices, devPath)
			logger.Debug("WirelessRemoved..")
		}
		for devType, devs := range m.devices {
			if m.isDeviceExists(devs, devPath) {
				m.devices[devType] = m.removeDevice(devs, devPath)
				break
			}
		}
		m.updatePropDevices()
	default:
		logger.Error("didn't support operation")
	}
}
func (m *Manager) addDevice(devs []*device, nmDev *nm.Device) []*device {
	dev := newDevice(nmDev)
	if m.isDeviceExists(devs, nmDev.Path) {
		// device maybe repeat added
		return devs
	}

	// connect signal DeviceStateChanged()
	nmDev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		dev.State = newState
		if m.DeviceStateChanged != nil {
			m.DeviceStateChanged(string(nmDev.Path), newState)
			m.updatePropDevices()
		}
	})
	devs = append(devs, dev)
	return devs
}
func (m *Manager) addWiredDevice(nmDev *nm.Device) {
	wiredDevice := newDevice(nmDev)
	if m.isDeviceExists(m.WiredDevices, nmDev.Path) {
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
	m.updatePropWiredDevices()
}
func (m *Manager) addWirelessDevice(nmDev *nm.Device) {
	wirelessDevice := newDevice(nmDev)
	if m.isDeviceExists(m.WirelessDevices, nmDev.Path) {
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
	if devWireless, err := nmNewDeviceWireless(nmDev.Path); err == nil {
		devWireless.ConnectAccessPointAdded(func(apPath dbus.ObjectPath) {
			if m.AccessPointAdded != nil {
				if ap, err := NewAccessPoint(apPath); err == nil {
					if len(ap.Ssid) == 0 {
						// ignore hidden access point
						return
					}
					apJSON, _ := marshalJSON(ap)
					m.AccessPointAdded(string(nmDev.Path), apJSON)
				}
			}
		})
		devWireless.ConnectAccessPointRemoved(func(apPath dbus.ObjectPath) {
			if m.AccessPointRemoved != nil {
				apJSON, _ := marshalJSON(accessPoint{Path: apPath})
				m.AccessPointRemoved(string(nmDev.Path), apJSON)
			}
		})
	}

	m.WirelessDevices = append(m.WirelessDevices, wirelessDevice)
	m.updatePropWirelessDevices()
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

func (m *Manager) removeDevice(devs []*device, path dbus.ObjectPath) []*device {
	i := m.getDeviceIndex(devs, path)
	if i < 0 {
		return devs
	}
	copy(devs[i:], devs[i+1:])
	devs[len(devs)-1] = nil
	devs = devs[:len(devs)-1]
	return devs
}

// TODO
func (m *Manager) enableDevice(devPath dbus.ObjectPath) (err error) {
	return
}
func (m *Manager) disableDevice(devPath dbus.ObjectPath) (err error) {
	return
}
