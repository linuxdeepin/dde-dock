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

package main

import (
	"dlib/dbus"
	"errors"
	"fmt"
)

func (grub *Grub2) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Grub2",
		"/com/deepin/daemon/Grub2",
		"com.deepin.daemon.Grub2",
	}
}

func (grub *Grub2) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			_LOGGER.Error("%v", err)
		}
	}()
	_LOGGER.Info("OnPropertiesChanged: " + name)
	switch name {
	case "DefaultEntry":
		grub.setDefaultEntry(grub.DefaultEntry)
	case "Timeout":
		grub.setTimeout(grub.Timeout)
	}

	grub.writeSettings()
	grub.notifyUpdate()
}

// Get entry titles in level one.
func (grub *Grub2) GetSimpleEntryTitles() ([]string, error) {
	entryTitles := make([]string, 0)
	for _, entry := range grub.entries {
		if entry.parentSubMenu == nil && entry.entryType == MENUENTRY {
			entryTitles = append(entryTitles, entry.getFullTitle())
		}
	}
	if len(entryTitles) == 0 {
		s := fmt.Sprintf("there is no menu entry in %s", _GRUB_MENU_FILE)
		_LOGGER.Error(s)
		return entryTitles, errors.New(s)
	}
	return entryTitles, nil
}
