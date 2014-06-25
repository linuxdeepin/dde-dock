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
	"io/ioutil"
	"os"
	"path"
	"pkg.linuxdeepin.com/lib/dbus"
)

var networkConfigFile = os.Getenv("HOME") + "/.config/deepin_network.json"

type config struct {
	configFile   string
	WiredEnabled bool
	VpnEnabled   bool
	Devices      map[string]*deviceConfig
}

type deviceConfig struct {
	Enabled        bool
	LastConnection string // uuid of last activated connection
}

func newConfig() (c *config) {
	c = &config{}
	c.setConfigFile(networkConfigFile)
	c.WiredEnabled = true
	c.VpnEnabled = true
	c.Devices = make(map[string]*deviceConfig)
	c.load()
	return
}

func newDeviceConfig() (d *deviceConfig) {
	d = &deviceConfig{}
	d.Enabled = true
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
	ensureDirExists(path.Dir(c.configFile))
	fileContent, _ := marshalJSON(c)
	err := ioutil.WriteFile(c.configFile, []byte(fileContent), 0644)
	if err != nil {
		logger.Error(err)
	}
}

func (c *config) getDeviceConfig(hwAddr string) (d *deviceConfig) {
	d, ok := c.Devices[hwAddr]
	if !ok {
		logger.Errorf("device config for %s not exists", hwAddr)
	}
	return
}

func (m *Manager) IsDeviceEnabled(devPath dbus.ObjectPath) (enabled bool, err error) {
	// TODO
	return
}
func (m *Manager) EnableDevice(devPath dbus.ObjectPath, enabled bool) (err error) {
	// TODO
	// hwAddr, err := nmGeneralGetDeviceHwAddr(devPath)
	if enabled {
		// devconf := m.config.getDeviceConfig(hwAddr)
	} else {
	}
	if m.DeviceEnabled != nil {
		m.DeviceEnabled(enabled)
	}
	return
}
