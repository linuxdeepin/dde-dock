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
	"os"
	"os/exec"
	"strings"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/utils"
)

const (
	themeDBusPath      = dbusPath + "/Theme"
	themeDBusInterface = dbusInterface + ".Theme"
)

func (*Theme) GetInterfaceName() string {
	return themeDBusInterface
}

func (theme *Theme) SetBackgroundSourceFile(sender dbus.Sender, filename string) *dbus.Error {
	theme.service.DelayAutoQuit()

	logger.Debugf("SetBackgroundSourceFile: %q", filename)
	err := theme.g.checkAuth(sender, polikitActionIdCommon)
	if err != nil {
		return dbusutil.ToError(err)
	}

	filename = utils.DecodeURI(filename)
	cmd := exec.Command(adjustThemeCmd, "-set-background", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}
	theme.emitSignalBackgroundChanged()
	return nil
}

func (theme *Theme) GetBackground() (string, *dbus.Error) {
	theme.service.DelayAutoQuit()

	theme.g.PropsMu.RLock()
	themeFile := theme.g.ThemeFile
	theme.g.PropsMu.RUnlock()

	if strings.Contains(themeFile, "/deepin/") {
		_, err := os.Stat(themeBgFile)
		if err != nil {
			if os.IsNotExist(err) {
				return fallbackThemeBgFile, nil
			}
			return "", dbusutil.ToError(err)
		} else {
			return themeBgFile, nil
		}
	} else if strings.Contains(themeFile, "/deepin-fallback/") {
		return fallbackThemeBgFile, nil
	} else {
		return "", nil
	}
}

func (theme *Theme) emitSignalBackgroundChanged() {
	err := theme.service.Emit(theme, "BackgroundChanged")
	if err != nil {
		logger.Warning(err)
	}
}
