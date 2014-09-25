package dock

import "dbus/com/deepin/daemon/dock"
import "pkg.linuxdeepin.com/lib/dbus"

// import "flag"
import "github.com/BurntSushi/xgb/xproto"
import "os"
import "path/filepath"
import "strings"

var (
	ENTRY_MANAGER = NewEntryManager()
)

type EntryManager struct {
	runtimeApps map[string]*RuntimeApp
	normalApps  map[string]*NormalApp
	appEntries  map[string]*AppEntry
}

func NewEntryManager() *EntryManager {
	m := &EntryManager{}
	m.runtimeApps = make(map[string]*RuntimeApp)
	m.normalApps = make(map[string]*NormalApp)
	m.appEntries = make(map[string]*AppEntry)
	return m
}

func (m *EntryManager) listenDockedApp() {
	if DOCKED_APP_MANAGER == nil {
		var err error
		DOCKED_APP_MANAGER, err = dock.NewDockedAppManager(
			"com.deepin.daemon.Dock",
			"/dde/dock/DockedAppManager",
		)
		if err != nil {
			logger.Warning("get DockedAppManager failed", err)
			return
		}
	}

	DOCKED_APP_MANAGER.ConnectDocked(func(id string) {
		if _, ok := m.normalApps[id]; ok {
			logger.Info(id, "is already docked")
			return
		}
		m.createNormalApp(id)
	})

	DOCKED_APP_MANAGER.ConnectUndocked(func(id string) {
		// undocked is operated on normal app
		logger.Info("ConnectUndocked: Undock", id)
		if app, ok := m.normalApps[id]; ok {
			logger.Info("destroy normal app")
			m.destroyNormalApp(app)
		}
	})
}

func (m *EntryManager) runtimeAppChanged(xids []xproto.Window) {
	logger.Debug("runtime app changed")
	willBeDestroied := make(map[string]*RuntimeApp)
	for _, app := range m.runtimeApps {
		willBeDestroied[app.Id] = app
	}

	// 1. create newfound RuntimeApps
	for _, xid := range xids {
		if isNormalWindow(xid) {
			appId := find_app_id_by_xid(xid,
				DisplayModeType(setting.GetDisplayMode()))
			if rApp, ok := m.runtimeApps[appId]; ok {
				logger.Debugf("%s is alread existed, attach xid: 0x%x", appId, xid)
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

func (m *EntryManager) mustGetEntry(nApp *NormalApp, rApp *RuntimeApp) *AppEntry {
	if rApp != nil {
		if e, ok := m.appEntries[rApp.Id]; ok {
			return e
		} else {
			e := NewAppEntryWithRuntimeApp(rApp)
			m.appEntries[rApp.Id] = e
			err := dbus.InstallOnSession(e)
			if err != nil {
				logger.Warning(err)
			}
			return e
		}
	} else if nApp != nil {
		if e, ok := m.appEntries[nApp.Id]; ok {
			return e
		} else {
			e := NewAppEntryWithNormalApp(nApp)
			m.appEntries[nApp.Id] = e
			err := dbus.InstallOnSession(e)
			if err != nil {
				logger.Warning(err)
			}
			return e
		}
	}
	panic("mustGetEntry: at least give one app")
}

func (m *EntryManager) destroyEntry(appId string) {
	if e, ok := m.appEntries[appId]; ok {
		e.detachNormalApp()
		e.detachRuntimeApp()
		dbus.ReleaseName(e)
		dbus.UnInstallObject(e)
		logger.Info("destroyEntry:", appId)
	}
	delete(m.appEntries, appId)
}

func (m *EntryManager) updateEntry(appId string, nApp *NormalApp, rApp *RuntimeApp) {
	switch {
	case nApp == nil && rApp == nil:
		m.destroyEntry(appId)
	case nApp == nil && rApp != nil:
		e := m.mustGetEntry(nApp, rApp)
		e.attachRuntimeApp(rApp)
		e.detachNormalApp()
	case nApp != nil && rApp != nil:
		e := m.mustGetEntry(nApp, rApp)
		e.attachNormalApp(nApp)
		e.attachRuntimeApp(rApp)
	case nApp != nil && rApp == nil:
		e := m.mustGetEntry(nApp, rApp)
		e.attachNormalApp(nApp)
		e.detachRuntimeApp()
	}
}

func (m *EntryManager) createRuntimeApp(xid xproto.Window) *RuntimeApp {
	appId := find_app_id_by_xid(xid,
		DisplayModeType(setting.GetDisplayMode()))
	if appId == "" {
		logger.Debug("get appid for", xid, "failed")
		return nil
	}

	if v, ok := m.runtimeApps[appId]; ok {
		return v
	}

	rApp := NewRuntimeApp(xid, appId)
	if rApp == nil {
		return nil
	}

	m.runtimeApps[appId] = rApp
	m.updateEntry(appId, m.mustGetEntry(nil, rApp).nApp, rApp)
	return rApp
}
func (m *EntryManager) destroyRuntimeApp(rApp *RuntimeApp) {
	delete(m.runtimeApps, rApp.Id)
	m.updateEntry(rApp.Id, m.mustGetEntry(nil, rApp).nApp, nil)
}
func (m *EntryManager) createNormalApp(id string) {
	logger.Info("createNormalApp for", id)
	if _, ok := m.normalApps[id]; ok {
		logger.Info("normal app for", id, "is exist")
		return
	}

	desktopId := id + ".desktop"
	nApp := NewNormalApp(desktopId)
	if nApp == nil {
		logger.Info("create from scratch file")
		newId := filepath.Join(
			os.Getenv("HOME"),
			".config/dock/scratch",
			desktopId,
		)
		nApp = NewNormalApp(newId)
		if nApp == nil {
			logger.Warning("create normal app failed:", id)
			return
		}
	}

	m.normalApps[id] = nApp
	m.updateEntry(id, nApp, m.mustGetEntry(nApp, nil).rApp)
}
func (m *EntryManager) destroyNormalApp(nApp *NormalApp) {
	delete(m.normalApps, nApp.Id)
	m.updateEntry(nApp.Id, nil, m.mustGetEntry(nApp, nil).rApp)
}

func initialize() {
	ENTRY_MANAGER.listenDockedApp()
	for _, id := range loadAll() {
		id = strings.ToLower(strings.Replace(id, "_", "-", -1))
		logger.Debug("load", id)
		ENTRY_MANAGER.createNormalApp(id)
	}
	initTrayManager()
}
