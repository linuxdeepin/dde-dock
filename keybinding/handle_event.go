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

package keybinding

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/keybind"
	"os/exec"
	"pkg.deepin.io/dde/daemon/keybinding/core"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/lib/dbus"
	"strings"
)

const (
	cmdDDEOSD = "/usr/lib/deepin-daemon/dde-osd"
)

func (m *Manager) handleKeyEvent(mod uint16, code int, pressed bool) {
	modStr := keybind.ModifierString(mod)
	codeStr := keybind.LookupString(m.xu, mod, xproto.Keycode(code))
	if code == 65 {
		codeStr = "space"
	}
	logger.Debug("Handle key event:", mod, modStr, code, codeStr, pressed)

	var accel = codeStr
	if len(modStr) != 0 {
		accel = modStr + "-" + codeStr
	}
	s := m.grabedList.GetByAccel(accel)
	if s == nil {
		logger.Debugf("'%s' not in grabed list", accel)
		return
	}

	switch s.Type {
	case shortcuts.KeyTypeSystem, shortcuts.KeyTypeCustom:
		if !pressed {
			return
		}
		logger.Debug("Exec action:", s.GetAction())
		go doAction(s.GetAction())
	case shortcuts.KeyTypeMedia:
		m.handleMediaEvent(modStr, codeStr, pressed)
	}
	return
}

func (m *Manager) handleMediaEvent(modStr, codeStr string, pressed bool) {
	signal := getMediakeySignal(modStr, codeStr)
	if len(signal) == 0 {
		return
	}

	//TODO: emit signal
	logger.Debug("Emit signal:", signal)
	if pressed {
		go doAction(cmdDDEOSD + " --" + signal)
	}
	dbus.Emit(m.media, signal, pressed)
}

func getMediakeySignal(modStr, codeStr string) string {
	switch codeStr {
	case "XF86MonBrightnessUp":
		return "BrightnessUp"
	case "XF86MonBrightnessDown":
		return "BrightnessDown"
	case "XF86AudioMute":
		return "AudioMute"
	case "XF86AudioLowerVolume":
		return "AudioDown"
	case "XF86AudioRaiseVolume":
		return "AudioUp"
	case "Num_Lock":
		// num_lock --> mod2
		if strings.Contains(modStr, "mod2") {
			return "NumLockOff"
		}
		return "NumLockOn"
	case "Caps_Lock":
		// caps_lock --> lock
		if strings.Contains(modStr, "lock") {
			return "CapsLockOff"
		}
		return "CapsLockOn"
	case "XF86TouchpadToggle":
		return "TouchpadToggle"
	case "XF86TouchpadOn":
		return "TouchpadOn"
	case "XF86TouchpadOff":
		return "TouchpadOff"
	case "XF86Display":
		return "SwitchMonitors"
	case "XF86PowerOff":
		return "PowerOff"
	case "XF86Sleep":
		return "PowerSleep"
	//case "p", "P":
	case "XF86AudioPlay":
		return "AudioPlay"
	case "XF86AudioPause":
		return "AudioPause"
	case "XF86AudioStop":
		return "AudioStop"
	case "XF86AudioPrev":
		return "AudioPrevious"
	case "XF86AudioNext":
		return "AudioNext"
	case "XF86AudioRewind":
		return "AudioRewind"
	case "XF86AudioForward":
		return "AudioForward"
	case "XF86AudioRepeat":
		return "AudioRepeat"
	case "XF86WWW":
		return "LaunchBrowser"
	case "XF86Mail":
		return "LaunchEmail"
	case "XF86Calculator":
		return "LaunchCalculator"
	}

	if len(modStr) != 0 {
		codeStr = modStr + "-" + codeStr
	}
	switch {
	case core.IsAccelEqual("mod4-p", codeStr):
		return "SwitchMonitors"
	case core.IsAccelEqual("mod4-space", codeStr):
		// TODO: check can switch layout
		return "SwitchLayout"
	}

	return ""
}

func doAction(cmd string) error {
	if len(cmd) == 0 {
		return nil
	}

	return exec.Command("/bin/sh", "-c", cmd).Run()
}
