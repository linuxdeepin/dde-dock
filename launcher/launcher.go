package launcher

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/howeyc/fsnotify"

	storeApi "dbus/com/deepin/store/api"

	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/dde/daemon/appinfo"
	"pkg.deepin.io/dde/daemon/launcher/category"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/dde/daemon/launcher/item"
	"pkg.deepin.io/dde/daemon/launcher/item/search"
	. "pkg.deepin.io/dde/daemon/launcher/utils"
	"pkg.deepin.io/lib/dbus"
	"gir/glib-2.0"
	"pkg.deepin.io/lib/utils"
)

const (
	launcherObject    string = "com.deepin.dde.daemon.Launcher"
	launcherPath      string = "/com/deepin/dde/daemon/Launcher"
	launcherInterface string = launcherObject

	AppDirName string = "applications"

	AppStatusCreated  string = "created"
	AppStatusModified string = "updated"
	AppStatusDeleted  string = "deleted"
)

// ItemChangedStatus stores item's current changed status.
type ItemChangedStatus struct {
	renamed, created, notRenamed, notCreated chan bool
	timeInstalled                            int64
}

// Launcher 为launcher的后端。
type Launcher struct {
	setting             Setting
	itemManager         ItemManager
	categoryManager     CategoryManager
	cancelMutex         sync.Mutex
	cancelSearchingChan chan struct{}
	pinyinObj           PinYin
	store               *storeApi.DStoreDesktop
	appMonitor          *fsnotify.Watcher

	// ItemChanged当launcher中的item有改变后触发。
	ItemChanged func(
		status string,
		itemInfo ItemInfoExport,
		categoryID CategoryID,
	)
	// UninstallSuccess在卸载程序成功后触发。
	UninstallSuccess func(ItemID)
	// UninstallFailed在卸载程序失败后触发。
	UninstallFailed func(ItemID, string)

	// SendToDesktopSuccess在发送到桌面成功后触发。
	SendToDesktopSuccess func(ItemID)
	// SendToDesktopFailed在发送到桌面失败后触发。
	SendToDesktopFailed func(ItemID, string)

	// RemoveFromDesktopSuccess在从桌面移除成功后触发。
	RemoveFromDesktopSuccess func(ItemID)
	// RemoveFromDesktopFailed在从桌面移除失败后触发。
	RemoveFromDesktopFailed func(ItemID, string)

	// SearchDone在搜索结束后触发。
	SearchDone func([]ItemID)

	// NewAppLaunched在新安装程序被标记为已启动后被触发。（废弃，不够直观，使用新信号NewAppMarkedAsLaunched）
	NewAppLaunched func(ItemID)
	// NewAppMarkedAsLaunched在新安装程序被标记为已启动后被触发。
	NewAppMarkedAsLaunched func(ItemID)
}

// NewLauncher creates a new launcher object.
func NewLauncher() *Launcher {
	launcher := &Launcher{
		cancelSearchingChan: make(chan struct{}),
	}
	return launcher
}

func (self *Launcher) setSetting(s Setting) {
	self.setting = s
}

func (self *Launcher) setCategoryManager(cm CategoryManager) {
	self.categoryManager = cm
}

func (self *Launcher) setItemManager(im ItemManager) {
	self.itemManager = im
}

func (self *Launcher) setPinYinObject(pinyinObj PinYin) {
	self.pinyinObj = pinyinObj
}

func (self *Launcher) setStoreAPI(s *storeApi.DStoreDesktop) {
	self.store = s
}

// GetDBusInfo returns launcher's dbus info.
func (self *Launcher) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		launcherObject,
		launcherPath,
		launcherInterface,
	}
}

// RequestUninstall 请求卸载程。
func (self *Launcher) RequestUninstall(id string, purge bool) {
	go func(id ItemID) {
		logger.Info("uninstall", id)
		err := self.itemManager.UninstallItem(id, purge, time.Minute*20)
		if err == nil {
			dbus.Emit(self, "UninstallSuccess", id)
			return
		}

		dbus.Emit(self, "UninstallFailed", id, err.Error())
	}(ItemID(id))
}

// RequestSendToDesktop 请求将程序发送到桌面。
func (self *Launcher) RequestSendToDesktop(id string) bool {
	itemID := ItemID(id)
	if filepath.IsAbs(id) {
		dbus.Emit(self, "SendToDesktopFailed", itemID, "app id is expected")
		return false
	}

	if err := self.itemManager.SendItemToDesktop(itemID); err != nil {
		dbus.Emit(self, "SendToDesktopFailed", itemID, err.Error())
		return false
	}

	soundutils.PlaySystemSound(soundutils.EventIconToDesktop, "", false)
	dbus.Emit(self, "SendToDesktopSuccess", itemID)
	return true
}

// RequestRemoveFromDesktop 请求将程序从桌面移除。
func (self *Launcher) RequestRemoveFromDesktop(id string) bool {
	itemID := ItemID(id)
	if filepath.IsAbs(id) {
		dbus.Emit(self, "RemoveFromDesktopFailed", itemID, "app id is expected")
		return false
	}

	if err := self.itemManager.RemoveItemFromDesktop(itemID); err != nil {
		dbus.Emit(self, "RemoveFromDesktopFailed", itemID, err.Error())
		return false
	}

	dbus.Emit(self, "RemoveFromDesktopSuccess", itemID)
	return true
}

// IsItemOnDesktop 判断程序是否已经在桌面上。
func (self *Launcher) IsItemOnDesktop(id string) bool {
	itemID := ItemID(id)
	if filepath.IsAbs(id) {
		return false
	}

	return self.itemManager.IsItemOnDesktop(itemID)
}

// GetCategoryInfo 获取分类id对应的分类信息。
// 包括：分类名，分类id，分类所包含的程序。
func (self *Launcher) GetCategoryInfo(cid int64) CategoryInfoExport {
	return NewCategoryInfoExport(self.categoryManager.GetCategory(CategoryID(cid)))
}

// GetAllCategoryInfos 获取所有分类的分类信息。
func (self *Launcher) GetAllCategoryInfos() []CategoryInfoExport {
	infos := []CategoryInfoExport{}
	ids := self.categoryManager.GetAllCategory()
	for _, id := range ids {
		infos = append(infos, NewCategoryInfoExport(self.categoryManager.GetCategory(id)))
	}

	return infos
}

// GetItemInfo 获取id对应的item信息。
// 包括：item的path，item的Name，item的id，item的icon，item的分类id，item的安装时间
func (self *Launcher) GetItemInfo(id string) ItemInfoExport {
	return NewItemInfoExport(self.itemManager.GetItem(ItemID(id)))
}

// GetAllItemInfos 获取所有item的信息。
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

	id := item.GenID(name)

	if status == AppStatusCreated && self.itemManager.HasItem(id) {
		status = AppStatusModified
	}
	logger.Info("start emitItemChanged", name, "Status:", status)

	if status != AppStatusDeleted {
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
		itemInfo := item.New(app)
		if info[name].timeInstalled != 0 {
			itemInfo.SetTimeInstalled(info[name].timeInstalled)
		}

		dbPath, _ := category.GetDBPath(category.SoftwareCenterDataDir, category.CategoryNameDBPath)
		self.categoryManager.LoadAppCategoryInfo(dbPath, "")
		defer self.categoryManager.FreeAppCategoryInfo()

		cid, err := self.categoryManager.QueryID(app)
		if err != nil {
			itemInfo.SetCategoryID(category.OthersID)
		}
		itemInfo.SetCategoryID(cid)

		self.itemManager.AddItem(itemInfo)
		self.categoryManager.AddItem(itemInfo.ID(), itemInfo.CategoryID())
	}

	if !self.itemManager.HasItem(id) {
		logger.Warning("has no such a item", id)
		return
	}

	item := self.itemManager.GetItem(id)
	cid := item.CategoryID()
	itemInfo := NewItemInfoExport(item)

	if status == AppStatusDeleted {
		self.itemManager.MarkLaunched(id)
		self.categoryManager.RemoveItem(id, cid)
		self.itemManager.RemoveItem(id)
	} else {
		self.categoryManager.AddItem(id, cid)
	}

	logger.Info("emit ItemChanged signal", status, dbus.Emit(self, "ItemChanged", status, itemInfo, cid))
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
		// logger.Info("renamed")
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
					self.emitItemChanged(name, AppStatusDeleted, info)
				}
			}
		}()
		info[name].renamed <- true
	} else if ev.IsCreate() {
		self.emitItemChanged(name, AppStatusCreated, info)
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
				self.emitItemChanged(name, AppStatusModified, info)
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
			self.emitItemChanged(name, AppStatusDeleted, info)
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
	var dirs []string
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
		self.store.ConnectNewDesktopAdded(func(desktopID string, timeInstalled int32) {
			self.emitItemChanged(desktopID, AppStatusCreated, map[string]ItemChangedStatus{
				desktopID: ItemChangedStatus{
					timeInstalled: int64(timeInstalled),
				},
			})
		})
	}
}

// RecordRate 记录程序的使用频率。（废弃，统一用词，请使用新接口RecordFrequency）
func (self *Launcher) RecordRate(id string) {
}

// RecordFrequency 记录程序的使用频率。
func (self *Launcher) RecordFrequency(id string) {
	self.RecordRate(id)
}

// GetAllFrequency 获取所有的使用频率信息。
// 包括：item的id与使用频率。
func (self *Launcher) GetAllFrequency() (infos []FrequencyExport) {
	f, err := appinfo.GetFrequencyRecordFile()
	if err != nil {
		return
	}

	frequency := self.itemManager.GetAllFrequency(f)
	f.Free()

	for id, rate := range frequency {
		infos = append(infos, FrequencyExport{Frequency: rate, ID: id})
	}

	return
}

// GetAllTimeInstalled 获取所有程序的安装时间。
// 包括：item的id与安装时间。
func (self *Launcher) GetAllTimeInstalled() []TimeInstalledExport {
	infos := []TimeInstalledExport{}
	times, err := self.itemManager.GetAllTimeInstalled()
	if err != nil {
		logger.Warning("GetAllTimeInstalled error:", err)
	}

	for id, t := range times {
		infos = append(infos, TimeInstalledExport{Time: t, ID: id})
	}

	return infos
}

type FreqGetter struct {
	f *glib.KeyFile
}

func NewFreqGetter(f *glib.KeyFile) *FreqGetter {
	getter := &FreqGetter{f: f}
	return getter
}

func (getter *FreqGetter) GetFrequency(id string) uint64 {
	if getter.f == nil {
		return 0
	}
	return appinfo.GetFrequency(id, getter.f)
}

// Search 搜索给定的关键字。
func (self *Launcher) Search(key string) {
	self.cancelMutex.Lock()
	defer self.cancelMutex.Unlock()

	close(self.cancelSearchingChan)
	self.cancelSearchingChan = make(chan struct{})
	go func(cancelChan chan struct{}) {
		resultChan := make(chan search.Result)
		recordFile, _ := appinfo.GetFrequencyRecordFile()
		if recordFile != nil {
			defer recordFile.Free()
		}
		freqGetter := NewFreqGetter(recordFile)
		transaction, err := search.NewTransaction(self.pinyinObj, resultChan, cancelChan, 0)
		if err != nil {
			return
		}

		transaction.SetFreqGetter(freqGetter)
		dataSet := self.itemManager.GetAllItems()
		go func() {
			transaction.Search(key, dataSet)
			close(resultChan)
		}()

		resultMap := map[ItemID]search.Result{}
		for result := range resultChan {
			select {
			case <-cancelChan:
				return
			default:
				resultMap[result.ID] = result
			}
		}

		var res search.ResultList
		for _, data := range resultMap {
			res = append(res, data)
		}

		sort.Sort(res)

		logger.Debug("search result", res)
		itemIDs := []ItemID{}
		for _, data := range res {
			itemIDs = append(itemIDs, data.ID)
		}
		dbus.Emit(self, "SearchDone", itemIDs)
	}(self.cancelSearchingChan)
}

// MarkLaunched 将程序标记为已启动过。
func (self *Launcher) MarkLaunched(id string) {
	err := self.itemManager.MarkLaunched(ItemID(id))
	if err != nil {
		logger.Warning("MarkLaunched error:", err)
		return
	}

	dbus.Emit(self, "NewAppLaunched", ItemID(id))
}

// GetAllNewInstalledApps 获取所有新安装的程序。
func (self *Launcher) GetAllNewInstalledApps() []ItemID {
	ids, err := self.itemManager.GetAllNewInstalledApps()
	if err != nil {
		logger.Info("GetAllNewInstalledApps", err)
	}
	return ids
}

func (self *Launcher) destroy() {
	if self.setting != nil {
		self.setting.Destroy()
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
