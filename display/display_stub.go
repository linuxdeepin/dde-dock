package main

import "dlib/dbus"
import "fmt"
import "github.com/BurntSushi/xgb/xproto"

func (dpy *Display) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Display",
		"/com/deepin/daemon/Display",
		"com.deepin.daemon.Display",
	}
}

func (dpy *Display) setPropWidth(v uint16) {
	if dpy.Width != v {
		dpy.Width = v
		dbus.NotifyChange(dpy, "Width")
	}
}

func (dpy *Display) setPropHeight(v uint16) {
	if dpy.Height != v {
		dpy.Height = v
		dbus.NotifyChange(dpy, "Height")
	}
}

func (dpy *Display) setPropOutputs(v []*Output) {
	dpy.Outputs = v
	dbus.NotifyChange(dpy, "Outputs")
}

func (dpy *Display) setPropRotation(v uint16) {
	if dpy.Rotation != v {
		dpy.Rotation = v
		dbus.NotifyChange(dpy, "Rotation")
	}
}

func (dpy *Display) setPropReflect(v uint16) {
	if dpy.Reflect != v {
		dpy.Reflect = v
		dbus.NotifyChange(dpy, "Reflect")
	}
}

func (dpy *Display) setPropPrimaryOutput(v *Output) {
	if dpy.PrimaryOutput != v {
		dpy.PrimaryOutput = v
		dbus.NotifyChange(dpy, "PrimaryOutput")
		if v != nil {
			fmt.Println("SetPropPrimaryOutput to ", v.Name)
		} else {
			fmt.Println("SetPropPrimaryOutput to None")
		}
	}
}

func (dpy *Display) setPropPrimaryRect(v xproto.Rectangle) {
	if dpy.PrimaryRect != v {
		fmt.Println("SetPropPrimaryRect to ", v)
		dpy.PrimaryRect = v
		dbus.NotifyChange(dpy, "PrimaryRect")

		if dpy.PrimaryChanged != nil {
			dpy.PrimaryChanged(dpy.PrimaryRect)
		}
	}
}

func (dpy *Display) setPropDisplayMode(v int16) {
	if dpy.DisplayMode != v {
		dpy.DisplayMode = v
		dbus.NotifyChange(dpy, "DisplayMode")
	}
}

func (dpy *Display) setPropBuiltinOutput(v *Output) {
	if dpy.BuiltinOutput != v {
		dpy.BuiltinOutput = v
		dbus.NotifyChange(dpy, "BuiltinOutput")
	}
}
