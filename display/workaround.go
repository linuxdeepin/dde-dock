package main

import "dbus/com/deepin/daemon/keybinding"

var __keepMediakeyManagerAlive interface{}

func (dpy *Display) workaroundBacklight() {
	mediaKeyManager, err := keybinding.NewMediaKey("com.deepin.daemon.KeyBinding", "/com/deepin/daemon/MediaKey")
	if err != nil {
		Logger.Error("Can't connect to /com/deepin/daemon/MediaKey", err)
		return
	}
	__keepMediakeyManagerAlive = mediaKeyManager

	workaround := func(m *Monitor) {
		for name, op := range GetDisplayInfo().outputNames {
			if ok, backlight := supportedBacklight(xcon, op); ok {
				dpy.setPropBrightness(name, backlight)
			}
		}
	}

	mediaKeyManager.ConnectBrightnessUp(func(onPress bool) {
		for _, m := range dpy.Monitors {
			workaround(m)
		}
	})
	mediaKeyManager.ConnectBrightnessDown(func(onPress bool) {
		for _, m := range dpy.Monitors {
			workaround(m)
		}
	})

	mediaKeyManager.ConnectSwitchMonitors(func(onPress bool) {
		if !onPress {
			if int(dpy.DisplayMode) >= len(dpy.Monitors) {
				dpy.SwitchMode(DisplayModeMirrors)
			} else {
				dpy.SwitchMode(dpy.DisplayMode + 1)
			}
		}
	})
}
func init() {
	GetDisplay().workaroundBacklight()
}
