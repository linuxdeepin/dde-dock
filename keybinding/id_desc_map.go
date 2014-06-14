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

package keybinding

import . "dlib/gettext"

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
	11:  "disable-touchpad",
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
	11:  16,
}

var SystemNameDescMap map[string]string

func initSystemNameDescMap() {
	SystemNameDescMap = map[string]string{
		"launcher":                     Tr("Launcher"),
		"show-desktop":                 Tr("Show desktop"),
		"lock-screen":                  Tr("Lock screen"),
		"file-manager":                 Tr("File manager"),
		"switch-applications":          Tr("Switch applications"),
		"switch-applications-backward": Tr("Reverse switch applications"),
		"prev-key":                     Tr("3D switch applications"),
		"next-key":                     Tr("3D reverse switch applications"),
		"show-dock":                    Tr("Show/Hide the dock"),
		"screenshot":                   Tr("Screenshot"),
		"screenshot-full-screen":       Tr("Full screenshot"),
		"screenshot-window":            Tr("Window screenshot"),
		"screenshot-delayed":           Tr("Delay screenshot"),
		"terminal":                     Tr("Terminal"),
		"terminal-quake":               Tr("Terminal Quake Window"),
		"logout":                       Tr("Logout"),
		"disable-touchpad":             Tr("Disable Touchpad"),
	}
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

var MediaNameDescMap map[string]string

func initMediaNameDescMap() {
	MediaNameDescMap = map[string]string{
		"calculator":  Tr("Calculator"),
		"email":       Tr("Email client"),
		"www":         Tr("Web browser"),
		"media":       Tr("Media player"),
		"play":        Tr("Play/Pause"),
		"pause":       Tr("Pause"),
		"stop":        Tr("Stop"),
		"volume-mute": Tr("Mute"),
		"volume-down": Tr("Volume down"),
		"volume-up":   Tr("Volume up"),
		"previous":    Tr("Previous"),
		"next":        Tr("Next"),
		"eject":       Tr("Eject"),
	}
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

var WindowNameDescMap map[string]string

func initWindowNameDescMap() {
	WindowNameDescMap = map[string]string{
		"close":                Tr("Close window"),
		"maximize":             Tr("Maximize window"),
		"unmaximize":           Tr("Restore window"),
		"minimize":             Tr("Minimize window"),
		"begin-move":           Tr("Move window"),
		"begin-resize":         Tr("Resize window"),
		"toggle-shaded":        Tr("Switch furl state"),
		"activate-window-menu": Tr("Activate window menu"),
	}
}

var WorkspaceIdNameMap = map[int32]string{
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

var WorkspaceIdIndexMap = map[int32]int32{
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

var WorkspaceNameDescMap map[string]string

func initWorkspaceNameDescMap() {
	WorkspaceNameDescMap = map[string]string{
		"switch-to-workspace-1":     Tr("Switch to workspace 1"),
		"switch-to-workspace-2":     Tr("Switch to workspace 2"),
		"switch-to-workspace-3":     Tr("Switch to workspace 3"),
		"switch-to-workspace-4":     Tr("Switch to workspace 4"),
		"switch-to-workspace-left":  Tr("Switch to left workspace"),
		"switch-to-workspace-right": Tr("Switch to right workspace"),
		"switch-to-workspace-up":    Tr("Switch to up workspace"),
		"switch-to-workspace-down":  Tr("Switch to down workspace"),
		"put-viewport-1-key":        Tr("Move to workspace 1"),
		"put-viewport-2-key":        Tr("Move to workspace 2"),
		"put-viewport-3-key":        Tr("Move to workspace 3"),
		"put-viewport-4-key":        Tr("Move to workspace 4"),
		"move-to-workspace-left":    Tr("Move to left workspace"),
		"move-to-workspace-right":   Tr("Move to right workspace"),
		"move-to-workspace-up":      Tr("Move to up workspace"),
		"move-to-workspace-down":    Tr("Move to down workspace"),
	}
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
	11:  "disable-touchpad",
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
