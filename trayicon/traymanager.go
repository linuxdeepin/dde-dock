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
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/damage"
	"github.com/BurntSushi/xgb/xfixes"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"

	"errors"
	"github.com/BurntSushi/xgbutil/xwindow"
	"pkg.deepin.io/lib/dbus"
	"sync"
)

const (
	OpcodeSystemTrayRequestDock uint32 = iota
	OpcodeSystemTrayBeginMessage
	OpcodeSystemTrayCancelMessage
)

// TrayManager为系统托盘的管理器。
type TrayManager struct {
	owner  xproto.Window // the manager selection owner window
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
		visual: visualId,
		icons:  make(map[xproto.Window]*TrayIcon),
	}
	err := m.init()
	// TODO
	if err != nil {
		logger.Warning(err)
	}
	return m
}

func (m *TrayManager) init() error {
	win, err := createOwnerWindow(m.visual)
	if err != nil {
		return err
	}
	logger.Debug("create owner window", win)
	m.owner = win
	err = m.acquireSystemTraySelection()
	if err != nil {
		return err
	}

	m.sendClientMsgMANAGER()

	xfixes.SelectSelectionInput(
		XU.Conn(),
		XU.RootWin(),
		_NET_SYSTEM_TRAY_S0,
		xfixes.SelectionEventMaskSelectionClientClose,
	)
	go m.eventHandleLoop()

	dbus.InstallOnSession(m)
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

// to notify tray icon applications
func (m *TrayManager) sendClientMsgMANAGER() error {
	cm, err := xevent.NewClientMessage(
		32, // Format
		XU.RootWin(),
		ATOM_MANAGER,             // message type
		xproto.TimeCurrentTime,   // data[0]
		int(_NET_SYSTEM_TRAY_S0), // data[1]
		int(m.owner),             // data[2]
	)

	if err != nil {
		panic(err)
	}

	logger.Debug("send clientMsg MANAGER")
	// !!! ewmh.ClientEvent not use EventMaskStructureNotify.
	return xevent.SendRootEvent(XU, cm,
		uint32(xproto.EventMaskStructureNotify))
}

func getSystemTraySelectionOwner() (xproto.Window, error) {
	reply, err := xproto.GetSelectionOwner(XU.Conn(), _NET_SYSTEM_TRAY_S0).Reply()
	if err != nil {
		return 0, err
	}
	return reply.Owner, nil
}

func createOwnerWindow(visual xproto.Visualid) (xproto.Window, error) {
	win, _ := xwindow.Generate(XU)

	err := xproto.CreateWindowChecked(XU.Conn(),
		0, win.Id, XU.RootWin(), 0, 0, 1, 1, 0, xproto.WindowClassInputOnly, visual, 0, nil).Check()
	if err != nil {
		return 0, err
	}

	win.Listen(xproto.EventMaskStructureNotify)
	return win.Id, nil
}

func (m *TrayManager) acquireSystemTraySelection() error {
	currentOwner, err := getSystemTraySelectionOwner()
	if err != nil {
		return err
	}
	logger.Debug("currentOwner is ", currentOwner)
	if currentOwner != 0 && currentOwner != m.owner {
		return errors.New("Another System tray application is running")
	}

	err = xproto.SetSelectionOwnerChecked(
		XU.Conn(),
		m.owner,
		_NET_SYSTEM_TRAY_S0,
		xproto.TimeCurrentTime,
	).Check()
	if err != nil {
		return err
	}

	xprop.ChangeProp32(
		XU,
		m.owner,
		"_NET_SYSTEM_TRAY_VISUAL",
		"VISUALID",
		uint(m.visual),
	)
	xprop.ChangeProp32(
		XU,
		m.owner,
		"_NET_SYSTEM_TRAY_ORIENTAION",
		"CARDINAL",
		0,
	)

	logger.Debug("acquire selection successful")
	return nil
}

func (m *TrayManager) eventHandleLoop() {
	for {
		ev, err := XU.Conn().WaitForEvent()
		if ev == nil && err == nil {
			logger.Warning("Both event and error are nil. Exiting...")
			return
		}

		if err != nil {
			logger.Warning(err)
		}
		if ev != nil {
			m.handleXEvent(ev)
		}

	}
}

func (m *TrayManager) handleXEvent(e xgb.Event) {
	switch ev := e.(type) {
	case xproto.ClientMessageEvent:
		if ev.Type == _NET_SYSTEM_TRAY_OPCODE {
			opcode := ev.Data.Data32[1]
			logger.Debug("system tray opcode", opcode)

			if opcode == OpcodeSystemTrayRequestDock {
				win := xproto.Window(ev.Data.Data32[2])
				logger.Debug("ClientMessageEvent: System tray request dock", win)
				m.addIcon(win)
			}
		}
	//case xproto.SelectionClearEvent:
	//logger.Debug("SelectionClearEvent")
	//m.Unmanage()
	case xfixes.SelectionNotifyEvent:
		logger.Debug("SelectionNotifyEvent")
		m.Manage()

		// tray icon events:
	case damage.NotifyEvent:
		m.handleTrayDamage(xproto.Window(ev.Drawable))

	case xproto.MapNotifyEvent:
		logger.Debug("MapNotifyEvent", ev.Window)
		//m.addIcon(ev.Window)
	case xproto.UnmapNotifyEvent:
		logger.Debug("UnmapNotifyEvent", ev.Window)
		//m.removeIcon(ev.Window)
	case xproto.DestroyNotifyEvent:
		logger.Debug("DestroyNotifyEvent", ev.Window)
		m.removeIcon(ev.Window)
	}
}
