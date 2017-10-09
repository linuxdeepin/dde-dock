/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/lib/dbus"
)

func (m *Manager) doGrabScreen(ss *shortcuts.ShortcutManager) error {
	conn, err := x.NewConn()
	if err != nil {
		return err
	}

	err = grabKbdAndMouse(conn)
	if err != nil {
		return err
	}

	// Temporarily disable record
	ss.EnableRecord(false)
	defer ss.EnableRecord(true)

loop:
	for {
		event := conn.WaitForEvent()
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
			accel := key.ToAccel(m.keySymbols)
			dbus.Emit(m, "KeyEvent", true, accel.String())
			if accel.IsGood() {
				logger.Debug("good accel", accel)
				m.grabScreenPressedAccel = &accel
			} else {
				logger.Debug("bad accel", accel)
				m.grabScreenPressedAccel = nil
			}
		case x.KeyReleaseEventCode:
			event, _ := x.NewKeyReleaseEvent(event)
			logger.Debug(event)
			if m.grabScreenPressedAccel != nil {
				dbus.Emit(m, "KeyEvent", false, m.grabScreenPressedAccel.String())
				m.grabScreenPressedAccel = nil
			} else {
				dbus.Emit(m, "KeyEvent", false, "")
			}

			break loop
		case x.ButtonPressEventCode:
			dbus.Emit(m, "KeyEvent", true, "")
			break loop
		case x.ButtonReleaseEventCode:
			dbus.Emit(m, "KeyEvent", false, "")
			break loop
		}
	}

	ungrabKbdAndMouse(conn)
	conn.Close()
	return nil
}

func grabKbdAndMouse(conn *x.Conn) error {
	rootWin := conn.GetDefaultScreen().Root
	err := keybind.GrabKeyboard(conn, rootWin)
	if err != nil {
		return err
	}

	// Ignore mouse grab error
	const pointerEventMask = x.EventMaskButtonRelease | x.EventMaskButtonPress
	mousebind.GrabPointer(conn, rootWin, pointerEventMask, x.None, x.None)
	return nil
}

func ungrabKbdAndMouse(conn *x.Conn) {
	keybind.UngrabKeyboard(conn)
	mousebind.UngrabPointer(conn)
}
