/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package grub2

// Setup grub2 environment, regenerate configure and theme if need, don't depends on dbus
func (grub *Grub2) Setup(gfxmode string) {
	runWithoutDbus = true

	// do not call grub.readEntries() here, for that
	// "/boot/grub/grub.cfg" may not exists
	grub.readSettings()
	grub.fixSettings()
	grub.fixSettingDistro()

	// setup gfxmode
	if len(gfxmode) == 0 {
		grub.setSettingGfxmode(grub.config.Resolution)
	} else {
		grub.setSettingGfxmode(gfxmode)
	}

	// write settings
	grub.writeSettings()
	// reset NeedUpdate flag for that will run update-grub always
	// after setup grub
	grub.config.NeedUpdate = false
	grub.config.save()

	// setup theme and generate theme background
	grub.SetupTheme(gfxmode)
}

func (grub *Grub2) SetupTheme(gfxmode string) {
	runWithoutDbus = true
	grub.config.loadOrSaveConfig()
	if len(gfxmode) == 0 {
		gfxmode = grub.config.Resolution
	}
	w, h := parseGfxmode(gfxmode)
	err := doGenerateThemeBackground(w, h)
	if err != nil {
		logger.Warning(err)
	}
}
