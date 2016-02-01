/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package bluetooth

import (
	"pkg.deepin.io/lib/utils"
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
	logger.Info("load bluetooth config file:", c.core.GetConfigFile())
	c.Adapters = make(map[string]*adapterConfig)
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
