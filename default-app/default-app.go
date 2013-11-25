package main

import (
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
