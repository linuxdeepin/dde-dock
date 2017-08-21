/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	mmdbus "dbus/org/freedesktop/modemmanager1"
	nmdbus "dbus/org/freedesktop/networkmanager"
	"fmt"
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus"
	"strings"
	"time"
)

type device struct {
	nmDev         *nmdbus.Device
	nmDevWired    *nmdbus.DeviceWired
	nmDevWireless *nmdbus.DeviceWireless
	mmDevModem    *mmdbus.Modem
	nmDevType     uint32
	id            string
	udi           string

	Path          dbus.ObjectPath
	State         uint32
	Interface     string
	HwAddress     string
	PermHwAddress string
	Driver        string
	Managed       bool

	// Vendor is the device vendor ID and product ID, if failed, use
	// interface name instead. BTW, we use Vendor instead of
	// Description as the name to keep compatible with the old
	// front-end code.
	Vendor string

	// Unique connection uuid for this device, works for wired,
	// wireless and modem devices, for wireless device the unique uuid
	// will be the connection uuid of hotspot mode.
	UniqueUuid string

	UsbDevice bool // not works for mobile device(modem)

	// used for wireless device
	ActiveAp       dbus.ObjectPath
	SupportHotspot bool

	// used for mobile device
	MobileNetworkType   string
	MobileSignalQuality uint32
}

func (m *Manager) initDeviceManage() {
	m.devicesLock.Lock()
	m.devices = make(map[string][]*device)
	m.devicesLock.Unlock()

	m.accessPointsLock.Lock()
	m.accessPoints = make(map[dbus.ObjectPath][]*accessPoint)
	m.accessPointsLock.Unlock()

	nmManager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		m.addDevice(path)
	})
	nmManager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		notifyDeviceRemoved(path)
		m.removeDevice(path)
	})
	for _, path := range nmGetDevices() {
		m.addDevice(path)
	}
}

func (m *Manager) newDevice(devPath dbus.ObjectPath) (dev *device, err error) {
	nmDev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}

	// ignore virtual network interfaces
	if isVirtualDeviceIfc(nmDev) {
		err = fmt.Errorf("ignore virtual network interface which driver is %s %s", nmDev.Driver.Get(), devPath)
		logger.Info(err)
		return
	}

	devType := nmDev.DeviceType.Get()
	if !isDeviceTypeValid(devType) {
		err = fmt.Errorf("ignore invalid device type %d", devType)
		logger.Info(err)
		return
	}

	dev = &device{
		nmDev:     nmDev,
		nmDevType: devType,
		udi:       nmDev.Udi.Get(),
		Path:      nmDev.Path,
		State:     nmDev.State.Get(),
		Interface: nmDev.Interface.Get(),
		Driver:    nmDev.Driver.Get(),
	}
	dev.Managed = nmGeneralIsDeviceManaged(devPath)
	dev.Vendor = nmGeneralGetDeviceDesc(devPath)
	dev.UsbDevice = nmGeneralIsUsbDevice(devPath)
	dev.id, _ = nmGeneralGetDeviceIdentifier(devPath)
	dev.UniqueUuid = nmGeneralGetDeviceUniqueUuid(devPath)

	// dispatch for different device types
	switch dev.nmDevType {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		if nmDevWired, err := nmNewDeviceWired(dev.Path); err == nil {
			dev.nmDevWired = nmDevWired
			// for mac address clone
			dev.nmDevWired.HwAddress.ConnectChanged(func() {
				v := dev.nmDevWired.HwAddress.Get()
				if v == dev.HwAddress {
					return
				}
				dev.HwAddress = v
				m.setPropDevices()
			})
			dev.HwAddress = dev.nmDevWired.HwAddress.Get()
			dev.PermHwAddress = dev.nmDevWired.PermHwAddress.Get()
		}
		if nmHasSystemSettingsModifyPermission() {
			m.ensureWiredConnectionExists(dev.Path, true)
		}
	case nm.NM_DEVICE_TYPE_WIFI:
		if nmDevWireless, err := nmNewDeviceWireless(dev.Path); err == nil {
			dev.nmDevWireless = nmDevWireless
			dev.HwAddress = nmDevWireless.HwAddress.Get()
			dev.PermHwAddress = nmDevWireless.PermHwAddress.Get()

			// connect property, about wireless active access point
			dev.nmDevWireless.ActiveAccessPoint.ConnectChanged(func() {
				if !m.isDeviceExists(devPath) {
					return
				}
				m.devicesLock.Lock()
				defer m.devicesLock.Unlock()
				dev.ActiveAp = nmDevWireless.ActiveAccessPoint.Get()
				m.setPropDevices()
			})
			dev.ActiveAp = nmDevWireless.ActiveAccessPoint.Get()
			dev.SupportHotspot = isWirelessDeviceSuportHotspot(nmDev.Interface.Get())

			dev.nmDevWireless.HwAddress.ConnectChanged(func() {
				v := dev.nmDevWireless.HwAddress.Get()
				if v == dev.HwAddress {
					return
				}
				dev.HwAddress = v
				m.setPropDevices()
			})
			// connect signals AccessPointAdded() and AccessPointRemoved()
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
	case nm.NM_DEVICE_TYPE_MODEM:
		if len(dev.id) == 0 {
			// some times, modem device will not be identified
			// successful for battery issue, so check and ignore it
			// here
			err = fmt.Errorf("modem device is not properly identified, please re-plugin it")
			return
		}
		m.ensureMobileConnectionExists(dev.Path, false)
		go func() {
			// disable autoconnect property for mobile devices
			// notice: sleep is necessary seconds before setting dbus values
			// FIXME: seems network-manager will restore Autoconnect's value some times
			time.Sleep(3 * time.Second)
			nmSetDeviceAutoconnect(dev.Path, false)
		}()
		if mmDevModem, err := mmNewModem(dbus.ObjectPath(dev.udi)); err == nil {
			dev.mmDevModem = mmDevModem

			// connect properties
			dev.mmDevModem.AccessTechnologies.ConnectChanged(func() {
				if !m.isDeviceExists(devPath) {
					return
				}
				m.devicesLock.Lock()
				defer m.devicesLock.Unlock()
				dev.MobileNetworkType = mmDoGetModemMobileNetworkType(mmDevModem.AccessTechnologies.Get())
				m.setPropDevices()
			})
			dev.MobileNetworkType = mmDoGetModemMobileNetworkType(mmDevModem.AccessTechnologies.Get())

			dev.mmDevModem.SignalQuality.ConnectChanged(func() {
				if !m.isDeviceExists(devPath) {
					return
				}
				m.devicesLock.Lock()
				defer m.devicesLock.Unlock()
				dev.MobileSignalQuality = mmDoGetModemDeviceSignalQuality(mmDevModem.SignalQuality.Get())
				m.setPropDevices()
			})
			dev.MobileSignalQuality = mmDoGetModemDeviceSignalQuality(mmDevModem.SignalQuality.Get())
		}
	}

	// connect signals
	dev.nmDev.ConnectStateChanged(func(newState, oldState, reason uint32) {
		logger.Infof("device state changed, %d => %d, reason[%d] %s", oldState, newState, reason, deviceErrorTable[reason])
		if !m.isDeviceExists(devPath) {
			return
		}

		m.devicesLock.Lock()
		defer m.devicesLock.Unlock()
		if reason == nm.NM_DEVICE_STATE_REASON_REMOVED && !isNmDeviceObjectExists(dev.Path) {
			// check if device dbus object exists for that if device
			// was set to unmanaged it will notify device removed, too
			return
		}
		dev.State = newState
		dev.Managed = nmGeneralIsDeviceManaged(dev.Path)

		// need get device vendor again for that some usb device may
		// not ready before
		dev.Interface = dev.nmDev.Interface.Get()
		dev.Vendor = nmGeneralGetDeviceDesc(dev.Path)
		dev.UsbDevice = nmGeneralIsUsbDevice(dev.Path)

		m.setPropDevices()

		m.config.updateDeviceConfig(dev.Path)
		m.config.syncDeviceState(dev.Path)
	})
	dev.State = dev.nmDev.State.Get()

	m.config.addDeviceConfig(devPath)
	m.switchHandler.initDeviceState(devPath)

	return
}
func (m *Manager) destroyDevice(dev *device) {
	// destroy object to reset all property connects
	if dev.nmDevWired != nil {
		nmDestroyDeviceWired(dev.nmDevWired)
	}
	if dev.nmDevWireless != nil {
		nmDestroyDeviceWireless(dev.nmDevWireless)
	}
	if dev.mmDevModem != nil {
		mmDestroyModem(dev.mmDevModem)
	}
	nmDestroyDevice(dev.nmDev)
}

func (m *Manager) clearDevices() {
	m.devicesLock.Lock()
	defer m.devicesLock.Unlock()
	for _, devs := range m.devices {
		for _, dev := range devs {
			m.destroyDevice(dev)
		}
	}
	m.devices = make(map[string][]*device)
	m.setPropDevices()
}

func (m *Manager) addDevice(devPath dbus.ObjectPath) {
	if m.isDeviceExists(devPath) {
		logger.Warning("device already exists", devPath)
		return
	}

	m.devicesLock.Lock()
	defer m.devicesLock.Unlock()
	dev, err := m.newDevice(devPath)
	if err != nil {
		return
	}
	logger.Infof("add device %#v", dev)
	devType := getCustomDeviceType(dev.nmDevType)
	m.devices[devType] = append(m.devices[devType], dev)
	m.setPropDevices()
}

func (m *Manager) removeDevice(devPath dbus.ObjectPath) {
	if !m.isDeviceExists(devPath) {
		return
	}
	devType, i := m.getDeviceIndex(devPath)

	m.devicesLock.Lock()
	defer m.devicesLock.Unlock()
	m.devices[devType] = m.doRemoveDevice(m.devices[devType], i)
	m.setPropDevices()
}
func (m *Manager) doRemoveDevice(devs []*device, i int) []*device {
	logger.Infof("remove device %#v", devs[i])
	m.destroyDevice(devs[i])
	copy(devs[i:], devs[i+1:])
	devs[len(devs)-1] = nil
	devs = devs[:len(devs)-1]
	return devs
}

func (m *Manager) getDevice(devPath dbus.ObjectPath) (dev *device) {
	devType, i := m.getDeviceIndex(devPath)
	if i < 0 {
		return
	}

	m.devicesLock.Lock()
	defer m.devicesLock.Unlock()
	return m.devices[devType][i]
}
func (m *Manager) isDeviceExists(devPath dbus.ObjectPath) bool {
	_, i := m.getDeviceIndex(devPath)
	if i >= 0 {
		return true
	}
	return false
}
func (m *Manager) getDeviceIndex(devPath dbus.ObjectPath) (devType string, index int) {
	m.devicesLock.Lock()
	defer m.devicesLock.Unlock()
	for t, devs := range m.devices {
		for i, dev := range devs {
			if dev.Path == devPath {
				return t, i
			}
		}
	}
	return "", -1
}

func (m *Manager) IsDeviceEnabled(devPath dbus.ObjectPath) (enabled bool, err error) {
	enabled = m.config.getDeviceEnabled(devPath)
	return
}

func (m *Manager) EnableDevice(devPath dbus.ObjectPath, enabled bool) (err error) {
	return m.switchHandler.enableDevice(devPath, enabled)
}

// SetDeviceManaged set target device managed or unmnaged from
// NetworkManager, and a little difference with other interface is
// that devPathOrIfc could be a device DBus path or the device
// interface name.
func (m *Manager) SetDeviceManaged(devPathOrIfc string, managed bool) (err error) {
	var devPath dbus.ObjectPath
	if strings.HasPrefix(devPathOrIfc, "/org/freedesktop/NetworkManager/Devices") {
		devPath = dbus.ObjectPath(devPathOrIfc)
	} else {
		m.devicesLock.Lock()
		defer m.devicesLock.Unlock()
	out:
		for _, devs := range m.devices {
			for _, dev := range devs {
				if dev.Interface == devPathOrIfc {
					devPath = dev.Path
					break out
				}
			}
		}
	}
	if len(devPath) > 0 {
		err = nmSetDeviceManaged(devPath, managed)
	} else {
		err = fmt.Errorf("invalid device identifier: %s", devPathOrIfc)
		logger.Error(err)
	}
	return
}

// ListDeviceConnections return the available connections for the device
func (m *Manager) ListDeviceConnections(devPath dbus.ObjectPath) ([]dbus.ObjectPath, error) {
	nmDev, err := nmNewDevice(devPath)
	if err != nil {
		return nil, err
	}
	defer nmDestroyDevice(nmDev)

	// ignore virtual network interfaces
	if isVirtualDeviceIfc(nmDev) {
		err = fmt.Errorf("ignore virtual network interface which driver is %s %s", nmDev.Driver.Get(), devPath)
		logger.Info(err)
		return nil, err
	}

	devType := nmDev.DeviceType.Get()
	if !isDeviceTypeValid(devType) {
		err = fmt.Errorf("ignore invalid device type %d", devType)
		logger.Info(err)
		return nil, err
	}

	return nmDev.AvailableConnections.Get(), nil
}
