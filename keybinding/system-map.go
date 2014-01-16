/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package main

var SystemIdNameMap = map[int32]string{
	0:   "launcher",
	620: "show-desktop",
	2:   "lock-screen",
	10:  "file-manager",
	621: "switch-applications",
	622: "switch-applications-backward",
	800: "prev-key",
	801: "next-key",
	3:   "show-dock",
	6:   "screenshot",
	7:   "screenshot-full-screen",
	8:   "screenshot-window",
	9:   "screenshot-delayed",
	1:   "terminal",
	5:   "terminal-quake",
	4:   "logout",
}

var SystemIdIndexMap = map[int32]int32{
	0:   0,
	620: 1,
	2:   2,
	10:  3,
	621: 4,
	622: 5,
	800: 6,
	801: 7,
	3:   8,
	6:   9,
	7:   10,
	8:   11,
	9:   12,
	1:   13,
	5:   14,
	4:   15,
}

var SystemNameDescMap = map[string]string{
	"launcher":                     "Launcher",
	"show-desktop":                 "Show desktop",
	"lock-screen":                  "Lock screen",
	"file-manager":                 "File manager",
	"switch-applications":          "Switch applications",
	"switch-applications-backward": "Reverse switch applications",
	"prev-key":                     "Switch applications with 3D effect",
	"next-key":                     "Reverse switch applications with 3D effect",
	"show-dock":                    "Show/Hide the dock",
	"screenshot":                   "Take a screenshot",
	"screenshot-full-screen":       "Take a screenshot of full screen",
	"screenshot-window":            "Take a screenshot of a window",
	"screenshot-delayed":           "Take a screenshot delayed",
	"terminal":                     "Terminal",
	"terminal-quake":               "Terminal Quake Window",
	"logout":                       "Logout",
}

var MediaIdNameMap = map[int32]string{
	300: "calculator",
	302: "email",
	303: "www",
	304: "media",
	305: "play",
	306: "pause",
	307: "stop",
	309: "volume-mute",
	308: "volume-down",
	310: "volume-up",
	311: "previous",
	312: "next",
	301: "eject",
}

var MediaIdIndexMap = map[int32]int32{
	300: 0,
	302: 1,
	303: 2,
	304: 3,
	305: 4,
	306: 5,
	307: 6,
	309: 7,
	308: 8,
	310: 9,
	311: 10,
	312: 11,
	301: 12,
}

var MediaNameDescMap = map[string]string{
	"calculator":  "Launch calculator",
	"email":       "Launch email client",
	"www":         "Launch web browser",
	"media":       "Launch media player",
	"play":        "Play (or play/pause)",
	"pause":       "Pause playback",
	"stop":        "Stop playback",
	"volume-mute": "Volume mute",
	"volume-down": "Volume down",
	"volume-up":   "Volume up",
	"previous":    "Previous track",
	"next":        "Next track",
	"eject":       "Eject",
}

var WindowIdNameMap = map[int32]string{
	603: "close",
	604: "maximize",
	607: "unmaximize",
	605: "minimize",
	601: "begin-move",
	602: "begin-resize",
	606: "toggle-shaded",
	600: "activate-window-menu",
}

var WindowIdIndexMap = map[int32]int32{
	603: 0,
	604: 1,
	607: 2,
	605: 3,
	601: 4,
	602: 5,
	606: 6,
	600: 7,
}

var WindowNameDescMap = map[string]string{
	"close":                "Close window",
	"maximize":             "Maximize window",
	"unmaximize":           "Restore window",
	"minimize":             "Minimize window",
	"begin-move":           "Move window",
	"begin-resize":         "Resize window",
	"toggle-shaded":        "Toggle shaded state",
	"activate-window-menu": "Activate the window menu",
}

var WorkSpaceIdNameMap = map[int32]string{
	608: "switch-to-workspace-1",
	609: "switch-to-workspace-2",
	610: "switch-to-workspace-3",
	611: "switch-to-workspace-4",
	613: "switch-to-workspace-left",
	614: "switch-to-workspace-right",
	615: "switch-to-workspace-up",
	612: "switch-to-workspace-down",
	900: "put-viewport-1-key",
	901: "put-viewport-2-key",
	902: "put-viewport-3-key",
	903: "put-viewport-4-key",
	617: "move-to-workspace-left",
	618: "move-to-workspace-right",
	619: "move-to-workspace-up",
	616: "move-to-workspace-down",
}

var WorkSpaceIdIndexMap = map[int32]int32{
	608: 0,
	609: 1,
	610: 2,
	611: 3,
	613: 4,
	614: 5,
	615: 6,
	612: 7,
	900: 8,
	901: 9,
	902: 10,
	903: 11,
	617: 12,
	618: 13,
	619: 14,
	616: 15,
}

var WorkSpaceNameDescMap = map[string]string{
	"switch-to-workspace-1":     "Switch to workspace 1",
	"switch-to-workspace-2":     "Switch to workspace 2",
	"switch-to-workspace-3":     "Switch to workspace 3",
	"switch-to-workspace-4":     "Switch to workspace 4",
	"switch-to-workspace-left":  "Switch to workspace left",
	"switch-to-workspace-right": "Switch to workspace right",
	"switch-to-workspace-up":    "Switch to workspace up",
	"switch-to-workspace-down":  "Switch to workspace down",
	"put-viewport-1-key":        "Move to workspace 1",
	"put-viewport-2-key":        "Move to workspace 2",
	"put-viewport-3-key":        "Move to workspace 3",
	"put-viewport-4-key":        "Move to workspace 4",
	"move-to-workspace-left":    "Move to workspace left",
	"move-to-workspace-right":   "Move to workspace right",
	"move-to-workspace-up":      "Move to workspace up",
	"move-to-workspace-down":    "Move to workspace down",
}

var IdNameMap = map[int32]string{
	0:   "launcher",
	1:   "terminal",
	2:   "lock-screen",
	3:   "show-dock",
	4:   "logout",
	5:   "terminal-quake",
	6:   "screenshot",
	7:   "screenshot-full-screen",
	8:   "screenshot-window",
	9:   "screenshot-delayed",
	10:  "file-manager",
	620: "show-desktop",
	621: "switch-applications",
	622: "switch-applications-backward",
	800: "prev-key",
	801: "next-key",
	300: "calculator",
	301: "eject",
	302: "email",
	303: "www",
	304: "media",
	305: "play",
	306: "pause",
	307: "stop",
	308: "volume-down",
	309: "volume-mute",
	310: "volume-up",
	311: "previous",
	312: "next",
	600: "activate-window-menu",
	601: "begin-move",
	602: "begin-resize",
	603: "close",
	604: "maximize",
	605: "minimize",
	606: "toggle-shaded",
	607: "unmaximize",
	608: "switch-to-workspace-1",
	609: "switch-to-workspace-2",
	610: "switch-to-workspace-3",
	611: "switch-to-workspace-4",
	612: "switch-to-workspace-down",
	613: "switch-to-workspace-left",
	614: "switch-to-workspace-right",
	615: "switch-to-workspace-up",
	900: "put-viewport-1-key",
	901: "put-viewport-2-key",
	902: "put-viewport-3-key",
	903: "put-viewport-4-key",
	616: "move-to-workspace-down",
	617: "move-to-workspace-left",
	618: "move-to-workspace-right",
	619: "move-to-workspace-up",
}
