/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

type EntryType int

const (
	MENUENTRY EntryType = iota
	SUBMENU
)

type Entry struct {
	entryType     EntryType
	title         string
	num           int
	parentSubMenu *Entry
}

func (entry *Entry) getFullTitle() string {
	if entry.parentSubMenu != nil {
		return entry.parentSubMenu.getFullTitle() + ">" + entry.title
	} else {
		return entry.title
	}
}
