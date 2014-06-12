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
			logger.Errorf("%v", err)
		}
	}()
	switch name {
	case "ItemColor":
		if theme.ItemColor == oldv.(string) {
			return
		}
		theme.updatePropItemColor(theme.ItemColor)
		theme.customTheme()
	case "SelectedItemColor":
		if theme.SelectedItemColor == oldv.(string) {
			return
		}
		theme.updatePropSelectedItemColor(theme.SelectedItemColor)
		theme.customTheme()
	}
}

func (theme *Theme) updatePropUpdating(value bool) {
	theme.Updating = value
	dbus.NotifyChange(theme, "Updating")
}

func (theme *Theme) updatePropBackground(value string) {
	theme.background = value
	// generate background thumbnail
	theme.Background = graphic.GenerateCacheFilePath("grub background" + theme.background)
	graphic.ThumbnailImage(theme.background, theme.Background, 300, 300, graphic.PNG)
	dbus.NotifyChange(theme, "Background")
}

func (theme *Theme) updatePropItemColor(value string) {
	itemColor := value
	if len(itemColor) == 0 {
		// set a default value to avoid empty string
		itemColor = theme.tplJSONData.DarkScheme.ItemColor
	}
	theme.ItemColor = itemColor
	theme.tplJSONData.CurrentScheme.ItemColor = itemColor
	dbus.NotifyChange(theme, "ItemColor")
}

func (theme *Theme) updatePropSelectedItemColor(value string) {
	selectedItemColor := value
	if len(selectedItemColor) == 0 {
		// set a default value to avoid empty string
		selectedItemColor = theme.tplJSONData.DarkScheme.SelectedItemColor
	}
	theme.SelectedItemColor = selectedItemColor
	theme.tplJSONData.CurrentScheme.SelectedItemColor = selectedItemColor
	dbus.NotifyChange(theme, "SelectedItemColor")
}

// SetBackgroundSourceFile setup the background source file, then
// generate the background to fit the screen resolution, support png
// and jpeg image format.
func (theme *Theme) SetBackgroundSourceFile(imageFile string) {
	go func() {
		theme.doSetBackgroundSourceFile(imageFile)
	}()
}

func (theme *Theme) doSetBackgroundSourceFile(imageFile string) bool {
	theme.updateLock.Lock()
	defer theme.updateLock.Unlock()
	theme.updatePropUpdating(true)
	screenWidth, screenHeight := getPrimaryScreenBestResolution()
	grub2extDoSetThemeBackgroundSourceFile(imageFile, screenWidth, screenHeight)
	theme.updatePropBackground(theme.background)
	theme.updatePropUpdating(false)

	// set item color through background's dominant color
	_, _, v, _ := graphic.GetDominantColorOfImage(theme.bgSrcFile)
	if v < 0.5 {
		// background is dark
		theme.tplJSONData.CurrentScheme = theme.tplJSONData.DarkScheme
		logger.Info("background is dark, use the dark theme scheme")
	} else {
		// background is bright
		theme.tplJSONData.CurrentScheme = theme.tplJSONData.BrightScheme
		logger.Info("background is bright, so use the bright theme scheme")
	}
	theme.updatePropItemColor(theme.tplJSONData.CurrentScheme.ItemColor)
	theme.updatePropSelectedItemColor(theme.tplJSONData.CurrentScheme.SelectedItemColor)
	theme.customTheme()

	logger.Info("update background sucess")
	return true
}
