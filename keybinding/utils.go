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

package keybinding

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	wm "github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
	"pkg.deepin.io/dde/daemon/keybinding/util"
	"pkg.deepin.io/lib/strv"

	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbus1"
)

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

func resetKWin(wmObj *wm.Wm) error {
	accels, err := util.GetAllKWinAccels(wmObj)
	if err != nil {
		return err
	}
	for _, accel := range accels {
		if !strv.Strv(accel.Keystrokes).Equal(accel.DefaultKeystrokes) &&
			len(accel.DefaultKeystrokes) > 0 && accel.DefaultKeystrokes[0] != "" {
			accelJson, err := util.MarshalJSON(&util.KWinAccel{
				Id:         accel.Id,
				Keystrokes: accel.DefaultKeystrokes,
			})
			if err != nil {
				logger.Warning(err)
				continue
			}
			ok, err := wmObj.SetAccel(0, accelJson)
			if !ok {
				logger.Warning("wm.SetAccel failed, id: ", accel.Id)
				continue
			}
			if err != nil {
				logger.Warning("failed to set accel:", err, accel.Id)
				continue
			}
		}
	}
	return nil
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

func getRfkillWlanState() (int, error) {
	dir := "/sys/class/rfkill"
	fileInfoList, err := ioutil.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	for _, fileInfo := range fileInfoList {
		typeFile := filepath.Join(dir, fileInfo.Name(), "type")
		typeBytes, err := readTinyFile(typeFile)
		if err != nil {
			continue
		}
		if bytes.Equal(bytes.TrimSpace(typeBytes), []byte("wlan")) {
			stateFile := filepath.Join(dir, fileInfo.Name(), "state")
			stateBytes, err := readTinyFile(stateFile)
			if err != nil {
				return 0, err
			}
			stateBytes = bytes.TrimSpace(stateBytes)
			state, err := strconv.Atoi(string(stateBytes))
			if err != nil {
				return 0, err
			}
			return state, nil

		}
	}
	return 0, errors.New("not found rfkill with type wlan")
}

func readTinyFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := make([]byte, 8)
	n, err := f.Read(buf)

	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func shouldUseDDEKwin() bool {
	_, err := os.Stat("/usr/bin/kwin_no_scale")
	return err == nil
}
