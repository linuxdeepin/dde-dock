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

package inputdevices

import (
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest = "com.deepin.daemon.InputDevices"
	dbusPath = "/com/deepin/daemon/InputDevices"
	dbusIFC  = dbusDest

	kbdDBusPath = "/com/deepin/daemon/InputDevice/Keyboard"
	kbdDBusIFC  = "com.deepin.daemon.InputDevice.Keyboard"

	mouseDBusPath     = "/com/deepin/daemon/InputDevice/Mouse"
	mouseDBusIFC      = "com.deepin.daemon.InputDevice.Mouse"
	trackPointDBusIFC = "com.deepin.daemon.InputDevice.TrackPoint"

	tpadDBusPath = "/com/deepin/daemon/InputDevice/TouchPad"
	tpadDBusIFC  = "com.deepin.daemon.InputDevice.TouchPad"

	wacomDBusPath = "/com/deepin/daemon/InputDevice/Wacom"
	wacomDBusIFC  = "com.deepin.daemon.InputDevice.Wacom"
)

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (*Keyboard) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: kbdDBusPath,
		Interface:  kbdDBusIFC,
	}
}

func (*Mouse) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: mouseDBusPath,
		Interface:  mouseDBusIFC,
	}
}

func (*TrackPoint) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: mouseDBusPath,
		Interface:  trackPointDBusIFC,
	}
}

func (m *Mouse) setPropExist(exist bool) {
	if exist == m.Exist {
		return
	}

	m.Exist = exist
	dbus.NotifyChange(m, "Exist")
}

func (tp *TrackPoint) setPropExist(exist bool) {
	if exist == tp.Exist {
		return
	}

	tp.Exist = exist
	dbus.NotifyChange(tp, "Exist")
}

func (*Touchpad) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: tpadDBusPath,
		Interface:  tpadDBusIFC,
	}
}

func (tpad *Touchpad) setPropExist(exist bool) {
	if exist == tpad.Exist {
		return
	}

	tpad.Exist = exist
	dbus.NotifyChange(tpad, "Exist")
}

func (*Wacom) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: wacomDBusPath,
		Interface:  wacomDBusIFC,
	}
}

func (w *Wacom) setPropExist(exist bool) {
	if exist == w.Exist {
		return
	}

	w.Exist = exist
	dbus.NotifyChange(w, "Exist")
}

func (w *Wacom) setPropMapOutput(output string) bool {
	if output == w.MapOutput {
		return false
	}

	w.MapOutput = output
	dbus.NotifyChange(w, "MapOutput")
	return true
}

func setPropString(obj dbus.DBusObject, handler *string, prop, v string) {
	if *handler == v {
		return
	}
	*handler = v
	dbus.NotifyChange(obj, prop)
}
