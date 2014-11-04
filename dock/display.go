package dock

import (
	"dbus/com/deepin/daemon/display"
	"github.com/BurntSushi/xgb/xproto"
)

var displayRect xproto.Rectangle

func initDisplay() bool {
	dpy, err := display.NewDisplay(
		"com.deepin.daemon.Display",
		"/com/deepin/daemon/Display",
	)
	if err != nil {
		logger.Error("connect to display failed:", err)
		return false
	}
	// to avoid get PrimaryRect failed
	defer func() {
		if r := recover(); r != nil {
			logger.Warning("Recovered in initDisplay", r)
		}
	}()
	setDisplayRect(dpy.PrimaryRect.Get())
	dpy.ConnectPrimaryChanged(func(rect []interface{}) {
		setDisplayRect(rect)

		for _, app := range ENTRY_MANAGER.runtimeApps {
			for _, winInfo := range app.xids {
				winInfo.OverlapDock = isWindowOverlapDock(winInfo.Xid)
			}
		}
		hideModemanager.UpdateState()
	})
	return true
}

func setDisplayRect(rect []interface{}) {
	if len(rect) != 4 {
		return
	}
	displayRect.X, _ = rect[0].(int16)
	displayRect.Y, _ = rect[1].(int16)
	displayRect.Width, _ = rect[2].(uint16)
	displayRect.Height, _ = rect[3].(uint16)
}
