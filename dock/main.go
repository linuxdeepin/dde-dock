package dock

import (
	"dbus/com/deepin/api/xmousearea"
	"dbus/com/deepin/daemon/display"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
	"os"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

var (
	logger                                     = log.NewLogger("daemon/dock")
	region              *Region                = nil
	setting             *Setting               = nil
	hideModemanager     *HideStateManager      = nil
	dpy                 *display.Display       = nil
	dockProperty        *DockProperty          = nil
	entryProxyerManager *EntryProxyerManager   = nil
	DOCKED_APP_MANAGER  *DockedAppManager      = nil
	areaImp             *xmousearea.XMouseArea = nil
	mouseArea           *XMouseAreaProxyer     = nil
)

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("dock", daemon, logger)
	return daemon
}

func (d *Daemon) Stop() error {
	if dockProperty != nil {
		dockProperty.destroy()
		dockProperty = nil
	}

	if DOCKED_APP_MANAGER != nil {
		DOCKED_APP_MANAGER.destroy()
		DOCKED_APP_MANAGER = nil
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
	return nil
}

func (d *Daemon) startFailed(args ...interface{}) {
	logger.Error(args...)
	d.Stop()
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

func (d *Daemon) Start() error {
	if dockProperty != nil {
		return nil
	}

	logger.BeginTracing()

	initDeepin()

	if logger.GetLogLevel() == log.LevelDebug {
		os.Setenv("G_MESSAGES_DEBUG", "all")
	}

	if !initDisplay() {
		d.Stop()
		logger.Info("initialize display failed")
		return nil
	}
	logger.Info("initialize display done")

	var err error

	XU, err = xgbutil.NewConn()
	if err != nil {
		d.startFailed(err)
		return err
	}

	TrayXU, err = xgbutil.NewConn()
	if err != nil {
		d.startFailed(err)
		return err
	}

	initAtom()
	logger.Info("initialize atoms done")

	dockProperty = NewDockProperty()
	err = dbus.InstallOnSession(dockProperty)
	if err != nil {
		d.startFailed("register dbus interface failed:", err)
		return err
	}
	logger.Info("initialize dock property done")

	entryProxyerManager = NewEntryProxyerManager()
	err = dbus.InstallOnSession(entryProxyerManager)
	if err != nil {
		d.startFailed("register dbus interface failed:", err)
		return err
	}
	entryProxyerManager.watchEntries()
	logger.Info("initialize entry proxyer manager done")

	DOCKED_APP_MANAGER = NewDockedAppManager()
	err = dbus.InstallOnSession(DOCKED_APP_MANAGER)
	if err != nil {
		d.startFailed("register dbus interface failed:", err)
		return err
	}
	logger.Info("initialize docked app manager done")

	setting = NewSetting()
	if setting == nil {
		d.startFailed("get setting failed")
	}
	err = dbus.InstallOnSession(setting)
	if err != nil {
		d.startFailed("register dbus interface failed:", err)
		return err
	}
	logger.Info("initialize settings done")

	dockProperty.updateDockHeight(DisplayModeType(setting.GetDisplayMode()))

	hideModemanager =
		NewHideStateManager(HideModeType(setting.GetHideMode()))
	err = dbus.InstallOnSession(hideModemanager)
	if err != nil {
		d.startFailed("register dbus interface failed:", err)
		return err
	}
	logger.Info("initialize hide mode manager done")

	hideModemanager.UpdateState()

	clientManager := NewClientManager()
	err = dbus.InstallOnSession(clientManager)
	if err != nil {
		d.startFailed("register dbus interface failed:", err)
		return err
	}
	go clientManager.listenRootWindow()
	logger.Info("initialize client manager done")

	region = NewRegion()
	err = dbus.InstallOnSession(region)
	if err != nil {
		d.startFailed("register dbus interface failed:", err)
		return err
	}

	areaImp, err = xmousearea.NewXMouseArea(
		"com.deepin.api.XMouseArea",
		"/com/deepin/api/XMouseArea",
	)
	mouseArea, err = NewXMouseAreaProxyer(areaImp, err)

	if err != nil {
		d.startFailed("register xmouse area failed:", err)
		return err
	}

	err = dbus.InstallOnSession(mouseArea)
	if err != nil {
		d.startFailed(err)
		return err
	}

	dbus.Emit(mouseArea, "InvalidId")
	mouseArea.connectMotionInto(func(_, _ int32, id string) {
		if mouseAreaTimer != nil {
			mouseAreaTimer.Stop()
			mouseAreaTimer = nil
		}
		mouseAreaTimer = time.AfterFunc(TOGGLE_HIDE_TIME, func() {
			logger.Debug("MouseIn:", id)
			mouseAreaTimer = nil
			hideModemanager.UpdateState()
		})
	})

	mouseArea.connectMotionOut(func(_, _ int32, id string) {
		if mouseAreaTimer != nil {
			mouseAreaTimer.Stop()
			mouseAreaTimer = nil
		}
		logger.Debug("MouseOut:", id)
		hideModemanager.UpdateState()
	})

	initialize()
	logger.Info("initialize done")
	return nil
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Name() string {
	return "dock"
}

func initialize() {
	for _, id := range DOCKED_APP_MANAGER.DockedAppList() {
		id = normalizeAppID(id)
		logger.Debug("load", id)
		ENTRY_MANAGER.createNormalApp(id)
	}
	initTrayManager()
}
