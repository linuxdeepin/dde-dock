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

package bluetooth

import (
	"strings"

	"pkg.deepin.io/lib/utils"
)

type config struct {
	core utils.Config

	Adapters map[string]*adapterConfig // use adapter hardware address as key
	Devices  map[string]*deviceConfig  // use adapter address/device address as key
}

type adapterConfig struct {
	Powered bool
}

type deviceConfig struct {
	Connected bool
}

func newConfig() (c *config) {
	c = &config{}
	c.core.SetConfigName("bluetooth")
	logger.Info("load bluetooth config file:", c.core.GetConfigFile())
	c.Adapters = make(map[string]*adapterConfig)
	c.Devices = make(map[string]*deviceConfig)
	c.load()
	return
}

func (c *config) load() {
	c.core.Load(c)
}

func (c *config) save() {
	c.core.Save(c)
}

func newAdapterConfig() (ac *adapterConfig) {
	ac = &adapterConfig{Powered: true}
	return
}

func (c *config) clearSpareConfig(b *Bluetooth) {
	var adapterAddresses []string
	// key is adapter address
	var adapterDevicesMap = make(map[string][]*device)

	b.adaptersLock.Lock()
	for _, adapter := range b.adapters {
		adapterAddresses = append(adapterAddresses, adapter.address)
	}
	b.adaptersLock.Unlock()

	for _, adapterAddr := range adapterAddresses {
		adapterDevicesMap[adapterAddr] = b.getAdapterDevices(adapterAddr)
	}

	c.core.Lock()
	// remove spare adapters
	for addr := range c.Adapters {
		if !isStringInArray(addr, adapterAddresses) {
			delete(c.Adapters, addr)
		}
	}

	// remove spare devices
	for addr := range c.Devices {
		addrParts := strings.SplitN(addr, "/", 2)
		if len(addrParts) != 2 {
			delete(c.Devices, addr)
			continue
		}
		adapterAddr := addrParts[0]
		deviceAddr := addrParts[1]

		devices := adapterDevicesMap[adapterAddr]
		var foundDevice bool
		for _, device := range devices {
			if device.Address == deviceAddr {
				foundDevice = true
				break
			}
		}

		if !foundDevice {
			delete(c.Devices, addr)
			continue
		}
	}

	c.core.Unlock()
}

func (c *config) addAdapterConfig(address string) {
	if c.isAdapterConfigExists(address) {
		return
	}

	c.core.Lock()
	c.Adapters[address] = newAdapterConfig()
	c.core.Unlock()
	c.save()
}

func (c *config) removeAdapterConfig(address string) {
	c.core.Lock()
	delete(c.Adapters, address)
	c.core.Unlock()
	c.save()
}

func (c *config) getAdapterConfig(address string) (ac *adapterConfig, ok bool) {
	c.core.Lock()
	defer c.core.Unlock()
	ac, ok = c.Adapters[address]
	return
}

func (c *config) isAdapterConfigExists(address string) (ok bool) {
	c.core.Lock()
	defer c.core.Unlock()
	_, ok = c.Adapters[address]
	return
}

func (c *config) getAdapterConfigPowered(address string) (powered bool) {
	c.core.Lock()
	defer c.core.Unlock()
	if ac, ok := c.Adapters[address]; ok {
		return ac.Powered
	}
	return false
}

func (c *config) setAdapterConfigPowered(address string, powered bool) {
	c.core.Lock()
	if ac, ok := c.Adapters[address]; ok {
		ac.Powered = powered
	}
	c.core.Unlock()
	c.save()
	return
}

func newDeviceConfig() (ac *deviceConfig) {
	ac = &deviceConfig{Connected: false}
	return
}

func (c *config) isDeviceConfigExist(address string) (ok bool) {
	c.core.Lock()
	defer c.core.Unlock()
	_, ok = c.Devices[address]
	return
}

func (c *config) addDeviceConfig(address string) {
	if c.isDeviceConfigExist(address) {
		return
	}
	c.core.Lock()
	c.Devices[address] = newDeviceConfig()
	c.core.Unlock()
	c.save()
}

func (c *config) getDeviceConfig(address string) (dc *deviceConfig, ok bool) {
	c.core.Lock()
	defer c.core.Unlock()
	dc, ok = c.Devices[address]
	return
}

func (c *config) removeDeviceConfig(address string) {
	c.core.Lock()
	delete(c.Devices, address)
	c.core.Unlock()
	c.save()
}

func (c *config) getDeviceConfigConnected(address string) (connected bool) {
	dc, ok := c.getDeviceConfig(address)
	if !ok {
		return
	}

	c.core.Lock()
	defer c.core.Unlock()
	return dc.Connected
}

func (c *config) setDeviceConfigConnected(address string, connected bool) {
	dc, ok := c.getDeviceConfig(address)
	if !ok {
		return
	}

	c.core.Lock()
	dc.Connected = connected
	c.core.Unlock()

	c.save()
	return
}
