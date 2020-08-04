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

package trayicon

import (
	"bytes"
	"errors"
	"sync"
	"time"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/composite"
	"github.com/linuxdeepin/go-x11-client/ext/damage"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	OpcodeSystemTrayRequestDock uint32 = iota
	OpcodeSystemTrayBeginMessage
	OpcodeSystemTrayCancelMessage
)

//go:generate dbusutil-gen -type TrayManager,StatusNotifierWatcher traymanager.go status-notifier-watcher.go

// TrayManager为系统托盘的管理器。
type TrayManager struct {
	service *dbusutil.Service
	owner   x.Window // the manager selection owner window
	visual  x.VisualID
	icons   map[x.Window]*TrayIcon
	mutex   sync.Mutex

	damageNotifyEventHandler DamageNotifyEventHandler

	// 目前已有系统托盘窗口的id。
	PropsMu sync.RWMutex
	// dbusutil-gen: equal=nil
	TrayIcons []uint32

	// nolint
	signals *struct {
		// Inited when tray manager is initialized.
		Inited struct{}
		// Added信号会在系统过盘图标增加时被触发。
		// Removed信号会在系统过盘图标被移除时被触发。
		// Changed信号会在系统托盘图标改变后被触发。
		Added, Removed, Changed struct {
			id uint32
		}
	}

	// nolint
	methods *struct {
		Manage             func() `out:"ok"`
		GetName            func() `in:"win" out:"name"`
		EnableNotification func() `in:"win,enabled"`
	}
}

type DamageNotifyEventHandler struct {
	mu           sync.Mutex
	queuedWinIds []x.Window
	timer        *time.Timer
	timerStarted bool
	manager      *TrayManager
}

func (handler *DamageNotifyEventHandler) process(winId x.Window) {
	handler.mu.Lock()
	var found bool
	for _, winId0 := range handler.queuedWinIds {
		if winId0 == winId {
			found = true
		}
	}
	if !found {
		handler.queuedWinIds = append(handler.queuedWinIds, winId)
	}

	if !handler.timerStarted {
		handler.timerStarted = true

		handler.timer = time.AfterFunc(60*time.Millisecond, func() {
			handler.mu.Lock()
			m := handler.manager
			for _, winId := range handler.queuedWinIds {
				m.handleDamageNotifyEvent(winId)
			}
			handler.queuedWinIds = nil
			handler.timerStarted = false
			handler.mu.Unlock()
		})
	}

	handler.mu.Unlock()
}

func NewTrayManager(service *dbusutil.Service) *TrayManager {
	visualId := findRGBAVisualID()

	m := &TrayManager{
		service: service,
		visual:  visualId,
		icons:   make(map[x.Window]*TrayIcon),
	}
	m.damageNotifyEventHandler.manager = m
	err := m.init()
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

	go m.eventHandleLoop()

	return nil
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

func (m *TrayManager) handleDamageNotifyEvent(xid x.Window) {
	m.mutex.Lock()
	icon, ok := m.icons[xid]
	m.mutex.Unlock()
	if !ok {
		return
	}

	icon.mu.Lock()
	if !icon.notify {
		// ignore event
		icon.mu.Unlock()
		return
	}
	icon.mu.Unlock()

	newData, err := icon.getPixmapData()
	if err != nil {
		logger.Warning(err)
		return
	}
	if !bytes.Equal(icon.data, newData) {
		icon.data = newData
		err := m.service.Emit(m, "Changed", uint32(xid))
		if err != nil {
			logger.Warning(err)
		}
		logger.Debugf("handleDamageNotifyEvent %v changed", xid)
	} else {
		logger.Debugf("handleDamageNotifyEvent %v no changed", xid)
	}
}

func sendClientMessage(win, dest x.Window, msgType x.Atom, pArray *[5]uint32) error {
	var data x.ClientMessageData
	data.SetData32(pArray)
	event := x.ClientMessageEvent{
		Format: 32,
		Window: win,
		Type:   msgType,
		Data:   data,
	}
	w := x.NewWriter()
	x.WriteClientMessageEvent(w, &event)
	const evMask = x.EventMaskStructureNotify
	return x.SendEventChecked(XConn, false, dest, evMask, w.Bytes()).Check(XConn)
}

// to notify tray icon applications
func (m *TrayManager) sendClientMsgMANAGER() error {
	screen := XConn.GetDefaultScreen()
	array := [5]uint32{
		x.CurrentTime,
		uint32(XA_NET_SYSTEM_TRAY_S0),
		uint32(m.owner),
	}
	logger.Debug("send clientMsg MANAGER")
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
	winId, err := XConn.AllocID()
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
		[]uint32{x.EventMaskStructureNotify},
	).Check(XConn)
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
		32, w.Bytes())

	w = x.NewWriter()
	w.Write4b(0)
	x.ChangeProperty(XConn,
		x.PropModeReplace,
		m.owner,                       // window
		XA_NET_SYSTEM_TRAY_ORIENTAION, // property
		x.AtomCardinal,                // type
		32,
		w.Bytes())

	logger.Debug("acquire selection successful")
	return nil
}

func (m *TrayManager) eventHandleLoop() {
	damageExtData := XConn.GetExtensionData(damage.Ext())
	damageFirstEvent := damageExtData.FirstEvent

	eventChan := make(chan x.GenericEvent, 500)
	XConn.AddEventChan(eventChan)

	for ev := range eventChan {
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
			m.damageNotifyEventHandler.process(x.Window(event.Drawable))
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

func (m *TrayManager) addIcon(win x.Window) {
	m.checkValid()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.icons[win]
	if ok {
		logger.Debugf("addIcon failed: %v existed", win)
		return
	}
	damageId, err := XConn.AllocID()
	if err != nil {
		logger.Debug("addIcon failed, new damage id failed:", err)
		return
	}
	d := damage.Damage(damageId)

	icon := NewTrayIcon(win)
	icon.damage = d

	err = damage.CreateChecked(XConn, d, x.Drawable(win), damage.ReportLevelRawRectangles).Check(XConn)
	if err != nil {
		logger.Debug("addIcon failed, damage create failed:", err)
		return
	}

	composite.RedirectWindow(XConn, win, composite.RedirectAutomatic)

	const valueMask = x.CWBackPixel | x.CWEventMask
	x.ChangeWindowAttributes(XConn, win, valueMask, []uint32{
		0, // background pixel
		x.EventMaskVisibilityChange | x.EventMaskStructureNotify, // event mask
	})

	err = m.service.Emit(m, "Added", uint32(win))
	if err != nil {
		logger.Warning(err)
	}
	logger.Infof("Add tray icon %v name: %q", win, icon.getName())
	m.icons[win] = icon
	m.updateTrayIcons()
}

func (m *TrayManager) removeIcon(win x.Window) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	icon, ok := m.icons[win]
	if !ok {
		logger.Debugf("removeIcon failed: %v not exist", win)
		return
	}
	// NOTE: no need to destroy the damage, the window of icon is destroyed.

	err := XConn.FreeID(uint32(icon.damage))
	if err != nil {
		logger.Warning(err)
	}

	delete(m.icons, win)
	err = m.service.Emit(m, "Removed", uint32(win))
	if err != nil {
		logger.Warning(err)
	}
	logger.Debugf("remove tray icon %v", win)
	m.updateTrayIcons()
}

func (m *TrayManager) updateTrayIcons() {
	var icons []uint32
	for _, icon := range m.icons {
		icons = append(icons, uint32(icon.win))
	}
	m.PropsMu.Lock()
	m.setPropTrayIcons(icons)
	m.PropsMu.Unlock()
}
