package main

import "dbus/com/deepin/daemon/keybinding"
import "strings"

func (dpy *Display) workaroundBacklight() {
	mediaKeyManager, err := keybinding.NewMediaKey("/com/deepin/daemon/MediaKey")
	if err != nil {
		Logger.Error("Can't connect to /com/deepin/daemon/MediaKey", err)
		return
	}

	workaround := func(m *Monitor) {
		names := strings.Split(m.Name, joinSeparator)
		for i, op := range m.outputs {
			if ok, backlight := supportedBacklight(X, op); ok {
				m.setPropBrightness(names[i], backlight)
			}
		}
	}

	mediaKeyManager.ConnectBrightnessUp(func(press bool) {
		if !press {
			for _, m := range dpy.Monitors {
				workaround(m)
			}
		}
	})
	mediaKeyManager.ConnectBrightnessDown(func(press bool) {
		if !press {
			for _, m := range dpy.Monitors {
				workaround(m)
			}
		}
	})
}
