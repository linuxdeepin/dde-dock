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

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"pkg.linuxdeepin.com/lib/dbus"
	"sync"
)

var networkConfigFile = os.Getenv("HOME") + "/.config/deepin_network.json"

type config struct {
	configFile string
	saveLock   sync.Mutex

	WiredEnabled        bool
	VpnEnabled          bool
	LastWirelessEnabled bool
	LastWwanEnabled     bool
	LastWiredEnabled    bool
	LastVpnEnabled      bool

	Devices        map[string]*deviceConfig // config for each device
	VpnConnections map[string]*vpnConfig    // config for each vpn connection
}

type deviceConfig struct {
	Enabled     bool
	LastEnabled bool

	// uuid of last activated connection, don't need to save it to
	// configuration file for that NetworkManager will decide which
	// connection to activate for each device after login, so just
	// follow NetworkManager's choice.
	lastConnectionUuid string
}

type vpnConfig struct {
	AutoConnect bool

	// don't need to save activated state
	activated     bool
	lastActivated bool
}

func newConfig() (c *config) {
	c = &config{}
	c.setConfigFile(networkConfigFile)
	c.WiredEnabled = true
	c.VpnEnabled = true
	c.Devices = make(map[string]*deviceConfig)
	c.VpnConnections = make(map[string]*vpnConfig)
	c.LastWirelessEnabled = true
	c.LastWwanEnabled = true
	c.LastWiredEnabled = true
	c.LastVpnEnabled = true
	c.load()
	c.clearSpareConfig()
	return
}

func newDeviceConfig() (d *deviceConfig) {
	d = &deviceConfig{}
	d.Enabled = true
	d.LastEnabled = d.Enabled
	return
}

func newVpnConfig() (v *vpnConfig) {
	v = &vpnConfig{}
	v.AutoConnect = false
	v.activated = false
	v.lastActivated = v.activated
	return
}

func (c *config) setConfigFile(file string) {
	c.configFile = file
}

func (c *config) load() {
	if isFileExists(c.configFile) {
		fileContent, err := ioutil.ReadFile(c.configFile)
		if err != nil {
			logger.Error(err)
			return
		}
		unmarshalJSON(string(fileContent), c)
	} else {
		c.save()
	}
}
func (c *config) save() {
	c.saveLock.Lock()
	defer c.saveLock.Unlock()
	ensureDirExists(path.Dir(c.configFile))
	fileContent, _ := marshalJSON(c)
	err := ioutil.WriteFile(c.configFile, []byte(fileContent), 0644)
	if err != nil {
		logger.Error(err)
	}
}

func (c *config) clearSpareConfig() {
	// remove spare device and vpn config
	devIds := nmGetDeviceIdentifiers()
	for id, _ := range c.Devices {
		if !isStringInArray(id, devIds) {
			c.removeDeviceConfig(id)
		}
	}
	vpnUuids := nmGetSpecialConnectionUuids(NM_SETTING_VPN_SETTING_NAME)
	for uuid, _ := range c.VpnConnections {
		if !isStringInArray(uuid, vpnUuids) {
			c.removeVpnConfig(uuid)
		}
	}
}

func (c *config) setLastGlobalSwithes(enabled bool) {
	c.LastWirelessEnabled = enabled
	c.LastWwanEnabled = enabled
	c.LastWiredEnabled = enabled
	c.LastVpnEnabled = enabled
	c.save()
}
func (c *config) setLastWirelessEnabled(enabled bool) {
	if c.LastWirelessEnabled != enabled {
		c.LastWirelessEnabled = enabled
		c.save()
	}
}
func (c *config) setLastWwanEnabled(enabled bool) {
	if c.LastWwanEnabled != enabled {
		c.LastWwanEnabled = enabled
		c.save()
	}
}
func (c *config) setLastWiredEnabled(enabled bool) {
	if c.LastWiredEnabled != enabled {
		c.LastWiredEnabled = enabled
		c.save()
	}
}
func (c *config) setLastVpnEnabled(enabled bool) {
	if c.LastVpnEnabled != enabled {
		c.LastVpnEnabled = enabled
		c.save()
	}
}

func (c *config) setWiredEnabled(enabled bool) {
	if c.WiredEnabled != enabled {
		c.LastWiredEnabled = c.WiredEnabled
		c.WiredEnabled = enabled
		c.save()
	}
}
func (c *config) setVpnEnabled(enabled bool) {
	if c.VpnEnabled != enabled {
		c.LastVpnEnabled = c.VpnEnabled
		c.VpnEnabled = enabled
		c.save()
	}
}

func (c *config) removeConnection(uuid string) {
	for _, devConfig := range c.Devices {
		if devConfig.lastConnectionUuid == uuid {
			devConfig.lastConnectionUuid = ""
		}
	}
	c.removeVpnConfig(uuid)
	c.save()
}

// deviceConfig related functions
func (c *config) isDeviceConfigExists(devId string) (ok bool) {
	_, ok = c.Devices[devId]
	return
}
func (c *config) getDeviceConfigByPath(devPath dbus.ObjectPath) (d *deviceConfig, err error) {
	devId, err := nmGetDeviceIdentifier(devPath)
	if err != nil {
		return
	}
	return c.getDeviceConfig(devId)
}
func (c *config) getDeviceConfig(devId string) (d *deviceConfig, err error) {
	if !c.isDeviceConfigExists(devId) {
		err = fmt.Errorf("device config for %s not exists", devId)
		logger.Error(err)
		return
	}
	d, _ = c.Devices[devId]
	return
}
func (c *config) addDeviceConfig(devPath dbus.ObjectPath) {
	devId, err := nmGetDeviceIdentifier(devPath)
	if err != nil {
		return
	}
	if !c.isDeviceConfigExists(devId) {
		devConfig := newDeviceConfig()
		devConfig.lastConnectionUuid, _ = nmGetDeviceActiveConnectionUuid(devPath)
		c.Devices[devId] = devConfig
		c.save()
	} else {
		c.updateDeviceConfig(devPath)
	}
}
func (c *config) removeDeviceConfig(devId string) {
	if !c.isDeviceConfigExists(devId) {
		logger.Errorf("device config for %s not exists", devId)
	}
	delete(c.Devices, devId)
	c.save()
}
func (c *config) updateDeviceConfig(devPath dbus.ObjectPath) {
	devConfig, err := c.getDeviceConfigByPath(devPath)
	if err != nil {
		return
	}
	devState, _ := nmGetDeviceState(devPath)
	if isDeviceStateActivated(devState) {
		devConfig.lastConnectionUuid, _ = nmGetDeviceActiveConnectionUuid(devPath)
		c.setDeviceEnabled(devPath, true)
	}
	logger.Debugf("updateDeviceConfig %s %#v", devPath, devConfig) // TODO test
}
func (c *config) setDeviceEnabled(devPath dbus.ObjectPath, enabled bool) {
	devConfig, err := c.getDeviceConfigByPath(devPath)
	if err != nil {
		return
	}
	if devConfig.Enabled == enabled {
		return
	}
	devConfig.Enabled = enabled
	if manager.DeviceEnabled != nil {
		logger.Debug("DeviceEnabled", devPath, enabled) //  TODO test
		manager.DeviceEnabled(string(devPath), enabled)
	}
}
func (c *config) setDeviceLastConnectionUuid(devPath dbus.ObjectPath, uuid string) {
	devConfig, err := c.getDeviceConfigByPath(devPath)
	if err != nil {
		return
	}
	if devConfig.lastConnectionUuid != uuid {
		devConfig.lastConnectionUuid = uuid
		c.save()
	}
}
func (c *config) setDeviceLastEnabled(devPath dbus.ObjectPath, enabled bool) {
	devConfig, err := c.getDeviceConfigByPath(devPath)
	if err != nil {
		return
	}
	if devConfig.LastEnabled != enabled {
		devConfig.LastEnabled = enabled
		c.save()
	}
}
func (c *config) setAllDeviceLastEnabled(enabled bool) {
	for _, devConfig := range c.Devices {
		devConfig.LastEnabled = enabled
	}
	c.save()
}
func (c *config) saveDeviceState(devPath dbus.ObjectPath) {
	devConfig, err := c.getDeviceConfigByPath(devPath)
	if err != nil {
		return
	}
	if devConfig.LastEnabled != devConfig.Enabled {
		devConfig.LastEnabled = devConfig.Enabled
		c.save()
	}
}

// vpnConfig related functions
func (c *config) isVpnConfigExists(uuid string) (ok bool) {
	_, ok = c.VpnConnections[uuid]
	return
}
func (c *config) getVpnConfig(uuid string) (v *vpnConfig, err error) {
	if !c.isVpnConfigExists(uuid) {
		err = fmt.Errorf("vpn config for %s not exists", uuid)
		logger.Error(err)
		return
	}
	v, _ = c.VpnConnections[uuid]
	return
}
func (c *config) addVpnConfig(uuid string) {
	if !c.isVpnConfigExists(uuid) {
		vpnConfig := newVpnConfig()
		c.VpnConnections[uuid] = vpnConfig
		c.save()
	}
}
func (c *config) removeVpnConfig(uuid string) {
	if c.isVpnConfigExists(uuid) {
		delete(c.VpnConnections, uuid)
		c.save()
	}
}
func (c *config) setVpnConnectionActivated(uuid string, activated bool) {
	vpnConfig, err := c.getVpnConfig(uuid)
	if err != nil {
		return
	}
	if vpnConfig.activated != activated {
		vpnConfig.activated = activated
		c.save()
	}
}

// Manager related functions
func (m *Manager) IsDeviceEnabled(devPath dbus.ObjectPath) (enabled bool, err error) {
	devConfig, err := m.config.getDeviceConfigByPath(devPath)
	if err != nil {
		enabled = true // return true as default
		return
	}
	enabled = devConfig.Enabled
	return
}

func (m *Manager) restoreDeviceState(devPath dbus.ObjectPath) (err error) {
	devConfig, err := m.config.getDeviceConfigByPath(devPath)
	if err != nil {
		return
	}
	err = m.EnableDevice(devPath, devConfig.LastEnabled)
	return
}
func (m *Manager) disconnectAndSaveDeviceState(devPath dbus.ObjectPath) (err error) {
	m.config.saveDeviceState(devPath)
	err = m.EnableDevice(devPath, false)
	return
}

func (m *Manager) EnableDevice(devPath dbus.ObjectPath, enabled bool) (err error) {
	if enabled && m.trunOnGlobalDeviceSwitchIfNeed(devPath) {
		return
	}
	devConfig, err := m.config.getDeviceConfigByPath(devPath)
	if err != nil {
		return
	}
	logger.Debugf("EnableDevice %s %v %#v", devPath, enabled, devConfig) // TODO test

	if enabled {
		// active last connection if device is disconnected
		if len(devConfig.lastConnectionUuid) > 0 {
			if _, err := nmGetConnectionByUuid(devConfig.lastConnectionUuid); err == nil {
				devState, err := nmGetDeviceState(devPath)
				if err == nil {
					if isDeviceStateAvailable(devState) && !isDeviceStateActivated(devState) {
						err = m.ActivateConnection(devConfig.lastConnectionUuid, devPath)
					}
				}
			}
		}
	} else {
		err = m.doDisconnectDevice(devPath)
	}

	m.config.setDeviceEnabled(devPath, enabled)
	m.config.save()
	return
}

func (m *Manager) restoreVpnConnectionState(uuid string) (err error) {
	vpnConfig, err := m.config.getVpnConfig(uuid)
	if err != nil {
		return
	}
	if vpnConfig.lastActivated {
		err = m.ActivateConnection(uuid, "/")
	} else {
		err = m.DeactivateConnection(uuid)
	}
	m.config.save()
	return
}
func (m *Manager) deactivateVpnConnection(uuid string) (err error) {
	vpnConfig, err := m.config.getVpnConfig(uuid)
	if err != nil {
		return
	}
	vpnConfig.lastActivated = vpnConfig.activated
	err = m.DeactivateConnection(uuid)
	m.config.save()
	return
}

func (m *Manager) trunOnGlobalDeviceSwitchIfNeed(devPath dbus.ObjectPath) (need bool) {
	// if global device switch is off, turn it on, and only keep
	// current device alive
	need = (m.generalGetGlobalDeviceEnabled(devPath) == false)
	if !need {
		return
	}
	m.config.setAllDeviceLastEnabled(false)
	m.config.setDeviceLastEnabled(devPath, true)
	m.generalSetGlobalDeviceEnabled(devPath, true)
	return
}

func (m *Manager) generalGetGlobalDeviceEnabled(devPath dbus.ObjectPath) (enabled bool) {
	devType, _ := nmGetDeviceType(devPath)
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		enabled = m.WiredEnabled
	case NM_DEVICE_TYPE_WIFI:
		enabled = m.WirelessEnabled
	case NM_DEVICE_TYPE_MODEM:
		enabled = m.WwanEnabled
	}
	return
}
func (m *Manager) generalSetGlobalDeviceEnabled(devPath dbus.ObjectPath, enabled bool) {
	devType, _ := nmGetDeviceType(devPath)
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		m.setWiredEnabled(enabled)
	case NM_DEVICE_TYPE_WIFI:
		m.setWirelessEnabled(enabled)
	case NM_DEVICE_TYPE_MODEM:
		m.setWwanEnabled(enabled)
	}
}
