/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package inputdevices

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus/property"
	"pkg.deepin.io/lib/xdg/basedir"
)

const (
	gsSchemaInputDevices = "com.deepin.dde.inputdevices"
	gsKeyWheelSpeed      = "wheel-speed"
	imWheelBin           = "imwheel"
)

type devicePathInfo struct {
	Path string
	Type string
}
type devicePathInfos []*devicePathInfo

type Manager struct {
	Infos             devicePathInfos
	settings          *gio.Settings
	imWheelConfigFile string
	WheelSpeed        *property.GSettingsUintProperty `access:"readwrite"`

	kbd        *Keyboard
	mouse      *Mouse
	trackPoint *TrackPoint
	tpad       *Touchpad
	wacom      *Wacom
}

func NewManager() *Manager {
	var m = new(Manager)
	m.imWheelConfigFile = filepath.Join(basedir.GetUserHomeDir(), ".imwheelrc")

	m.Infos = devicePathInfos{
		&devicePathInfo{
			Path: "com.deepin.daemon.InputDevice.Keyboard",
			Type: "keyboard",
		},
		&devicePathInfo{
			Path: "com.deepin.daemon.InputDevice.Mouse",
			Type: "mouse",
		},
		&devicePathInfo{
			Path: "com.deepin.daemon.InputDevice.TrackPoint",
			Type: "trackpoint",
		},
		&devicePathInfo{
			Path: "com.deepin.daemon.InputDevice.TouchPad",
			Type: "touchpad",
		},
	}

	m.settings = gio.NewSettings(gsSchemaInputDevices)
	m.WheelSpeed = property.NewGSettingsUintProperty(m, "WheelSpeed", m.settings, gsKeyWheelSpeed)

	m.kbd = getKeyboard()
	m.wacom = getWacom()
	m.tpad = getTouchpad()
	m.mouse = getMouse()
	m.trackPoint = getTrackPoint()

	return m
}

func (m *Manager) setWheelSpeed(inInit bool) {
	speed := m.settings.GetUint(gsKeyWheelSpeed)
	// speed range is [1,100]
	logger.Debug("setWheelSpeed", speed)

	var err error
	shouldWrite := true
	if inInit {
		if _, err := os.Stat(m.imWheelConfigFile); err == nil {
			shouldWrite = false
		}
	}

	if shouldWrite {
		err = writeImWheelConfig(m.imWheelConfigFile, speed)
		if err != nil {
			logger.Warning("failed to write imwheel config file:", err)
			return
		}
	}

	err = controlImWheel(speed)
	if err != nil {
		logger.Warning("failed to control imwheel:", err)
		return
	}
}

func controlImWheel(speed uint32) error {
	if speed == 1 {
		// quit
		return exec.Command(imWheelBin, "-k", "-q").Run()
	}
	// restart
	return exec.Command(imWheelBin, "-k").Run()
}

func writeImWheelConfig(file string, speed uint32) error {
	logger.Debugf("writeImWheelConfig file:%q, speed: %d", file, speed)

	const header = `# written by ` + dbusDest + `
".*"
Control_L,Up,Control_L|Button4
Control_R,Up,Control_R|Button4
Control_L,Down,Control_L|Button5
Control_R,Down,Control_R|Button5
Shift_L,Up,Shift_L|Button4
Shift_R,Up,Shift_R|Button4
Shift_L,Down,Shift_L|Button5
Shift_R,Down,Shift_R|Button5
`
	fh, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fh.Close()
	writer := bufio.NewWriter(fh)

	_, err = writer.Write([]byte(header))
	if err != nil {
		return err
	}

	//  Delay Before Next KeyPress Event
	delay := 240000 / speed
	_, err = fmt.Fprintf(writer, "None,Up,Button4,%d,0,%d\n", speed, delay)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(writer, "None,Down,Button5,%d,0,%d\n", speed, delay)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return fh.Sync()
}

func (m *Manager) init() {
	m.kbd.init()
	m.kbd.handleGSettings()
	m.wacom.init()
	m.wacom.handleGSettings()
	m.tpad.init()
	m.tpad.handleGSettings()
	m.mouse.init()
	m.mouse.handleGSettings()
	m.trackPoint.init()
	m.trackPoint.handleGSettings()

	m.setWheelSpeed(true)
	m.handleGSettings()
}
