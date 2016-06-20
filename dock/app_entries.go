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
	"errors"
	"fmt"
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
