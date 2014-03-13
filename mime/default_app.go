package main

import (
        "dlib/gio-2.0"
        "dlib/glib-2.0"
        "github.com/howeyc/fsnotify"
        "io/ioutil"
        "os"
        "os/user"
        "strings"
)

type DefaultApps struct {
        DefaultAppChanged func()
}

type AppInfo struct {
        ID      string
        Name    string
        Exec    string
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

        _TerminalGSettings = gio.NewSettings(_TERMINAL_SCHEMA)
        mimeWatcher        *fsnotify.Watcher
)

func NewDAppInfo(gioApp *gio.AppInfo) AppInfo {
        dappInfo := AppInfo{}
        if gioApp == nil {
                logObject.Info("gioApp is nil in NewDAppInfo")
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
                exec := _TerminalGSettings.GetString("exec")
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

                if _TerminalGSettings.SetString("exec", appInfo.Exec) {
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
                                logObject.Info("%v", err)
                                return false
                        }
                        break
                }
        }

        return true
}

func (dapp *DefaultApps) listenMimeCacheFile() {
        var (
                err      error
                userInfo *user.User
        )

        mimeWatcher, err = fsnotify.NewWatcher()
        if err != nil {
                logObject.Info("New Watcher Failed: %v\n", err)
                panic(err)
        }

        userInfo, err = user.Current()
        if err != nil {
                logObject.Info("Get current user failed: %v\n", err)
                panic(err)
        }

        mimeFile := userInfo.HomeDir + "/" + MIME_CACHE_FILE
        // TODO: if mimeFile not exist, touch it
        err = mimeWatcher.Watch(mimeFile)
        if err != nil {
                logObject.Info("Watch '%s' Failed: %s\n",
                        MIME_CACHE_FILE, err)
                panic(err)
        }

        go func() {
                defer mimeWatcher.Close()
                for {
                        select {
                        case ev := <-mimeWatcher.Event:
                                logObject.Info("Watch Event: %v\n", ev)
                                if ev.IsDelete() {
                                        mimeWatcher.Watch(mimeFile)
                                } else {
                                        dapp.DefaultAppChanged()
                                }
                        case err := <-mimeWatcher.Error:
                                logObject.Info("Watch Error: %v\n", err)
                        }
                }
        }()
}

func NewDefaultApps() *DefaultApps {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Info("Recover Error in NewDefaultApps: %v\n",
                                err)
                }
        }()

        dapp := &DefaultApps{}

        var err error
        mimeWatcher, err = fsnotify.NewWatcher()
        if err != nil {
                logObject.Info("Create mime file watcher failed: %v\n", err)
                panic(err)
        }

        dapp.listenMimeCacheFile()

        return dapp
}

func NewAppInfoByID(id string) (AppInfo, bool) {
        appInfo := AppInfo{}
        keyFile := glib.NewKeyFile()
        defer keyFile.Free()
        lang := GetLocalLang()

        _, err1 := keyFile.LoadFromFile(_DESKTOP_PATH+id, glib.KeyFileFlagsNone)
        if err1 != nil {
                logObject.Info("Load File Failed: %v\n", err1)
                return AppInfo{}, false
        }

        name, err := keyFile.GetString(_DESKTOP_ENTRY, "Name["+lang+"]")
        if err != nil {
                name, err = keyFile.GetString(_DESKTOP_ENTRY, "Name")
        }

        exec, err2 := keyFile.GetString(_DESKTOP_ENTRY, _EXEC)
        if err2 != nil {
                logObject.Info("Get Exec Failed: %v\n", err2)
                return AppInfo{}, false
        }

        appInfo.ID = id
        appInfo.Name = name
        appInfo.Exec = exec

        return appInfo, true
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
                logObject.Info("Get Desktop Entry List Failed")
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
                logObject.Info("KeyFile Load File Failed: %v\n", err)
                return false
        }

        categories, err := keyFile.GetString(_DESKTOP_ENTRY, _CATEGORY)
        if err != nil {
                logObject.Info("KeyFile Get String Failed: %v\n", err)
                return false
        }

        if strings.Contains(categories, _TERMINAL_EMULATOR) {
                execName, err := keyFile.GetString(_DESKTOP_ENTRY, _EXEC)
                if err != nil {
                        logObject.Info("KeyFile Get String Failed: %v\n",
                                err)
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
                logObject.Info("Read Dir Failed: %v\n", err)
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
