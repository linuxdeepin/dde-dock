package launcher

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/howeyc/fsnotify"

	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/item"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/item/search"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/utils"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/glib-2.0"
	"pkg.linuxdeepin.com/lib/utils"
)

const (
	launcherObject    string = "com.deepin.dde.daemon.Launcher"
	launcherPath      string = "/com/deepin/dde/daemon/Launcher"
	launcherInterface string = launcherObject

	AppDirName string = "applications"

	SoftwareStatusCreated  string = "created"
	SoftwareStatusModified string = "updated"
	SoftwareStatusDeleted  string = "deleted"
)

type ItemChangedStatus struct {
	renamed, created, notRenamed, notCreated chan bool
}

type Launcher struct {
	setting             SettingInterface
	itemManager         ItemManagerInterface
	categoryManager     CategoryManagerInterface
	cancelSearchingChan chan struct{}
	pinyinObj           PinYinInterface

	ItemChanged func(
		status string,
		itemInfo ItemInfoExport,
		categoryId CategoryId,
	)
	UninstallSuccess func(ItemId)
	UninstallFailed  func(ItemId, string)

	SendToDesktopSuccess func(ItemId)
	SendToDesktopFailed  func(ItemId, string)

	RemoveFromDesktopSuccess func(ItemId)
	RemoveFromDesktopFailed  func(ItemId, string)

	SearchDone func([]ItemId)
}

func NewLauncher() *Launcher {
	launcher = &Launcher{
		cancelSearchingChan: make(chan struct{}),
	}
	return launcher
}

func (self *Launcher) setSetting(s SettingInterface) {
	self.setting = s
}

func (self *Launcher) setCategoryManager(cm CategoryManagerInterface) {
	self.categoryManager = cm
}

func (self *Launcher) setItemManager(im ItemManagerInterface) {
	self.itemManager = im
}

func (self *Launcher) setPinYinObject(pinyinObj PinYinInterface) {
	self.pinyinObj = pinyinObj
}

func (self *Launcher) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		launcherObject,
		launcherPath,
		launcherInterface,
	}
}

func (self *Launcher) RequestUninstall(id string, purge bool) {
	go func(id ItemId) {
		logger.Warning("uninstall", id)
		err := self.itemManager.UninstallItem(id, purge, time.Minute*20)
		if err == nil {
			dbus.Emit(self, "UninstallSuccess", id)
			return
		}

		dbus.Emit(self, "UninstallFailed", id, err.Error())
	}(ItemId(id))
}

func (self *Launcher) RequestSendToDesktop(id string) bool {
	itemId := ItemId(id)
	if filepath.IsAbs(id) {
		dbus.Emit(self, "SendToDesktopFailed", itemId, "app id is expected")
		return false
	}

	if err := self.itemManager.SendItemToDesktop(itemId); err != nil {
		dbus.Emit(self, "SendToDesktopFailed", itemId, err.Error())
		return false
	}

	dbus.Emit(self, "SendToDesktopSuccess", itemId)
	return true
}

func (self *Launcher) RequestRemoveFromDesktop(id string) bool {
	itemId := ItemId(id)
	if filepath.IsAbs(id) {
		dbus.Emit(self, "RemoveFromDesktopFailed", itemId, "app id is expected")
		return false
	}

	if err := self.itemManager.RemoveItemFromDesktop(itemId); err != nil {
		dbus.Emit(self, "RemoveFromDesktopFailed", itemId, err.Error())
		return false
	}

	dbus.Emit(self, "RemoveFromDesktopSuccess", itemId)
	return true
}

func (self *Launcher) IsItemOnDesktop(id string) bool {
	itemId := ItemId(id)
	if filepath.IsAbs(id) {
		return false
	}

	return self.itemManager.IsItemOnDesktop(itemId)
}

func (self *Launcher) GetCategoryInfo(cid int64) CategoryInfoExport {
	return NewCategoryInfoExport(self.categoryManager.GetCategory(CategoryId(cid)))
}

func (self *Launcher) GetAllCategoryInfos() []CategoryInfoExport {
	infos := []CategoryInfoExport{}
	ids := self.categoryManager.GetAllCategory()
	for _, id := range ids {
		infos = append(infos, NewCategoryInfoExport(self.categoryManager.GetCategory(id)))
	}

	return infos
}

func (self *Launcher) GetItemInfo(id string) ItemInfoExport {
	return NewItemInfoExport(self.itemManager.GetItem(ItemId(id)))
}

func (self *Launcher) GetAllItemInfos() []ItemInfoExport {
	items := self.itemManager.GetAllItems()
	infos := []ItemInfoExport{}
	for _, item := range items {
		infos = append(infos, NewItemInfoExport(item))
	}

	return infos
}

func (self *Launcher) emitItemChanged(name, status string, info map[string]ItemChangedStatus) {
	defer delete(info, name)
	id := GenId(name)

	logger.Info(name, "Status:", status)
	if status != SoftwareStatusDeleted {
		logger.Info(name)
		<-time.After(time.Second * 10)
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
		itemInfo := NewItem(app)
		self.itemManager.AddItem(itemInfo)
		self.categoryManager.AddItem(itemInfo.Id(), itemInfo.GetCategoryId())
	}
	if !self.itemManager.HasItem(id) {
		logger.Info("get item failed")
		return
	}

	item := self.itemManager.GetItem(id)
	logger.Info("emit ItemChanged signal")
	dbus.Emit(self, "ItemChanged", status, NewItemInfoExport(item), item.GetCategoryId())

	cid := self.itemManager.GetItem(id).GetCategoryId()
	if status == SoftwareStatusDeleted {
		self.categoryManager.RemoveItem(id, cid)
		self.itemManager.RemoveItem(id)
	} else {
		self.categoryManager.AddItem(id, cid)
	}
	logger.Info(name, status, "successful")
}

func (self *Launcher) itemChangedHandler(ev *fsnotify.FileEvent, name string, info map[string]ItemChangedStatus) {
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
				self.emitItemChanged(name, SoftwareStatusDeleted, info)
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
				self.emitItemChanged(name, SoftwareStatusCreated, info)
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
				self.emitItemChanged(name, SoftwareStatusModified, info)
			default:
				logger.Info("modify created")
				self.emitItemChanged(name, SoftwareStatusCreated, info)
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
		self.emitItemChanged(name, SoftwareStatusDeleted, info)
	}
}

func (self *Launcher) eventHandler(watcher *fsnotify.Watcher) {
	var info = map[string]ItemChangedStatus{}
	for {
		select {
		case ev := <-watcher.Event:
			name := path.Clean(ev.Name)
			basename := path.Base(name)
			matched, _ := path.Match(`[^#.]*.desktop`, basename)
			if basename == "kde4" {
				if ev.IsCreate() {
					watcher.Watch(name)
				} else if ev.IsDelete() {
					watcher.RemoveWatch(name)
				}
			}
			if matched {
				self.itemChangedHandler(ev, name, info)
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
		if utils.IsFileExist(applicationsDir) {
			dirs = append(dirs, applicationsDir)
		}
		applicationsDirForKde := path.Join(applicationsDir, "kde4")
		if utils.IsFileExist(applicationsDirForKde) {
			dirs = append(dirs, applicationsDirForKde)
		}
	}

	userDataDir := path.Join(glib.GetUserDataDir(), AppDirName)
	dirs = append(dirs, userDataDir)
	if !utils.IsFileExist(userDataDir) {
		os.MkdirAll(userDataDir, DirDefaultPerm)
	}
	userDataDirForKde := path.Join(userDataDir, "kde4")
	if utils.IsFileExist(userDataDirForKde) {
		dirs = append(dirs, userDataDirForKde)
	}
	return dirs
}

func (self *Launcher) listenItemChanged() {
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

	go self.eventHandler(watcher)
}

func (self *Launcher) RecordRate(id string) {
	f, err := GetFrequencyRecordFile()
	if err != nil {
		logger.Warning("Open frequency record file failed:", err)
		return
	}
	defer f.Free()
	self.itemManager.SetRate(ItemId(id), self.itemManager.GetRate(ItemId(id), f)+1, f)
}

func (self *Launcher) GetAllFrequency() (infos []FrequencyExport) {
	f, err := GetFrequencyRecordFile()
	frequency := self.itemManager.GetAllFrequency(f)

	for id, rate := range frequency {
		infos = append(infos, FrequencyExport{Frequency: rate, Id: id})
	}

	if err != nil {
		logger.Warning("Open frequency record file failed:", err)
		return
	}
	f.Free()

	return
}

func (self *Launcher) GetAllTimeInstalled() []TimeInstalledExport {
	infos := []TimeInstalledExport{}
	for id, t := range self.itemManager.GetAllTimeInstalled() {
		infos = append(infos, TimeInstalledExport{Time: t, Id: id})
	}

	return infos
}

func (self *Launcher) Search(key string) {
	close(self.cancelSearchingChan)
	self.cancelSearchingChan = make(chan struct{})
	go func() {
		resultChan := make(chan SearchResult)
		transaction, err := NewSearchTransaction(self.pinyinObj, resultChan, self.cancelSearchingChan, 0)
		if err != nil {
			return
		}

		dataSet := self.itemManager.GetAllItems()
		go func() {
			transaction.Search(key, dataSet)
			close(resultChan)
		}()

		select {
		case <-self.cancelSearchingChan:
			return
		default:
			resultMap := map[ItemId]SearchResult{}
			for result := range resultChan {
				resultMap[result.Id] = result
			}

			var res SearchResultList
			for _, data := range resultMap {
				res = append(res, data)
			}

			sort.Sort(res)

			itemIds := []ItemId{}
			for _, data := range res {
				itemIds = append(itemIds, data.Id)
			}
			dbus.Emit(self, "SearchDone", itemIds)
		}
	}()
}

func (self *Launcher) destroy() {
	if self.setting != nil {
		self.setting.destroy()
		launcher.setting = nil
	}
	dbus.UnInstallObject(self)
}
