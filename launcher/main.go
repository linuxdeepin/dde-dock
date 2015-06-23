package launcher

import (
	"database/sql"
	storeApi "dbus/com/deepin/store/api"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/category"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/item"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/item/search"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/item/softwarecenter"
	. "pkg.linuxdeepin.com/dde-daemon/loader"
	. "pkg.linuxdeepin.com/lib/gettext"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	. "pkg.linuxdeepin.com/lib/initializer"
	"pkg.linuxdeepin.com/lib/log"
	"sync"
)

type Daemon struct {
	*ModuleBase
	launcher *Launcher
}

func NewLauncherDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = NewModuleBase("launcher", daemon, logger)
	return daemon
}

//TODO
func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Stop() error {
	if d.launcher == nil {
		return nil
	}

	d.launcher.destroy()
	d.launcher = nil

	logger.EndTracing()
	return nil
}

func loadItemsInfo(im *ItemManager, cm *CategoryManager) {
	timeInfo, _ := im.GetAllTimeInstalled()

	appChan := make(chan *gio.AppInfo)
	go func() {
		allApps := gio.AppInfoGetAll()
		for _, app := range allApps {
			appChan <- app
		}
		close(appChan)
	}()

	dbPath, _ := GetDBPath(SoftwareCenterDataDir, CategoryNameDBPath)
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
				item := NewItem(desktopApp)
				cid, err := QueryCategoryId(desktopApp, db)
				if err != nil {
					item.SetCategoryId(OtherID)
				}
				item.SetCategoryId(cid)
				item.SetTimeInstalled(timeInfo[item.Id()])

				im.AddItem(item)
				cm.AddItem(item.Id(), item.GetCategoryId())

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

func (d *Daemon) Start() error {
	if d.launcher != nil {
		return nil
	}

	logger.BeginTracing()

	InitI18n()

	// DesktopAppInfo.ShouldShow does not know deepin.
	gio.DesktopAppInfoSetDesktopEnv("Deepin")

	err := NewInitializer().Init(func(interface{}) (interface{}, error) {
		return NewSoftwareCenter()
	}).InitOnSessionBus(func(soft interface{}) (interface{}, error) {
		d.launcher = NewLauncher()

		im := NewItemManager(soft.(SoftwareCenterInterface))
		cm := NewCategoryManager()

		d.launcher.setItemManager(im)
		d.launcher.setCategoryManager(cm)

		loadItemsInfo(im, cm)

		store, err := storeApi.NewDStoreDesktop("com.deepin.store.Api", "/com/deepin/store/Api")
		if err == nil {
			d.launcher.setStoreApi(store)
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
		setting, err := NewSetting(coreSetting)
		d.launcher.setSetting(setting)
		return setting, err
	}).GetError()

	if err != nil {
		d.Stop()
	}

	return err
}
