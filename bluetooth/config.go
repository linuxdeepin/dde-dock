/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/lib/utils"
)

type config struct {
	core utils.Config

	Adapters map[string]*adapterConfig // use adapter hardware address as key
	Devices  map[string]*deviceConfig  // use device dpath as key
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
	go c.clearSpareConfig()
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

func (c *config) clearSpareConfig() {
	var addresses []string
	apathes := bluezGetAdapters()
	for _, apath := range apathes {
		addresses = append(addresses, bluezGetAdapterAddress(apath))
	}

	c.core.Lock()
	defer c.core.Unlock()
	for address, _ := range c.Adapters {
		if !isStringInArray(address, addresses) {
			delete(c.Adapters, address)
		}
	}
}

func (c *config) addAdapterConfig(address string) {
	if c.isAdapterConfigExists(address) {
		return
	}

	c.core.Lock()
	defer c.core.Unlock()
	c.Adapters[address] = newAdapterConfig()
}

func (c *config) removeAdapterConfig(address string) {
	if !c.isAdapterConfigExists(address) {
		logger.Errorf("config for adapter %s not exists", address)
		return
	}

	c.core.Lock()
	defer c.core.Unlock()
	delete(c.Adapters, address)
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
func (c *config) setDeviceConfigConnected(address string, conneted bool) {
	dc, ok := c.getDeviceConfig(address)
	if !ok {
		return
	}

	c.core.Lock()
	dc.Connected = conneted
	c.core.Unlock()

	c.save()
	return
}
