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

func ListMediaShortcuts() Shortcuts {
	gs := newMediakeyGSetting()
	defer gs.Unref()

	keys := gs.ListKeys()
	idMap := getMediaIdNameMap()
	var list Shortcuts
	for _, k := range keys {
		tmp := Shortcut{
			Id:     k,
			Type:   KeyTypeMedia,
			Name:   getNameFromMap(k, idMap),
			Accels: gs.GetStrv(k),
		}
		list = append(list, &tmp)
	}
	return list
}

func resetMediaAccels() {
	gs := newMediakeyGSetting()
	defer gs.Unref()

	for _, key := range gs.ListKeys() {
		gs.Reset(key)
	}
}

func disableMediaAccels(key string) {
	gs := newMediakeyGSetting()
	defer gs.Unref()

	gs.SetStrv(key, []string{})
}

func addMediaAccel(key, accel string) {
	gs := newMediakeyGSetting()
	defer gs.Unref()

	list := gs.GetStrv(key)
	ret, added := addAccelToList(accel, list)
	if !added {
		return
	}
	gs.SetStrv(key, ret)
}

func delMediaAccel(key, accel string) {
	gs := newMediakeyGSetting()
	defer gs.Unref()

	list := gs.GetStrv(key)
	ret, deleted := delAccelFromList(accel, list)
	if !deleted {
		return
	}

	gs.SetStrv(key, ret)
}

func getMediaIdNameMap() map[string]string {
	var idMap = map[string]string{
		"calculator":  "Calculator",
		"eject":       "Eject",
		"email":       "Email client",
		"www":         "Web broswer",
		"media":       "Media player",
		"play":        "Play/Pause",
		"pause":       "Pause",
		"stop":        "Stop",
		"previous":    "Previous",
		"next":        "Next",
		"volume-mute": "Mute",
		"volume-down": "Volume down",
		"volume-up":   "Volume up",
	}
	return idMap
}
