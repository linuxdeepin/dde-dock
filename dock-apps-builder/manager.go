package main

import "dlib/dbus"
import "dlib/logger"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/xwindow"
import "github.com/BurntSushi/xgbutil/xevent"
import "github.com/BurntSushi/xgb/xproto"
import "github.com/BurntSushi/xgbutil/xprop"
import "github.com/BurntSushi/xgbutil/ewmh"

var (
	XU, _                = xgbutil.NewConn()
	_NET_CLIENT_LIST, _  = xprop.Atm(XU, "_NET_CLIENT_LIST")
	ATOM_WINDOW_ICON, _  = xprop.Atm(XU, "_NET_WM_ICON")
	ATOM_WINDOW_NAME, _  = xprop.Atm(XU, "_NET_WM_NAME")
	ATOM_WINDOW_STATE, _ = xprop.Atm(XU, "_NET_WM_STATE")
	ATOM_WINDOW_TYPE, _  = xprop.Atm(XU, "_NET_WM_WINDOW_TYPE")
	MANAGER              = initManager()
	LOGGER               = logger.NewLogger("com.deepin.daemon.DockAppsBuilder")
)

func listenerRootWindow() {
	xwindow.New(XU, XU.RootWin()).Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		if ev.Atom == _NET_CLIENT_LIST {
			list, err := ewmh.ClientListGet(XU)
			if err != nil {
				LOGGER.Warning("Can't Get _NET_CLIENT_LIST", err)
			}
			MANAGER.runtimeAppChangged(list)
		}
	}).Connect(XU, XU.RootWin())
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

func (m *Manager) runtimeAppChangged(xids []xproto.Window) {
	willBeDestroied := make(map[string]*RuntimeApp)
	for _, app := range m.runtimeApps {
		willBeDestroied[app.Id] = app
	}
	for _, xid := range xids {
		appId := find_app_id_by_xid(xid)
		if _, ok := m.runtimeApps[appId]; ok {
			willBeDestroied[appId] = nil
		} else {
			m.createRuntimeApp(xid)
		}
	}
	for _, app := range willBeDestroied {
		m.destroyRuntimeApp(app)
	}
}

func (m *Manager) destroyEntry(appId string) {
}
func (m *Manager) mustGetEntry(appId string) *AppEntry {
	return nil
}

func (m *Manager) updateEntry(appId string, nApp *NormalApp, rApp *RuntimeApp) {
	if nApp == nil && rApp == nil {
		m.destroyEntry(appId)
	}
	m.mustGetEntry(appId)
	if _, ok := m.appEntries[appId]; !ok {
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
	if nApp, ok := m.normalApps[rApp.Id]; ok {
		m.updateEntry(appId, nApp, rApp)
	} else {
		// an undocked RuntimeApp
		m.updateEntry(appId, nil, rApp)
	}
}
func (m *Manager) destroyRuntimeApp(app *RuntimeApp) {
	m.updateEntry(app.Id, m.mustGetEntry(app.Id).nApp, nil)
	app.Destroy()
}

func main() {
	for _, id := range loadAll() {
		if e := NewAppEntry(id + ".desktop"); e != nil {
			dbus.InstallOnSession(e)
		}
	}
	listenerRootWindow()
	go xevent.Main(XU)
	dbus.Wait()
}
