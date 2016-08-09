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
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xwindow"
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
	logger.Info("TrayManager Manage")
	m.destroyOwnerWindow()

	win, _ := xwindow.Generate(TrayXU)
	m.owner = win.Id
	logger.Debug("owner window id", m.owner)

	xproto.CreateWindowChecked(TrayXU.Conn(), 0, m.owner, TrayXU.RootWin(), 0, 0, 1, 1, 0, xproto.WindowClassInputOnly, m.visual, 0, nil)
	TrayXU.Sync()
	win.Listen(xproto.EventMaskStructureNotify)
	return m.tryOwner()
}

// Unmanage移除系统托盘图标的管理权限。
func (m *TrayManager) Unmanage() bool {
	logger.Info("TrayManager Unmanage")
	reply, err := m.getSelectionOwner()
	if err != nil {
		logger.Warning("get selection owner failed:", err)
		return false
	}
	if reply.Owner != m.owner {
		logger.Warning("not selection owner")
		return false
	}

	m.destroyOwnerWindow()
	trayicons := m.TrayIcons
	logger.Debug("m.TrayIcons:", trayicons)
	for _, icon := range trayicons {
		m.removeIcon(xproto.Window(icon))
	}
	logger.Debug("removeIcon done")

	timeStamp, _ := ewmh.WmUserTimeGet(TrayXU, m.owner)
	return xproto.SetSelectionOwnerChecked(
		TrayXU.Conn(),
		0,
		_NET_SYSTEM_TRAY_S0,
		xproto.Timestamp(timeStamp),
	).Check() == nil
}

// RetryManager方法尝试获取系统托盘图标的权利全，并出发Added信号。
func (m *TrayManager) RetryManager() {
	m.Unmanage()
	m.Manage()

	logger.Debug("emit Added signal for m.TrayIcons", m.TrayIcons)
	for _, icon := range m.TrayIcons {
		dbus.Emit(m, "Added", icon)
	}
}

// GetName返回传入的系统图标的窗口id的窗口名。
func (m *TrayManager) GetName(xid uint32) string {
	icon, ok := m.icons[xproto.Window(xid)]
	if !ok {
		return ""
	}
	return icon.getName()
}

// EnableNotification设置对应id的窗口是否可以通知。
func (m *TrayManager) EnableNotification(xid uint32, enable bool) {
	icon, ok := m.icons[xproto.Window(xid)]
	if !ok {
		return
	}
	icon.notify = enable
}
