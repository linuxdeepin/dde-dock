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
	ActiveAp  dbus.ObjectPath // used for wireless device
}

func (m *Manager) newDevice(devPath dbus.ObjectPath) (dev *device) {
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
	dev.id, _ = nmGeneralGetDeviceIdentifier(devPath)

	// add device config
	m.config.addDeviceConfig(devPath)

	// connect signal, about device state
	dev.nmDev.ConnectStateChanged(func(newState, oldState, reason uint32) {
		dev.State = newState
		m.config.updateDeviceConfig(dev.Path)
		m.config.syncDeviceState(dev.Path)
		if m.DeviceStateChanged != nil { // TODO
			m.DeviceStateChanged(string(dev.Path), dev.State)
			m.setPropDevices()
		}
	})

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
				m.setPropDevices()
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
	devType := getCustomDeviceType(dev.nmDevType)
	m.devices[devType] = m.doAddDevice(m.devices[devType], dev)
	m.setPropDevices()
}
func (m *Manager) doAddDevice(devs []*device, dev *device) []*device {
	if m.isDeviceExists(devs, dev.Path) {
		// device maybe repeat added
		return devs
	}
	devs = append(devs, dev)
	return devs
}

func (m *Manager) removeDevice(path dbus.ObjectPath) {
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
