/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package appearance

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"sync"

	"gir/gio-2.0"
	ddaemon "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.daemon"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.greeter"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/notify"
	"pkg.deepin.io/lib/utils"
)

var pamEnvFile = os.Getenv("HOME") + "/.pam_environment"

func doGetScaleFactor() float64 {
	var s = gio.NewSettings("com.deepin.xsettings")
	scale := s.GetDouble("scale-factor")
	s.Unref()
	return scale
}

func doSetScaleFactor(scale float64) {
	sendNotify(gettext.Tr("Display scaling"),
		gettext.Tr("Setting display scaling"), "dialog-window-scale")
	var s = gio.NewSettings("com.deepin.xsettings")
	defer s.Unref()
	v := s.GetDouble("scale-factor")
	if scale != v {
		s.SetDouble("scale-factor", scale)
	}

	// for qt
	err := writeKeyToEnvFile("QT_SCALE_FACTOR", fmt.Sprintf("%v", scale), pamEnvFile)
	if err != nil {
		logger.Warning("Failed to set qt scale factor:", err)
	}

	// for java
	doSetJAVAScale(scale)

	// if 1.7 < scale < 2, window scale = 2
	windowScale := int32(math.Trunc((scale+0.3)*10) / 10)
	if windowScale < 1 {
		windowScale = 1
	}
	oldWindowScale := s.GetInt("window-scale")
	if oldWindowScale != windowScale {
		s.SetInt("window-scale", windowScale)
	}

	// set cursor size for deepin-metacity
	gsWrapGDI := gio.NewSettings("com.deepin.wrap.gnome.desktop.interface")
	cursorSize := s.GetInt("gtk-cursor-theme-size")
	gsWrapGDI.SetInt("cursor-size", cursorSize*windowScale)
	gsWrapGDI.Unref()

	doScaleGreeter(scale)
	go func() {
		doScalePlymouth(uint32(windowScale))
		sendNotify(gettext.Tr("Set successfully"),
			gettext.Tr("View by logging out after set display scaling"), "dialog-window-scale")
		setScaleStatus(false)
	}()
}

func doScaleGreeter(scale float64) {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		logger.Warning(err)
		return
	}
	setter := greeter.NewGreeter(systemConn)
	err = setter.SetScaleFactor(0, scale)
	if err != nil {
		logger.Warning("Failed to set greeter scale:", err)
	}
}

func doScalePlymouth(scale uint32) {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		logger.Warning(err)
		return
	}
	setter := ddaemon.NewDaemon(systemConn)
	err = setter.ScalePlymouth(0, scale)
	if err != nil {
		logger.Warning("Failed to scale plymouth:", err)
	}
}

func doSetJAVAScale(scale float64) {
	var envName = "_JAVA_OPTIONS"
	var scaleKey = "-Dswt.autoScale="

	value := os.Getenv(envName)
	if strings.Contains(value, scaleKey) {
		list1 := strings.Split(value, scaleKey)
		value = list1[0]

		list2 := strings.Split(list1[1], " ")
		value += strings.Join(list2[1:], " ")
	}

	value += fmt.Sprintf(" %s%d", scaleKey, int(scale*100))
	err := writeKeyToEnvFile(envName, value, pamEnvFile)
	if err != nil {
		logger.Warning("Failed to set java scale:", value, err)
	}
}

func writeKeyToEnvFile(key, value, filename string) error {
	if !utils.IsFileExist(filename) {
		return ioutil.WriteFile(filename, []byte(key+"="+value+"\n"), 0644)
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var lines = strings.Split(string(content), "\n")
	var idx = -1
	for i, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}
		line = strings.TrimSpace(line)
		if !strings.Contains(line, key+"=") {
			continue
		}

		if line == key+"="+value {
			return nil
		}
		idx = i
		break
	}

	if idx != -1 {
		lines[idx] = key + "=" + value
	} else {
		if lines[len(lines)-1] == "" {
			lines[len(lines)-1] = key + "=" + value
		} else {
			lines = append(lines, key+"="+value)
		}
		lines = append(lines, "")
	}
	return ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
}

var (
	_isScaling   bool = false
	_scaleLocker sync.Mutex
)

func setScaleStatus(status bool) {
	_scaleLocker.Lock()
	_isScaling = status
	_scaleLocker.Unlock()
}

func getScaleStatus() bool {
	_scaleLocker.Lock()
	defer _scaleLocker.Unlock()
	return _isScaling
}

var _notifier *notify.Notification

func sendNotify(summary, body, icon string) {
	if _notifier == nil {
		notify.Init("dde-daemon")
		_notifier = notify.NewNotification(summary, body, icon)
	} else {
		_notifier.Update(summary, body, icon)
	}
	err := _notifier.Show()
	if err != nil {
		logger.Warning("Failed to send notify:", summary, body)
	}
}
