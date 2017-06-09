/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"pkg.deepin.io/dde/daemon/keybinding/keybind"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/dde/daemon/keybinding/xrecord"
	"pkg.deepin.io/lib/dbus"
)

func (m *Manager) doGrabScreen() error {
	xu, err := xgbutil.NewConn()
	if err != nil {
		return err
	}
	keybind.Initialize(xu)
	mousebind.Initialize(xu)

	err = grabKbdAndMouse(xu)
	if err != nil {
		return err
	}

	// Disable xrecord
	xrecord.Enable(false)
	defer xrecord.Enable(true)

	xevent.ButtonPressFun(
		func(x *xgbutil.XUtil, e xevent.ButtonPressEvent) {
			dbus.Emit(m, "KeyEvent", true, "")
			ungrabKbdAndMouse(xu)
			xevent.Quit(xu)
		}).Connect(xu, xu.RootWin())

	xevent.ButtonReleaseFun(
		func(x *xgbutil.XUtil, e xevent.ButtonReleaseEvent) {
			dbus.Emit(m, "KeyEvent", false, "")
			ungrabKbdAndMouse(xu)
			xevent.Quit(xu)
		}).Connect(xu, xu.RootWin())

	xevent.KeyPressFun(
		func(x *xgbutil.XUtil, ev xevent.KeyPressEvent) {
			logger.Debug(ev)
			mods := shortcuts.GetConcernedModifiers(ev.State)
			logger.Debug("event mods:", shortcuts.Modifiers(ev.State))
			key := shortcuts.Key{
				Mods: mods,
				Code: shortcuts.Keycode(ev.Detail),
			}
			logger.Debug("event key:", key)
			accel := key.ToAccel(x)
			dbus.Emit(m, "KeyEvent", true, accel.String())
			if accel.IsGood() {
				logger.Debug("good accel", accel)
				m.grabScreenPressedAccel = &accel
			} else {
				logger.Debug("bad accel", accel)
				m.grabScreenPressedAccel = nil
			}
		}).Connect(xu, xu.RootWin())

	xevent.KeyReleaseFun(
		func(x *xgbutil.XUtil, ev xevent.KeyReleaseEvent) {
			logger.Debug(ev)
			if m.grabScreenPressedAccel != nil {
				dbus.Emit(m, "KeyEvent", false, m.grabScreenPressedAccel.String())
				m.grabScreenPressedAccel = nil
			} else {
				dbus.Emit(m, "KeyEvent", false, "")
			}

			ungrabKbdAndMouse(xu)
			xevent.Quit(xu)
		}).Connect(xu, xu.RootWin())

	xevent.Main(xu)
	return nil
}

func grabKbdAndMouse(xu *xgbutil.XUtil) error {
	err := keybind.GrabKeyboard(xu, xu.RootWin())
	if err != nil {
		return err
	}

	// Ignore mouse grab error
	mousebind.GrabChecked(xu, xu.RootWin(), 0, 1, false)
	mousebind.GrabChecked(xu, xu.RootWin(), 0, 2, false)
	mousebind.GrabChecked(xu, xu.RootWin(), 0, 3, false)
	return nil
}

func ungrabKbdAndMouse(xu *xgbutil.XUtil) {
	keybind.UngrabKeyboard(xu)
	mousebind.Ungrab(xu, xu.RootWin(), 0, 1)
	mousebind.Ungrab(xu, xu.RootWin(), 0, 2)
	mousebind.Ungrab(xu, xu.RootWin(), 0, 3)
}
