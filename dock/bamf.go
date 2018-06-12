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
	x "github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/lib/dbus1"
)

const (
	bamfDBusServiceName   = "org.ayatana.bamf"
	bamfDBusObjPathPrefix = "/org/ayatana/bamf"
	bamfMatcherObjPath    = bamfDBusObjPathPrefix + "/matcher"
	bamfMatcherInterface  = bamfDBusServiceName + ".matcher"
	bamfAppInterface      = bamfDBusServiceName + ".application"
)

func getDesktopFromWindowByBamf(win x.Window) (string, error) {
	bus, err := dbus.SessionBus()
	if err != nil {
		return "", err
	}
	matcher := bus.Object(bamfDBusServiceName, bamfMatcherObjPath)
	var applicationObjPathStr string
	err = matcher.Call(bamfMatcherInterface+".ApplicationForXid", 0,
		uint32(win)).Store(&applicationObjPathStr)
	if err != nil {
		return "", err
	}
	applicationObjPath := dbus.ObjectPath(applicationObjPathStr)
	if !applicationObjPath.IsValid() {
		return "", nil
	}
	application := bus.Object(bamfDBusServiceName, applicationObjPath)
	var desktopFile string
	err = application.Call(bamfAppInterface+".DesktopFile", 0).Store(&desktopFile)
	if err != nil {
		return "", err
	}
	return desktopFile, nil
}
