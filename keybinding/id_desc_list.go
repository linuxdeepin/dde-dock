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

import . "pkg.deepin.io/lib/gettext"

type idDescInfo struct {
	Id   int32
	Name string
	Desc string
}

var systemIdDescList []idDescInfo

/**
 * Id Range
 * 600 ~ 649: org.compiz.core.gschema.xml
 * 650 ~ 699: org.compiz.move.gschema.xml
 * 700 ~ 749: org.compiz.resize.gschema.xml
 * 750 ~ 799: org.compiz.vpswitch.gschema.xml
 * 800 ~ 849: org.compiz.put.gschema.xml
 * 850 ~ 899: org.compiz.wall.gschema.xml
 * 900 ~ 949: org.compiz.shift.gschema.xml
 * 950 ~ 999: org.compiz.switcher.gschema.xml
 **/
func initSystemIdDescList() {
	systemIdDescList = []idDescInfo{
		idDescInfo{0, "launcher", Tr("Launcher")},
		idDescInfo{600, "show-desktop", Tr("Show desktop")},
		idDescInfo{2, "lock-screen", Tr("Lock screen")},
		idDescInfo{10, "file-manager", Tr("File manager")},
		idDescInfo{950, "switch-applications", Tr("Switch applications")},
		idDescInfo{951, "switch-applications-backward", Tr("Reverse switch applications")},
		idDescInfo{900, "prev-key", Tr("3D switch applications")},
		idDescInfo{901, "next-key", Tr("3D reverse switch applications")},
		idDescInfo{3, "show-dock", Tr("Show/Hide the dock")},
		idDescInfo{6, "screenshot", Tr("Screenshot")},
		idDescInfo{7, "screenshot-full-screen", Tr("Full screenshot")},
		idDescInfo{8, "screenshot-window", Tr("Window screenshot")},
		idDescInfo{9, "screenshot-delayed", Tr("Delay screenshot")},
		idDescInfo{1, "terminal", Tr("Terminal")},
		idDescInfo{5, "terminal-quake", Tr("Terminal Quake Window")},
		idDescInfo{4, "logout", Tr("Logout")},
		idDescInfo{12, "disable-touchpad", Tr("Disable Touchpad")},
		idDescInfo{13, "deepin-translator", Tr("Deepin Translator")},
		idDescInfo{11, "switch-layout", Tr("Switch Layout")},
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
		idDescInfo{601, "close", Tr("Close window")},
		idDescInfo{602, "maximize", Tr("Maximize window")},
		idDescInfo{603, "unmaximize", Tr("Restore window")},
		idDescInfo{604, "minimize", Tr("Minimize window")},
		idDescInfo{650, "begin-move", Tr("Move window")},
		idDescInfo{700, "begin-resize", Tr("Resize window")},
		idDescInfo{605, "toggle-shaded", Tr("Switch furl state")},
		idDescInfo{606, "activate-window-menu", Tr("Activate window menu")},
	}
}

var workspaceIdDescList []idDescInfo

func initWorkspaceIdDescList() {
	workspaceIdDescList = []idDescInfo{
		idDescInfo{750, "switch-to-workspace-1", Tr("Switch to workspace 1")},
		idDescInfo{751, "switch-to-workspace-2", Tr("Switch to workspace 2")},
		idDescInfo{752, "switch-to-workspace-3", Tr("Switch to workspace 3")},
		idDescInfo{753, "switch-to-workspace-4", Tr("Switch to workspace 4")},
		idDescInfo{850, "switch-to-workspace-left", Tr("Switch to left workspace")},
		idDescInfo{851, "switch-to-workspace-right", Tr("Switch to right workspace")},
		idDescInfo{852, "switch-to-workspace-up", Tr("Switch to upper workspace")},
		idDescInfo{853, "switch-to-workspace-down", Tr("Switch to lower workspace")},
		idDescInfo{800, "put-viewport-1-key", Tr("Move to workspace 1")},
		idDescInfo{801, "put-viewport-2-key", Tr("Move to workspace 2")},
		idDescInfo{802, "put-viewport-3-key", Tr("Move to workspace 3")},
		idDescInfo{803, "put-viewport-4-key", Tr("Move to workspace 4")},
		idDescInfo{854, "move-to-workspace-left", Tr("Move to left workspace")},
		idDescInfo{855, "move-to-workspace-right", Tr("Move to right workspace")},
		idDescInfo{856, "move-to-workspace-up", Tr("Move to upper workspace")},
		idDescInfo{857, "move-to-workspace-down", Tr("Move to lower workspace")},
	}
}
