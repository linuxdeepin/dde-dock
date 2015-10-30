package launcher

import (
	"database/sql"
	storeApi "dbus/com/deepin/store/api"
	"pkg.deepin.io/dde/daemon/launcher/category"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/dde/daemon/launcher/item"
	"pkg.deepin.io/dde/daemon/launcher/item/dstore"
	. "pkg.deepin.io/dde/daemon/launcher/item/search"
	. "pkg.deepin.io/dde/daemon/loader"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gio-2.0"
	. "pkg.deepin.io/lib/initializer"
	"pkg.deepin.io/lib/log"
	"sync"
)

// Daemon is the module interface's implementation.
type Daemon struct {
	*ModuleBase
	launcher *Launcher
}

// NewLauncherDaemon creates a new launcher daemon module.
func NewLauncherDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = NewModuleBase("launcher", daemon, logger)
	return daemon
}

// GetDependencies returns the dependencies of this module.
func (d *Daemon) GetDependencies() []string {
	return []string{}
}

// Stop stops the launcher module.
func (d *Daemon) Stop() error {
	if d.launcher == nil {
		return nil
	}

	d.launcher.destroy()
	d.launcher = nil

	logger.EndTracing()
	return nil
}

func loadItemsInfo(im *item.Manager, cm *category.Manager) {
	timeInfo, _ := im.GetAllTimeInstalled()

	appChan := make(chan *gio.AppInfo)
	go func() {
		allApps := gio.AppInfoGetAll()
		for _, app := range allApps {
			appChan <- app
		}
		close(appChan)
	}()

	dbPath, _ := category.GetDBPath(category.SoftwareCenterDataDir, category.CategoryNameDBPath)
	db, err := sql.Open("sqlite3", dbPath)

	var wg sync.WaitGroup
	const N = 20
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			for app := range appChan {
				if !app.ShouldShow() {
					app.Unref()
					continue
				}

				desktopApp := gio.ToDesktopAppInfo(app)
				newItem := item.New(desktopApp)
				cid, err := category.QueryID(desktopApp, db)
				newItem.SetCategoryID(cid)
				if err != nil {
				}
				newItem.SetTimeInstalled(timeInfo[newItem.ID()])

				im.AddItem(newItem)
				cm.AddItem(newItem.ID(), newItem.CategoryID())

				app.Unref()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if err == nil {
		db.Close()
	}
}

// Start starts the launcher module.
func (d *Daemon) Start() error {
	if d.launcher != nil {
		return nil
	}

	logger.BeginTracing()

	InitI18n()

	// DesktopAppInfo.ShouldShow does not know deepin.
	gio.DesktopAppInfoSetDesktopEnv("Deepin")

	err := NewInitializer().Init(func(interface{}) (interface{}, error) {
		return dstore.New()
	}).InitOnSessionBus(func(soft interface{}) (interface{}, error) {
		d.launcher = NewLauncher()

		im := item.NewManager(soft.(DStore))
		cm := category.NewManager()

		d.launcher.setItemManager(im)
		d.launcher.setCategoryManager(cm)

		loadItemsInfo(im, cm)

		store, err := storeApi.NewDStoreDesktop("com.deepin.store.Api", "/com/deepin/store/Api")
		if err == nil {
			d.launcher.setStoreAPI(store)
		}

		names := []string{}
		for _, item := range im.GetAllItems() {
			names = append(names, item.Name())
		}

		pinyinObj, err := NewPinYinSearchAdapter(names)
		d.launcher.setPinYinObject(pinyinObj)

		d.launcher.listenItemChanged()

		return d.launcher, nil
	}).InitOnSessionBus(func(interface{}) (interface{}, error) {
		coreSetting := gio.NewSettings("com.deepin.dde.launcher")
		setting, err := NewSettings(coreSetting)
		d.launcher.setSetting(setting)
		return setting, err
	}).GetError()

	if err != nil {
		d.Stop()
	}

	return err
}
