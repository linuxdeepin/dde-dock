package main

import (
	"fmt"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/dde/daemon/system/gesture"
	"pkg.deepin.io/lib/dbus1"
)

func (*Daemon) SetLongPressDuration(duration uint32) *dbus.Error {
	epath := dbusPath + ".SetLongPressDuration"
	if duration < 1 {
		return dbus.NewError(epath,
			[]interface{}{fmt.Errorf("invalid duration: %d", duration)})
	}
	var m = loader.GetModule("gesture")
	if m == nil {
		return dbus.NewError(epath,
			[]interface{}{"Not found module 'gesture'"})
	}
	m.(*gesture.Daemon).SetLongPressDuration(int(duration))
	return nil
}
