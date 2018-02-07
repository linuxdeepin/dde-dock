/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package dock

import (
	"github.com/BurntSushi/xgb/xproto"
	"pkg.deepin.io/lib/dbus"
)

const (
	bamfDBusDest          = "org.ayatana.bamf"
	bamfDBusObjPathPrefix = "/org/ayatana/bamf"
	bamfMatcherObjPath    = bamfDBusObjPathPrefix + "/matcher"
	bamfMatcherIfc        = bamfDBusDest + ".matcher"
	bamfAppIfc            = bamfDBusDest + ".application"
)

func _getDestkopFromWindowByBamf(xid uint32) (string, error) {
	bus, err := dbus.SessionBus()
	if err != nil {
		return "", err
	}
	matcher := bus.Object(bamfDBusDest, bamfMatcherObjPath)
	var applicationObjPathStr string
	err = matcher.Call(bamfMatcherIfc+".ApplicationForXid", 0, xid).Store(&applicationObjPathStr)
	if err != nil {
		return "", err
	}
	applicationObjPath := dbus.ObjectPath(applicationObjPathStr)
	if !applicationObjPath.IsValid() {
		return "", nil
	}
	application := bus.Object(bamfDBusDest, applicationObjPath)
	var desktopFile string
	err = application.Call(bamfAppIfc+".DesktopFile", 0).Store(&desktopFile)
	if err != nil {
		return "", err
	}
	return desktopFile, nil
}

func getDesktopFromWindowByBamf(win xproto.Window) string {
	desktopFile, err := _getDestkopFromWindowByBamf(uint32(win))
	if err != nil {
		logger.Warning(err)
		return ""
	}
	return desktopFile
}
