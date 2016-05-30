/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"gir/gio-2.0"
	"github.com/BurntSushi/xgbutil/ewmh"
	. "pkg.deepin.io/lib/gettext"
)

func (entry *AppEntry) updateMenu() {
	logger.Debug("Update menu")
	menu := NewMenu()
	menu.AppendItem(entry.getMenuItemLaunch())

	desktopActionMenuItems := entry.getMenuItemDesktopActions()
	menu.AppendItem(desktopActionMenuItems...)

	if entry.isActive() {
		menu.AppendItem(entry.getMenuItemCloseAll())
	}

	// menu item dock or undock
	logger.Info(entry.Id, "Item docked?", entry.isDocked)
	if entry.isDocked {
		menu.AppendItem(entry.getMenuItemUndock())
	} else {
		menu.AppendItem(entry.getMenuItemDock())
	}

	entry.coreMenu = menu
	menuJSON := menu.GenerateJSON()
	entry.setData(FieldMenu, menuJSON)
}

func (entry *AppEntry) getMenuItemDesktopActions() []*MenuItem {
	ai := entry.appInfo
	if ai == nil {
		return nil
	}

	var menuItems []*MenuItem
	for _, actionName := range ai.ListActions() {
		//NOTE: don't directly use 'actionName' with closure in an forloop
		actionNameCopy := actionName
		menuItem := NewMenuItem(
			ai.GetActionName(actionName),
			func(timestamp uint32) {
				logger.Debug("desktop app info launch action:", actionNameCopy)
				ai.LaunchAction(actionNameCopy,
					gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
			}, true)
		menuItems = append(menuItems, menuItem)
	}
	return menuItems
}

func (entry *AppEntry) launchApp(timestamp uint32) {
	var appInfo *gio.AppInfo
	if entry.appInfo.DesktopAppInfo != nil {
		logger.Debug("Has AppInfo")
		appInfo = (*gio.AppInfo)(entry.appInfo.DesktopAppInfo)
	} else {
		logger.Debug("No AppInfo", entry.exec)
		var err error = nil
		appInfo, err = gio.AppInfoCreateFromCommandline(
			entry.exec,
			"",
			gio.AppInfoCreateFlagsNone,
		)
		if err != nil {
			logger.Warning("Launch App Falied: ", err)
			return
		}

		defer appInfo.Unref()
	}

	if appInfo == nil {
		logger.Warning("create app info to run program failed")
		return
	}

	_, err := appInfo.Launch(
		make([]*gio.File, 0), gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
	if err != nil {
		logger.Warning("Launch App Failed: ", err)
	}
}

func (entry *AppEntry) getMenuItemLaunch() *MenuItem {
	var itemName string
	if entry.isActive() {
		itemName = entry.getDisplayName()
	} else {
		itemName = Tr("_Run")
	}
	logger.Debugf("getMenuItemLaunch, itemName: %q", itemName)
	return NewMenuItem(itemName, entry.launchApp, true)
}

func (entry *AppEntry) getMenuItemCloseAll() *MenuItem {
	return NewMenuItem(Tr("_Close All"), func(timestamp uint32) {
		logger.Debug("Close All")
		for win, _ := range entry.windows {
			ewmh.CloseWindow(XU, win)
		}
	}, true)
}

func (entry *AppEntry) getMenuItemDock() *MenuItem {
	return NewMenuItem(Tr("_Dock"), func(uint32) {
		logger.Debug("menu action dock entry")
		entry.entryManager.dockEntry(entry)
	}, true)
}

func (entry *AppEntry) getMenuItemUndock() *MenuItem {
	return NewMenuItem(Tr("_Undock"), func(uint32) {
		logger.Debug("menu action undock entry")
		entry.entryManager.undockEntry(entry)
	}, true)
}
