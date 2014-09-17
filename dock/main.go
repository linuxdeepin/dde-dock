package dock

import (
	"dbus/com/deepin/api/xmousearea"
	"dbus/com/deepin/daemon/display"
	"os"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
	"time"
)

var (
	logger                            = log.NewLogger("com.deepin.daemon.Dock")
	region          *Region           = nil
	setting         *Setting          = nil
	hideModemanager *HideStateManager = nil
	dpy             *display.Display  = nil
	dockProperty    *DockProperty     = nil
)

func Stop() {
	logger.EndTracing()
	display.DestroyDisplay(dpy)
}
func Start() {
	logger.BeginTracing()

	initDeepin()

	if logger.GetLogLevel() == log.LevelDebug {
		os.Setenv("G_MESSAGES_DEBUG", "all")
	}

	if !initDisplay() {
		Stop()
		return
	}

	var err error
	dockProperty = NewDockProperty()
	dbus.InstallOnSession(dockProperty)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}

	m := NewEntryProxyerManager()
	err = dbus.InstallOnSession(m)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		Stop()
		return
	}

	m.watchEntries()

	d := NewDockedAppManager()
	err = dbus.InstallOnSession(d)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		Stop()
		return
	}

	setting = NewSetting()
	err = dbus.InstallOnSession(setting)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		Stop()
		return
	}

	dockProperty.updateDockHeight(DisplayModeType(setting.GetDisplayMode()))

	hideModemanager =
		NewHideStateManager(HideModeType(setting.GetHideMode()))
	err = dbus.InstallOnSession(hideModemanager)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		Stop()
		return
	}
	hideModemanager.UpdateState()

	cm := NewClientManager()
	err = dbus.InstallOnSession(cm)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		Stop()
		return
	}
	go cm.listenRootWindow()

	region = NewRegion()
	err = dbus.InstallOnSession(region)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		Stop()
		return
	}

	area, err := NewXMouseAreaProxyer(xmousearea.NewXMouseArea(
		"com.deepin.api.XMouseArea",
		"/com/deepin/api/XMouseArea",
	))

	if err != nil {
		logger.Errorf("register xmouse area failed:", err)
		Stop()
		return
	}

	dbus.InstallOnSession(area)

	area.connectMotionInto(func(_, _ int32, id string) {
		logger.Info("MouseIn:", id)
		if mouseAreaTimer != nil {
			mouseAreaTimer.Stop()
			mouseAreaTimer = nil
		}
		mouseAreaTimer = time.AfterFunc(TOGGLE_HIDE_TIME, func() {
			logger.Info("MouseIn:", id)
			mouseAreaTimer = nil
			hideModemanager.UpdateState()
		})
	})

	area.connectMotionOut(func(_, _ int32, id string) {
		if mouseAreaTimer != nil {
			mouseAreaTimer.Stop()
			mouseAreaTimer = nil
		}
		logger.Info("MouseOut:", id)
		hideModemanager.UpdateState()
	})

	initialize()
}
