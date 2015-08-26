/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
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

package shortcuts

import (
	"pkg.deepin.io/dde/daemon/keybinding/core"
	"pkg.deepin.io/lib/gio-2.0"
)

const (
	systemSchema   = "com.deepin.dde.keybinding.system"
	mediakeySchema = "com.deepin.dde.keybinding.mediakey"
	wmSchema       = "com.deepin.wrap.gnome.desktop.wm.keybindings"
)

func newSystemGSetting() *gio.Settings {
	return gio.NewSettings(systemSchema)
}

func newWMGSetting() *gio.Settings {
	return gio.NewSettings(wmSchema)
}

func newMediakeyGSetting() *gio.Settings {
	return gio.NewSettings(mediakeySchema)
}

func getNameFromMap(id string, m map[string]string) string {
	name, ok := m[id]
	if !ok {
		return id
	}
	return name
}

func isAccelInList(accel string, list []string) bool {
	for _, v := range list {
		if core.IsAccelEqual(v, accel) {
			return true
		}
	}

	return false
}

func filterNilStr(list []string) []string {
	var ret []string
	for _, k := range list {
		if len(k) == 0 {
			continue
		}
		ret = append(ret, k)
	}
	return ret
}

func addAccelToList(accel string, list []string) ([]string, bool) {
	if isAccelInList(accel, list) {
		return list, false
	}

	list = append(list, accel)
	return list, true
}

func delAccelFromList(accel string, list []string) ([]string, bool) {
	var (
		ret   []string
		found bool
	)
	for _, v := range list {
		if core.IsAccelEqual(accel, v) {
			found = true
			continue
		}
		ret = append(ret, v)
	}
	return ret, found
}

func isStrInList(s string, list []string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}
