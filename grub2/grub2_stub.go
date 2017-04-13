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
		newv := grub.FixSettingsAlways
		if oldvBool == newv {
			return
		}
		grub.setPropFixSettingsAlways(newv)
		grub.config.setFixSettingsAlways(newv)
	case "EnableTheme":
		oldvBool, _ := oldv.(bool)
		newv := grub.EnableTheme
		if oldvBool == newv {
			return
		}
		grub.setPropEnableTheme(newv)
		grub.config.setEnableTheme(newv)
		grub.setEnableTheme(newv)
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
		newv := grub.Timeout
		if oldvInt == newv {
			return
		}
		grub.setPropTimeout(newv)
		grub.setSettingTimeout(newv)
		grub.notifyUpdate()
	case "Resolution":
		oldvStr, _ := oldv.(string)
		newv := grub.Resolution
		if oldvStr == newv {
			return
		}
		grub.setPropResolution(newv)
		grub.setSettingGfxmode(newv)
		grub.notifyUpdate()
	}
}

func (grub *Grub2) setPropFixSettingsAlways(value bool) {
	grub.FixSettingsAlways = value
	dbus.NotifyChange(grub, "FixSettingsAlways")
}

func (grub *Grub2) setPropEnableTheme(value bool) {
	grub.EnableTheme = value
	dbus.NotifyChange(grub, "EnableTheme")
}

func (grub *Grub2) setPropDefaultEntry(value string) {
	grub.setSettingDefaultEntry(value)
	grub.DefaultEntry = grub.getSettingDefaultEntry()
	dbus.NotifyChange(grub, "DefaultEntry")
}

func (grub *Grub2) setPropTimeout(value int32) {
	grub.Timeout = value
	dbus.NotifyChange(grub, "Timeout")
}

func (grub *Grub2) setPropResolution(value string) {
	grub.Resolution = value
	dbus.NotifyChange(grub, "Resolution")
}

func (grub *Grub2) setPropUpdating(value bool) {
	grub.Updating = value
	dbus.NotifyChange(grub, "Updating")
}
