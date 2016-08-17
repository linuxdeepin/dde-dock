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
	"github.com/BurntSushi/xgb/damage"
	"github.com/BurntSushi/xgb/xfixes"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"pkg.deepin.io/lib/dbus"
	"sync"
)

const (
	OpCodeSystemTrayRequestDock   uint32 = 0
	OpCodeSystemTrayBeginMessage  uint32 = 1
	OpCodeSystemTrayCancelMessage uint32 = 2
)

// TrayManager为系统托盘的管理器。
type TrayManager struct {
	owner  xproto.Window
	visual xproto.Visualid
	icons  map[xproto.Window]*TrayIcon
	mutex  sync.Mutex

	// 目前已有系统托盘窗口的id。
	TrayIcons []uint32

	// Signals:
	// Removed信号会在系统过盘图标被移除时被触发。
	Removed func(id uint32)
	// Added信号会在系统过盘图标增加时被触发。
	Added func(id uint32)
	// Changed信号会在系统托盘图标改变后被触发。
	Changed func(id uint32)
	// Inited when tray manager is initialized.
	Inited func()
}

func NewTrayManager() *TrayManager {
	visualId := findRGBAVisualID()

	m := &TrayManager{
		owner:  0,
		visual: visualId,
		icons:  make(map[xproto.Window]*TrayIcon),
	}
	m.init()
	return m
}

func (m *TrayManager) init() error {
	m.Manage()
	dbus.InstallOnSession(m)

	xfixes.SelectSelectionInput(
		TrayXU.Conn(),
		TrayXU.RootWin(),
		_NET_SYSTEM_TRAY_S0,
		xfixes.SelectionEventMaskSelectionClientClose,
	)
	go m.startListener()
	dbus.Emit(m, "Inited")
	return nil
}

func (m *TrayManager) destroy() {

}

func (m *TrayManager) checkValid() {
	for _, id := range m.TrayIcons {
		xid := xproto.Window(id)
		if isValidWindow(xid) {
			continue
		}

		m.removeIcon(xid)
	}
}

func (m *TrayManager) handleTrayDamage(xid xproto.Window) {
	icon, ok := m.icons[xid]
	if !ok {
		return
	}
	md5 := icon2md5(xid)
	if !md5Equal(icon.md5, md5) {
		icon.md5 = md5
		dbus.Emit(m, "Changed", uint32(xid))
		logger.Debugf("handleTrayDamage %v name: %q changed %v", xid, icon.getName(), md5)
	}
}

func (m *TrayManager) destroyOwnerWindow() {
	if m.owner != 0 {
		xproto.DestroyWindow(TrayXU.Conn(), m.owner)
	}
	m.owner = 0
}

func (m *TrayManager) requireManageTrayIcons() {
	mstype, err := xprop.Atm(TrayXU, "MANAGER")
	if err != nil {
		logger.Warning("Get MANAGER Failed")
		return
	}

	timeStamp, _ := ewmh.WmUserTimeGet(TrayXU, m.owner)
	cm, err := xevent.NewClientMessage(
		32,
		TrayXU.RootWin(),
		mstype,
		int(timeStamp),
		int(_NET_SYSTEM_TRAY_S0),
		int(m.owner),
	)

	if err != nil {
		logger.Warning("Send MANAGER Request failed:", err)
		return
	}

	// !!! ewmh.ClientEvent not use EventMaskStructureNotify.
	xevent.SendRootEvent(TrayXU, cm,
		uint32(xproto.EventMaskStructureNotify))
}

func (m *TrayManager) getSelectionOwner() (*xproto.GetSelectionOwnerReply, error) {
	_trayInstance := xproto.GetSelectionOwner(TrayXU.Conn(), _NET_SYSTEM_TRAY_S0)
	return _trayInstance.Reply()
}

func (m *TrayManager) tryOwner() bool {
	// Make a check, the tray application MUST be 1.
	reply, err := m.getSelectionOwner()
	if err != nil {
		logger.Error(err)
		return false
	}
	if reply.Owner != 0 {
		logger.Warning("Another System tray application is running")
		return false
	}

	timeStamp, _ := ewmh.WmUserTimeGet(TrayXU, m.owner)
	err = xproto.SetSelectionOwnerChecked(
		TrayXU.Conn(),
		m.owner,
		_NET_SYSTEM_TRAY_S0,
		xproto.Timestamp(timeStamp),
	).Check()
	if err != nil {
		logger.Warning("Set Selection Owner failed: ", err)
		return false
	}

	//owner the _NET_SYSTEM_TRAY_Sn
	logger.Debug("Required _NET_SYSTEM_TRAY_S0 successful")

	m.requireManageTrayIcons()

	xprop.ChangeProp32(
		TrayXU,
		m.owner,
		"_NET_SYSTEM_TRAY_VISUAL",
		"VISUALID",
		uint(m.visual),
	)
	xprop.ChangeProp32(
		TrayXU,
		m.owner,
		"_NET_SYSTEM_TRAY_ORIENTAION",
		"CARDINAL",
		0,
	)
	reply, err = m.getSelectionOwner()
	if err != nil {
		logger.Warning("get selection owner failed:", err)
		return false
	}
	return reply.Owner != 0
}

// TODO
var isListened bool = false

func (m *TrayManager) startListener() {
	// to avoid creating too much listener when SelectionNotifyEvent occurs.
	if isListened {
		return
	}
	isListened = true

	for {
		e, err := TrayXU.Conn().WaitForEvent()
		if err != nil || e == nil {
			logger.Warning("WaitForEvent error:", err)
			m.checkValid()
			continue
		}

		switch ev := e.(type) {
		case xproto.ClientMessageEvent:
			// logger.Info("ClientMessageEvent")
			if ev.Type == _NET_SYSTEM_TRAY_OPCODE {
				// timeStamp = ev.Data.Data32[0]
				opCode := ev.Data.Data32[1]
				// logger.Info("TRAY_OPCODE")

				switch opCode {
				case OpCodeSystemTrayRequestDock:
					logger.Debug("System tray request dock")
					xid := xproto.Window(ev.Data.Data32[2])
					m.addIcon(xid)
				case OpCodeSystemTrayBeginMessage:
				case OpCodeSystemTrayCancelMessage:
				}
			}
		case damage.NotifyEvent:
			m.handleTrayDamage(xproto.Window(ev.Drawable))
		case xproto.DestroyNotifyEvent:
			logger.Debug("DestroyNotifyEvent", ev.Window)
			m.removeIcon(ev.Window)
		case xproto.SelectionClearEvent:
			logger.Debug("SelectionClearEvent")
			m.Unmanage()
		case xfixes.SelectionNotifyEvent:
			logger.Debug("SelectionNotifyEvent")
			m.Manage()
		case xproto.UnmapNotifyEvent:
			logger.Debug("UnmapNotifyEvent", ev.Window)
			m.removeIcon(ev.Window)
		case xproto.MapNotifyEvent:
			logger.Debug("MapNotifyEvent", ev.Window)
			m.addIcon(ev.Window)
		}
	}
}
