package main

import (
	"dlib/gio-2.0"
	"fmt"
)

type DefaultApps struct{}

type AppInfo struct {
	ID          string
	Name        string
	DisplayName string
	Desc        string
	Exec        string
	Command     string
}

func NewDAppInfo(gioApp *gio.AppInfo) *AppInfo {
	dappInfo := AppInfo{}

	dappInfo.ID = gioApp.GetId()
	dappInfo.Name = gioApp.GetName()
	dappInfo.DisplayName = gioApp.GetDisplayName()
	dappInfo.Desc = gioApp.GetDescription()
	dappInfo.Exec = gioApp.GetExecutable()
	dappInfo.Command = gioApp.GetCommandline()
	return &dappInfo
}

func (dapp *DefaultApps) AppsListViaType(typeName string) []*AppInfo {
	var defaultAppsList []*AppInfo
	gioAppsList := gio.AppInfoGetAllForType(typeName)
	for _, gioApp := range gioAppsList {
		defaultAppsList = append(defaultAppsList, NewDAppInfo(gioApp))
	}
	return defaultAppsList
}

func (dapp *DefaultApps) DefaultAppViaType(typeName string) *AppInfo {
	gioApp := gio.AppInfoGetDefaultForType(typeName, false)
	return NewDAppInfo(gioApp)
}

func (dapp *DefaultApps) SetDefaultAppViaType(typeName, appID string) bool {
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
