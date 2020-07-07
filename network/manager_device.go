/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package network

import (
	"fmt"
	"net"
	"strings"
	"time"

	mmdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.modemmanager1"
	nmdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.networkmanager"

	"pkg.deepin.io/dde/daemon/network/nm"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

type device struct {
	nmDev      *nmdbus.Device
	mmDevModem *mmdbus.Modem
	nmDevType  uint32
	id         string
	udi        string

	Path          dbus.ObjectPath
	State         uint32
	Interface     string
	ClonedAddress string
	HwAddress     string
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

	InterfaceFlags uint32

	quitFlagsCheckChan chan struct{}
}

const (
	nmInterfaceFlagUP uint32 = 1 << iota
	nmInterfaceFlagLowerUP

	nmInterfaceFlagCarrier = 0x10000
)

func (m *Manager) initDeviceManage() {
	m.devicesLock.Lock()
	m.devices = make(map[string][]*device)
	m.devicesLock.Unlock()

	m.accessPointsLock.Lock()
	m.accessPoints = make(map[dbus.ObjectPath][]*accessPoint)
	m.accessPointsLock.Unlock()

	_, err := nmManager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		m.addDevice(path)
	})
	if err != nil {
		logger.Warning(err)
	}
	_, err = nmManager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		notifyDeviceRemoved(path)
		m.removeDevice(path)
	})
	if err != nil {
		logger.Warning(err)
	}
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
		driver, _ := nmDev.Driver().Get(0)
		err = fmt.Errorf("ignore virtual network interface which driver is %s %s", driver, devPath)
		logger.Info(err)
		return
	}

	devType, _ := nmDev.DeviceType().Get(0)
	if !isDeviceTypeValid(devType) {
		err = fmt.Errorf("ignore invalid device type %d", devType)
		logger.Info(err)
		return
	}

	dev = &device{
		nmDev:     nmDev,
		nmDevType: devType,
		Path:      nmDev.Path_(),
	}
	dev.udi, _ = nmDev.Udi().Get(0)
	dev.State, _ = nmDev.State().Get(0)
	dev.Interface, _ = nmDev.Interface().Get(0)
	dev.Driver, _ = nmDev.Driver().Get(0)

	dev.Managed = nmGeneralIsDeviceManaged(devPath)
	dev.Vendor = nmGeneralGetDeviceDesc(devPath)
	dev.UsbDevice = nmGeneralIsUsbDevice(devPath)
	dev.id, _ = nmGeneralGetDeviceIdentifier(devPath)
	dev.UniqueUuid = nmGeneralGetDeviceUniqueUuid(devPath)

	nmDev.InitSignalExt(m.sysSigLoop, true)

	// dispatch for different device types
	switch dev.nmDevType {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		// for mac address clone
		nmDevWired := nmDev.Wired()
		err = nmDevWired.HwAddress().ConnectChanged(func(hasValue bool, value string) {
			if !hasValue {
				return
			}
			if value == dev.ClonedAddress {
				return
			}
			dev.ClonedAddress = value
			m.updatePropDevices()
		})
		if err != nil {
			logger.Warning(err)
		}
		dev.ClonedAddress, _ = nmDevWired.HwAddress().Get(0)
		dev.HwAddress, _ = nmDevWired.PermHwAddress().Get(0)

		if dev.HwAddress == "" {
			dev.HwAddress = dev.ClonedAddress
		}

		if nmHasSystemSettingsModifyPermission() {
			carrierChanged := func(hasValue, value bool) {
				if !hasValue || !value {
					return
				}

				logger.Info("wired plugin", dev.Path)
				logger.Debug("ensure wired connection exists", dev.Path)
				_, _, err = m.ensureWiredConnectionExists(dev.Path, true)
				if err != nil {
					logger.Warning(err)
				}
			}

			nmDev.Wired().Carrier().ConnectChanged(carrierChanged)

			carrier, _ := nmDev.Wired().Carrier().Get(0)
			carrierChanged(true, carrier)
		} else {
			logger.Debug("do not have modify permission")
		}
	case nm.NM_DEVICE_TYPE_WIFI:
		nmDevWireless := nmDev.Wireless()
		dev.ClonedAddress, _ = nmDevWireless.HwAddress().Get(0)
		dev.HwAddress, _ = nmDevWireless.PermHwAddress().Get(0)

		// connect property, about wireless active access point
		err = nmDevWireless.ActiveAccessPoint().ConnectChanged(func(hasValue bool,
			value dbus.ObjectPath) {
			if !hasValue {
				return
			}
			if !m.isDeviceExists(devPath) {
				return
			}
			m.devicesLock.Lock()
			defer m.devicesLock.Unlock()
			dev.ActiveAp = value
			m.updatePropDevices()

			// Re-active connection if wireless 'ActiveAccessPoint' not equal active connection 'SpecificObject'
			// such as wifi roaming, but the active connection state is activated
			err := m.wirelessReActiveConnection(nmDev)
			if err != nil {
				logger.Warning("Failed to re-active connection:", err)
			}
		})
		if err != nil {
			logger.Warning(err)
		}
		dev.ActiveAp, _ = nmDevWireless.ActiveAccessPoint().Get(0)
		permHwAddress, _ := nmDevWireless.PermHwAddress().Get(0)
		dev.SupportHotspot = isWirelessDeviceSupportHotspot(permHwAddress)

		err = nmDevWireless.HwAddress().ConnectChanged(func(hasValue bool, value string) {
			if !hasValue {
				return
			}
			if value == dev.ClonedAddress {
				return
			}
			dev.ClonedAddress = value
			m.updatePropDevices()
		})
		if err != nil {
			logger.Warning(err)
		}
		// connect signals AccessPointAdded() and AccessPointRemoved()
		_, err = nmDevWireless.ConnectAccessPointAdded(func(apPath dbus.ObjectPath) {
			m.addAccessPoint(dev.Path, apPath)
		})
		if err != nil {
			logger.Warning(err)
		}

		_, err = nmDevWireless.ConnectAccessPointRemoved(func(apPath dbus.ObjectPath) {
			m.removeAccessPoint(dev.Path, apPath)
		})
		if err != nil {
			logger.Warning(err)
		}
		for _, apPath := range nmGetAccessPoints(dev.Path) {
			m.addAccessPoint(dev.Path, apPath)
		}
	case nm.NM_DEVICE_TYPE_MODEM:
		if len(dev.id) == 0 {
			// some times, modem device will not be identified
			// successful for battery issue, so check and ignore it
			// here
			err = fmt.Errorf("modem device is not properly identified, please re-plugin it")
			return
		}
		go func() {
			// disable autoconnect property for mobile devices
			// notice: sleep is necessary seconds before setting dbus values
			// FIXME: seems network-manager will restore Autoconnect's value some times
			time.Sleep(3 * time.Second)
			nmSetDeviceAutoconnect(dev.Path, false)
		}()
		if mmDevModem, err := mmNewModem(dbus.ObjectPath(dev.udi)); err == nil {
			mmDevModem.InitSignalExt(m.sysSigLoop, true)
			dev.mmDevModem = mmDevModem

			// connect properties
			err = dev.mmDevModem.AccessTechnologies().ConnectChanged(func(hasValue bool,
				value uint32) {
				if !m.isDeviceExists(devPath) {
					return
				}
				if !hasValue {
					return
				}
				m.devicesLock.Lock()
				defer m.devicesLock.Unlock()
				dev.MobileNetworkType = mmDoGetModemMobileNetworkType(value)
				m.updatePropDevices()
			})
			if err != nil {
				logger.Warning(err)
			}
			accessTech, _ := mmDevModem.AccessTechnologies().Get(0)
			dev.MobileNetworkType = mmDoGetModemMobileNetworkType(accessTech)

			err = dev.mmDevModem.SignalQuality().ConnectChanged(func(hasValue bool,
				value mmdbus.ModemSignalQuality) {
				if !m.isDeviceExists(devPath) {
					return
				}
				if !hasValue {
					return
				}

				m.devicesLock.Lock()
				defer m.devicesLock.Unlock()
				dev.MobileSignalQuality = value.Quality
				m.updatePropDevices()
			})
			if err != nil {
				logger.Warning(err)
			}
			dev.MobileSignalQuality = mmDoGetModemDeviceSignalQuality(mmDevModem)
		}
	}

	// connect signals
	_, err = dev.nmDev.ConnectStateChanged(func(newState uint32, oldState uint32, reason uint32) {
		logger.Debugf("device state changed, %d => %d, reason[%d] %s",
			oldState, newState, reason, deviceErrorTable[reason])

		if !m.isDeviceExists(devPath) {
			return
		}

		dev.State = newState
		m.devicesLock.Lock()
		m.updatePropDevices()
		m.devicesLock.Unlock()

	})
	if err != nil {
		logger.Warning(err)
	}

	err = dev.nmDev.Interface().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}

		dev.Interface = value
		m.devicesLock.Lock()
		m.updatePropDevices()
		m.devicesLock.Unlock()
	})
	if err != nil {
		logger.Warning(err)
	}

	err = dev.nmDev.Managed().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}

		dev.Managed = value
		m.devicesLock.Lock()
		m.updatePropDevices()
		m.devicesLock.Unlock()
	})
	if err != nil {
		logger.Warning(err)
	}

	// TODO: NetworkManager 升级 1.22 后，直接使用 NetworkManager 的 InterfaceFlags 属性
	dev.InterfaceFlags = m.getInterfaceFlags(dev)
	ticker := time.NewTicker(1 * time.Second)
	dev.quitFlagsCheckChan = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				newFlags := m.getInterfaceFlags(dev)
				if dev.InterfaceFlags != newFlags {
					dev.InterfaceFlags = newFlags

					m.devicesLock.Lock()
					m.updatePropDevices()
					m.devicesLock.Unlock()
				}

			case <-dev.quitFlagsCheckChan:
				ticker.Stop()
				return
			}
		}
	}()

	return
}

func (m *Manager) getInterfaceFlags(dev *device) uint32 {
	interfaceInfo, err := net.InterfaceByName(dev.Interface)
	if err != nil {
		logger.Warning("failed to get interface info:", err)
		return 0
	}

	var flags uint32
	if interfaceInfo.Flags&net.FlagUp != 0 {
		flags |= nmInterfaceFlagUP
	}

	return flags
}

func (m *Manager) destroyDevice(dev *device) {
	close(dev.quitFlagsCheckChan)

	// destroy object to reset all property connects
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
	m.updatePropDevices()
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
	logger.Debug("add device", devPath)
	devType := getCustomDeviceType(dev.nmDevType)
	m.devices[devType] = append(m.devices[devType], dev)
	m.updatePropDevices()
}

func (m *Manager) removeDevice(devPath dbus.ObjectPath) {
	if !m.isDeviceExists(devPath) {
		return
	}
	devType, i := m.getDeviceIndex(devPath)

	m.devicesLock.Lock()
	defer m.devicesLock.Unlock()
	m.devices[devType] = m.doRemoveDevice(m.devices[devType], i)
	m.updatePropDevices()
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

func (m *Manager) IsDeviceEnabled(devPath dbus.ObjectPath) (bool, *dbus.Error) {
	b, err := m.sysNetwork.IsDeviceEnabled(0, string(devPath))
	return b, dbusutil.ToError(err)
}

func (m *Manager) EnableDevice(devPath dbus.ObjectPath, enabled bool) *dbus.Error {
	err := m.enableDevice(devPath, enabled)
	return dbusutil.ToError(err)
}

func (m *Manager) enableDevice(devPath dbus.ObjectPath, enabled bool) (err error) {
	cpath, err := m.sysNetwork.EnableDevice(0, string(devPath), enabled)
	if err != nil {
		return
	}
	if enabled {
		var uuid string
		uuid, err = nmGetConnectionUuid(cpath)
		if err != nil {
			return
		}
		m.ActivateConnection(uuid, devPath)
	}

	m.stateHandler.locker.Lock()
	defer m.stateHandler.locker.Unlock()
	dsi, ok := m.stateHandler.devices[devPath]
	if !ok {
		return
	}
	dsi.enabled = enabled
	return
}

// SetDeviceManaged set target device managed or unmnaged from
// NetworkManager, and a little difference with other interface is
// that devPathOrIfc could be a device DBus path or the device
// interface name.
func (m *Manager) SetDeviceManaged(devPathOrIfc string, managed bool) *dbus.Error {
	err := m.setDeviceManaged(devPathOrIfc, managed)
	return dbusutil.ToError(err)
}

func (m *Manager) setDeviceManaged(devPathOrIfc string, managed bool) (err error) {
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
func (m *Manager) ListDeviceConnections(devPath dbus.ObjectPath) ([]dbus.ObjectPath, *dbus.Error) {
	paths, err := m.listDeviceConnections(devPath)
	return paths, dbusutil.ToError(err)
}

func (m *Manager) listDeviceConnections(devPath dbus.ObjectPath) ([]dbus.ObjectPath, error) {
	nmDev, err := nmNewDevice(devPath)
	if err != nil {
		return nil, err
	}

	// ignore virtual network interfaces
	if isVirtualDeviceIfc(nmDev) {
		driver, _ := nmDev.Driver().Get(0)
		err = fmt.Errorf("ignore virtual network interface which driver is %s %s", driver, devPath)
		logger.Info(err)
		return nil, err
	}

	devType, _ := nmDev.DeviceType().Get(0)
	if !isDeviceTypeValid(devType) {
		err = fmt.Errorf("ignore invalid device type %d", devType)
		logger.Info(err)
		return nil, err
	}

	availableConnections, _ := nmDev.AvailableConnections().Get(0)
	return availableConnections, nil
}

// RequestWirelessScan request all wireless devices re-scan access point list.
func (m *Manager) RequestWirelessScan() *dbus.Error {
	m.devicesLock.Lock()
	defer m.devicesLock.Unlock()

	if devices, ok := m.devices[deviceWifi]; ok {
		for _, dev := range devices {
			err := dev.nmDev.RequestScan(0, nil)
			if err != nil {
				logger.Debug(err)
			}
		}
	}
	return nil
}

func (m *Manager) wirelessReActiveConnection(nmDev *nmdbus.Device) error {
	wireless := nmDev.Wireless()
	apPath, err := wireless.ActiveAccessPoint().Get(0)
	if err != nil {
		return err
	}
	if apPath == "/" {
		logger.Debug("Invalid active access point path:", nmDev.Path_())
		return nil
	}
	connPath, err := nmDev.ActiveConnection().Get(0)
	if err != nil {
		return err
	}
	if connPath == "/" {
		logger.Debug("Invalid active connection path:", nmDev.Path_())
		return nil
	}

	connObj, err := nmNewActiveConnection(connPath)
	if err != nil {
		return err
	}

	// check network connectivity state
	state, err := connObj.State().Get(0)
	if err != nil {
		return err
	}
	if state != nm.NM_ACTIVE_CONNECTION_STATE_ACTIVATED {
		logger.Debug("[Inactive] re-active connection not activated:", connPath, nmDev.Path_(), state)
		return nil
	}

	spePath, err := connObj.SpecificObject().Get(0)
	if err != nil {
		return err
	}
	if spePath == "/" {
		logger.Debug("Invalid specific access point path:", connObj.Path_(), nmDev.Path_())
		return nil
	}

	if string(apPath) == string(spePath) {
		logger.Debug("[NONE] re-active connection not changed:", connPath, spePath, nmDev.Path_())
		return nil
	}

	ip4Path, _ := connObj.Ip4Config().Get(0)
	if m.checkGatewayConnectivity(ip4Path) {
		logger.Debug("Network is connectivity, don't re-active")
		return nil
	}

	settingsPath, err := connObj.Connection().Get(0)
	if err != nil {
		return err
	}
	logger.Debug("[DO] re-active connection:", settingsPath, connPath, spePath, nmDev.Path_())
	_, err = nmActivateConnection(settingsPath, nmDev.Path_())
	return err
}

func (m *Manager) checkGatewayConnectivity(ipPath dbus.ObjectPath) bool {
	if ipPath == "/" {
		return false
	}

	addr, mask, gateways, domains := nmGetIp4ConfigInfo(ipPath)
	logger.Debugf("The active connection ip4 info: address(%s), mask(%s), gateways(%v), domains(%v)",
		addr, mask, gateways, domains)
	// check whether the gateway is connected by ping
	for _, gw := range gateways {
		if len(gw) == 0 {
			continue
		}
		if !m.doPing(gw, 3) {
			return false
		}
	}
	return true
}

func (m *Manager) doPing(addr string, retries int) bool {
	for i := 0; i < retries; i++ {
		err := m.sysNetwork.Ping(0, addr)
		if err == nil {
			return true
		}
		logger.Warning("Failed to ping gateway:", i, addr, err)
		time.Sleep(time.Millisecond * 500)
	}
	return false
}
