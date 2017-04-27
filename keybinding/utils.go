/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

import (
	"encoding/json"
	"gir/gio-2.0"
	"io/ioutil"
	"os"
	"os/exec"
	"pkg.deepin.io/lib/dbus"
	"strings"
)

func canShowCapsOSD() bool {
	s := gio.NewSettings("com.deepin.dde.keyboard")
	defer s.Unref()
	return s.GetBoolean("capslock-toggle")
}

func getPowerButtonPressedExec() string {
	s := gio.NewSettings("com.deepin.dde.power")
	defer s.Unref()
	return s.GetString("power-button-pressed-exec")
}

func resetGSettings(gs *gio.Settings) {
	for _, key := range gs.ListKeys() {
		userVal := gs.GetUserValue(key)
		if userVal != nil {
			// TODO unref userVal
			logger.Debug("reset gsettings key", key)
			gs.Reset(key)
		}
	}
}

func doMarshal(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func execCmd(cmd string) error {
	if len(cmd) == 0 {
		logger.Debug("cmd is empty")
		return nil
	}

	logger.Debugf("len environ: %d, exec cmd: %q", len(os.Environ()), cmd)
	return exec.Command("/bin/sh", "-c", cmd).Run()
}

func showOSD(signal string) {
	logger.Debug("show OSD", signal)
	sessionDBus, _ := dbus.SessionBus()
	go sessionDBus.Object("com.deepin.dde.osd", "/").Call("com.deepin.dde.osd.ShowOSD", 0, signal)
}

const sessionManagerDest = "com.deepin.SessionManager"
const sessionManagerObjPath = "/com/deepin/SessionManager"

func systemSuspend() {
	sessionDBus, _ := dbus.SessionBus()
	go sessionDBus.Object(sessionManagerDest, sessionManagerObjPath).Call(sessionManagerDest+".RequestSuspend", 0)
}

func queryCommandByMime(mime string) string {
	app := gio.AppInfoGetDefaultForType(mime, false)
	if app == nil {
		return ""
	}
	defer app.Unref()

	return app.GetExecutable()
}

const (
	ibmHotkeyFile = "/proc/acpi/ibm/hotkey"
)

var driverSupportedHotkey = func() func() bool {
	var (
		init      bool = false
		supported bool = false
	)

	return func() bool {
		if !init {
			init = true
			supported = checkIBMHotkey(ibmHotkeyFile)
		}
		return supported
	}
}()

func checkIBMHotkey(file string) bool {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		list := strings.Split(line, ":")
		if len(list) != 2 {
			continue
		}

		if list[0] != "status" {
			continue
		}

		if strings.TrimSpace(list[1]) == "enabled" {
			return true
		}
	}
	return false
}
