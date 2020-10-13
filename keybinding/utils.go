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
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/dpms"

	wm "github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
	"pkg.deepin.io/dde/daemon/keybinding/util"
	"pkg.deepin.io/lib/strv"

	"pkg.deepin.io/gir/gio-2.0"
	dbus "pkg.deepin.io/lib/dbus1"
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
		//logger.Debug("resetKwin each accel:", accel.Id, accel.Keystrokes, accel.DefaultKeystrokes)
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
			logger.Debug("resetKwin SetAccel", accelJson)
			ok, err := wmObj.SetAccel(0, accelJson)
			// 目前 wm 的实现，调用 SetAccel 如果遇到冲突情况，会导致目标快捷键被清空。
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

func (m *Manager) systemSuspend() {
	can, err := m.sessionManager.CanSuspend(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if !can {
		logger.Info("can not suspend")
		return
	}

	logger.Debug("suspend")
	err = m.sessionManager.RequestSuspend(0)
	if err != nil {
		logger.Warning("failed to suspend:", err)
	}
}

func (m *Manager) systemHibernate() {
	if os.Getenv("POWER_CAN_SLEEP") == "0" {
		logger.Info("can not Hibernate, env POWER_CAN_SLEEP == 0")
		return
	}
	can, err := m.sessionManager.CanHibernate(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if !can {
		logger.Info("can not Hibernate")
		return
	}

	logger.Debug("Hibernate")
	err = m.sessionManager.RequestHibernate(0)
	if err != nil {
		logger.Warning("failed to Hibernate:", err)
	}
}

func (m *Manager) systemShutdown() {
	can, err := m.sessionManager.CanShutdown(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if !can {
		logger.Info("can not Shutdown")
		return
	}

	logger.Debug("Shutdown")
	err = m.sessionManager.RequestShutdown(0)
	if err != nil {
		logger.Warning("failed to Shutdown:", err)
	}
}

func (m *Manager) systemTurnOffScreen() {
	const settingKeyScreenBlackLock = "screen-black-lock"
	logger.Info("DPMS Off")
	var err error
	var useWayland bool
	if len(os.Getenv("WAYLAND_DISPLAY")) != 0 {
		useWayland = true
	} else {
		useWayland = false
	}
	if m.gsPower.GetBoolean(settingKeyScreenBlackLock) {
		m.doLock(true)
	}

	if useWayland {
		err = exec.Command("dde_wldpms", "-s", "Off").Run()
	} else {
		xConn, err := x.NewConn()
		if err != nil {
			logger.Error(err)
		}
		err = dpms.ForceLevelChecked(xConn, dpms.DPMSModeOff).Check(xConn)
	}
	if err != nil {
		logger.Warning("Set DPMS off error:", err)
	}
}

func (m *Manager) systemLogout() {
	can, err := m.sessionManager.CanLogout(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if !can {
		logger.Info("can not logout")
		return
	}

	logger.Debug("logout")
	err = m.sessionManager.RequestLogout(0)
	if err != nil {
		logger.Warning("failed to logout:", err)
	}
}

func (m *Manager) systemAway() {
	err := m.sessionManager.RequestLock(0)
	if err != nil {
		logger.Warning(err)
	}
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

func (m *Manager) doLock(autoStartAuth bool) {
	logger.Info("Lock Screen")
	err := m.lockFront.ShowAuth(0, autoStartAuth)
	if err != nil {
		logger.Warning("failed to call lockFront ShowAuth:", err)
	}
}
