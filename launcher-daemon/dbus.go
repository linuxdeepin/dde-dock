package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/howeyc/fsnotify"

	"dlib/dbus"
	"dlib/gio-2.0"
	"dlib/glib-2.0"
)

const (
	launcherObject    string = "com.deepin.dde.daemon.Launcher"
	launcherPath      string = "/com/deepin/dde/daemon/Launcher"
	launcherInterface string = launcherObject

	AppDirName     string      = "applications"
	DirDefaultPerm os.FileMode = 0755

	SOFTWARE_STATUS_CREATED  string = "created"
	SOFTWARE_STATUS_MODIFIED string = "updated"
	SOFTWARE_STATUS_DELETED  string = "deleted"
)

type ItemChangedStatus struct {
	renamed, created, notRenamed, notCreated chan bool
}

type LauncherDBus struct {
	ItemChanged func(
		status string,
		itemInfo ItemInfo,
		categoryId CategoryId,
	)
}

func (d *LauncherDBus) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		launcherObject,
		launcherPath,
		launcherInterface,
	}
}

func (d *LauncherDBus) CategoryInfos() CategoryInfosResult {
	return getCategoryInfos()
}

func (d *LauncherDBus) ItemInfos(id int32) []ItemInfo {
	return getItemInfos(CategoryId(id))
}

func (d *LauncherDBus) emitItemChanged(name, status string, info map[string]ItemChangedStatus) {
	defer delete(info, name)
	id := genId(name)

	logger.Info("Status:", status)
	if status != SOFTWARE_STATUS_DELETED {
		logger.Info(name)
		app := gio.NewDesktopAppInfoFromFilename(name)
		for count := 0; app == nil; count++ {
			<-time.After(time.Millisecond * 200)
			app = gio.NewDesktopAppInfoFromFilename(name)
			if app == nil && count == 20 {
				logger.Info("create DesktopAppInfo failed")
				return
			}
		}
		defer app.Unref()
		if !app.ShouldShow() {
			logger.Info(app.GetFilename(), "should NOT show")
			return
		}
		itemTable[id] = &ItemInfo{}
		itemTable[id].init(app)
	}
	if _, ok := itemTable[id]; !ok {
		logger.Info("get item from itemTable failed")
		return
	}
	d.ItemChanged(status, *itemTable[id], itemTable[id].getCategoryId())

	if status == SOFTWARE_STATUS_DELETED {
		itemTable[id].destroy()
		delete(itemTable, id)
	} else {
		cid := itemTable[id].getCategoryId()
		fmt.Printf("add id to category#%d\n", cid)
		categoryTable[cid].items[id] = true
	}
	logger.Info(status, "successful")
}

func (d *LauncherDBus) itemChangedHandler(ev *fsnotify.FileEvent, name string, info map[string]ItemChangedStatus) {
	if _, ok := info[name]; !ok {
		info[name] = ItemChangedStatus{
			make(chan bool),
			make(chan bool),
			make(chan bool),
			make(chan bool),
		}
	}
	if ev.IsRename() {
		logger.Info("renamed")
		select {
		case <-info[name].renamed:
		default:
		}
		go func() {
			select {
			case <-info[name].notRenamed:
				return
			case <-time.After(time.Second):
				<-info[name].renamed
				d.emitItemChanged(name, SOFTWARE_STATUS_DELETED, info)
			}
		}()
		info[name].renamed <- true
	} else if ev.IsCreate() {
		go func() {
			select {
			case <-info[name].renamed:
				// logger.Info("not renamed")
				info[name].notRenamed <- true
				info[name].renamed <- true
			default:
				// logger.Info("default")
			}
			select {
			case <-info[name].notCreated:
				return
			case <-time.After(time.Second):
				<-info[name].created
				logger.Info("create")
				d.emitItemChanged(name, SOFTWARE_STATUS_CREATED, info)
			}
		}()
		info[name].created <- true
	} else if ev.IsModify() && !ev.IsAttrib() {
		go func() {
			select {
			case <-info[name].created:
				info[name].notCreated <- true
			}
			select {
			case <-info[name].renamed:
				d.emitItemChanged(name, SOFTWARE_STATUS_MODIFIED, info)
			default:
				logger.Info("modify created")
				d.emitItemChanged(name, SOFTWARE_STATUS_CREATED, info)
			}
		}()
	} else if ev.IsAttrib() {
		go func() {
			select {
			case <-info[name].renamed:
				<-info[name].created
				info[name].notCreated <- true
			default:
			}
		}()
	} else if ev.IsDelete() {
		d.emitItemChanged(name, SOFTWARE_STATUS_DELETED, info)
	}
}

func (d *LauncherDBus) eventHandler(watcher *fsnotify.Watcher) {
	var info = map[string]ItemChangedStatus{}
	for {
		select {
		case ev := <-watcher.Event:
			name := path.Clean(ev.Name)
			basename := path.Base(name)
			matched, _ := path.Match(`[^#.]*.desktop`, basename)
			if matched {
				d.itemChangedHandler(ev, name, info)
			}
		case <-watcher.Error:
		}
	}
}

func getApplicationsDirs() []string {
	dirs := make([]string, 0)
	dataDirs := glib.GetSystemDataDirs()
	for _, dir := range dataDirs {
		applicationsDir := path.Join(dir, AppDirName)
		if exist(applicationsDir) {
			dirs = append(dirs, applicationsDir)
		}
	}

	userDataDir := path.Join(glib.GetUserDataDir(), AppDirName)
	dirs = append(dirs, userDataDir)
	if !exist(userDataDir) {
		os.MkdirAll(userDataDir, DirDefaultPerm)
	}
	return dirs
}

func (d *LauncherDBus) listenItemChanged() {
	dirs := getApplicationsDirs()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	// FIXME: close watcher.
	for _, dir := range dirs {
		logger.Info("monitor:", dir)
		watcher.Watch(dir)
	}

	go d.eventHandler(watcher)
}

func (d *LauncherDBus) Search(key string) []ItemId {
	return search(key)
}

func (d *LauncherDBus) IsOnDesktop(name string) bool {
	return isOnDesktop(name)
}

func (d *LauncherDBus) SendToDesktop(name string) {
	sendToDesktop(name)
}

func (d *LauncherDBus) GetFavors() FavorItemList {
	return getFavors()
}

func (d *LauncherDBus) SaveFavors(items FavorItemList) bool {
	return saveFavors(items)
}

func (d *LauncherDBus) GetAppId(path string) string {
	return string(genId(path))
}

func initDBus() {
	launcherDbus := LauncherDBus{}
	dbus.InstallOnSession(&launcherDbus)
	launcherDbus.listenItemChanged()
}
