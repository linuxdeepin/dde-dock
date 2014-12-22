package launcher

import (
	"database/sql"
	storeApi "dbus/com/deepin/store/api"
	"errors"
	"sync"
	// . "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/category"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/item"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/item/search"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/item/softwarecenter"
	"pkg.linuxdeepin.com/lib/dbus"
	. "pkg.linuxdeepin.com/lib/gettext"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/log"
)

var logger = log.NewLogger("dde-daemon/launcher-daemon")
var launcher *Launcher = nil

func Stop() {
	if launcher != nil {
		launcher.destroy()
		launcher = nil
	}

	logger.EndTracing()
}

func startFailed(err error) {
	logger.Error(err)
	Stop()
}

func Start() {
	var err error

	logger.BeginTracing()

	InitI18n()

	// DesktopAppInfo.ShouldShow does not know deepin.
	gio.DesktopAppInfoSetDesktopEnv("Deepin")

	soft, err := NewSoftwareCenter()
	if err != nil {
		startFailed(err)
		return
	}

	im := NewItemManager(soft)
	cm := NewCategoryManager()

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

	launcher = NewLauncher()
	launcher.setItemManager(im)
	launcher.setCategoryManager(cm)

	store, err := storeApi.NewDStoreDesktop("com.deepin.store.Api", "/com/deepin/store/Api")
	if err == nil {
		launcher.setStoreApi(store)
	}

	names := []string{}
	for _, item := range im.GetAllItems() {
		names = append(names, item.Name())
	}
	pinyinObj, err := NewPinYinSearchAdapter(names)
	launcher.setPinYinObject(pinyinObj)

	launcher.listenItemChanged()

	err = dbus.InstallOnSession(launcher)
	if err != nil {
		startFailed(err)
		return
	}

	coreSetting := gio.NewSettings("com.deepin.dde.launcher")
	if coreSetting == nil {
		startFailed(errors.New("get schema failed"))
		return
	}
	setting, err := NewSetting(coreSetting)
	if err != nil {
		startFailed(err)
		return
	}
	err = dbus.InstallOnSession(setting)
	if err != nil {
		startFailed(err)
		return
	}
	launcher.setSetting(setting)
}
