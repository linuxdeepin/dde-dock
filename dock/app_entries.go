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
	"strings"
)

type AppEntries []*AppEntry

func (entries AppEntries) GetFirstByInnerId(id string) *AppEntry {
	for _, entry := range entries {
		if entry.innerId == id {
			return entry
		}
	}
	return nil
}

func (entries AppEntries) Insert(entry *AppEntry, index int) AppEntries {
	// append
	if index < 0 || index >= len(entries) {
		return append(entries, entry)
	}
	// insert
	return append(entries[:index],
		append([]*AppEntry{entry}, entries[index:]...)...)
}

func (entries AppEntries) Remove(entry *AppEntry) AppEntries {
	index := entries.IndexOf(entry)
	if index != -1 {
		return append(entries[:index], entries[index+1:]...)
	}
	return entries
}

func (entries AppEntries) IndexOf(entry *AppEntry) int {
	var index int = -1
	for i, v := range entries {
		if v.Id == entry.Id {
			index = i
		}
	}
	return index
}

func (entries AppEntries) Move(index, newIndex int) (AppEntries, error) {
	if index == newIndex {
		return nil, errors.New("index == newIndex")
	}

	entriesLength := len(entries)
	if 0 <= index && index < entriesLength &&
		0 <= newIndex && newIndex < entriesLength {

		entry := entries[index]
		// remove entry at index
		removed := append(entries[:index], entries[index+1:]...)
		// insert entry at newIndex
		return append(removed[:newIndex],
			append([]*AppEntry{entry}, removed[newIndex:]...)...), nil
	}
	return nil, fmt.Errorf("Index out of bounds, index: %v, newIndex: %v, len: %v", index, newIndex, entriesLength)
}

func (entries AppEntries) FilterDocked() AppEntries {
	var dockedEntries AppEntries
	for _, entry := range entries {
		if entry.appInfo != nil && entry.IsDocked == true {
			dockedEntries = append(dockedEntries, entry)
		}
	}
	return dockedEntries
}

func (entries AppEntries) GetByWindowPid(pid uint) *AppEntry {
	for _, entry := range entries {
		for _, winInfo := range entry.windows {
			if winInfo.pid == pid {
				return entry
			}
		}
	}
	return nil
}

func (entries AppEntries) GetByAppId(id string) *AppEntry {
	for _, entry := range entries {
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

func (entries AppEntries) GetByDesktopFilePath(desktopFilePath string) (*AppEntry, error) {
	// same file
	for _, entry := range entries {
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
		return nil, errors.New("Invalid desktopFilePath")
	}
	hash := appInfo.innerId
	for _, entry := range entries {
		if entry.appInfo == nil {
			continue
		}
		if entry.appInfo.innerId == hash {
			return entry, nil
		}
	}
	return nil, nil
}
