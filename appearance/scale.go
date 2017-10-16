/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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
	"dbus/com/deepin/daemon/greeter"
	"fmt"
	"gir/gio-2.0"
	"io/ioutil"
	"math"
	"os"
	"pkg.deepin.io/lib/utils"
	"strings"
)

var pamEnvFile = os.Getenv("HOME") + "/.pam_environment"

func doGetScaleFactor() float64 {
	var s = gio.NewSettings("com.deepin.xsettings")
	scale := s.GetDouble("scale-factor")
	s.Unref()
	return scale
}

func doSetScaleFactor(scale float64) {
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

	// if 1.7 < scale < 2, window scale = 2
	tmp := int32(math.Trunc((scale+0.3)*10) / 10)
	if tmp < 1 {
		tmp = 1
	}
	window := s.GetInt("window-scale")
	if window != tmp {
		s.SetInt("window-scale", tmp)
	}

	doSetGreeterScale(scale)
}

func doSetGreeterScale(scale float64) {
	setter, err := greeter.NewGreeter("com.deepin.daemon.Greeter", "/com/deepin/daemon/Greeter")
	if err != nil {
		logger.Warning("Failed to create greeter setter connection:", err)
		return
	}

	err = setter.SetScaleFactor(scale)
	if err != nil {
		logger.Warning("Failed to set greeter scale:", err)
	}
	setter = nil
}

func writeKeyToEnvFile(key, value, filename string) error {
	if !utils.IsFileExist(filename) {
		return ioutil.WriteFile(filename, []byte(key+"="+value), 0644)
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
		lines[len(lines)-1] = key + "=" + value
	}
	return ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
}
