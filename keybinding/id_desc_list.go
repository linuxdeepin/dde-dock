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
	Id   int32
	Name string
	Desc string
}

var systemIdDescList []idDescInfo

func initSystemIdDescList() {
	systemIdDescList = []idDescInfo{
		idDescInfo{0, "launcher", Tr("Launcher")},
		idDescInfo{620, "show-desktop", Tr("Show desktop")},
		idDescInfo{2, "lock-screen", Tr("Lock screen")},
		idDescInfo{10, "file-manager", Tr("File manager")},
		idDescInfo{621, "switch-applications", Tr("Switch applications")},
		idDescInfo{622, "switch-applications-backward", Tr("Reverse switch applications")},
		idDescInfo{800, "prev-key", Tr("3D switch applications")},
		idDescInfo{801, "next-key", Tr("3D reverse switch applications")},
		idDescInfo{3, "show-dock", Tr("Show/Hide the dock")},
		idDescInfo{6, "screenshot", Tr("Screenshot")},
		idDescInfo{7, "screenshot-full-screen", Tr("Full screenshot")},
		idDescInfo{8, "screenshot-window", Tr("Window screenshot")},
		idDescInfo{9, "screenshot-delayed", Tr("Delay screenshot")},
		idDescInfo{1, "terminal", Tr("Terminal")},
		idDescInfo{5, "terminal-quake", Tr("Terminal Quake Window")},
		idDescInfo{4, "logout", Tr("Logout")},
		idDescInfo{12, "deepin-translator", Tr("Disable Touchpad")},
		idDescInfo{13, "switch-layout", Tr("Deepin Translator")},
		idDescInfo{11, "disable-touchpad", Tr("Switch Layout")},
	}
}

var mediaIdDescList []idDescInfo

func initMediaIdDescList() {
	mediaIdDescList = []idDescInfo{
		idDescInfo{300, "calculator", Tr("Calculator")},
		idDescInfo{302, "email", Tr("Email client")},
		idDescInfo{303, "www", Tr("Web browser")},
		idDescInfo{304, "media", Tr("Media player")},
		idDescInfo{305, "play", Tr("Play/Pause")},
		idDescInfo{306, "pause", Tr("Pause")},
		idDescInfo{307, "stop", Tr("Stop")},
		idDescInfo{309, "volume-mute", Tr("Mute")},
		idDescInfo{308, "volume-down", Tr("Volume down")},
		idDescInfo{310, "volume-up", Tr("Volume up")},
		idDescInfo{311, "previous", Tr("Previous")},
		idDescInfo{312, "next", Tr("Next")},
		idDescInfo{301, "eject", Tr("Eject")},
	}
}

var windowIdDescList []idDescInfo

func initWindowIdDescList() {
	windowIdDescList = []idDescInfo{
		idDescInfo{603, "close", Tr("Close window")},
		idDescInfo{604, "maximize", Tr("Maximize window")},
		idDescInfo{607, "unmaximize", Tr("Restore window")},
		idDescInfo{605, "minimize", Tr("Minimize window")},
		idDescInfo{601, "begin-move", Tr("Move window")},
		idDescInfo{602, "begin-resize", Tr("Resize window")},
		idDescInfo{606, "toggle-shaded", Tr("Switch furl state")},
		idDescInfo{600, "activate-window-menu", Tr("Activate window menu")},
	}
}

var workspaceIdDescList []idDescInfo

func initWorkspaceIdDescList() {
	workspaceIdDescList = []idDescInfo{
		idDescInfo{608, "switch-to-workspace-1", Tr("Switch to workspace 1")},
		idDescInfo{609, "switch-to-workspace-2", Tr("Switch to workspace 2")},
		idDescInfo{610, "switch-to-workspace-3", Tr("Switch to workspace 3")},
		idDescInfo{611, "switch-to-workspace-4", Tr("Switch to workspace 4")},
		idDescInfo{613, "switch-to-workspace-left", Tr("Switch to left workspace")},
		idDescInfo{614, "switch-to-workspace-right", Tr("Switch to right workspace")},
		idDescInfo{615, "switch-to-workspace-up", Tr("Switch to upper workspace")},
		idDescInfo{612, "switch-to-workspace-down", Tr("Switch to lower workspace")},
		idDescInfo{900, "put-viewport-1-key", Tr("Move to workspace 1")},
		idDescInfo{901, "put-viewport-2-key", Tr("Move to workspace 2")},
		idDescInfo{902, "put-viewport-3-key", Tr("Move to workspace 3")},
		idDescInfo{903, "put-viewport-4-key", Tr("Move to workspace 4")},
		idDescInfo{617, "move-to-workspace-left", Tr("Move to left workspace")},
		idDescInfo{618, "move-to-workspace-right", Tr("Move to right workspace")},
		idDescInfo{619, "move-to-workspace-up", Tr("Move to upper workspace")},
		idDescInfo{616, "move-to-workspace-down", Tr("Move to lower workspace")},
	}
}
