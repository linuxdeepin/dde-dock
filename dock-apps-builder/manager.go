package main

import "dbus/com/deepin/daemon/dock"
import "dlib"
import "dlib/dbus"
import "dlib/logger"
import "flag"
import "fmt"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/xwindow"
import "github.com/BurntSushi/xgbutil/xevent"
import "github.com/BurntSushi/xgb/xproto"
import "github.com/BurntSushi/xgbutil/xprop"
import "github.com/BurntSushi/xgbutil/ewmh"
import "os"
import "path/filepath"

var (
	XU, _                 = xgbutil.NewConn()
	TrayXU, _             = xgbutil.NewConn()
	_NET_CLIENT_LIST, _   = xprop.Atm(XU, "_NET_CLIENT_LIST")
	_NET_ACTIVE_WINDOW, _ = xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	ATOM_WINDOW_ICON, _   = xprop.Atm(XU, "_NET_WM_ICON")
	ATOM_WINDOW_NAME, _   = xprop.Atm(XU, "_NET_WM_NAME")
	ATOM_WINDOW_STATE, _  = xprop.Atm(XU, "_NET_WM_STATE")
	ATOM_WINDOW_TYPE, _   = xprop.Atm(XU, "_NET_WM_WINDOW_TYPE")
	MANAGER               = initManager()
	LOGGER                = logger.NewLogger("com.deepin.daemon.DockAppsBuilder")
)

func listenRootWindow() {
	var update = func() {
		list, err := ewmh.ClientListGet(XU)
		if err != nil {
			LOGGER.Warning("Can't Get _NET_CLIENT_LIST", err)
		}
		MANAGER.runtimeAppChangged(list)
	}

	xwindow.New(XU, XU.RootWin()).Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case _NET_CLIENT_LIST:
			update()
		case _NET_ACTIVE_WINDOW:
			if activedWindow, err := ewmh.ActiveWindowGet(XU); err == nil {
				appId := find_app_id_by_xid(activedWindow)
				if rApp, ok := MANAGER.runtimeApps[appId]; ok {
					rApp.setLeader(activedWindow)
				}
			}
		}
	}).Connect(XU, XU.RootWin())
	update()
}

type Manager struct {
	runtimeApps map[string]*RuntimeApp
	normalApps  map[string]*NormalApp
	appEntries  map[string]*AppEntry
}

func initManager() *Manager {
	m := &Manager{}
	m.runtimeApps = make(map[string]*RuntimeApp)
	m.normalApps = make(map[string]*NormalApp)
	m.appEntries = make(map[string]*AppEntry)
	return m
}

func (m *Manager) listenDockedApp() {
	// TODO:
	if DOCKED_APP_MANAGER == nil {
		var err error
		DOCKED_APP_MANAGER, err = dock.NewDockedAppManager(
			"com.deepin.daemon.Dock",
			"/dde/dock/DockedAppManager",
		)
		if err != nil {
			LOGGER.Warning("get DockedAppManager failed", err)
			return
		}
	}

	DOCKED_APP_MANAGER.ConnectDocked(func(id string) {
		if _, ok := m.appEntries[id]; ok {
			LOGGER.Info(id, "is docked")
			return
		}

		m.createNormalApp(id + ".desktop")
	})

	DOCKED_APP_MANAGER.ConnectUndocked(func(id string) {
		LOGGER.Info("TODO: handle undockde")
	})
}

func (m *Manager) runtimeAppChangged(xids []xproto.Window) {
	willBeDestroied := make(map[string]*RuntimeApp)
	for _, app := range m.runtimeApps {
		willBeDestroied[app.Id] = app
	}

	// 1. create newfound RuntimeApps
	for _, xid := range xids {
		if isNormalWindow(xid) {
			appId := find_app_id_by_xid(xid)
			if rApp, ok := m.runtimeApps[appId]; ok {
				willBeDestroied[appId] = nil
				rApp.attachXid(xid)
			} else {
				m.createRuntimeApp(xid)
			}
		}
	}
	// 2. destroy disappeared RuntimeApps since last runtimeAppChanged point
	for _, app := range willBeDestroied {
		if app != nil {
			m.destroyRuntimeApp(app)
		}
	}
}

func (m *Manager) mustGetEntry(nApp *NormalApp, rApp *RuntimeApp) *AppEntry {
	if rApp != nil {
		if e, ok := m.appEntries[rApp.Id]; ok {
			return e
		} else {
			e := NewAppEntryWithRuntimeApp(rApp)
			m.appEntries[rApp.Id] = e
			dbus.InstallOnSession(e)
			return e
		}
	} else if nApp != nil {
		if e, ok := m.appEntries[nApp.Id]; ok {
			return e
		} else {
			e := NewAppEntryWithNormalApp(nApp)
			m.appEntries[nApp.Id] = e
			dbus.InstallOnSession(e)
			return e
		}
	}
	panic("mustGetEntry: at least give one app")
}

func (m *Manager) destroyEntry(appId string) {
	if e, ok := m.appEntries[appId]; ok {
		e.detachNormalApp()
		e.detachRuntimeApp()
		dbus.ReleaseName(e)
		dbus.UnInstallObject(e)
		fmt.Println("destroyEntry:", appId)
	}
	delete(m.appEntries, appId)
}

func (m *Manager) updateEntry(appId string, nApp *NormalApp, rApp *RuntimeApp) {
	switch {
	case nApp == nil && rApp == nil:
		m.destroyEntry(appId)
	case nApp == nil && rApp != nil:
		e := m.mustGetEntry(nApp, rApp)
		e.attachRuntimeApp(rApp)
		e.detachNormalApp()
	case nApp != nil && rApp != nil:
		e := m.mustGetEntry(nApp, rApp)
		e.attachNoramlApp(nApp)
		e.attachRuntimeApp(rApp)
	case nApp != nil && rApp == nil:
		e := m.mustGetEntry(nApp, rApp)
		e.attachNoramlApp(nApp)
		e.detachRuntimeApp()
	}
}

func (m *Manager) createRuntimeApp(xid xproto.Window) {
	appId := find_app_id_by_xid(xid)
	if _, ok := m.runtimeApps[appId]; ok {
		return
	}

	//TODO: xid 未改变但appId改变的情况， 比如nautils/libreoffice就会动态改变

	rApp := NewRuntimeApp(xid, appId)
	if rApp == nil {
		return
	}

	m.runtimeApps[appId] = rApp
	m.updateEntry(appId, m.mustGetEntry(nil, rApp).nApp, rApp)
}
func (m *Manager) destroyRuntimeApp(rApp *RuntimeApp) {
	delete(m.runtimeApps, rApp.Id)
	m.updateEntry(rApp.Id, m.mustGetEntry(nil, rApp).nApp, nil)
}
func (m *Manager) createNormalApp(id string) {
	fmt.Println("createNormalApp")
	if _, ok := m.normalApps[id]; ok {
		return
	}

	nApp := NewNormalApp(id)
	if nApp == nil {
		newId := filepath.Join(
			os.Getenv("HOME"),
			".config/dock/scratch",
			id,
		)
		nApp = NewNormalApp(newId)
		if nApp == nil {
			return
		}
	}

	m.normalApps[id] = nApp
	m.updateEntry(id, nApp, m.mustGetEntry(nApp, nil).rApp)
}
func (m *Manager) destroyNormalApp(nApp *NormalApp) {
	delete(m.normalApps, nApp.Id)
	m.updateEntry(nApp.Id, nil, m.mustGetEntry(nApp, nil).rApp)
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "d", false, "start debug")
	flag.Parse()
	if debug {
		LOGGER.SetLogLevel(logger.LEVEL_DEBUG)
	}

	dlib.InitI18n()
	initDeepin()
	for _, id := range loadAll() {
		// LOGGER.Debug(id)
		MANAGER.createNormalApp(id + ".desktop")
	}
	listenRootWindow()
	initTrayManager()
	go xevent.Main(XU)
	dbus.Wait()
}
