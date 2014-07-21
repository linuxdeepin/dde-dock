/**
 * Copyright (c) 2011 ~ 2014 Deepin}, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author,      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer,  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License}, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful},
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not}, see <http,//www.gnu.org/licenses/>.
 **/

package keybinding

import . "pkg.linuxdeepin.com/lib/gettext"

type idDescInfo struct {
	Id    int32
	Name  string //gsetting key
	Desc  string
	Index int32
}

var systemIdDescList []idDescInfo

func initSystemIdDescList() {
	systemIdDescList = []idDescInfo{
		idDescInfo{0, "launcher", Tr("Launcher"), 0},
		idDescInfo{620, "show-desktop", Tr("Show desktop"), 1},
		idDescInfo{2, "lock-screen", Tr("Lock screen"), 2},
		idDescInfo{10, "file-manager", Tr("File manager"), 3},
		idDescInfo{621, "switch-applications", Tr("Switch applications"), 4},
		idDescInfo{622, "switch-applications-backward", Tr("Reverse switch applications"), 5},
		idDescInfo{800, "prev-key", Tr("3D switch applications"), 6},
		idDescInfo{801, "next-key", Tr("3D reverse switch applications"), 7},
		idDescInfo{3, "show-dock", Tr("Show/Hide the dock"), 8},
		idDescInfo{6, "screenshot", Tr("Screenshot"), 9},
		idDescInfo{7, "screenshot-full-screen", Tr("Full screenshot"), 10},
		idDescInfo{8, "screenshot-window", Tr("Window screenshot"), 11},
		idDescInfo{9, "screenshot-delayed", Tr("Delay screenshot"), 12},
		idDescInfo{1, "terminal", Tr("Terminal"), 13},
		idDescInfo{5, "terminal-quake", Tr("Terminal Quake Window"), 14},
		idDescInfo{4, "logout", Tr("Logout"), 15},
		idDescInfo{12, "deepin-translator", Tr("Disable Touchpad"), 16},
		idDescInfo{13, "switch-layout", Tr("Deepin Translator"), 17},
		idDescInfo{11, "disable-touchpad", Tr("Switch Layout"), 18},
	}
}

var mediaIdDescList []idDescInfo

func initMediaIdDescList() {
	mediaIdDescList = []idDescInfo{
		idDescInfo{300, "calculator", Tr("Calculator"), 0},
		idDescInfo{302, "email", Tr("Email client"), 1},
		idDescInfo{303, "www", Tr("Web browser"), 2},
		idDescInfo{304, "media", Tr("Media player"), 3},
		idDescInfo{305, "play", Tr("Play/Pause"), 4},
		idDescInfo{306, "pause", Tr("Pause"), 5},
		idDescInfo{307, "stop", Tr("Stop"), 6},
		idDescInfo{309, "volume-mute", Tr("Mute"), 7},
		idDescInfo{308, "volume-down", Tr("Volume down"), 8},
		idDescInfo{310, "volume-up", Tr("Volume up"), 9},
		idDescInfo{311, "previous", Tr("Previous"), 10},
		idDescInfo{312, "next", Tr("Next"), 11},
		idDescInfo{301, "eject", Tr("Eject"), 12},
	}
}

var windowIdDescList []idDescInfo

func initWindowIdDescList() {
	windowIdDescList = []idDescInfo{
		idDescInfo{603, "close", Tr("Close window"), 0},
		idDescInfo{604, "maximize", Tr("Maximize window"), 1},
		idDescInfo{607, "unmaximize", Tr("Restore window"), 2},
		idDescInfo{605, "minimize", Tr("Minimize window"), 3},
		idDescInfo{601, "begin-move", Tr("Move window"), 4},
		idDescInfo{602, "begin-resize", Tr("Resize window"), 5},
		idDescInfo{606, "toggle-shaded", Tr("Switch furl state"), 6},
		idDescInfo{600, "activate-window-menu", Tr("Activate window menu"), 7},
	}
}

var workspaceIdDescList []idDescInfo

func initWorkspaceIdDescList() {
	workspaceIdDescList = []idDescInfo{
		idDescInfo{608, "switch-to-workspace-1", Tr("Switch to workspace 1"), 0},
		idDescInfo{609, "switch-to-workspace-2", Tr("Switch to workspace 2"), 1},
		idDescInfo{610, "switch-to-workspace-3", Tr("Switch to workspace 3"), 2},
		idDescInfo{611, "switch-to-workspace-4", Tr("Switch to workspace 4"), 3},
		idDescInfo{613, "switch-to-workspace-left", Tr("Switch to left workspace"), 4},
		idDescInfo{614, "switch-to-workspace-right", Tr("Switch to right workspace"), 5},
		idDescInfo{615, "switch-to-workspace-up", Tr("Switch to upper workspace"), 6},
		idDescInfo{612, "switch-to-workspace-down", Tr("Switch to lower workspace"), 7},
		idDescInfo{900, "put-viewport-1-key", Tr("Move to workspace 1"), 8},
		idDescInfo{901, "put-viewport-2-key", Tr("Move to workspace 2"), 9},
		idDescInfo{902, "put-viewport-3-key", Tr("Move to workspace 3"), 10},
		idDescInfo{903, "put-viewport-4-key", Tr("Move to workspace 4"), 11},
		idDescInfo{617, "move-to-workspace-left", Tr("Move to left workspace"), 12},
		idDescInfo{618, "move-to-workspace-right", Tr("Move to right workspace"), 13},
		idDescInfo{619, "move-to-workspace-up", Tr("Move to upper workspace"), 14},
		idDescInfo{616, "move-to-workspace-down", Tr("Move to lower workspace"), 15},
	}
}
