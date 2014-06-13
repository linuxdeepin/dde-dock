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
	"dlib/dbus"
	"fmt"
)

const (
	dbusGrubDest = "com.deepin.daemon.Grub2"
	dbusGrubPath = "/com/deepin/daemon/Grub2"
	dbusGrubIfs  = "com.deepin.daemon.Grub2"
)

// GetDBusInfo implements interface of dbus.DBusObject.
func (grub *Grub2) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		dbusGrubDest,
		dbusGrubPath,
		dbusGrubIfs,
	}
}

// OnPropertiesChanged implements interface of dbus.DBusObject.
func (grub *Grub2) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("%v", err)
		}
	}()
	logger.Debug("OnPropertiesChanged: " + name)
	switch name {
	case "DefaultEntry":
		oldvStr, _ := oldv.(string)
		if grub.DefaultEntry == oldvStr {
			return
		}
		grub.updatePropDefaultEntry(grub.DefaultEntry)
	case "Timeout":
		oldvInt, _ := oldv.(int32)
		if grub.Timeout == oldvInt {
			return
		}
		grub.updatePropTimeout(grub.Timeout)
	}
	grub.writeSettings()
	grub.notifyUpdate()
}

func (grub *Grub2) updatePropDefaultEntry(value string) {
	grub.DefaultEntry = value
	grub.setSettingDefaultEntry(grub.DefaultEntry)
	dbus.NotifyChange(grub, "DefaultEntry")
}

func (grub *Grub2) updatePropTimeout(value int32) {
	grub.Timeout = value
	grub.setSettingTimeout(grub.Timeout)
	dbus.NotifyChange(grub, "Timeout")
}

func (grub *Grub2) updatePropUpdating(value bool) {
	grub.Updating = value
	dbus.NotifyChange(grub, "Updating")
}

// GetSimpleEntryTitles return entry titles in level one.
func (grub *Grub2) GetSimpleEntryTitles() ([]string, error) {
	entryTitles := make([]string, 0)
	for _, entry := range grub.entries {
		if entry.parentSubMenu == nil && entry.entryType == MENUENTRY {
			entryTitles = append(entryTitles, entry.getFullTitle())
		}
	}
	if len(entryTitles) == 0 {
		err := fmt.Errorf("there is no menu entry in %s", grubMenuFile)
		return entryTitles, err
	}
	return entryTitles, nil
}

// Reset reset all configuretion.
func (grub *Grub2) Reset() {
	simpleEntryTitles, _ := grub.GetSimpleEntryTitles()
	firstEntry := ""
	if len(simpleEntryTitles) > 0 {
		firstEntry = simpleEntryTitles[0]
	}
	grub.updatePropDefaultEntry(firstEntry)
	grub.updatePropTimeout(int32(10))
	grub.theme.reset()
}
