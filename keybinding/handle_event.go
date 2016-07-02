/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

import (
	"os/exec"
	"strings"

	"gir/gio-2.0"
	"pkg.deepin.io/dde/daemon/keybinding/core"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/lib/dbus"
)

const (
	cmdDDEOSD = "/usr/lib/deepin-daemon/dde-osd"
)

func (m *Manager) handleKeyEvent(mod uint16, code int, pressed bool) {
	modStr, codeStr, _ := core.LookupKeyEvent(mod, code)
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
	dbus.Emit(m.media, signal, pressed)
	if pressed {
		if !canShowOSD(id) {
			return
		}
		go showOSD(signal)
	}
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
	case "kbd-brightness-up":
		return "KbdBrightnessUp"
	case "kbd-brightness-down":
		return "KbdBrightnessDown"
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
	case "suspend":
		return "PowerSuspend"
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
	case "www", "home":
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
	case "eject":
		return "Eject"
	case "media":
		return "AudioMedia"
	}

	return ""
}

func doAction(cmd string) error {
	if len(cmd) == 0 {
		return nil
	}

	return exec.Command("/bin/sh", "-c", cmd).Run()
}

func showOSD(signal string) {
	sessionDBus, _ := dbus.SessionBus()
	sessionDBus.Object("com.deepin.dde.osd", "/").Call("com.deepin.dde.osd.ShowOSD", 0, signal)
}

func canShowOSD(id string) bool {
	switch id {
	case "mon-brightness-up", "mon-brightness-down", "volume-up", "volume-down":
		// Move to mpris
		return false
	case "capslock":
		return canShowCapsOSD()

	}
	return true
}

func canShowCapsOSD() bool {
	s := gio.NewSettings("com.deepin.dde.keyboard")
	defer s.Unref()
	return s.GetBoolean("capslock-toggle")
}
