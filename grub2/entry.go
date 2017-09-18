/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
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
