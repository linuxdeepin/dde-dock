/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package dock

import (
	"reflect"
	"sync"

	"github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/lib/appinfo/desktopappinfo"
	"pkg.deepin.io/lib/dbus1"
	. "pkg.deepin.io/lib/gettext"
)

func (entry *AppEntry) updateMenu() {
	logger.Debug("Update menu")
	menu := NewMenu()
	menu.AppendItem(entry.getMenuItemLaunch())

	desktopActionMenuItems := entry.getMenuItemDesktopActions()
	menu.AppendItem(desktopActionMenuItems...)
	hasWin := entry.hasWindow()
	if hasWin {
		menu.AppendItem(entry.getMenuItemAllWindows())
	}

	// menu item dock or undock
	logger.Debug(entry.Id, "Item docked?", entry.IsDocked)
	if entry.IsDocked {
		menu.AppendItem(entry.getMenuItemUndock())
	} else {
		menu.AppendItem(entry.getMenuItemDock())
	}

	if hasWin {
		menu.AppendItem(entry.getMenuItemForceQuit())
		if entry.hasAllowedCloseWindow() {
			menu.AppendItem(entry.getMenuItemCloseAll())
		}
	}
	entry.Menu.setMenu(menu)
}

func (entry *AppEntry) getMenuItemDesktopActions() []*MenuItem {
	ai := entry.appInfo
	if ai == nil {
		return nil
	}

	var items []*MenuItem
	launchAction := func(action desktopappinfo.DesktopAction) func(timestamp uint32) {
		return func(timestamp uint32) {
			logger.Debugf("launch action %+v", action)
			err := entry.manager.startManager.LaunchAppAction(dbus.FlagNoAutoStart,
				ai.GetFileName(), action.Section, timestamp)
			if err != nil {
				logger.Warning("launchAppAction failed:", err)
			}
		}
	}

	for _, action := range ai.GetActions() {
		item := NewMenuItem(action.Name, launchAction(action), true)
		items = append(items, item)
	}
	return items
}

func (entry *AppEntry) launchApp(timestamp uint32) {
	logger.Debug("launchApp timestamp:", timestamp)
	if entry.appInfo != nil {
		logger.Debug("Has AppInfo")
		entry.manager.launch(entry.appInfo.GetFileName(), timestamp, nil)
	} else {
		// TODO
		logger.Debug("not supported")
	}
}

func (entry *AppEntry) getMenuItemLaunch() *MenuItem {
	var itemName string
	if entry.hasWindow() {
		itemName = entry.getName()
	} else {
		itemName = Tr("Open")
	}
	logger.Debugf("getMenuItemLaunch, itemName: %q", itemName)
	return NewMenuItem(itemName, entry.launchApp, true)
}

func (entry *AppEntry) getMenuItemCloseAll() *MenuItem {
	return NewMenuItem(Tr("Close All"), func(timestamp uint32) {
		logger.Debug("Close All")
		entry.PropsMu.RLock()
		winIds := entry.getAllowedCloseWindows()
		entry.PropsMu.RUnlock()

		for _, win := range winIds {
			err := closeWindow(win, x.Timestamp(timestamp))
			if err != nil {
				logger.Warningf("failed to close window %d: %v", win, err)
			}
		}
	}, true)
}

func (entry *AppEntry) getMenuItemForceQuit() *MenuItem {
	return NewMenuItem(Tr("Force Quit"), func(timestamp uint32) {
		logger.Debug("Force Quit")
		entry.ForceQuit()
	}, true)
}

func (entry *AppEntry) getMenuItemDock() *MenuItem {
	return NewMenuItem(Tr("Dock"), func(uint32) {
		logger.Debug("menu action dock entry")
		entry.RequestDock()
	}, true)
}

func (entry *AppEntry) getMenuItemUndock() *MenuItem {
	return NewMenuItem(Tr("Undock"), func(uint32) {
		logger.Debug("menu action undock entry")
		entry.RequestUndock()
	}, true)
}

func (entry *AppEntry) getMenuItemAllWindows() *MenuItem {
	menuItem := NewMenuItem(Tr("All Windows"), func(uint32) {
		logger.Debug("menu action all windows")
		entry.PresentWindows()
	}, true)
	menuItem.hint = menuItemHintShowAllWindows
	return menuItem
}

type AppEntryMenu struct {
	manager *Manager
	cache   string
	is3DWM  bool
	dirty   bool
	menu    *Menu
	mu      sync.Mutex
}

func (m *AppEntryMenu) setMenu(menu *Menu) {
	m.mu.Lock()
	m.menu = menu
	m.dirty = true
	m.mu.Unlock()
}

func (m *AppEntryMenu) getMenu() *Menu {
	m.mu.Lock()
	ret := m.menu
	m.mu.Unlock()
	return ret
}

func (*AppEntryMenu) SetValue(val interface{}) (changed bool, err *dbus.Error) {
	// read only
	return
}

func (m *AppEntryMenu) GetValue() (val interface{}, err *dbus.Error) {
	is3DWM := m.manager.is3DWM()
	m.mu.Lock()
	if m.dirty || m.cache == "" || m.is3DWM != is3DWM {
		items := make([]*MenuItem, 0, len(m.menu.Items))
		for _, item := range m.menu.Items {
			if is3DWM || item.hint != menuItemHintShowAllWindows {
				items = append(items, item)
			}
		}
		menu := NewMenu()
		menu.Items = items
		m.cache = menu.GenerateJSON()
		m.dirty = false
		m.is3DWM = is3DWM
	}
	val = m.cache
	m.mu.Unlock()
	return
}

func (*AppEntryMenu) SetNotifyChangedFunc(func(val interface{})) {
}

func (*AppEntryMenu) GetType() reflect.Type {
	return reflect.TypeOf("")
}
