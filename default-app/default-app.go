package main

import (
	"dlib/gio-2.0"
	"dlib/glib-2.0"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type DefaultApps struct{}

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
)

var (
	_TerminalBlacklist = []string{"guake"}

	_TerminalGSettings = gio.NewSettings(_TERMINAL_SCHEMA)
)

func NewDAppInfo(gioApp *gio.AppInfo) *AppInfo {
	dappInfo := AppInfo{}

	dappInfo.ID = gioApp.GetId()
	dappInfo.Name = gioApp.GetDisplayName()
	dappInfo.Exec = gioApp.GetExecutable()
	return &dappInfo
}

func (dapp *DefaultApps) AppsListViaType(typeName string) []*AppInfo {
	var defaultAppsList []*AppInfo

	if typeName == "terminal" {
		lists := GetTerminalList()
		if lists == nil {
			return nil
		}

		for _, v := range lists {
			defaultAppsList = append(defaultAppsList,
				NewAppInfoByID(v))
		}

		return defaultAppsList
	}

	gioAppsList := gio.AppInfoGetAllForType(typeName)
	for _, gioApp := range gioAppsList {
		defaultAppsList = append(defaultAppsList, NewDAppInfo(gioApp))
	}
	return defaultAppsList
}

func (dapp *DefaultApps) DefaultAppViaType(typeName string) *AppInfo {
	if typeName == "terminal" {
		exec := _TerminalGSettings.GetString("exec")
		terminalApps := dapp.AppsListViaType(typeName)

		for _, v := range terminalApps {
			if exec == v.Exec {
				return v
			}
		}

		return nil
	}

	gioApp := gio.AppInfoGetDefaultForType(typeName, false)
	return NewDAppInfo(gioApp)
}

func (dapp *DefaultApps) SetDefaultAppViaType(typeName, appID string) bool {
	if typeName == "terminal" {
		appInfo := NewAppInfoByID(appID)
		if appInfo == nil {
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
				fmt.Println(err)
				return false
			}
			break
		}
	}

	return true
}

func NewAppInfoByID(id string) *AppInfo {
	appInfo := &AppInfo{}
	keyFile := glib.NewKeyFile()
	lang := GetLocalLang()

	_, err1 := keyFile.LoadFromFile(_DESKTOP_PATH+id, glib.KeyFileFlagsNone)
	if err1 != nil {
		fmt.Println("Load File Failed:", err1)
		return nil
	}

	name, err := keyFile.GetString(_DESKTOP_ENTRY, "Name["+lang+"]")
	if err != nil {
		name, err = keyFile.GetString(_DESKTOP_ENTRY, "Name")
	}

	exec, err2 := keyFile.GetString(_DESKTOP_ENTRY, _EXEC)
	if err2 != nil {
		fmt.Println("Get Exec Failed:", err2)
		return nil
	}

	appInfo.ID = id
	appInfo.Name = name
	appInfo.Exec = exec

	return appInfo
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
		fmt.Println("Get Desktop Entry List Failed")
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
	_, err := keyFile.LoadFromFile(fileName, glib.KeyFileFlagsNone)
	if err != nil {
		fmt.Println("KeyFile Load File Failed:", err)
		return false
	}

	categories, err := keyFile.GetString(_DESKTOP_ENTRY, _CATEGORY)
	if err != nil {
		fmt.Println("KeyFile Get String Failed:", err)
		return false
	}

	if strings.Contains(categories, _TERMINAL_EMULATOR) {
		execName, err := keyFile.GetString(_DESKTOP_ENTRY, _EXEC)
		if err != nil {
			fmt.Println("KeyFile Get String Failed:", err)
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
		fmt.Println("Read Dir Failed:", err)
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
