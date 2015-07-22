package mime

import (
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"os"
	"path"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
	"sync"
)

type DefaultApps struct {
	DefaultAppChanged func()

	watcher *dutils.WatchProxy
}

type AppInfo struct {
	ID   string
	Name string
	Exec string
}

const (
	_TERMINAL_SCHEMA = "com.deepin.desktop.default-applications.terminal"
	_DESKTOP_PATH    = "/usr/share/applications/"

	_TERMINAL_EMULATOR   = "TerminalEmulator"
	_DESKTOP_ENTRY       = "Desktop Entry"
	_X_TERMINAL_EMULATOR = "x-terminal-emulator"
	_CATEGORY            = "Categories"
	_EXEC                = "Exec"

	MIME_CACHE_FILE = ".local/share/applications/mimeapps.list"
)

var (
	_TerminalBlacklist = []string{"guake"}
	mimeWatcher        *fsnotify.Watcher
)

var _TerminalGSettings = func() func() *gio.Settings {
	var terminalGSettings *gio.Settings
	var initTerminalGSettings sync.Once

	return func() *gio.Settings {
		initTerminalGSettings.Do(func() {
			terminalGSettings = gio.NewSettings(_TERMINAL_SCHEMA)
		})
		return terminalGSettings
	}
}()

func NewDAppInfo(gioApp *gio.AppInfo) AppInfo {
	dappInfo := AppInfo{}
	if gioApp == nil {
		logger.Info("gioApp is nil in NewDAppInfo")
		return dappInfo
	}

	dappInfo.ID = gioApp.GetId()
	dappInfo.Name = gioApp.GetDisplayName()
	dappInfo.Exec = gioApp.GetExecutable()
	return dappInfo
}

func (dapp *DefaultApps) AppsListViaType(typeName string) []AppInfo {
	var defaultAppsList []AppInfo

	if typeName == "terminal" {
		lists := GetTerminalList()
		if lists == nil {
			return nil
		}

		for _, v := range lists {
			app, ok := NewAppInfoByID(v)
			if !ok {
				continue
			}
			defaultAppsList = append(defaultAppsList,
				app)
		}

		return defaultAppsList
	}

	gioAppsList := gio.AppInfoGetAllForType(typeName)
	for _, gioApp := range gioAppsList {
		defaultAppsList = append(defaultAppsList, NewDAppInfo(gioApp))
	}
	return defaultAppsList
}

func (dapp *DefaultApps) DefaultAppViaType(typeName string) AppInfo {
	if typeName == "terminal" {
		exec := _TerminalGSettings().GetString("exec")
		terminalApps := dapp.AppsListViaType(typeName)

		for _, v := range terminalApps {
			if exec == v.Exec {
				return v
			}
		}

		return AppInfo{}
	}

	gioApp := gio.AppInfoGetDefaultForType(typeName, false)
	return NewDAppInfo(gioApp)
}

func (dapp *DefaultApps) SetDefaultAppViaType(typeName, appID string) bool {
	if typeName == "terminal" {
		appInfo, ok := NewAppInfoByID(appID)
		if !ok {
			return false
		}

		if _TerminalGSettings().SetString("exec", appInfo.Exec) {
			gio.SettingsSync()
			return true
		}
		return false
	}

	gio.AppInfoResetTypeAssociations(typeName)
	gioAppsList := gio.AppInfoGetAllForType(typeName)
	for _, gioApp := range gioAppsList {
		if gioApp.GetId() == appID {
			_, err := gioApp.SetAsDefaultForType(typeName)
			if err != nil {
				logger.Debug(err)
				return false
			}
			break
		}
	}

	return true
}

func (dapp *DefaultApps) handleMimeFileChanged(ev *fsnotify.FileEvent) {
	if ev == nil {
		return
	}

	if ev.IsDelete() {
		if dapp.watcher != nil {
			dapp.watcher.ResetFileListWatch()
		}
	} else {
		dbus.Emit(dapp, "DefaultAppChanged")
	}
}

func (dapp *DefaultApps) destroy() {
	dbus.UnInstallObject(dapp)
	if dapp.watcher != nil {
		dapp.watcher.EndWatch()
	}
}

func NewDefaultApps() *DefaultApps {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error in NewDefaultApps: %v",
				err)
		}
	}()

	dapp := &DefaultApps{}

	dapp.watcher = dutils.NewWatchProxy()
	if dapp.watcher != nil {
		dapp.watcher.SetFileList(getWatchFiles())
		dapp.watcher.SetEventHandler(dapp.handleMimeFileChanged)
		go dapp.watcher.StartWatch()
	}
	dapp.initConfigData()

	return dapp
}

func NewAppInfoByID(id string) (AppInfo, bool) {
	appInfo := AppInfo{}
	keyFile := glib.NewKeyFile()
	defer keyFile.Free()
	lang := GetLocalLang()

	_, err1 := keyFile.LoadFromFile(_DESKTOP_PATH+id, glib.KeyFileFlagsNone)
	if err1 != nil {
		logger.Debug("Load File Failed:", err1)
		return AppInfo{}, false
	}

	name, err := keyFile.GetString(_DESKTOP_ENTRY, "Name["+lang+"]")
	if err != nil {
		name, err = keyFile.GetString(_DESKTOP_ENTRY, "Name")
	}

	exec, err2 := keyFile.GetString(_DESKTOP_ENTRY, _EXEC)
	if err2 != nil {
		logger.Debug("Get Exec Failed:", err2)
		return AppInfo{}, false
	}

	appInfo.ID = id
	appInfo.Name = name
	appInfo.Exec = exec

	return appInfo, true
}

func getWatchFiles() []string {
	homeDir := os.Getenv("HOME")
	if len(homeDir) == 0 {
		return nil
	}

	mimeFile := path.Join(homeDir, MIME_CACHE_FILE)
	if dutils.IsFileExist(mimeFile) {
		return []string{mimeFile}
	}

	fp, err := os.Create(mimeFile)
	if err != nil {
		return nil
	}
	fp.Close()

	return []string{mimeFile}
}

func GetLocalLang() string {
	langStr := os.Getenv("LANG")
	array := strings.Split(langStr, ".")
	return array[0]
}

func GetTerminalList() []string {
	terminalList := []string{}
	entryList, err := GetDesktopEntryList()
	if err != nil {
		logger.Debug("Get Desktop Entry List Failed")
		return nil
	}

	for _, v := range entryList {
		if IsTerminalEmulator(_DESKTOP_PATH + v) {
			terminalList = append(terminalList, v)
		}
	}

	return terminalList
}

func IsTerminalEmulator(fileName string) bool {
	keyFile := glib.NewKeyFile()
	defer keyFile.Free()
	_, err := keyFile.LoadFromFile(fileName, glib.KeyFileFlagsNone)
	if err != nil {
		logger.Debug("KeyFile Load File Failed:", err)
		return false
	}

	categories, err := keyFile.GetString(_DESKTOP_ENTRY, _CATEGORY)
	if err != nil {
		return false
	}

	if strings.Contains(categories, _TERMINAL_EMULATOR) {
		execName, err := keyFile.GetString(_DESKTOP_ENTRY, _EXEC)
		if err != nil {
			return false
		}

		if strings.Contains(execName, _X_TERMINAL_EMULATOR) {
			return false
		}

		for _, v := range _TerminalBlacklist {
			if strings.Contains(execName, v) {
				return false
			}
		}

		return true
	}

	return false
}

func GetDesktopEntryList() ([]string, error) {
	entryList := []string{}

	desktops, err := ioutil.ReadDir(_DESKTOP_PATH)
	if err != nil {
		logger.Debug("Read Dir Failed:", err)
		return nil, err
	}

	for _, fileInfo := range desktops {
		if fileInfo.IsDir() {
			continue
		}

		entryList = append(entryList, fileInfo.Name())
	}

	return entryList, nil
}
