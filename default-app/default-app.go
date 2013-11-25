package main

import (
	"dlib/dbus"
	"dlib/gio-2.0"
	"fmt"
)

type DefaultApps struct{}

type DAppInfo struct {
	AppID          string
	AppName        string
	AppDisplayName string
	AppDesc        string
	AppExec        string
	AppCommand     string
}

const (
	_DEFAULT_APPS_DEST = "com.deepin.daemon.DefaultApps"
	_DEFAULT_APPS_PATH = "/com/deepin/daemon/DefaultApps"
	_DEFAULT_APPS_IFC  = "com.deepin.daemon.DefaultApps"

	_HTTP_CONTENT_TYPE     = "x-scheme-handler/http"
	_HTTPS_CONTENT_TYPE    = "x-scheme-handler/https"
	_MAIL_CONTENT_TYPE     = "x-scheme-handler/mailto"
	_CALENDAR_CONTENT_TYPE = "text/calendar"
	_EDITOR_CONTENT_TYPE   = "text/plain"
	_AUDIO_CONTENT_TYPE    = "audio/mpeg"
	_VIDEO_CONTENT_TYPE    = "video/mp4"
)

func NewDAppInfo(gioApp *gio.AppInfo) DAppInfo {
	dappInfo := DAppInfo{}

	dappInfo.AppID = gioApp.GetId()
	dappInfo.AppName = gioApp.GetName()
	dappInfo.AppDisplayName = gioApp.GetDisplayName()
	dappInfo.AppDesc = gioApp.GetDescription()
	dappInfo.AppExec = gioApp.GetExecutable()
	dappInfo.AppCommand = gioApp.GetCommandline()
	return dappInfo
}

func (dapp *DefaultApps) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_DEFAULT_APPS_DEST,
		_DEFAULT_APPS_PATH,
		_DEFAULT_APPS_IFC,
	}
}

func (dapp *DefaultApps) GetAppsListViaType(typeName string) []DAppInfo {
	var defaultAppsList []DAppInfo
	gioAppsList := gio.AppInfoGetAllForType(typeName)
	for _, gioApp := range gioAppsList {
		defaultAppsList = append(defaultAppsList, NewDAppInfo(gioApp))
	}
	return defaultAppsList
}

func (dapp *DefaultApps) GetDefaultAppViaType(typeName string,
	supportUris bool) DAppInfo {
	gioApp := gio.AppInfoGetDefaultForType(typeName, supportUris)
	return NewDAppInfo(gioApp)
}

func (dapp *DefaultApps) SetDefaultAppViaType(typeName,
	appID string) (bool, error) {
	gio.AppInfoResetTypeAssociations(typeName)
	gioAppsList := gio.AppInfoGetAllForType(typeName)
	for _, gioApp := range gioAppsList {
		if gioApp.GetId() == appID {
			success, err := gioApp.SetAsDefaultForType(typeName)
			if err != nil {
				fmt.Println(err)
				return success, err
			}
			break
		}
	}

	return true, nil
}

func main() {
	dapp := DefaultApps{}
	dbus.InstallOnSession(&dapp)
	fmt.Println(dapp.GetAppsListViaType(_HTTP_CONTENT_TYPE))
	fmt.Println(dapp.GetDefaultAppViaType(_HTTP_CONTENT_TYPE, false))
	select {}
}
