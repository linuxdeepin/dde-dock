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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"strconv"
	"strings"
	"time"
)

const (
	dbusDest = "com.deepin.daemon.helper.Backlight"
	dbusPath = "/com/deepin/daemon/helper/Backlight"
	dbusIFC  = "com.deepin.daemon.helper.Backlight"
)

const (
	DisplayBacklight byte = iota + 1
	KeyboardBacklight
)

var logger = log.NewLogger("backlight_helper")

type Manager struct{}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) SetBrightness(type_ byte, name string, value int32) error {
	filename, err := getBrightnessFilename(type_, name)
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = fh.WriteString(strconv.Itoa(int(value)))
	if err != nil {
		return err
	}

	return nil
}

func getBrightnessFilename(type_ byte, name string) (string, error) {
	// check type_
	var subsystem string
	switch type_ {
	case DisplayBacklight:
		subsystem = "backlight"
	case KeyboardBacklight:
		subsystem = "leds"
	default:
		return "", fmt.Errorf("invalid type %d", type_)
	}

	// check name
	if strings.ContainsRune(name, '/') || name == "" ||
		name == "." || name == ".." {
		return "", fmt.Errorf("invalid name %q", name)
	}

	return filepath.Join("/sys/class", subsystem, name, "brightness"), nil
}

func main() {
	m := &Manager{}
	err := dbus.InstallOnSystem(m)
	if err != nil {
		logger.Error("Install session bus failed:", err)
		return
	}
	dbus.SetAutoDestroyHandler(time.Second*10, nil)
	dbus.DealWithUnhandledMessage()
	err = dbus.Wait()
	if err != nil {
		logger.Error("Lost dbus connection:", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
