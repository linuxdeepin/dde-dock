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
	"errors"
	"sync"

	"pkg.deepin.io/lib/dbus"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/damage"
)

const (
	OpcodeSystemTrayRequestDock uint32 = iota
	OpcodeSystemTrayBeginMessage
	OpcodeSystemTrayCancelMessage
)

// TrayManager为系统托盘的管理器。
type TrayManager struct {
	owner  x.Window // the manager selection owner window
	visual x.VisualID
	icons  map[x.Window]*TrayIcon
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
		icons:  make(map[x.Window]*TrayIcon),
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

	go m.eventHandleLoop()

	dbus.InstallOnSession(m)
	dbus.Emit(m, "Inited")
	return nil
}

func (m *TrayManager) destroy() {

}

func (m *TrayManager) checkValid() {
	for _, id := range m.TrayIcons {
		xid := x.Window(id)
		if isValidWindow(xid) {
			continue
		}

		m.removeIcon(xid)
	}
}

func (m *TrayManager) handleTrayDamage(xid x.Window) {
	icon, ok := m.icons[xid]
	if !ok {
		return
	}
	dbus.Emit(m, "Changed", uint32(xid))
	logger.Debugf("handleTrayDamage %v name: %q", xid, icon.getName())
}

func sendClientMessage(win, dest x.Window, msgType x.Atom, pArray *[5]uint32) error {
	var data x.ClientMessageData
	data.SetData32(pArray)
	event := x.ClientMessageEvent{
		ResponseType: x.ClientMessageEventCode,
		Format:       32,
		Window:       win,
		Type:         msgType,
		Data:         data,
	}
	w := x.NewWriter()
	x.ClientMessageEventWrite(w, &event)
	const evMask = x.EventMaskSubstructureNotify | x.EventMaskSubstructureRedirect
	return x.SendEventChecked(XConn, x.False, dest, evMask, w.Bytes()).Check(XConn)
}

// to notify tray icon applications
func (m *TrayManager) sendClientMsgMANAGER() error {
	screen := XConn.GetDefaultScreen()
	array := [5]uint32{
		x.CurrentTime,
		uint32(XA_NET_SYSTEM_TRAY_S0),
		uint32(m.owner),
	}
	return sendClientMessage(screen.Root, screen.Root, XA_MANAGER, &array)
}

func getSystemTraySelectionOwner() (x.Window, error) {
	reply, err := x.GetSelectionOwner(XConn, XA_NET_SYSTEM_TRAY_S0).Reply(XConn)
	if err != nil {
		return 0, err
	}
	return reply.Owner, nil
}

func createOwnerWindow(visual x.VisualID) (x.Window, error) {
	winId, err := XConn.GenerateID()
	if err != nil {
		return 0, err
	}
	win := x.Window(winId)
	screen := XConn.GetDefaultScreen()
	err = x.CreateWindowChecked(XConn,
		0,
		win,         // window
		screen.Root, // parent
		0, 0, 1, 1, 0,
		x.WindowClassInputOnly,
		visual,
		x.CWEventMask,
		&x.CreateWindowValueList{
			EventMask: x.EventMaskStructureNotify,
		}).Check(XConn)
	if err != nil {
		return 0, err
	}
	return win, nil
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

	err = x.SetSelectionOwnerChecked(
		XConn,
		m.owner,
		XA_NET_SYSTEM_TRAY_S0,
		x.CurrentTime).Check(XConn)
	if err != nil {
		return err
	}

	w := x.NewWriter()
	w.Write4b(uint32(m.visual))
	x.ChangeProperty(XConn,
		x.PropModeReplace,
		m.owner,                   // window
		XA_NET_SYSTEM_TRAY_VISUAL, // property
		x.AtomVisualID,            // type
		32,
		1,
		w.Bytes())

	w = x.NewWriter()
	w.Write4b(0)
	x.ChangeProperty(XConn,
		x.PropModeReplace,
		m.owner, // window
		XA_NET_SYSTEM_TRAY_ORIENTAION, // property
		x.AtomCardinal,                // type
		32,
		1,
		w.Bytes())

	logger.Debug("acquire selection successful")
	return nil
}

func (m *TrayManager) eventHandleLoop() {
	damageExtData := XConn.GetExtensionData(damage.Ext())
	damageFirstEvent := damageExtData.FirstEvent

	for {
		ev := XConn.WaitForEvent()
		switch ev.GetEventCode() {
		case x.ClientMessageEventCode:
			event, _ := x.NewClientMessageEvent(ev)
			if event.Type == XA_NET_SYSTEM_TRAY_OPCODE {
				data32 := event.Data.GetData32()
				opcode := data32[1]
				logger.Debug("system tray opcode", opcode)

				if opcode == OpcodeSystemTrayRequestDock {
					win := x.Window(data32[2])
					logger.Debug("ClientMessageEvent: system tray request dock", win)
					m.addIcon(win)
				}
			}
		case damage.NotifyEventCode + damageFirstEvent:
			event, _ := damage.NewNotifyEvent(ev)
			m.handleTrayDamage(x.Window(event.Drawable))
		case x.MapNotifyEventCode:
			event, _ := x.NewMapNotifyEvent(ev)
			logger.Debug("MapNotifyEvent", event.Window)
		case x.UnmapNotifyEventCode:
			event, _ := x.NewUnmapNotifyEvent(ev)
			logger.Debug("UnmapNotifyEvent", event.Window)
		case x.DestroyNotifyEventCode:
			event, _ := x.NewDestroyNotifyEvent(ev)
			logger.Debug("DestroyNotifyEvent", event.Window)
			m.removeIcon(event.Window)

		default:
			logger.Debug(ev)
		}
	}
}
