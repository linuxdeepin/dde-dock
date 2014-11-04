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
	"pkg.linuxdeepin.com/lib/dbus"
	graphic "pkg.linuxdeepin.com/lib/gdkpixbuf"
	"pkg.linuxdeepin.com/lib/utils"
)

// GetDBusInfo implements interface of dbus.DBusObject.
func (theme *Theme) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Grub2",
		ObjectPath: "/com/deepin/daemon/Grub2/Theme",
		Interface:  "com.deepin.daemon.Grub2.Theme",
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
		theme.setPropItemColor(theme.ItemColor)
		theme.customTheme()
	case "SelectedItemColor":
		if theme.SelectedItemColor == oldv.(string) {
			return
		}
		theme.setPropSelectedItemColor(theme.SelectedItemColor)
		theme.customTheme()
	}
}

func (theme *Theme) setPropUpdating(value bool) {
	theme.Updating = value
	dbus.NotifyChange(theme, "Updating")
}

func (theme *Theme) setPropBackground(value string) {
	theme.Background = value
	dbus.NotifyChange(theme, "Background")
}

func (theme *Theme) setPropItemColor(value string) {
	itemColor := value
	if len(itemColor) == 0 {
		// set a default value to avoid empty string
		itemColor = theme.tplJSONData.DarkScheme.ItemColor
	}
	theme.ItemColor = itemColor
	theme.tplJSONData.CurrentScheme.ItemColor = itemColor
	dbus.NotifyChange(theme, "ItemColor")
}

func (theme *Theme) setPropSelectedItemColor(value string) {
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
func (theme *Theme) SetBackgroundSourceFile(imageFile string) (ok bool, err error) {
	imageFile = utils.DecodeURI(imageFile)
	ok = graphic.IsSupportedImage(imageFile)
	if ok {
		go func() {
			theme.doSetBackgroundSourceFile(imageFile)
		}()
	}
	return
}

func (theme *Theme) doSetBackgroundSourceFile(imageFile string) bool {
	theme.updateLock.Lock()
	defer theme.updateLock.Unlock()
	theme.setPropUpdating(true)
	screenWidth, screenHeight := parseCurrentGfxmode()
	grub2extDoSetThemeBackgroundSourceFile(imageFile, screenWidth, screenHeight)
	theme.setPropBackground(theme.bgThumbFile)
	theme.setPropUpdating(false)

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
	theme.setPropItemColor(theme.tplJSONData.CurrentScheme.ItemColor)
	theme.setPropSelectedItemColor(theme.tplJSONData.CurrentScheme.SelectedItemColor)
	theme.customTheme()

	logger.Info("update background sucess")
	return true
}
