package main

import (
	"dlib/dbus"
	"dlib/gio-2.0"
	"fmt"
)

type DefaultApps struct {
	AppID string
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

func (dapp *DefaultApps) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_DEFAULT_APPS_DEST,
		_DEFAULT_APPS_PATH,
		_DEFAULT_APPS_IFC,
	}
}

func (dapp *DefaultApps) GetAppsListViaType(typeName string) []string {
	var defaultAppsList []string
	gioAppsList := gio.AppInfoGetAllForType(typeName)
	for _, gioApp := range gioAppsList {
		defaultAppsList = append(defaultAppsList, gioApp.GetName())
	}
	return defaultAppsList
}

func (dapp *DefaultApps) GetDefaultAppViaType(typeName string,
	supportUris bool) string {
	gioApp := gio.AppInfoGetDefaultForType(typeName, supportUris)
	return gioApp.GetName()
}

func (dapp *DefaultApps) SetDefaultAppViaType(typeName string) bool {
	gio.AppInfoResetTypeAssociations(typeName)
	/*success := gio.*/
	return true
}

func main() {
	dapp := DefaultApps{AppID: "1"}
	dbus.InstallOnSession(&dapp)
	fmt.Println(dapp.GetAppsListViaType(_HTTP_CONTENT_TYPE))
	fmt.Println(dapp.GetDefaultAppViaType(_HTTP_CONTENT_TYPE, false))
	select {}
}
