package dock

import (
	"dbus/com/deepin/api/xmousearea"
	"dbus/com/deepin/daemon/display"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
	"os"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
	"time"
)

var (
	logger                                     = log.NewLogger("dde-daemon/dock")
	region              *Region                = nil
	setting             *Setting               = nil
	hideModemanager     *HideStateManager      = nil
	dpy                 *display.Display       = nil
	dockProperty        *DockProperty          = nil
	entryProxyerManager *EntryProxyerManager   = nil
	dockedAppManager    *DockedAppManager      = nil
	areaImp             *xmousearea.XMouseArea = nil
	mouseArea           *XMouseAreaProxyer     = nil
)

func Stop() {
	if dockProperty != nil {
		dockProperty.destroy()
		dockProperty = nil
	}

	if dockedAppManager != nil {
		dockedAppManager.destroy()
		dockedAppManager = nil
	}

	if region != nil {
		region.destroy()
		region = nil
	}

	if setting != nil {
		setting.destroy()
		setting = nil
	}

	if hideModemanager != nil {
		hideModemanager.destroy()
		hideModemanager = nil
	}

	if entryProxyerManager != nil {
		entryProxyerManager.destroy()
		entryProxyerManager = nil
	}

	if mouseArea != nil {
		mouseArea.destroy()
		xmousearea.DestroyXMouseArea(areaImp)
		mouseArea = nil
	}

	if dpy != nil {
		display.DestroyDisplay(dpy)
		dpy = nil
	}

	if XU != nil {
		XU.Conn().Close()
		XU = nil
	}

	if TrayXU != nil {
		TrayXU.Conn().Close()
		TrayXU = nil
	}

	logger.EndTracing()
}

func startFailed(args ...interface{}) {
	logger.Error(args...)
	Stop()
}

func initAtom() {
	_NET_SHOWING_DESKTOP, _ = xprop.Atm(XU, "_NET_SHOWING_DESKTOP")
	DEEPIN_SCREEN_VIEWPORT, _ = xprop.Atm(XU, "DEEPIN_SCREEN_VIEWPORT")
	_NET_CLIENT_LIST, _ = xprop.Atm(XU, "_NET_CLIENT_LIST")
	_NET_ACTIVE_WINDOW, _ = xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	ATOM_WINDOW_ICON, _ = xprop.Atm(XU, "_NET_WM_ICON")
	ATOM_WINDOW_NAME, _ = xprop.Atm(XU, "_NET_WM_NAME")
	ATOM_WINDOW_STATE, _ = xprop.Atm(XU, "_NET_WM_STATE")
	ATOM_WINDOW_TYPE, _ = xprop.Atm(XU, "_NET_WM_WINDOW_TYPE")
	ATOM_DOCK_APP_ID, _ = xprop.Atm(XU, "_DDE_DOCK_APP_ID")

	_NET_SYSTEM_TRAY_S0, _ = xprop.Atm(TrayXU, "_NET_SYSTEM_TRAY_S0")
	_NET_SYSTEM_TRAY_OPCODE, _ = xprop.Atm(TrayXU, "_NET_SYSTEM_TRAY_OPCODE")
}

func Start() {
	if dockProperty != nil {
		return
	}

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

	XU, err = xgbutil.NewConn()
	if err != nil {
		startFailed(err)
		return
	}

	TrayXU, err = xgbutil.NewConn()
	if err != nil {
		startFailed(err)
		return
	}

	initAtom()

	dockProperty = NewDockProperty()
	err = dbus.InstallOnSession(dockProperty)
	if err != nil {
		startFailed("register dbus interface failed:", err)
		return
	}

	entryProxyerManager = NewEntryProxyerManager()
	err = dbus.InstallOnSession(entryProxyerManager)
	if err != nil {
		startFailed("register dbus interface failed:", err)
		return
	}

	entryProxyerManager.watchEntries()

	dockedAppManager = NewDockedAppManager()
	err = dbus.InstallOnSession(dockedAppManager)
	if err != nil {
		startFailed("register dbus interface failed:", err)
		return
	}

	setting = NewSetting()
	if setting == nil {
		startFailed("get setting failed")
	}
	err = dbus.InstallOnSession(setting)
	if err != nil {
		startFailed("register dbus interface failed:", err)
		return
	}

	dockProperty.updateDockHeight(DisplayModeType(setting.GetDisplayMode()))

	hideModemanager =
		NewHideStateManager(HideModeType(setting.GetHideMode()))
	err = dbus.InstallOnSession(hideModemanager)
	if err != nil {
		startFailed("register dbus interface failed:", err)
		return
	}
	hideModemanager.UpdateState()

	clientManager := NewClientManager()
	err = dbus.InstallOnSession(clientManager)
	if err != nil {
		startFailed("register dbus interface failed:", err)
		return
	}
	go clientManager.listenRootWindow()

	region = NewRegion()
	err = dbus.InstallOnSession(region)
	if err != nil {
		startFailed("register dbus interface failed:", err)
		return
	}

	areaImp, err = xmousearea.NewXMouseArea(
		"com.deepin.api.XMouseArea",
		"/com/deepin/api/XMouseArea",
	)
	mouseArea, err = NewXMouseAreaProxyer(areaImp, err)

	if err != nil {
		startFailed("register xmouse area failed:", err)
		return
	}

	err = dbus.InstallOnSession(mouseArea)
	if err != nil {
		startFailed(err)
		return
	}

	dbus.Emit(mouseArea, "InvalidId")
	mouseArea.connectMotionInto(func(_, _ int32, id string) {
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

	mouseArea.connectMotionOut(func(_, _ int32, id string) {
		if mouseAreaTimer != nil {
			mouseAreaTimer.Stop()
			mouseAreaTimer = nil
		}
		logger.Info("MouseOut:", id)
		hideModemanager.UpdateState()
	})

	initialize()
}
