package inputdevices

import (
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest = "com.deepin.daemon.InputDevices"
	dbusPath = "/com/deepin/daemon/InputDevices"
	dbusIFC = dbusDest

	kbdDBusPath = "/com/deepin/daemon/InputDevice/Keyboard"
	kbdDBusIFC = "com.deepin.daemon.InputDevice.Keyboard"

	mouseDBusPath = "/com/deepin/daemon/InputDevice/Mouse"
	mouseDBusIFC  = "com.deepin.daemon.InputDevice.Mouse"

	tpadDBusPath = "/com/deepin/daemon/InputDevice/TouchPad"
	tpadDBusIFC = "com.deepin.daemon.InputDevice.TouchPad"

	wacomDBusPath = "/com/deepin/daemon/InputDevice/Wacom"
	wacomDBusIFC  = "com.deepin.daemon.InputDevice.Wacom"
)

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest: dbusDest,
		ObjectPath: dbusPath,
		Interface: dbusIFC,
	}
}

func (*Keyboard) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest: dbusDest,
		ObjectPath: kbdDBusPath,
		Interface: kbdDBusIFC,
	}
}

func (*Mouse) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: mouseDBusPath,
		Interface:  mouseDBusIFC,
	}
}

func (m *Mouse) setPropExist(exist bool) {
	if exist == m.Exist {
		return
	}

	m.Exist = exist
	dbus.NotifyChange(m, "Exist")
}

func (*Touchpad) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest: dbusDest,
		ObjectPath: tpadDBusPath,
		Interface: tpadDBusIFC,
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
