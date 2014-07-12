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

const configFile = "/var/cache/deepin/grub2.json"

type config struct {
	core       utils.Config
	NeedUpdate bool // mark to generate grub configuration
	// EnableTheme  bool
	// FixSettingsAlways bool
	DefaultEntry string
	Timeout      int32
	Resolution   string
}

func newConfig() (c *config) {
	c = &config{}
	c.NeedUpdate = true
	// TODO
	// c.EnableTheme = true
	c.DefaultEntry = "0"
	c.Timeout = 10
	c.Resolution = getPrimaryScreenBestResolutionStr()
	c.core.SetConfigFile(configFile)
	logger.Info("config file:", c.core.GetConfigFile())
	return
}
func (c *config) load() {
	err := c.core.Load(c)
	if err != nil {
		logger.Error(err)
	}
}
func (c *config) save() {
	fileContent, err := c.core.GetFileContentToSave(c)
	if err != nil {
		logger.Error(err)
	}
	grub2extDoWriteCacheConfig(string(fileContent))
}

func (c *config) setDefaultEntry(value string) {
	c.DefaultEntry = value
	c.save()
}
func (c *config) setTimeout(value int32) {
	c.Timeout = value
	c.save()
}
func (c *config) setResolution(value string) {
	c.Resolution = value
	c.save()
}
