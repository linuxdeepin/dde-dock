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

package soundeffect

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"sync"

	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/strv"
)

const (
	soundEffectSchema = "com.deepin.dde.sound-effect"
	settingKeyEnabled = "enabled"

	DBusServiceName = "com.deepin.daemon.SoundEffect"
	dbusPath        = "/com/deepin/daemon/SoundEffect"
	dbusInterface   = DBusServiceName
)

type Manager struct {
	service *dbusutil.Service
	setting *gio.Settings
	count   int
	countMu sync.Mutex
	names   strv.Strv

	Enabled gsprop.Bool `prop:"access:rw"`

	methods *struct {
		PlaySystemSound    func() `in:"name"`
		GetSystemSoundFile func() `in:"name" out:"file"`
		PlaySound          func() `in:"name"`
		GetSoundFile       func() `in:"name" out:"file"`
		EnableSound        func() `in:"name,enabled"`
		IsSoundEnabled     func() `in:"name" out:"enabled"`
		GetSoundEnabledMap func() `out:"result"`
	}
}

func NewManager(service *dbusutil.Service) *Manager {
	var m = new(Manager)

	m.service = service
	m.setting = gio.NewSettings(soundEffectSchema)
	return m
}

func (m *Manager) init() error {
	m.Enabled.Bind(m.setting, settingKeyEnabled)
	var err error
	m.names, err = getSoundNames()
	if err != nil {
		return err
	}
	logger.Debug(m.names)
	return nil
}

func getSoundNames() ([]string, error) {
	var result []string
	out, err := exec.Command("gsettings", "list-recursively", soundEffectSchema).Output()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		parts := bytes.Fields(scanner.Bytes())
		if len(parts) == 3 {
			if bytes.Equal(parts[1], []byte("enabled")) {
				// skip key enabled
				continue
			}
			key := string(parts[1])
			value := parts[2]
			if bytes.Equal(value, []byte("true")) || bytes.Equal(value, []byte("false")) {
				result = append(result, key)
			}
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return result, nil
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

// deprecated
func (m *Manager) PlaySystemSound(name string) *dbus.Error {
	return m.PlaySound(name)
}

func (m *Manager) PlaySound(name string) *dbus.Error {
	m.service.DelayAutoQuit()

	if name == "" {
		return nil
	}

	go func() {
		m.countMu.Lock()
		m.count++
		logger.Debug("start", m.count)
		m.countMu.Unlock()

		err := soundutils.PlaySystemSound(name, "")
		if err != nil {
			logger.Error(err)
		}

		m.countMu.Lock()
		logger.Debug("end", m.count)
		m.count--
		m.countMu.Unlock()
	}()
	return nil
}

// deprecated
func (m *Manager) GetSystemSoundFile(name string) (string, *dbus.Error) {
	return m.GetSoundFile(name)
}

func (m *Manager) GetSoundFile(name string) (string, *dbus.Error) {
	m.service.DelayAutoQuit()

	file := soundutils.GetSystemSoundFile(name)
	if file == "" {
		return "", dbusutil.ToError(errors.New("sound file not found"))
	}

	return file, nil
}

func (m *Manager) canQuit() bool {
	m.countMu.Lock()
	count := m.count
	m.countMu.Unlock()
	return count == 0
}

func (m *Manager) EnableSound(name string, enabled bool) *dbus.Error {
	if !m.names.Contains(name) {
		return dbusutil.ToError(errors.New("invalid sound event"))
	}

	// TODO
	if name == "desktop-login" {
		logger.Debug("sync config to sound-theme-player")
	}

	m.setting.SetBoolean(name, enabled)
	return nil
}

func (m *Manager) IsSoundEnabled(name string) (bool, *dbus.Error) {
	if !m.names.Contains(name) {
		return false, dbusutil.ToError(errors.New("invalid sound event"))
	}

	return m.setting.GetBoolean(name), nil
}

func (m *Manager) GetSoundEnabledMap() (map[string]bool, *dbus.Error) {
	result := make(map[string]bool, len(m.names))
	for _, name := range m.names {
		result[name] = m.setting.GetBoolean(name)
	}
	return result, nil
}
