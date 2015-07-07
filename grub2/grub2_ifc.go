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
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"strings"
)

const (
	DbusGrubDest = "com.deepin.daemon.Grub2"
	DbusGrubPath = "/com/deepin/daemon/Grub2"
	DbusGrubIfs  = "com.deepin.daemon.Grub2"
)

// GetDBusInfo implements interface of dbus.DBusObject.
func (grub *Grub2) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       DbusGrubDest,
		ObjectPath: DbusGrubPath,
		Interface:  DbusGrubIfs,
	}
}

// GetSimpleEntryTitles return entry titles only in level one.
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
	defaultEntry := defaultGrubDefaultEntry
	if len(simpleEntryTitles) > 0 {
		defaultEntry = simpleEntryTitles[0]
	}
	grub.setPropFixSettingsAlways(true)
	grub.setPropEnableTheme(true)
	grub.setPropResolution(getDefaultGfxmode())
	grub.setPropDefaultEntry(defaultEntry)
	grub.setPropTimeout(defaultGrubTimeoutInt)
	grub.theme.reset()
	grub.notifyUpdate()
}
