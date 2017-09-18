/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package grub2

import (
	"errors"
	"fmt"
	"pkg.deepin.io/lib/dbus"
	graphic "pkg.deepin.io/lib/gdkpixbuf"
	"pkg.deepin.io/lib/utils"
	"regexp"
)

// GetDBusInfo implements interface of dbus.DBusObject.
func (theme *Theme) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       DBusDest,
		ObjectPath: DBusObjPath + "/Theme",
		Interface:  DBusInterface + ".Theme",
	}
}

var colorReg = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)

func checkColor(v string) error {
	if colorReg.MatchString(v) {
		return nil
	}
	return fmt.Errorf("invalid color %q", v)
}

func (theme *Theme) SetItemColor(dbusMsg dbus.DMessage, v string) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	err = checkColor(v)
	if err != nil {
		return
	}

	if theme.ItemColor == v {
		return
	}
	theme.setPropItemColor(v)
	theme.setCustomTheme()
	return
}

func (theme *Theme) SetSelectedItemColor(dbusMsg dbus.DMessage, v string) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	err = checkColor(v)
	if err != nil {
		return
	}

	if theme.SelectedItemColor == v {
		return
	}
	theme.setPropSelectedItemColor(v)
	theme.setCustomTheme()
	return
}

// SetBackgroundSourceFile setup the background source file, then
// generate the background to fit the screen resolution, support png
// and jpeg image format.
func (theme *Theme) SetBackgroundSourceFile(dbusMsg dbus.DMessage, imageFile string) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	imageFile = utils.DecodeURI(imageFile)
	if graphic.IsSupportedImage(imageFile) {
		go theme.doSetBackgroundSourceFile(imageFile)
		return
	}
	return errors.New("unsupported image file")
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
