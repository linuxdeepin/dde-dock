/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package shortcuts

import (
	"pkg.deepin.io/lib/gettext"
)

func getSystemIdNameMap() map[string]string {
	var idNameMap = map[string]string{
		"launcher":               gettext.Tr("Launcher"),
		"terminal":               gettext.Tr("Terminal"),
		"deepin-screen-recorder": gettext.Tr("Screen Recorder"),
		"lock-screen":            gettext.Tr("Lock screen"),
		"show-dock":              gettext.Tr("Show/Hide the dock"),
		"logout":                 gettext.Tr("Shutdown interface"),
		"terminal-quake":         gettext.Tr("Terminal Quake Window"),
		"screenshot":             gettext.Tr("Screenshot"),
		"screenshot-fullscreen":  gettext.Tr("Full screenshot"),
		"screenshot-window":      gettext.Tr("Window screenshot"),
		"screenshot-delayed":     gettext.Tr("Delay screenshot"),
		"file-manager":           gettext.Tr("File manager"),
		"disable-touchpad":       gettext.Tr("Disable Touchpad"),
		"wm-switcher":            gettext.Tr("Switch window effects"),
		"turn-off-screen":        gettext.Tr("Fast Screen Off"),
		"system-monitor":         gettext.Tr("System Monitor"),
		"color-picker":           gettext.Tr("Deepin Picker"),
		"ai-assistant":           gettext.Tr("Desktop AI Assistant"),
		"text-to-speech":         gettext.Tr("Text to Speech"),
		"speech-to-text":         gettext.Tr("Speech to Text"),
		"clipboard":              gettext.Tr("Clipboard"),
		"translation":            gettext.Tr("Translation"),
	}
	return idNameMap
}

func getSpecialIdNameMap() map[string]string {
	var idNameMap = map[string]string{
		"switch-kbd-layout": gettext.Tr("Switch Layout"),
	}
	return idNameMap
}

func getWMIdNameMap() map[string]string {
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
		"minimize":                     gettext.Tr("Minimize window"),
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
		"expose-all-windows":           gettext.Tr("Display windows of all workspaces"),
		"expose-windows":               gettext.Tr("Display windows of current workspace"),
		"preview-workspace":            gettext.Tr("Display workspace"),
	}
	return idNameMap
}

func getMediaIdNameMap() map[string]string {
	var idNameMap = map[string]string{
		"messenger":           "Messenger",         // XF86Messenger
		"save":                "Save",              // XF86Save
		"new":                 "New",               // XF86New
		"wake-up":             "WakeUp",            // XF86WakeUp
		"audio-rewind":        "AudioRewind",       // XF86AudioRewind
		"audio-mute":          "AudioMute",         // XF86AudioMute
		"mon-brightness-up":   "MonBrightnessUp",   // XF86MonBrightnessUp
		"wlan":                "WLAN",              // XF86WLAN
		"audio-media":         "AudioMedia",        // XF86AudioMedia
		"reply":               "Reply",             // XF86Reply
		"favorites":           "Favorites",         // XF86Favorites
		"audio-play":          "AudioPlay",         // XF86AudioPlay
		"audio-mic-mute":      "AudioMicMute",      // XF86AudioMicMute
		"audio-pause":         "AudioPause",        // XF86AudioPause
		"audio-stop":          "AudioStop",         // XF86AudioStop
		"power-off":           "PowerOff",          // XF86PowerOff
		"documents":           "Documents",         // XF86Documents
		"game":                "Game",              // XF86Game
		"search":              "Search",            // XF86Search
		"audio-record":        "AudioRecord",       // XF86AudioRecord
		"display":             "Display",           // XF86Display
		"reload":              "Reload",            // XF86Reload
		"explorer":            "Explorer",          // XF86Explorer
		"calculator":          "Calculator",        // XF86Calculator
		"calendar":            "Calendar",          // XF86Calendar
		"forward":             "Forward",           // XF86Forward
		"cut":                 "Cut",               // XF86Cut
		"mon-brightness-down": "MonBrightnessDown", // XF86MonBrightnessDown
		"copy":                "Copy",              // XF86Copy
		"tools":               "Tools",             // XF86Tools
		"audio-raise-volume":  "AudioRaiseVolume",  // XF86AudioRaiseVolume
		"close":               "Close",             // XF86Close
		"www":                 "WWW",               // XF86WWW
		"home-page":           "HomePage",          // XF86HomePage
		"sleep":               "Sleep",             // XF86Sleep
		"audio-lower-volume":  "AudioLowerVolume",  // XF86AudioLowerVolume
		"audio-prev":          "AudioPrev",         // XF86AudioPrev
		"audio-next":          "AudioNext",         // XF86AudioNext
		"paste":               "Paste",             // XF86Paste
		"open":                "Open",              // XF86Open
		"send":                "Send",              // XF86Send
		"my-computer":         "MyComputer",        // XF86MyComputer
		"mail":                "Mail",              // XF86Mail
		"adjust-brightness":   "BrightnessAdjust",  // XF86BrightnessAdjust
		"log-off":             "LogOff",            // XF86LogOff
		"pictures":            "Pictures",          // XF86Pictures
		"terminal":            "Terminal",          // XF86Terminal
		"video":               "Video",             // XF86Video
		"music":               "Music",             // XF86Music
		"app-left":            "ApplicationLeft",   // XF86ApplicationLeft
		"app-right":           "ApplicationRight",  // XF86ApplicationRight
		"meeting":             "Meeting",           // XF86Meeting
		"switch-monitors":     gettext.Tr("Switch monitors"),
	}
	return idNameMap
}
