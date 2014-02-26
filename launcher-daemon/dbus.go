package main

import (
	"dbus/com/deepin/api/graph"
	"fmt"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/howeyc/fsnotify"

	"dlib/dbus"
	"dlib/gio-2.0"
	"dlib/glib-2.0"
)

const (
	launcherObject            string = "com.deepin.dde.daemon.Launcher"
	launcherPath              string = "/com/deepin/dde/daemon/Launcher"
	launcherInterface         string = launcherObject
	launcherCategoryInterface string = launcherObject + ".category"
	launcherConfigInterface   string = launcherObject + ".config"

	AppDirName     string      = "applications"
	DirDefaultPerm os.FileMode = 775

	DefaultBackgroundImage string = "/usr/share/background/default_background.jpg"
)

type ItemChangedStatus struct {
	renamed, created, notRenamed, notCreated chan bool
}

type LauncherDBus struct {
	ItemChanged func(
		status string,
		itemInfo ItemInfo,
		categoryIds []CategoryId,
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

	fmt.Println("Status:", status)
	if status != SOFTWARE_STATUS_DELETED {
		app := gio.NewDesktopAppInfoFromFilename(name)
		if app == nil {
			fmt.Println("create DesktopAppInfo failed")
			return
		}
		defer app.Unref()
		if !app.ShouldShow() {
			fmt.Println(app.GetFilename(), "should NOT show")
			return
		}
		itemTable[id] = &ItemInfo{}
		itemTable[id].init(app)
	}
	if _, ok := itemTable[id]; !ok {
		fmt.Println("get item from itemTable failed")
		return
	}
	d.ItemChanged(status, *itemTable[id], itemTable[id].getCategoryIds())

	if status == SOFTWARE_STATUS_DELETED {
		itemTable[id].destroy()
		delete(itemTable, id)
	} else {
		for _, cid := range itemTable[id].getCategoryIds() {
			fmt.Printf("add id to category#%d\n", cid)
			categoryTable[cid].items[id] = true
		}
	}
	fmt.Println(status, "successful")
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
				// fmt.Println("not renamed")
				info[name].notRenamed <- true
				info[name].renamed <- true
			default:
				// fmt.Println("default")
			}
			select {
			case <-info[name].notCreated:
				return
			case <-time.After(time.Second):
				<-info[name].created
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
		fmt.Println("monitor:", dir)
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

func (d *LauncherDBus) LoadHiddenApps() []ItemId {
	return getHiddenApps()
}

func (d *LauncherDBus) SaveHiddenApps(ids []string) bool {
	return saveHiddenApps(ids)
}

func (d *LauncherDBus) GetFavors() FavorItemList {
	return getFavors()
}

func (d *LauncherDBus) SaveFavors(items FavorItemList) bool {
	return saveFavors(items)
}

func (d *LauncherDBus) GetPackageNames(path string) []string {
	return getPackageNames(path)
}

func (d *LauncherDBus) GetBackgroundPict() string {
	errorHandler := func(e error) string {
		fmt.Println(e)
		return DefaultBackgroundImage
	}

	i, err := graph.NewGraph("/com/deepin/api/Graph")
	if err != nil {
		return errorHandler(err)
	}

	u, err := user.Current()
	if err != nil {
		return errorHandler(err)
	}

	// use "" to get current pict
	res, err := i.BackgroundBlurPictPath(u.Uid, "")
	if err != nil {
		return errorHandler(err)
	}

	fmt.Println(res)
	blurPath := res[1].(string)
	return blurPath
}

func (d *LauncherDBus) GetAppId(path string) string {
	return string(genId(path))
}

func initDBus() {
	launcherDbus := LauncherDBus{}
	dbus.InstallOnSession(&launcherDbus)
	launcherDbus.listenItemChanged()
}
