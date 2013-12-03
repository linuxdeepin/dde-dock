// +build ignore

package main

import (
	"dlib/dbus"
	"strconv"
)

type Display struct {
	Index      int32   `access:"read"`
	Name       string  `access:"read"`
	Builtin    int32   `access:"read"`
	Active     int32   //turn on or off
	Primary    int32   //primary screen or not
	Resolution int32   //to be changed
	Refresh    float64 //refresh rate:auto,...
	Rotation   int32   //to be changed:0,90,180,270
	Brightness float64
	Mirror     int32    //whether to mirror displays
	Position   [4]int32 //position of this screen:Above,Left,Below,Right
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
