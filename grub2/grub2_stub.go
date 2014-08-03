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
	"encoding/json"
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
	. "pkg.linuxdeepin.com/lib/gettext"
	"strings"
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
	case "FixSettingsAlways":
		oldvBool, _ := oldv.(bool)
		if oldvBool == grub.FixSettingsAlways {
			return
		}
		grub.setPropFixSettingsAlways(grub.FixSettingsAlways)
	case "EnableTheme":
		oldvBool, _ := oldv.(bool)
		if oldvBool == grub.EnableTheme {
			return
		}
		grub.setPropEnableTheme(grub.EnableTheme)
		grub.notifyUpdate()
	case "DefaultEntry":
		oldvStr, _ := oldv.(string)
		if oldvStr == grub.DefaultEntry {
			return
		}
		grub.setPropDefaultEntry(grub.DefaultEntry)
		grub.notifyUpdate()
	case "Timeout":
		oldvInt, _ := oldv.(int32)
		if oldvInt == grub.Timeout {
			return
		}
		grub.setPropTimeout(grub.Timeout)
		grub.notifyUpdate()
	case "Resolution":
		oldvStr, _ := oldv.(string)
		if oldvStr == grub.Resolution {
			return
		}
		grub.setPropResolution(grub.Resolution)
		grub.notifyUpdate()
	}
}

func (grub *Grub2) setPropFixSettingsAlways(value bool) {
	grub.FixSettingsAlways = value
	grub.config.setFixSettingsAlways(value)
	dbus.NotifyChange(grub, "FixSettingsAlways")
}

func (grub *Grub2) setPropEnableTheme(value bool) {
	grub.EnableTheme = value
	grub.config.setEnableTheme(value)
	if value {
		grub.setSettingTheme(themeMainFile)
	} else {
		grub.setSettingTheme("")
	}
	dbus.NotifyChange(grub, "EnableTheme")
}

func (grub *Grub2) setPropDefaultEntry(value string) {
	grub.DefaultEntry = value
	grub.setSettingDefaultEntry(value)
	dbus.NotifyChange(grub, "DefaultEntry")
}

func (grub *Grub2) setPropTimeout(value int32) {
	grub.Timeout = value
	grub.setSettingTimeout(value)
	dbus.NotifyChange(grub, "Timeout")
}

func (grub *Grub2) setPropResolution(value string) {
	grub.Resolution = value
	grub.setSettingGfxmode(value)
	dbus.NotifyChange(grub, "Resolution")
}

func (grub *Grub2) setPropUpdating(value bool) {
	grub.Updating = value
	dbus.NotifyChange(grub, "Updating")
}

// GetSimpleEntryTitles return entry titles in level one.
func (grub *Grub2) GetSimpleEntryTitles() ([]string, error) {
	entryTitles := make([]string, 0)
	for _, entry := range grub.entries {
		if entry.parentSubMenu == nil && entry.entryType == MENUENTRY {
			title := entry.getFullTitle()
			if !strings.Contains(title, "memtest86+") {
				entryTitles = append(entryTitles, title)
			}
		}
	}
	if len(entryTitles) == 0 {
		err := fmt.Errorf("there is no menu entry in %s", grubMenuFile)
		return entryTitles, err
	}
	return entryTitles, nil
}

func (grub *Grub2) GetAvailableResolutions() (modesJSON string, err error) {
	type mode struct{ Text, Value string }
	primaryResolution := getPrimaryScreenBestResolutionStr()
	appendModeUniq := func(modes []mode, r string) []mode {
		if r != primaryResolution {
			modes = append(modes, mode{Text: r, Value: r})
		}
		return modes
	}
	var modes []mode
	modes = append(modes, mode{Text: Tr("Auto"), Value: "auto"})
	modes = append(modes, mode{Text: primaryResolution, Value: primaryResolution})
	modes = appendModeUniq(modes, "1440x900")
	modes = appendModeUniq(modes, "1400x1050")
	modes = appendModeUniq(modes, "1280x800")
	modes = appendModeUniq(modes, "1280x720")
	modes = appendModeUniq(modes, "1366x768")
	modes = appendModeUniq(modes, "1024x768")
	modes = appendModeUniq(modes, "800x600")
	tmpByteArray, err := json.Marshal(modes)
	modesJSON = string(tmpByteArray)
	return
}

// Reset reset all configuretion.
func (grub *Grub2) Reset() {
	simpleEntryTitles, _ := grub.GetSimpleEntryTitles()
	firstEntry := ""
	if len(simpleEntryTitles) > 0 {
		firstEntry = simpleEntryTitles[0]
	}
	grub.setPropFixSettingsAlways(true)
	grub.setPropEnableTheme(true)
	grub.setPropResolution(getPrimaryScreenBestResolutionStr())
	grub.setPropDefaultEntry(firstEntry)
	grub.setPropTimeout(int32(10))
	grub.theme.reset()
}
