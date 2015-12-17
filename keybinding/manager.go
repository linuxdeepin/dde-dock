/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package keybinding

import (
	"encoding/json"
	"sort"
	"sync"

	"github.com/BurntSushi/xgbutil"
	"pkg.deepin.io/dde/daemon/keybinding/core"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/lib/dbus"
	"gir/gio-2.0"
)

const (
	systemSchema   = "com.deepin.dde.keybinding.system"
	mediakeySchema = "com.deepin.dde.keybinding.mediakey"
)

type Manager struct {
	Added   func(string, int32)
	Deleted func(string, int32)
	Changed func(string, int32)

	// (pressed, accel)
	KeyEvent func(bool, string)

	xu *xgbutil.XUtil

	sysSetting   *gio.Settings
	mediaSetting *gio.Settings

	media *Mediakey

	grabLocker sync.Mutex
	grabedList shortcuts.Shortcuts
}

func NewManager() (*Manager, error) {
	var m = Manager{}

	xu, err := core.Initialize()
	if err != nil {
		return nil, err
	}
	m.xu = xu

	m.sysSetting = gio.NewSettings(systemSchema)
	m.mediaSetting = gio.NewSettings(mediakeySchema)

	m.media = &Mediakey{}

	return &m, nil
}

func (m *Manager) destroy() {
	m.ungrabShortcuts(m.grabedList)
	m.grabedList = nil
	m.stopLoop()

	if m.sysSetting != nil {
		m.sysSetting.Unref()
		m.sysSetting = nil
	}

	if m.mediaSetting != nil {
		m.mediaSetting.Unref()
		m.mediaSetting = nil
	}
}

func (m *Manager) startLoop() {
	core.StartLoop()
}

func (m *Manager) stopLoop() {
	core.Finalize()
}

func (m *Manager) initGrabedList() {
	sysList := shortcuts.ListSystemShortcut()
	customList := shortcuts.ListCustomKey().GetShortcuts()
	mediaList := shortcuts.ListMediaShortcut()

	m.grabShortcuts(sysList)
	m.grabShortcuts(customList)
	m.grabShortcuts(mediaList)
}

func (m *Manager) listAll() shortcuts.Shortcuts {
	list := shortcuts.ListWMShortcut()
	list = append(list, shortcuts.ListMetacityShortcut()...)
	list = append(list, m.grabedList...)
	return list
}

func (m *Manager) addToGrabedList(s *shortcuts.Shortcut) {
	m.grabLocker.Lock()
	defer m.grabLocker.Unlock()
	m.grabedList = m.grabedList.Add(s.Id, s.Type)
}

func (m *Manager) deleteFromGrabedList(s *shortcuts.Shortcut) {
	m.grabLocker.Lock()
	defer m.grabLocker.Unlock()
	m.grabedList = m.grabedList.Delete(s.Id, s.Type)
}

func (m *Manager) grabShortcuts(list shortcuts.Shortcuts) {
	for _, s := range list {
		err := m.grabShortcut(s)
		if err != nil {
			logger.Warningf("Grab '%s' %v failed: %v",
				s.Id, s.Accels, err)
			continue
		}
	}
}

func (m *Manager) ungrabShortcuts(list shortcuts.Shortcuts) {
	for _, s := range list {
		m.ungrabShortcut(s)
	}
}

func (m *Manager) grabShortcut(s *shortcuts.Shortcut) error {
	err := m.grabAccels(s.Accels, s.Type, m.handleKeyEvent)
	if err != nil {
		return err
	}

	m.addToGrabedList(s)
	return nil
}

func (m *Manager) ungrabShortcut(s *shortcuts.Shortcut) {
	m.ungrabAccels(s.Accels, s.Type)
	m.deleteFromGrabedList(s)
}

func (m *Manager) grabAccels(accels []string, ty int32, cb core.HandleType) error {
	switch ty {
	case shortcuts.KeyTypeWM, shortcuts.KeyTypeMetacity:
		return nil
	}

	return core.GrabAccels(accels, cb)
}

func (m *Manager) ungrabAccels(accels []string, ty int32) {
	switch ty {
	case shortcuts.KeyTypeWM, shortcuts.KeyTypeMetacity:
		return
	}

	core.UngrabAccels(accels)
}

func (m *Manager) updateShortcutById(id string, ty int32) {
	old := m.grabedList.GetById(id, ty)
	if old == nil {
		return
	}

	new := shortcuts.ListAllShortcuts().GetById(id, ty)
	if new == nil {
		return
	}

	if isListEqual(old.Accels, new.Accels) {
		return
	}

	m.ungrabAccels(old.Accels, old.Type)
	m.grabAccels(new.Accels, new.Type, m.handleKeyEvent)
	m.updateGrabedList(id, ty)
	dbus.Emit(m, "Changed", id, ty)
}

func (m *Manager) updateGrabedList(id string, ty int32) {
	m.grabLocker.Lock()
	defer m.grabLocker.Unlock()
	switch ty {
	case shortcuts.KeyTypeSystem,
		shortcuts.KeyTypeMedia,
		shortcuts.KeyTypeCustom:
		m.grabedList = m.grabedList.Add(id, ty)
	}
}

func doMarshal(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func isListEqual(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return false
	}

	sort.Strings(l1)
	sort.Strings(l2)
	for i, v := range l1 {
		if v != l2[i] {
			return false
		}
	}
	return true
}
