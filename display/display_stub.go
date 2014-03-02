package main

import "github.com/BurntSushi/xgb/xproto"

import "dlib/dbus"

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

func (dpy *Display) setPropPrimaryRect(v xproto.Rectangle) {
	if dpy.PrimaryRect != v {
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
