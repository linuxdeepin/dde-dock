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

func ListWMShortcut() Shortcuts {
	s := newWMGSetting()
	defer s.Unref()
	return doListShortcut(s, wmIdNameMap(), KeyTypeWM)
}

func resetWMAccels() {
	s := newWMGSetting()
	defer s.Unref()
	doResetAccels(s)
}

func disableWMAccels(key string) {
	s := newWMGSetting()
	defer s.Unref()
	doDisableAccles(s, key)
}

func addWMAccel(key, accel string) {
	s := newWMGSetting()
	defer s.Unref()
	doAddAccel(s, key, accel)
}

func delWMAccel(key, accel string) {
	s := newWMGSetting()
	defer s.Unref()
	doDelAccel(s, key, accel)
}

func wmIdNameMap() map[string]string {
	var idNameMap = map[string]string{
		"switch-to-workspace-1":        "Switch to workspace 1",
		"switch-to-workspace-2":        "Switch to workspace 2",
		"switch-to-workspace-3":        "Switch to workspace 3",
		"switch-to-workspace-4":        "Switch to workspace 4",
		"switch-to-workspace-5":        "Switch to workspace 5",
		"switch-to-workspace-6":        "Switch to workspace 6",
		"switch-to-workspace-7":        "Switch to workspace 7",
		"switch-to-workspace-8":        "Switch to workspace 8",
		"switch-to-workspace-9":        "Switch to workspace 9",
		"switch-to-workspace-10":       "Switch to workspace 10",
		"switch-to-workspace-11":       "Switch to workspace 11",
		"switch-to-workspace-12":       "Switch to workspace 12",
		"switch-to-workspace-left":     gettext.Tr("Switch to left workspace"),
		"switch-to-workspace-right":    gettext.Tr("Switch to right workspace"),
		"switch-to-workspace-up":       gettext.Tr("Switch to upper workspace"),
		"switch-to-workspace-down":     gettext.Tr("Switch to lower workspace"),
		"switch-to-workspace-last":     "Switch to last workspace",
		"switch-group":                 gettext.Tr("Switch similar windows"),
		"switch-group-backward":        gettext.Tr("Switch similar windows in reverse"),
		"switch-applications":          gettext.Tr("Switch windows"),
		"switch-applications-backward": gettext.Tr("Switch windows in reverse"),
		"switch-windows":               "Switch windows",
		"switch-windows-backward":      "Reverse switch windows",
		"switch-panels":                "Switch system controls",
		"switch-panels-backward":       "Reverse switch system controls",
		"cycle-group":                  "Switch windows of an app directly",
		"cycle-group-backward":         "Reverse switch windows of an app directly",
		"cycle-windows":                "Switch windows directly",
		"cycle-windows-backward":       "Reverse switch windows directly",
		"cycle-panels":                 "Switch system controls directly",
		"cycle-panels-backward":        "Reverse switch system controls directly",
		"show-desktop":                 gettext.Tr("Show desktop"),
		"panel-main-menu":              "Show the activities overview",
		"panel-run-dialog":             "Show the run command prompt",
		// Don't use
		// "set-spew-mark":                gettext.Tr(""),
		"activate-window-menu":         "Activate window menu",
		"toggle-fullscreen":            "toggle-fullscreen",
		"toggle-maximized":             "Toggle maximization state",
		"toggle-above":                 "Toggle window always appearing on top",
		"maximize":                     gettext.Tr("Maximize window"),
		"unmaximize":                   gettext.Tr("Restore window"),
		"toggle-shaded":                "Switch furl state",
		"minimize":                     "Minimize window",
		"close":                        gettext.Tr("Close window"),
		"begin-move":                   gettext.Tr("Move window"),
		"begin-resize":                 gettext.Tr("Resize window"),
		"toggle-on-all-workspaces":     "Toggle window on all workspaces or one",
		"move-to-workspace-1":          "Move to workspace 1",
		"move-to-workspace-2":          "Move to workspace 2",
		"move-to-workspace-3":          "Move to workspace 3",
		"move-to-workspace-4":          "Move to workspace 4",
		"move-to-workspace-5":          "Move to workspace 5",
		"move-to-workspace-6":          "Move to workspace 6",
		"move-to-workspace-7":          "Move to workspace 7",
		"move-to-workspace-8":          "Move to workspace 8",
		"move-to-workspace-9":          "Move to workspace 9",
		"move-to-workspace-10":         "Move to workspace 10",
		"move-to-workspace-11":         "Move to workspace 11",
		"move-to-workspace-12":         "Move to workspace 12",
		"move-to-workspace-last":       "Move to last workspace",
		"move-to-workspace-left":       gettext.Tr("Move to left workspace"),
		"move-to-workspace-right":      gettext.Tr("Move to right workspace"),
		"move-to-workspace-up":         gettext.Tr("Move to upper workspace"),
		"move-to-workspace-down":       gettext.Tr("Move to lower workspace"),
		"move-to-monitor-left":         "Move to left monitor",
		"move-to-monitor-right":        "Move to right monitor",
		"move-to-monitor-up":           "Move to up monitor",
		"move-to-monitor-down":         "Move to down monitor",
		"raise-or-lower":               "Raise window if covered, otherwise lower it",
		"raise":                        "Raise window above other windows",
		"lower":                        "Lower window below other windows",
		"maximize-vertically":          "Maximize window vertically",
		"maximize-horizontally":        "Maximize window horizontally",
		"move-to-corner-nw":            "Move window to top left corner",
		"move-to-corner-ne":            "Move window to top right corner",
		"move-to-corner-sw":            "Move window to bottom left corner",
		"move-to-corner-se":            "Move window to bottom right corner",
		"move-to-side-n":               "Move window to top edge of screen",
		"move-to-side-s":               "Move window to bottom edge of screen",
		"move-to-side-e":               "Move window to right side of screen",
		"move-to-side-w":               "Move window to left side of screen",
		"move-to-center":               "Move window to center of screen",
		"switch-input-source":          "Binding to select the next input source",
		"switch-input-source-backward": "Binding to select the previous input source",
		"always-on-top":                "Set or unset window to appear always on top",
	}
	return idNameMap
}
