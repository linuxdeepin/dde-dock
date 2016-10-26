/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
