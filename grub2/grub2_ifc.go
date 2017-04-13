/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
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

// GetSimpleEntryTitles return entry titles only in level one and will
// filter out some useless entries such as sub-menus and "memtest86+".
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
	// reset and save config
	grub.config.reset()

	// Fix settings always
	grub.setPropFixSettingsAlways(true)
	//grub.config.setFixSettingsAlways(true)

	// enable theme
	grub.setPropEnableTheme(true)
	grub.setEnableTheme(true)

	// resolution
	defaultGfxmode := getDefaultGfxmode()
	grub.setPropResolution(defaultGfxmode)
	grub.setSettingGfxmode(defaultGfxmode)

	// timeout
	grub.setPropTimeout(defaultGrubTimeoutInt)
	grub.setSettingTimeout(defaultGrubTimeoutInt)

	// grub default entry
	grub.setPropDefaultEntry(defaultGrubDefaultEntry)
	grub.setSettingDefaultEntry(defaultGrubDefaultEntry)

	grub.theme.reset()
	grub.notifyUpdate()
}
