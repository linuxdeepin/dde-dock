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
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	x "github.com/linuxdeepin/go-x11-client"

	"pkg.deepin.io/lib/dbus1"
)

type AppEntries struct {
	items []*AppEntry
	mu    sync.RWMutex

	insertCb func(entry *AppEntry, index int)
	removeCb func(entry *AppEntry)
}

func (entries *AppEntries) GetValue() (val interface{}, err *dbus.Error) {
	entries.mu.RLock()
	result := make([]dbus.ObjectPath, len(entries.items))
	for idx, entry := range entries.items {
		result[idx] = dbus.ObjectPath(entryDBusObjPathPrefix + entry.Id)
	}
	entries.mu.RUnlock()
	return result, nil
}

func (entries *AppEntries) SetNotifyChangedFunc(func(val interface{})) {
}

func (entries *AppEntries) SetValue(val interface{}) (changed bool, err *dbus.Error) {
	// readonly
	return
}

func (entries *AppEntries) GetType() reflect.Type {
	return reflect.TypeOf([]dbus.ObjectPath{})
}

func (entries *AppEntries) GetByInnerId(id string) *AppEntry {
	entries.mu.RLock()
	for _, entry := range entries.items {
		if entry.innerId == id {
			entries.mu.RUnlock()
			return entry
		}
	}
	entries.mu.RUnlock()
	return nil
}

func (entries *AppEntries) Append(entry *AppEntry) {
	entries.Insert(entry, -1)
}

func (entries *AppEntries) Insert(entry *AppEntry, index int) {
	entries.mu.Lock()
	if index < 0 || index >= len(entries.items) {
		// append
		index = len(entries.items)
		entries.items = append(entries.items, entry)
	} else {
		// insert
		entries.items = append(entries.items[:index],
			append([]*AppEntry{entry}, entries.items[index:]...)...)
	}

	if entries.insertCb != nil {
		entries.insertCb(entry, index)
	}
	entries.mu.Unlock()
}

func (entries *AppEntries) Remove(entry *AppEntry) {
	entries.mu.Lock()
	index := entries.indexOf(entry)
	if index != -1 {
		entries.items = append(entries.items[:index], entries.items[index+1:]...)
		entries.removeCb(entry)
	}
	entries.mu.Unlock()
}

func (entries *AppEntries) indexOf(entry *AppEntry) int {
	index := -1
	for i, v := range entries.items {
		if v.Id == entry.Id {
			index = i
		}
	}
	return index
}

func (entries *AppEntries) IndexOf(entry *AppEntry) int {
	entries.mu.RLock()
	idx := entries.indexOf(entry)
	entries.mu.RUnlock()
	return idx
}

func (entries *AppEntries) Move(index, newIndex int) error {
	if index == newIndex {
		return errors.New("index == newIndex")
	}

	entries.mu.Lock()

	entriesLength := len(entries.items)
	if 0 <= index && index < entriesLength &&
		0 <= newIndex && newIndex < entriesLength {

		entry := entries.items[index]
		// remove entry at index
		removed := append(entries.items[:index], entries.items[index+1:]...)
		// insert entry at newIndex
		entries.items = append(removed[:newIndex],
			append([]*AppEntry{entry}, removed[newIndex:]...)...)

		entries.mu.Unlock()
		return nil
	}
	entries.mu.Unlock()
	return fmt.Errorf("index out of bounds, index: %v, newIndex: %v, len: %v", index, newIndex, entriesLength)
}

func (entries *AppEntries) FilterDocked() (dockedEntries []*AppEntry) {
	entries.mu.RLock()
	for _, entry := range entries.items {
		if entry.appInfo != nil && entry.IsDocked == true {
			dockedEntries = append(dockedEntries, entry)
		}
	}
	entries.mu.RUnlock()
	return dockedEntries
}

func (entries *AppEntries) GetByWindowPid(pid uint) *AppEntry {
	entries.mu.RLock()

	for _, entry := range entries.items {
		for _, winInfo := range entry.windows {
			if winInfo.pid == pid {
				entries.mu.RUnlock()
				return entry
			}
		}
	}

	entries.mu.RUnlock()
	return nil
}

func (entries *AppEntries) getByWindowId(winId x.Window) *AppEntry {
	entries.mu.RLock()
	for _, entry := range entries.items {
		entry.PropsMu.RLock()
		_, ok := entry.windows[winId]
		entry.PropsMu.RUnlock()
		if ok {
			entries.mu.RUnlock()
			return entry
		}
	}

	entries.mu.RUnlock()
	// not found
	return nil
}

func getByAppId(items []*AppEntry, id string) *AppEntry {
	for _, entry := range items {
		if entry.appInfo == nil {
			continue
		}

		eAppId := entry.appInfo.GetId()
		if strings.EqualFold(id, eAppId) {
			return entry
		}
	}
	return nil
}

func (entries *AppEntries) GetByAppId(id string) *AppEntry {
	entries.mu.RLock()
	e := getByAppId(entries.items, id)
	entries.mu.RUnlock()
	return e
}

func getByDesktopFilePath(entriesItems []*AppEntry, desktopFilePath string) (*AppEntry, error) {
	// same file
	for _, entry := range entriesItems {
		if entry.appInfo == nil {
			continue
		}
		file := entry.appInfo.GetFileName()
		if file == desktopFilePath {
			return entry, nil
		}
	}

	// hash equal
	appInfo := NewAppInfoFromFile(desktopFilePath)
	if appInfo == nil {
		return nil, errors.New("invalid desktopFilePath")
	}
	hash := appInfo.innerId
	for _, entry := range entriesItems {
		if entry.appInfo == nil {
			continue
		}
		if entry.appInfo.innerId == hash {
			return entry, nil
		}
	}
	return nil, nil
}

func (entries *AppEntries) GetByDesktopFilePath(desktopFilePath string) (*AppEntry, error) {
	entries.mu.RLock()
	e, err := getByDesktopFilePath(entries.items, desktopFilePath)
	entries.mu.RUnlock()
	return e, err
}
