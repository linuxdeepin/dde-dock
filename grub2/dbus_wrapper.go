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
	dbusGrub2ext "dbus/com/deepin/daemon/grub2ext"
)

// dbus api wrapper for grub2ext

func newDbusGrub2Ext() (grub2ext *dbusGrub2ext.Grub2Ext, err error) {
	grub2ext, err = dbusGrub2ext.NewGrub2Ext("com.deepin.daemon.Grub2Ext", "/com/deepin/daemon/Grub2Ext")
	if err != nil {
		logger.Error(err)
	}
	return
}

func grub2extDoWriteThemeMainFile(themeFileContent string) {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoWriteThemeMainFile(themeFileContent)
}

func grub2extDoGenerateGrubMenu() {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoGenerateGrubMenu()
}

func grub2extDoGenerateThemeBackground(screenWidth, screenHeight uint16) {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoGenerateThemeBackground(screenWidth, screenHeight)
}

func grub2extDoResetThemeBackground() {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoResetThemeBackground()
}

func grub2extDoSetThemeBackgroundSourceFile(imageFile string, screenWidth, screenHeight uint16) {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoSetThemeBackgroundSourceFile(imageFile, screenWidth, screenHeight)
}

func grub2extDoWriteConfig(fileContent string) {
	logger.Debug("grub2extDoWriteConfig")
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoWriteConfig(fileContent)
}

func grub2extDoWriteGrubSettings(fileContent string) {
	logger.Debug("grub2extDoWriteGrubSettings")
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoWriteGrubSettings(fileContent)
}

func grub2extDoWriteThemeTplFile(jsonContent string) {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoWriteThemeTplFile(jsonContent)
}
