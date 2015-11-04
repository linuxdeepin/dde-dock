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
	dutils "pkg.deepin.io/lib/utils"
)

const (
	systemSchema   = "com.deepin.dde.keybinding.system"
	mediakeySchema = "com.deepin.dde.keybinding.mediakey"
	wmSchema       = "com.deepin.wrap.gnome.desktop.wm.keybindings"
	metacitySchema = "com.deepin.wrap.gnome.metacity.keybindings"
	galaSchema     = "com.deepin.wrap.pantheon.desktop.gala.keybindings"
)

func doListShortcut(setting *gio.Settings, idMap map[string]string, ty int32) Shortcuts {
	var list = Shortcuts{}
	for _, key := range setting.ListKeys() {
		var tmp = Shortcut{
			Id:     key,
			Type:   ty,
			Name:   getNameFromMap(key, idMap),
			Accels: filterNilStr(setting.GetStrv(key)),
		}

		list = append(list, &tmp)
	}

	return list
}

func doResetAccels(setting *gio.Settings) {
	for _, key := range setting.ListKeys() {
		_, srcList := setting.GetDefaultValue(key).GetStrv()
		if isAccelsEqual(srcList, setting.GetStrv(key)) {
			continue
		}
		setting.Reset(key)
	}
}

func doDisableAccles(setting *gio.Settings, key string) {
	setting.SetStrv(key, []string{})
}

func doAddAccel(setting *gio.Settings, key, accel string) {
	accels := setting.GetStrv(key)
	list, added := addAccelToList(accel, accels)
	if !added {
		return
	}
	setting.SetStrv(key, list)
}

func doDelAccel(setting *gio.Settings, key, accel string) {
	accels := setting.GetStrv(key)
	list, deleted := delAccelFromList(accel, accels)
	if !deleted {
		return
	}
	setting.SetStrv(key, list)
}

func newSystemGSetting() *gio.Settings {
	return gio.NewSettings(systemSchema)
}

func newWMGSetting() *gio.Settings {
	return gio.NewSettings(wmSchema)
}

func newMediakeyGSetting() *gio.Settings {
	return gio.NewSettings(mediakeySchema)
}

func newMetacityGSetting() *gio.Settings {
	s, _ := dutils.CheckAndNewGSettings(metacitySchema)
	return s
}

func newGalaGSetting() *gio.Settings {
	s, _ := dutils.CheckAndNewGSettings(galaSchema)
	return s
}

func getNameFromMap(id string, m map[string]string) string {
	name, ok := m[id]
	if !ok {
		return id
	}
	return name
}

func isAccelsEqual(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return false
	}

	for _, v := range l1 {
		if !isAccelInList(v, l2) {
			return false
		}
	}

	return true
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

	var ret = []string{accel}
	ret = append(ret, list...)
	return ret, true
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
