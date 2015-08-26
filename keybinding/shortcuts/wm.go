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
	"pkg.deepin.io/lib/gettext"
)

func ListWMShortcuts() Shortcuts {
	gs := newWMGSetting()
	defer gs.Unref()

	keys := gs.ListKeys()
	idMap := wmIdNameMap()
	var list = Shortcuts{}
	for _, k := range keys {
		tmp := Shortcut{
			Id:     k,
			Type:   KeyTypeWM,
			Name:   getNameFromMap(k, idMap),
			Accels: filterNilStr(gs.GetStrv(k)),
		}

		list = append(list, &tmp)
	}

	return list
}

func resetWMAccels() {
	gs := newWMGSetting()
	defer gs.Unref()

	for _, key := range gs.ListKeys() {
		gs.Reset(key)
	}
}

func disableWMAccels(key string) {
	gs := newWMGSetting()
	defer gs.Unref()

	gs.SetStrv(key, []string{})
}

func addWMAccel(key, accel string) {
	gs := newWMGSetting()
	defer gs.Unref()

	list := gs.GetStrv(key)
	ret, added := addAccelToList(accel, list)
	if !added {
		return
	}
	gs.SetStrv(key, filterNilStr(ret))
}

func delWMAccel(key, accel string) {
	gs := newWMGSetting()
	defer gs.Unref()

	list := gs.GetStrv(key)
	ret, deleted := delAccelFromList(accel, list)
	if !deleted {
		return
	}

	gs.SetStrv(key, filterNilStr(ret))
}

func wmIdNameMap() map[string]string {
	var idNameMap = map[string]string{
		"switch-to-workspace-1":        gettext.Tr("Switch to workspace 1"),
		"switch-to-workspace-2":        gettext.Tr("Switch to workspace 2"),
		"switch-to-workspace-3":        gettext.Tr("Switch to workspace 3"),
		"switch-to-workspace-4":        gettext.Tr("Switch to workspace 4"),
		"switch-to-workspace-5":        gettext.Tr("Switch to workspace 5"),
		"switch-to-workspace-6":        gettext.Tr("Switch to workspace 6"),
		"switch-to-workspace-7":        gettext.Tr("Switch to workspace 7"),
		"switch-to-workspace-8":        gettext.Tr("Switch to workspace 8"),
		"switch-to-workspace-9":        gettext.Tr("Switch to workspace 9"),
		"switch-to-workspace-10":       gettext.Tr("Switch to workspace 10"),
		"switch-to-workspace-11":       gettext.Tr("Switch to workspace 11"),
		"switch-to-workspace-12":       gettext.Tr("Switch to workspace 12"),
		"switch-to-workspace-left":     gettext.Tr("Switch to left workspace"),
		"switch-to-workspace-right":    gettext.Tr("Switch to right workspace"),
		"switch-to-workspace-up":       gettext.Tr("Switch to up workspace"),
		"switch-to-workspace-down":     gettext.Tr("Switch to down workspace"),
		"switch-to-workspace-last":     gettext.Tr("Switch to last workspace"),
		"switch-group":                 gettext.Tr("Switch windows of an application"),
		"switch-group-backward":        gettext.Tr("Reverse switch windows of an application"),
		"switch-applications":          gettext.Tr("Switch applications"),
		"switch-applications-backward": gettext.Tr("Reverse switch applications"),
		"switch-windows":               gettext.Tr("Switch windows"),
		"switch-windows-backward":      gettext.Tr("Reverse switch windows"),
		"switch-panels":                gettext.Tr("Switch system controls"),
		"switch-panels-backward":       gettext.Tr("Reverse switch system controls"),
		"cycle-group":                  gettext.Tr("Switch windows of an app directly"),
		"cycle-group-backward":         gettext.Tr("Reverse switch windows of an app directly"),
		"cycle-windows":                gettext.Tr("Switch windows directly"),
		"cycle-windows-backward":       gettext.Tr("Reverse switch windows directly"),
		"cycle-panels":                 gettext.Tr("Switch system controls directly"),
		"cycle-panels-backward":        gettext.Tr("Reverse switch system controls directly"),
		"show-desktop":                 gettext.Tr("Show desktop"),
		"panel-main-menu":              gettext.Tr("Show the activities overview"),
		"panel-run-dialog":             gettext.Tr("Show the run command prompt"),
		// Don't use
		// "set-spew-mark":                gettext.Tr(""),
		"activate-window-menu":         gettext.Tr("Activate window menu"),
		"toggle-fullscreen":            gettext.Tr("toggle-fullscreen"),
		"toggle-maximized":             gettext.Tr("Toggle maximization state"),
		"toggle-above":                 gettext.Tr("Toggle window always appearing on top"),
		"maximize":                     gettext.Tr("Maximize window"),
		"unmaximize":                   gettext.Tr("Restore window"),
		"toggle-shaded":                gettext.Tr("Switch furl state"),
		"minimize":                     gettext.Tr("Minimize window"),
		"close":                        gettext.Tr("Close window"),
		"begin-move":                   gettext.Tr("Move window"),
		"begin-resize":                 gettext.Tr("Resize window"),
		"toggle-on-all-workspaces":     gettext.Tr("Toggle window on all workspaces or one"),
		"move-to-workspace-1":          gettext.Tr("Move to workspace 1"),
		"move-to-workspace-2":          gettext.Tr("Move to workspace 2"),
		"move-to-workspace-3":          gettext.Tr("Move to workspace 3"),
		"move-to-workspace-4":          gettext.Tr("Move to workspace 4"),
		"move-to-workspace-5":          gettext.Tr("Move to workspace 5"),
		"move-to-workspace-6":          gettext.Tr("Move to workspace 6"),
		"move-to-workspace-7":          gettext.Tr("Move to workspace 7"),
		"move-to-workspace-8":          gettext.Tr("Move to workspace 8"),
		"move-to-workspace-9":          gettext.Tr("Move to workspace 9"),
		"move-to-workspace-10":         gettext.Tr("Move to workspace 10"),
		"move-to-workspace-11":         gettext.Tr("Move to workspace 11"),
		"move-to-workspace-12":         gettext.Tr("Move to workspace 12"),
		"move-to-workspace-last":       gettext.Tr("Move to last workspace"),
		"move-to-workspace-left":       gettext.Tr("Move to left workspace"),
		"move-to-workspace-right":      gettext.Tr("Move to right workspace"),
		"move-to-workspace-up":         gettext.Tr("Move to up workspace"),
		"move-to-workspace-down":       gettext.Tr("Move to down workspace"),
		"move-to-monitor-left":         gettext.Tr("Move to left monitor"),
		"move-to-monitor-right":        gettext.Tr("Move to right monitor"),
		"move-to-monitor-up":           gettext.Tr("Move to up monitor"),
		"move-to-monitor-down":         gettext.Tr("Move to down monitor"),
		"raise-or-lower":               gettext.Tr("Raise window if covered, otherwise lower it"),
		"raise":                        gettext.Tr("Raise window above other windows"),
		"lower":                        gettext.Tr("Lower window below other windows"),
		"maximize-vertically":          gettext.Tr("Maximize window vertically"),
		"maximize-horizontally":        gettext.Tr("Maximize window horizontally"),
		"move-to-corner-nw":            gettext.Tr("Move window to top left corner"),
		"move-to-corner-ne":            gettext.Tr("Move window to top right corner"),
		"move-to-corner-sw":            gettext.Tr("Move window to bottom left corner"),
		"move-to-corner-se":            gettext.Tr("Move window to bottom right corner"),
		"move-to-side-n":               gettext.Tr("Move window to top edge of screen"),
		"move-to-side-s":               gettext.Tr("Move window to bottom edge of screen"),
		"move-to-side-e":               gettext.Tr("Move window to right side of screen"),
		"move-to-side-w":               gettext.Tr("Move window to left side of screen"),
		"move-to-center":               gettext.Tr("Move window to center of screen"),
		"switch-input-source":          gettext.Tr("Binding to select the next input source"),
		"switch-input-source-backward": gettext.Tr("Binding to select the previous input source"),
		"always-on-top":                gettext.Tr("Set or unset window to appear always on top"),
	}
	return idNameMap
}
