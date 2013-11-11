package main

import (
	"dlib/dbus"
	"strconv"
)

type Display struct {
	Index           int32    `access:"read"`
	Name            string   `access:"read"`
	Builtin         int32	 `access:"read"`
	Resolution      int32//to be changed
	Rotation        int32//to be changed
	Brightness      float64
	Mirror          int32    //whether to mirror displays
	Active          int32    //turn on or off
}

func (display *Display) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Display",
		"/com/deepin/daemon/Display" + strconv.FormatInt(int64(display.Index), 10),
		"com.deepin.daemon.Display",
	}
}

func main() {
	dbus.InstallOnSession(&Display{})
	select {}
}
