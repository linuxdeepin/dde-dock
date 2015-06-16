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

package bluetooth

import (
	"pkg.linuxdeepin.com/lib/utils"
)

type config struct {
	core utils.Config

	Adapters map[string]*adapterConfig // use adapter hardware address as key
}
type adapterConfig struct {
	Powered bool
}

func newConfig() (c *config) {
	c = &config{}
	c.core.SetConfigName("bluetooth")
	logger.Info("config file:", c.core.GetConfigFile())
	c.Adapters = make(map[string]*adapterConfig)
	c.load()
	c.clearSpareConfig()
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
	c.core.Lock()
	defer c.core.Unlock()
	var addresses []string
	apathes := bluezGetAdapters()
	for _, apath := range apathes {
		addresses = append(addresses, bluezGetAdapterAddress(apath))
	}
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
	ac, ok := c.getAdapterConfig(address)
	if !ok {
		return
	}

	c.core.Lock()
	powered = ac.Powered
	c.core.Unlock()
	return
}
func (c *config) setAdapterConfigPowered(address string, powered bool) {
	ac, ok := c.getAdapterConfig(address)
	if !ok {
		return
	}

	c.core.Lock()
	ac.Powered = powered
	c.core.Unlock()

	c.save()
	return
}
