package main

import (
	"fmt"
)

const (
	HiddenAppsConf      string = "launcher/apps.ini"
	HiddenAppsGroupName string = "HiddenApps"
	HiddenAppsKeyName   string = "app_ids"
)

// TODO: add a signal for content changed.

func getHiddenApps() []ItemId {
	file, err := configFile(HiddenAppsConf)
	defer file.Free()
	if err != nil {
		fmt.Println(err)
	}

	ids := make([]ItemId, 0)
	_, list, err := file.GetStringList(HiddenAppsGroupName,
		HiddenAppsKeyName)
	if err != nil {
		return ids
	}
	for _, id := range uniqueStringList(list) {
		ids = append(ids, ItemId(id))
	}
	return ids
}

func saveHiddenApps(ids []string) bool {
	file, err := configFile(HiddenAppsConf)
	defer file.Free()
	if err != nil {
		fmt.Println(fmt.Errorf("saveHiddenApps: %s", err))
		return false
	}
	file.SetStringList(HiddenAppsGroupName, HiddenAppsKeyName,
		uniqueStringList(ids))
	if err = saveKeyFile(file, configFilePath(HiddenAppsConf)); err != nil {
		fmt.Println(fmt.Errorf("saveHiddenApps: %s", err))
	}
	return true
}
