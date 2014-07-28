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

package keybinding

var compizKeysMap = map[string]string{
	//org.compiz.core.gschema.xml	(s)
	"show-desktop":         "show-desktop-key",
	"close":                "close-window-key",
	"maximize":             "maximize-window-key",
	"unmaximize":           "unmaximize-window-key",
	"minimize":             "minimize-window-key",
	"toggle-shaded":        "toggle-window-shaded-key",
	"activate-window-menu": "window-menu-key",
	//org.compiz.move.gschema.xml	(s)
	"begin-move": "initiate-key",
	//org.compiz.resize.gschema.xml	(s)
	"begin-resize": "initiate-key",
	//org.compiz.vpswitch.gschema.xml	(s)
	"switch-to-workspace-1": "switch-to-1-key",
	"switch-to-workspace-2": "switch-to-2-key",
	"switch-to-workspace-3": "switch-to-3-key",
	"switch-to-workspace-4": "switch-to-4-key",
	//org.compiz.put.gschema.xml	(s)
	"put-viewport-1-key": "put-viewport-1-key",
	"put-viewport-2-key": "put-viewport-2-key",
	"put-viewport-3-key": "put-viewport-3-key",
	"put-viewport-4-key": "put-viewport-4-key",
	//org.compiz.wall.gschema.xml	(s)
	"switch-to-workspace-left":  "left-key",
	"switch-to-workspace-right": "right-key",
	"switch-to-workspace-up":    "up-key",
	"switch-to-workspace-down":  "down-key",
	"move-to-workspace-left":    "left-window-key",
	"move-to-workspace-right":   "right-window-key",
	"move-to-workspace-up":      "up-window-key",
	"move-to-workspace-down":    "down-window-key",
	//org.compiz.shift.gschema.xml	(s)
	"next-key": "next-key",
	"prev-key": "prev-key",
	//org.compiz.switcher.gschema.xml	(s)
	"switch-applications":          "next-key",
	"switch-applications-backward": "prev-key",
}
