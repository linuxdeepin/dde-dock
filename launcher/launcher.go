package launcher

import (
	"database/sql"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/howeyc/fsnotify"

	storeApi "dbus/com/deepin/store/api"
	. "pkg.deepin.io/dde-daemon/launcher/category"
	. "pkg.deepin.io/dde-daemon/launcher/interfaces"
	. "pkg.deepin.io/dde-daemon/launcher/item"
	. "pkg.deepin.io/dde-daemon/launcher/item/search"
	. "pkg.deepin.io/dde-daemon/launcher/utils"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/utils"
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
	timeInstalled                            int64
}

// Launcher为launcher的后端。
type Launcher struct {
	setting             SettingInterface
	itemManager         ItemManagerInterface
	categoryManager     CategoryManagerInterface
	cancelSearchingChan chan struct{}
	pinyinObj           PinYinInterface
	store               *storeApi.DStoreDesktop
	appMonitor          *fsnotify.Watcher

	// ItemChanged当launcher中的item有改变后触发。
	ItemChanged func(
		status string,
		itemInfo ItemInfoExport,
		categoryId CategoryId,
	)
	// UninstallSuccess在卸载程序成功后触发。
	UninstallSuccess func(ItemId)
	// UninstallFailed在卸载程序失败后触发。
	UninstallFailed func(ItemId, string)

	// SendToDesktopSuccess在发送到桌面成功后触发。
	SendToDesktopSuccess func(ItemId)
	// SendToDesktopFailed在发送到桌面失败后触发。
	SendToDesktopFailed func(ItemId, string)

	// RemoveFromDesktopSuccess在从桌面移除成功后触发。
	RemoveFromDesktopSuccess func(ItemId)
	// RemoveFromDesktopFailed在从桌面移除失败后触发。
	RemoveFromDesktopFailed func(ItemId, string)

	// SearchDone在搜索结束后触发。
	SearchDone func([]ItemId)

	// NewAppLaunched在新安装程序被标记为已启动后被触发。（废弃，不够直观，使用新信号NewAppMarkedAsLaunched）
	NewAppLaunched func(ItemId)
	// NewAppMarkedAsLaunched在新安装程序被标记为已启动后被触发。
	NewAppMarkedAsLaunched func(ItemId)
}

func NewLauncher() *Launcher {
	launcher := &Launcher{
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

func (self *Launcher) setStoreApi(s *storeApi.DStoreDesktop) {
	self.store = s
}

func (self *Launcher) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		launcherObject,
		launcherPath,
		launcherInterface,
	}
}

// RequestUninstall请求卸载程。
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

// RequestSendToDesktop请求将程序发送到桌面。
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

// RequestRemoveFromDesktop请求将程序从桌面移除。
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

// IsItemOnDesktop判断程序是否已经在桌面上。
func (self *Launcher) IsItemOnDesktop(id string) bool {
	itemId := ItemId(id)
	if filepath.IsAbs(id) {
		return false
	}

	return self.itemManager.IsItemOnDesktop(itemId)
}

// GetCategoryInfo获取分类id对应的分类信息。
// 包括：分类名，分类id，分类所包含的程序。
func (self *Launcher) GetCategoryInfo(cid int64) CategoryInfoExport {
	return NewCategoryInfoExport(self.categoryManager.GetCategory(CategoryId(cid)))
}

// GetAllCategoryInfosu获取所有分类的分类信息。
func (self *Launcher) GetAllCategoryInfos() []CategoryInfoExport {
	infos := []CategoryInfoExport{}
	ids := self.categoryManager.GetAllCategory()
	for _, id := range ids {
		infos = append(infos, NewCategoryInfoExport(self.categoryManager.GetCategory(id)))
	}

	return infos
}

// GetItemInfo获取id对应的item信息。
// 包括：item的path，item的Name，item的id，item的icon，item的分类id，item的安装时间
func (self *Launcher) GetItemInfo(id string) ItemInfoExport {
	return NewItemInfoExport(self.itemManager.GetItem(ItemId(id)))
}

// GetAllItemInfos获取所有item的信息。
func (self *Launcher) GetAllItemInfos() []ItemInfoExport {
	items := self.itemManager.GetAllItems()
	infos := []ItemInfoExport{}
	for _, item := range items {
		infos = append(infos, NewItemInfoExport(item))
	}

	return infos
}

func (self *Launcher) emitItemChanged(name, status string, info map[string]ItemChangedStatus) {
	if info != nil {
		defer delete(info, name)
	}

	id := GenId(name)

	if status == SoftwareStatusCreated && self.itemManager.HasItem(id) {
		status = SoftwareStatusModified
	}
	logger.Info(name, "Status:", status)

	if status != SoftwareStatusDeleted {
		// cannot use float number here. the total wait time is about 12s.
		maxDuration := time.Second + time.Second/2
		waitDuration := time.Millisecond * 0
		deltaDuration := time.Millisecond * 100

		app := CreateDesktopAppInfo(name)
		for app == nil && waitDuration < maxDuration {
			<-time.After(waitDuration)
			app = CreateDesktopAppInfo(name)
			waitDuration += deltaDuration
		}

		if app == nil {
			logger.Infof("create DesktopAppInfo for %q failed", name)
			return
		}
		defer app.Unref()
		if !app.ShouldShow() {
			logger.Info(app.GetFilename(), "should NOT show")
			return
		}
		itemInfo := NewItem(app)
		if info[name].timeInstalled != 0 {
			itemInfo.SetTimeInstalled(info[name].timeInstalled)
		}

		dbPath, _ := GetDBPath(SoftwareCenterDataDir, CategoryNameDBPath)
		db, err := sql.Open("sqlite3", dbPath)
		if err == nil {
			defer db.Close()
			cid, err := QueryCategoryId(app, db)
			if err != nil {
				itemInfo.SetCategoryId(OtherID)
			}
			itemInfo.SetCategoryId(cid)
		}

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
		self.itemManager.MarkLaunched(id)
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
			0,
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
				if true {
					self.emitItemChanged(name, SoftwareStatusDeleted, info)
				}
			}
		}()
		info[name].renamed <- true
	} else if ev.IsCreate() {
		self.emitItemChanged(name, SoftwareStatusCreated, info)
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
		if true {
			self.emitItemChanged(name, SoftwareStatusDeleted, info)
		}
	}
}

func (self *Launcher) eventHandler(watcher *fsnotify.Watcher) {
	var info = map[string]ItemChangedStatus{}
	for {
		select {
		case ev := <-watcher.Event:
			name := path.Clean(ev.Name)
			basename := path.Base(name)
			if basename == "kde4" {
				if ev.IsCreate() {
					watcher.Watch(name)
				} else if ev.IsDelete() {
					watcher.RemoveWatch(name)
				}
			}
			matched, _ := path.Match(`[^#.]*.desktop`, basename)
			if matched {
				go self.itemChangedHandler(ev, name, info)
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

	self.appMonitor = watcher
	for _, dir := range dirs {
		logger.Info("monitor:", dir)
		watcher.Watch(dir)
	}

	go self.eventHandler(watcher)

	if self.store != nil {
		self.store.ConnectNewDesktopAdded(func(desktopId string, timeInstalled int32) {
			self.emitItemChanged(desktopId, SoftwareStatusCreated, map[string]ItemChangedStatus{
				desktopId: ItemChangedStatus{
					timeInstalled: int64(timeInstalled),
				},
			})
		})
	}
}

// RecordRate记录程序的使用频率。（废弃，统一用词，请使用新接口RecordFrequency）
func (self *Launcher) RecordRate(id string) {
	f, err := GetFrequencyRecordFile()
	if err != nil {
		logger.Warning("Open frequency record file failed:", err)
		return
	}
	defer f.Free()
	self.itemManager.SetRate(ItemId(id), self.itemManager.GetRate(ItemId(id), f)+1, f)
}

// RecordFrequency记录程序的使用频率。
func (self *Launcher) RecordFrequency(id string) {
	self.RecordRate(id)
}

// GetAllFrequency获取所有的使用频率信息。
// 包括：item的id与使用频率。
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

// GetAllTimeInstalled获取所有程序的安装时间。
// 包括：item的id与安装时间。
func (self *Launcher) GetAllTimeInstalled() []TimeInstalledExport {
	infos := []TimeInstalledExport{}
	times, err := self.itemManager.GetAllTimeInstalled()
	if err != nil {
		logger.Info(err)
	}

	for id, t := range times {
		infos = append(infos, TimeInstalledExport{Time: t, Id: id})
	}

	return infos
}

// Search搜索给定的关键字。
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

// MarkLaunched将程序标记为已启动过。
func (self *Launcher) MarkLaunched(id string) {
	err := self.itemManager.MarkLaunched(ItemId(id))
	if err != nil {
		logger.Info(err)
		return
	}

	dbus.Emit(self, "NewAppLaunched", ItemId(id))
}

// GetAllNewInstalledApps获取所有新安装的程序。
func (self *Launcher) GetAllNewInstalledApps() []ItemId {
	ids, err := self.itemManager.GetAllNewInstalledApps()
	if err != nil {
		logger.Info("GetAllNewInstalledApps", err)
	}
	return ids
}

func (self *Launcher) destroy() {
	if self.setting != nil {
		self.setting.destroy()
		self.setting = nil
	}
	if self.store != nil {
		storeApi.DestroyDStoreDesktop(self.store)
		self.store = nil
	}
	if self.appMonitor != nil {
		self.appMonitor.Close()
		self.appMonitor = nil
	}
	dbus.UnInstallObject(self)
}
