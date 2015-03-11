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
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/utils"
)

type config struct {
	core utils.Config

	Adapters map[dbus.ObjectPath]*adapterConfig // use adapter dbus path as key
}
type adapterConfig struct {
	Powered bool
}

func newConfig() (c *config) {
	c = &config{}
	c.core.SetConfigName("bluetooth")
	logger.Info("config file:", c.core.GetConfigFile())
	c.Adapters = make(map[dbus.ObjectPath]*adapterConfig)
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
	apathes := bluezGetAdapters()
	for apath, _ := range c.Adapters {
		if !isDBusPathInArray(apath, apathes) {
			delete(c.Adapters, apath)
		}
	}
}

func (c *config) addAdapterConfig(apath dbus.ObjectPath) {
	if c.isAdapterConfigExists(apath) {
		return
	}

	c.core.Lock()
	defer c.core.Unlock()
	c.Adapters[apath] = newAdapterConfig()
}
func (c *config) removeAdapterConfig(apath dbus.ObjectPath) {
	if !c.isAdapterConfigExists(apath) {
		logger.Errorf("config for adapter %s not exists", apath)
		return
	}

	c.core.Lock()
	defer c.core.Unlock()
	delete(c.Adapters, apath)
}
func (c *config) getAdapterConfig(apath dbus.ObjectPath) (ac *adapterConfig, ok bool) {
	c.core.Lock()
	defer c.core.Unlock()
	ac, ok = c.Adapters[apath]
	return
}
func (c *config) isAdapterConfigExists(apath dbus.ObjectPath) (ok bool) {
	c.core.Lock()
	defer c.core.Unlock()
	_, ok = c.Adapters[apath]
	return

}
func (c *config) getAdapterConfigPowered(apath dbus.ObjectPath) (powered bool) {
	ac, ok := c.getAdapterConfig(apath)
	if !ok {
		return
	}

	c.core.Lock()
	powered = ac.Powered
	c.core.Unlock()
	return
}
func (c *config) setAdapterConfigPowered(apath dbus.ObjectPath, powered bool) {
	ac, ok := c.getAdapterConfig(apath)
	if !ok {
		return
	}

	c.core.Lock()
	ac.Powered = powered
	c.core.Unlock()

	c.save()
	return
}
