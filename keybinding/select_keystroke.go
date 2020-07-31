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

package keybinding

import (
	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/keybind"
	"github.com/linuxdeepin/go-x11-client/util/mousebind"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
)

func (m *Manager) selectKeystroke() error {
	conn, err := x.NewConn()
	if err != nil {
		return err
	}
	err = grabKbdAndMouse(conn)
	if err != nil {
		logger.Warning("failed to grab keyboard and mouse:", err)
		return err
	}

	// Temporarily disable record
	m.shortcutManager.EnableRecord(false)
	defer m.shortcutManager.EnableRecord(true)

	eventChan := make(chan x.GenericEvent, 500)
	conn.AddEventChan(eventChan)

loop:
	for event := range eventChan {
		switch event.GetEventCode() {
		case x.KeyPressEventCode:
			event, _ := x.NewKeyPressEvent(event)
			logger.Debug(event)
			mods := shortcuts.GetConcernedModifiers(event.State)
			logger.Debug("event mods:", shortcuts.Modifiers(event.State))
			key := shortcuts.Key{
				Mods: mods,
				Code: shortcuts.Keycode(event.Detail),
			}
			logger.Debug("event key:", key)
			ks := key.ToKeystroke(m.keySymbols)
			m.emitSignalKeyEvent(true, ks.String())
			if ks.IsGood() {
				logger.Debug("good keystroke", ks)
				m.grabScreenKeystroke = ks
			} else {
				logger.Debug("bad keystroke", ks)
				m.grabScreenKeystroke = nil
			}
		case x.KeyReleaseEventCode:
			event, _ := x.NewKeyReleaseEvent(event)
			logger.Debug(event)
			if m.grabScreenKeystroke != nil {
				m.emitSignalKeyEvent(false, m.grabScreenKeystroke.String())
				m.grabScreenKeystroke = nil
			} else {
				m.emitSignalKeyEvent(false, "")
			}

			break loop
		case x.ButtonPressEventCode:
			m.emitSignalKeyEvent(true, "")
		case x.ButtonReleaseEventCode:
			m.emitSignalKeyEvent(false, "")
			break loop
		}
	}

	ungrabKbdAndMouse(conn)
	conn.Close()
	logger.Debug("end selectKeystroke")
	return nil
}

func (m *Manager) emitSignalKeyEvent(pressed bool, keystroke string) {
	err := m.service.Emit(m, "KeyEvent", pressed, keystroke)
	if err != nil {
		logger.Warning("emit KeyEvent Failed:", err)
	}
}

func grabKbdAndMouse(conn *x.Conn) error {
	rootWin := conn.GetDefaultScreen().Root
	err := keybind.GrabKeyboard(conn, rootWin)
	if err != nil {
		return err
	}

	// Ignore mouse grab error
	const pointerEventMask = x.EventMaskButtonRelease | x.EventMaskButtonPress
	err = mousebind.GrabPointer(conn, rootWin, pointerEventMask, x.None, x.None)
	if err != nil {
		err1 := keybind.UngrabKeyboard(conn)
		if err1 != nil {
			logger.Warning("ungrabKeyboard Failed:", err)
		}
		return err
	}
	return nil
}

func ungrabKbdAndMouse(conn *x.Conn) {
	err := keybind.UngrabKeyboard(conn)
	if err != nil {
		logger.Warning("ungrabKeyboard Failed", err)
	}
	err = mousebind.UngrabPointer(conn)
	if err != nil {
		logger.Warning("ungrabPointer Failed", err)
	}
}
