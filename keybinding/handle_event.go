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
	"os/exec"
	"strings"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/keybind"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/lib/dbus"
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

		if s.Id == "switch-layout" {
			m.handleMediaEvent(s.Id, s.Type, modStr, pressed)
			break
		}

		logger.Debug("Exec action:", s.GetAction())
		go doAction(s.GetAction())
	case shortcuts.KeyTypeMedia:
		m.handleMediaEvent(s.Id, s.Type, modStr, pressed)
	}
	return
}

func (m *Manager) handleMediaEvent(id string, ty int32, modStr string, pressed bool) {
	signal := getMediakeySignal(id, ty, modStr)
	if len(signal) == 0 {
		return
	}

	logger.Debug("Emit signal:", signal)
	if pressed {
		go doAction(cmdDDEOSD + " --" + signal)
	}
	dbus.Emit(m.media, signal, pressed)
}

func getMediakeySignal(id string, ty int32, modStr string) string {
	if ty != shortcuts.KeyTypeSystem &&
		ty != shortcuts.KeyTypeMedia {
		return ""
	}

	switch id {
	case "mon-brightness-up":
		return "BrightnessUp"
	case "mon-brightness-down":
		return "BrightnessDown"
	case "volume-mute":
		return "AudioMute"
	case "volume-down":
		return "AudioDown"
	case "volume-up":
		return "AudioUp"
	case "numlock":
		// num_lock --> mod2
		if strings.Contains(modStr, "mod2") {
			return "NumLockOff"
		}
		return "NumLockOn"
	case "capslock":
		// caps_lock --> lock
		if strings.Contains(modStr, "lock") {
			return "CapsLockOff"
		}
		return "CapsLockOn"
	case "touchpad-toggle":
		return "TouchpadToggle"
	case "touchpad-on":
		return "TouchpadOn"
	case "touchpad-off":
		return "TouchpadOff"
	case "display":
		return "SwitchMonitors"
	case "power-off":
		return "PowerOff"
	case "sleep":
		return "PowerSleep"
	//case "p", "P":
	case "play":
		return "AudioPlay"
	case "pause":
		return "AudioPause"
	case "stop":
		return "AudioStop"
	case "previous":
		return "AudioPrevious"
	case "next":
		return "AudioNext"
	case "audio-rewind":
		return "AudioRewind"
	case "audio-forward":
		return "AudioForward"
	case "audio-repeat":
		return "AudioRepeat"
	case "www":
		return "LaunchBrowser"
	case "email":
		return "LaunchEmail"
	case "calculator":
		return "LaunchCalculator"
	case "switch-monitors":
		return "SwitchMonitors"
	// system type
	case "switch-layout":
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
