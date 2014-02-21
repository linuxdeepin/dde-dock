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
	"dlib/graph"
)

func (theme *Theme) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Grub2",
		"/com/deepin/daemon/Grub2/Theme",
		"com.deepin.daemon.Grub2.Theme",
	}
}

func (theme *Theme) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logError("%v", err)
		}
	}()
	switch name {
	case "ItemColor":
		theme.setItemColor(theme.ItemColor)
	case "SelectedItemColor":
		theme.setSelectedItemColor(theme.SelectedItemColor)
	}
}

// Set the background source file, then generate the background
// to fit the screen resolution, support png and jpeg image format
func (theme *Theme) SetBackgroundSourceFile(imageFile string) uint32 {
	_UPDATE_THEME_BACKGROUND_ID++
	go func() {
		id := _UPDATE_THEME_BACKGROUND_ID
		ok := theme.doSetBackgroundSourceFile(imageFile)
		if theme.BackgroundUpdated != nil {
			theme.BackgroundUpdated(id, ok)
		}
	}()
	return _UPDATE_THEME_BACKGROUND_ID
}

func (theme *Theme) doSetBackgroundSourceFile(imageFile string) bool {
	// check image size
	w, h, err := graph.GetImageSize(imageFile)
	if err != nil {
		return false
	}
	if w < 800 || h < 600 {
		logError("image size is too small") // TODO
		return false
	}

	// backup background source file
	_, err = copyFile(imageFile, theme.bgSrcFile)
	if err != nil {
		return false
	}

	theme.generateBackground()

	// set item color through background's dominant color
	_, _, v := graph.GetDominantColorOfImage(theme.bgSrcFile)
	if v < 0.4 {
		// background is dark
		theme.tplJsonData.CurrentScheme = theme.tplJsonData.BrightScheme
	} else {
		// background is bright
		theme.tplJsonData.CurrentScheme = theme.tplJsonData.DarkScheme
	}
	theme.ItemColor = theme.tplJsonData.CurrentScheme.ItemColor
	theme.SelectedItemColor = theme.tplJsonData.CurrentScheme.SelectedItemColor
	dbus.NotifyChange(theme, "ItemColor")
	dbus.NotifyChange(theme, "SelectedItemColor")
	theme.customTheme()

	logInfo("update background sucess")
	return true
}
