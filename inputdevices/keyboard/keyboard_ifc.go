/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
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

package keyboard

import (
	"strings"
)

func (kbd *Keyboard) Reset() {
	for _, key := range kbd.settings.ListKeys() {
		kbd.settings.Reset(key)
	}
}

func (kbd *Keyboard) LayoutList() map[string]string {
	return kbd.layoutDescMap
}

//func (kbd *Keyboard) GetLayoutLocale(layout string) string {
func (kbd *Keyboard) GetLayoutDesc(layout string) string {
	if len(layout) == 0 || !strings.Contains(layout, kbdKeyLayoutDelim) {
		return ""
	}

	desc, ok := kbd.layoutDescMap[layout]
	if !ok {
		return ""
	}

	return desc
}

func (kbd *Keyboard) AddLayoutOption(option string) {
	if len(option) == 0 {
		return
	}

	list := kbd.settings.GetStrv(kbdKeyUserLayoutList)
	if isStrInList(option, list) {
		return
	}

	list = append(list, option)
	kbd.settings.SetStrv(kbdKeyUserLayoutList, list)
}

func (kbd *Keyboard) DeleteLayoutOption(option string) {
	if len(option) == 0 {
		return
	}

	list := kbd.settings.GetStrv(kbdKeyUserLayoutList)
	var tmpList []string
	for _, v := range list {
		if v == option {
			continue
		}
		tmpList = append(tmpList, v)
	}
	if len(list) == len(tmpList) {
		return
	}

	kbd.settings.SetStrv(kbdKeyUserLayoutList, tmpList)
}

func (kbd *Keyboard) ClearLayoutOption() {
	kbd.settings.SetStrv(kbdKeyUserLayoutList, []string{})
	setOptionList([]string{})
}

func (kbd *Keyboard) AddUserLayout(layout string) {
	if len(layout) == 0 {
		return
	}

	kbd.addUserLayoutToList(layout)
}

func (kbd *Keyboard) DeleteUserLayout(layout string) {
	if len(layout) == 0 {
		return
	}

	kbd.deleteUserLayoutFromList(layout)
}
