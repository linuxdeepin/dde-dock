/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package grub2

// EntryType is used to define entry's type in '/boot/grub/grub.cfg'.
type EntryType int

const (
	MENUENTRY EntryType = iota
	SUBMENU
)

// Entry is a struct to store each entry's data/
type Entry struct {
	entryType     EntryType
	title         string
	num           int
	parentSubMenu *Entry
}

func (entry *Entry) getFullTitle() string {
	if entry.parentSubMenu != nil {
		return entry.parentSubMenu.getFullTitle() + ">" + entry.title
	}
	return entry.title
}
