/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package trayicon

import (
	x "github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/lib/dbus"
)

func (*TrayManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.dde.TrayManager",
		ObjectPath: "/com/deepin/dde/TrayManager",
		Interface:  "com.deepin.dde.TrayManager",
	}
}

// Manage方法获取系统托盘图标的管理权。
func (m *TrayManager) Manage() bool {
	logger.Debug("call Manage by dbus")

	err := m.sendClientMsgMANAGER()
	if err != nil {
		logger.Warning(err)
		return false
	}
	return true
}

// GetName返回传入的系统图标的窗口id的窗口名。
func (m *TrayManager) GetName(win uint32) string {
	icon, ok := m.icons[x.Window(win)]
	if !ok {
		return ""
	}
	return icon.getName()
}

// EnableNotification设置对应id的窗口是否可以通知。
func (m *TrayManager) EnableNotification(win uint32, enable bool) {
	icon, ok := m.icons[x.Window(win)]
	if !ok {
		return
	}
	icon.notify = enable
}
