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
	"dlib/graphic"
)

// GetDBusInfo implements interface of dbus.DBusObject.
func (theme *Theme) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Grub2",
		"/com/deepin/daemon/Grub2/Theme",
		"com.deepin.daemon.Grub2.Theme",
	}
}

// OnPropertiesChanged implements interface of dbus.DBusObject.
func (theme *Theme) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("%v", err)
		}
	}()
	switch name {
	case "ItemColor":
		if theme.ItemColor == oldv.(string) {
			return
		}
		theme.setItemColor(theme.ItemColor)
	case "SelectedItemColor":
		if theme.SelectedItemColor == oldv.(string) {
			return
		}
		theme.setSelectedItemColor(theme.SelectedItemColor)
	}
}

// SetBackgroundSourceFile setup the background source file, then
// generate the background to fit the screen resolution, support png
// and jpeg image format.
func (theme *Theme) SetBackgroundSourceFile(imageFile string) uint32 {
	updateThemeBackgroundID++
	go func() {
		id := updateThemeBackgroundID
		ok := theme.doSetBackgroundSourceFile(imageFile)
		if theme.BackgroundUpdated != nil {
			theme.BackgroundUpdated(id, ok)
		}
	}()
	return updateThemeBackgroundID
}

func (theme *Theme) doSetBackgroundSourceFile(imageFile string) bool {
	// check image size
	w, h, err := graphic.GetImageSize(imageFile)
	if err != nil {
		return false
	}
	if w < 800 || h < 600 {
		logger.Error("image size is too small") // TODO
		return false
	}

	screenWidth, screenHeight := getPrimaryScreenBestResolution()
	grub2ext.DoSetThemeBackgroundSourceFile(imageFile, screenWidth, screenHeight)
	dbus.NotifyChange(theme, "Background")

	// set item color through background's dominant color
	_, _, v := graphic.GetDominantColorOfImage(theme.bgSrcFile)
	if v < 0.5 {
		// background is dark
		theme.tplJSONData.CurrentScheme = theme.tplJSONData.DarkScheme
		logger.Info("background is dark, use the dark theme scheme")
	} else {
		// background is bright
		theme.tplJSONData.CurrentScheme = theme.tplJSONData.BrightScheme
		logger.Info("background is bright, so use the bright theme scheme")
	}
	theme.ItemColor = theme.tplJSONData.CurrentScheme.ItemColor
	theme.SelectedItemColor = theme.tplJSONData.CurrentScheme.SelectedItemColor
	theme.setItemColor(theme.ItemColor)
	theme.setSelectedItemColor(theme.SelectedItemColor)

	logger.Info("update background sucess")
	return true
}
