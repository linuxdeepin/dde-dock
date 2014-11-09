/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package network

import nm "dbus/org/freedesktop/networkmanager"
import "pkg.linuxdeepin.com/lib/dbus"

type device struct {
	nmDev         *nm.Device
	nmDevWireless *nm.DeviceWireless
	nmDevType     uint32
	id            string

	Path      dbus.ObjectPath
	State     uint32
	HwAddress string
	Managed   bool
	Vendor    string
	UsbDevice bool
	ActiveAp  dbus.ObjectPath // used for wireless device
}

func (m *Manager) initDeviceManage() {
	m.devices = make(map[string][]*device)
	m.accessPoints = make(map[dbus.ObjectPath][]*accessPoint)
	nmManager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		m.addDevice(path)
	})
	nmManager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		m.removeDevice(path)
	})
	for _, path := range nmGetDevices() {
		m.addDevice(path)
	}
}

func (m *Manager) newDevice(devPath dbus.ObjectPath) (dev *device) {
	m.devicesLocker.Lock()
	defer m.devicesLocker.Unlock()

	nmDev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	dev = &device{
		nmDev:     nmDev,
		nmDevType: nmDev.DeviceType.Get(),
		Path:      nmDev.Path,
		State:     nmDev.State.Get(),
	}
	dev.Managed = nmGeneralIsDeviceManaged(devPath)
	dev.Vendor = udevGetDeviceVendor(nmDev.Udi.Get())
	dev.UsbDevice = udevIsUsbDevice(nmDev.Udi.Get())
	dev.id, _ = nmGeneralGetDeviceIdentifier(devPath)

	// add device config
	m.config.addDeviceConfig(devPath)

	// connect signal
	dev.nmDev.ConnectStateChanged(func(newState, oldState, reason uint32) {
		m.devicesLocker.Lock()
		defer m.devicesLocker.Unlock()

		if reason == NM_DEVICE_STATE_REASON_REMOVED {
			return
		}
		dev.State = newState
		dev.Managed = nmGeneralIsDeviceManaged(dev.Path)
		m.setPropDevices()
		dbus.Emit(m, "DeviceStateChanged", string(dev.Path), dev.State)

		m.config.updateDeviceConfig(dev.Path)
		m.config.syncDeviceState(dev.Path)
	})
	dev.State = dev.nmDev.State.Get()

	// dispatch for different device types
	switch dev.nmDevType {
	case NM_DEVICE_TYPE_ETHERNET:
		if nmDevWired, err := nmNewDeviceWired(dev.Path); err == nil {
			dev.HwAddress = nmDevWired.HwAddress.Get()
		}
	case NM_DEVICE_TYPE_WIFI:
		logger.Debug("add wireless device:", dev.Path)
		if nmDevWireless, err := nmNewDeviceWireless(dev.Path); err == nil {
			dev.HwAddress = nmDevWireless.HwAddress.Get()
			dev.nmDevWireless = nmDevWireless

			// connect property, about wireless active access point
			dev.nmDevWireless.ActiveAccessPoint.ConnectChanged(func() {
				m.devicesLocker.Lock()
				defer m.devicesLocker.Unlock()

				dev.ActiveAp = nmDevWireless.ActiveAccessPoint.Get()
				m.setPropDevices()
			})
			dev.ActiveAp = nmDevWireless.ActiveAccessPoint.Get()

			// connect signal AccessPointAdded() and AccessPointRemoved()
			dev.nmDevWireless.ConnectAccessPointAdded(func(apPath dbus.ObjectPath) {
				m.addAccessPoint(dev.Path, apPath)
			})
			dev.nmDevWireless.ConnectAccessPointRemoved(func(apPath dbus.ObjectPath) {
				m.removeAccessPoint(dev.Path, apPath)
			})
			for _, apPath := range nmGetAccessPoints(dev.Path) {
				m.addAccessPoint(dev.Path, apPath)
			}
		}
	}

	return
}
func (m *Manager) destroyDevice(dev *device) {
	// remove device config
	m.config.removeDeviceConfig(dev.id)

	// destroy object to reset all property connects
	if dev.nmDevWireless != nil {
		// nmDevWireless is optional, so check if is nil before
		// destroy it
		nmDestroyDeviceWireless(dev.nmDevWireless)
	}
	nmDestroyDevice(dev.nmDev)
}

func (m *Manager) addDevice(devPath dbus.ObjectPath) {
	dev := m.newDevice(devPath)

	m.devicesLocker.Lock()
	defer m.devicesLocker.Unlock()

	devType := getCustomDeviceType(dev.nmDevType)
	m.devices[devType] = m.doAddDevice(m.devices[devType], dev)
	m.setPropDevices()
}
func (m *Manager) doAddDevice(devs []*device, dev *device) []*device {
	if m.isDeviceExists(devs, dev.Path) {
		// maybe device repeat added
		return devs
	}
	devs = append(devs, dev)
	return devs
}

func (m *Manager) removeDevice(path dbus.ObjectPath) {
	m.devicesLocker.Lock()
	defer m.devicesLocker.Unlock()

	for devType, devs := range m.devices {
		if m.isDeviceExists(devs, path) {
			m.devices[devType] = m.doRemoveDevice(devs, path)
			break
		}
	}
	m.setPropDevices()
}
func (m *Manager) doRemoveDevice(devs []*device, path dbus.ObjectPath) []*device {
	i := m.getDeviceIndex(devs, path)
	if i < 0 {
		return devs
	}

	m.destroyDevice(devs[i])
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
