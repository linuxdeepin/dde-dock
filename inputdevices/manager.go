/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

	"pkg.deepin.io/dde/daemon/common/dsync"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
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
	Infos      devicePathInfos // readonly
	WheelSpeed gsprop.Uint     `prop:"access:rw"`

	settings          *gio.Settings
	imWheelConfigFile string

	kbd        *Keyboard
	mouse      *Mouse
	trackPoint *TrackPoint
	tpad       *Touchpad
	wacom      *Wacom

	sessionSigLoop *dbusutil.SignalLoop
	syncConfig     *dsync.Config
}

func NewManager(service *dbusutil.Service) *Manager {
	var m = new(Manager)
	m.imWheelConfigFile = filepath.Join(basedir.GetUserHomeDir(), ".imwheelrc")

	m.Infos = devicePathInfos{
		&devicePathInfo{
			Path: kbdDBusInterface,
			Type: "keyboard",
		},
		&devicePathInfo{
			Path: mouseDBusInterface,
			Type: "mouse",
		},
		&devicePathInfo{
			Path: trackPointDBusInterface,
			Type: "trackpoint",
		},
		&devicePathInfo{
			Path: touchPadDBusInterface,
			Type: "touchpad",
		},
	}

	m.settings = gio.NewSettings(gsSchemaInputDevices)
	m.WheelSpeed.Bind(m.settings, gsKeyWheelSpeed)

	m.kbd = newKeyboard(service)
	m.wacom = newWacom(service)

	m.tpad = newTouchpad(service)

	m.mouse = newMouse(service, m.tpad)

	m.trackPoint = newTrackPoint(service)

	m.sessionSigLoop = dbusutil.NewSignalLoop(service.Conn(), 10)
	m.syncConfig = dsync.NewConfig("peripherals", &syncConfig{m: m},
		m.sessionSigLoop, dbusPath, logger)

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
	return exec.Command(imWheelBin, "-k", "-b", "4 5").Run()
}

func writeImWheelConfig(file string, speed uint32) error {
	logger.Debugf("writeImWheelConfig file:%q, speed: %d", file, speed)

	const header = `# written by ` + dbusServiceName + `
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

	m.sessionSigLoop.Start()
}
