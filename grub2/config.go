/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
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

package grub2

import (
	"pkg.linuxdeepin.com/lib/utils"
)

const ConfigFileDefault = "/var/cache/deepin/grub2.json"

var configFile = ConfigFileDefault

// refactor code, split fields, grub, efi
type config struct {
	core              utils.Config
	NeedUpdate        bool // mark to generate grub configuration
	FixSettingsAlways bool
	EnableTheme       bool
	DefaultEntry      string
	Timeout           string
	Resolution        string
}

func newConfig() (c *config) {
	c = &config{}
	c.NeedUpdate = true
	c.FixSettingsAlways = true
	c.EnableTheme = true
	c.DefaultEntry = defaultGrubDefaultEntry
	c.Timeout = defaultGrubTimeout
	c.Resolution = getDefaultGfxmode()
	c.core.SetConfigFile(configFile)
	return
}

func (c *config) loadOrSaveConfig() {
	// do not merge this function to load() for permission issue
	if c.core.IsConfigFileExists() {
		c.load()
	} else {
		c.save()
	}
}

func (c *config) load() {
	logger.Info("config file:", c.core.GetConfigFile())
	err := c.core.Load(c)
	if err != nil {
		logger.Error(err)
	}
}
func (c *config) save() {
	if runWithoutDbus {
		c.doSaveWithoutDbus()
	} else {
		c.doSave()
	}
}
func (c *config) doSave() {
	fileContent, err := c.core.GetFileContentToSave(c)
	if err != nil {
		logger.Error(err)
	}
	grub2extDoWriteConfig(string(fileContent))
}
func (c *config) doSaveWithoutDbus() {
	fileContent, err := c.core.GetFileContentToSave(c)
	if err != nil {
		logger.Error(err)
	}
	ge := NewGrub2Ext()
	ge.DoWriteConfig(string(fileContent))
}

func (c *config) setFixSettingsAlways(value bool) {
	c.FixSettingsAlways = value
	c.save()
}
func (c *config) setEnableTheme(value bool) {
	c.EnableTheme = value
	c.save()
}

func (c *config) doSetDefaultEntry(value string) {
	c.DefaultEntry = value
}
func (c *config) doSetTimeout(value string) {
	c.Timeout = value
}
func (c *config) doSetResolution(value string) {
	c.Resolution = value
}
