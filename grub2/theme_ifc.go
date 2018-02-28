/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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
	"regexp"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	graphic "pkg.deepin.io/lib/gdkpixbuf"
	"pkg.deepin.io/lib/utils"
)

func (theme *Theme) GetDBusExportInfo() dbusutil.ExportInfo {
	return dbusutil.ExportInfo{
		Path:      DBusObjPath + "/Theme",
		Interface: DBusInterface + ".Theme",
	}
}

var colorReg = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)

func checkColor(v string) error {
	if colorReg.MatchString(v) {
		return nil
	}
	return fmt.Errorf("invalid color %q", v)
}

func (theme *Theme) SetItemColor(sender dbus.Sender, color string) *dbus.Error {
	theme.service.DelayAutoQuit()

	err := theme.g.checkAuth(sender)
	if err != nil {
		return dbusutil.ToError(err)
	}

	err = checkColor(color)
	if err != nil {
		return dbusutil.ToError(err)
	}

	theme.PropsMu.Lock()
	defer theme.PropsMu.Unlock()

	if theme.setPropItemColor(color) {
		err = theme.setCustomTheme()
		if err != nil {
			return dbusutil.ToError(err)
		}
	}
	return nil
}

func (theme *Theme) SetSelectedItemColor(sender dbus.Sender, color string) *dbus.Error {
	theme.service.DelayAutoQuit()

	err := theme.g.checkAuth(sender)
	if err != nil {
		return dbusutil.ToError(err)
	}

	err = checkColor(color)
	if err != nil {
		return dbusutil.ToError(err)
	}

	theme.PropsMu.Lock()
	defer theme.PropsMu.Unlock()

	if theme.setPropSelectedItemColor(color) {
		err = theme.setCustomTheme()
		if err != nil {
			return dbusutil.ToError(err)
		}
	}
	return nil
}

// SetBackgroundSourceFile setup the background source file, then
// generate the background to fit the screen resolution, support png
// and jpeg image format.
func (theme *Theme) SetBackgroundSourceFile(sender dbus.Sender, filename string) *dbus.Error {
	theme.service.DelayAutoQuit()

	logger.Debugf("SetBackgroundSourceFile: %q", filename)
	err := theme.g.checkAuth(sender)
	if err != nil {
		return dbusutil.ToError(err)
	}

	filename = utils.DecodeURI(filename)
	if graphic.IsSupportedImage(filename) {
		go theme.doSetBackgroundSourceFile(filename)
		return nil
	}
	return dbusutil.ToError(errors.New("unsupported image file"))
}

func (theme *Theme) GetBackground() (string, *dbus.Error) {
	theme.service.DelayAutoQuit()
	return theme.bgFile, nil
}

func (theme *Theme) emitSignalBackgroundChanged() {
	theme.service.Emit(theme, "BackgroundChanged")
}
