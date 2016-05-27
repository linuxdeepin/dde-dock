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
	"pkg.deepin.io/lib/dbus"
)

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
	grub.setSettingDefaultEntry(value)
	grub.DefaultEntry = grub.getSettingDefaultEntry()
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
