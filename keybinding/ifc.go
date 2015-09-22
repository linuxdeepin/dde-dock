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
	"fmt"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/lib/dbus"
)

// Reset reset all shortcut
func (m *Manager) Reset() {
	shortcuts.Reset()
	m.ungrabShortcuts(m.grabedList)
	m.initGrabedList()
}

// List list all shortcut
func (m *Manager) List() string {
	ret, _ := doMarshal(m.listAll())
	return ret
}

// Add add custom shortcut
//
// name: accel name
// action: accel command line
// accel: the binded accel
// ret0: the shortcut id
// ret1: whether accel conflict, if true, ret0 is conflict id
// ret2: error info
func (m *Manager) Add(name, action, accel string) (string, bool, error) {
	logger.Debugf("Add custom key: %s %s %s", name, action, accel)
	avaliable, conflict := m.CheckAvaliable(accel)
	if !avaliable {
		return conflict, true, nil
	}

	id, err := shortcuts.AddCustomKey(name, action, []string{accel})
	if err != nil {
		return "", false, err
	}

	list := shortcuts.ListCustomKey().GetShortcuts()
	info := list.GetById(id, shortcuts.KeyTypeCustom)
	if info == nil {
		return "", false, fmt.Errorf("Add custom accel failed")
	}
	dbus.Emit(m, "Added", id, shortcuts.KeyTypeCustom)

	err = m.grabShortcut(info)
	if err != nil {
		return "", false, err
	}
	return id, false, nil
}

// Delete delete shortcut by id and type
//
// id: the specail id
// ty: the special type
// ret0: error info
func (m *Manager) Delete(id string, ty int32) error {
	logger.Debugf("Delete '%s' type '%v'", id, ty)
	if ty != shortcuts.KeyTypeCustom {
		return fmt.Errorf("Invalid shortcut type '%v'", ty)
	}

	s := m.grabedList.GetById(id, ty)
	if s == nil {
		return fmt.Errorf("Invalid shortcut id '%s'", id)
	}

	err := shortcuts.DeleteCustomKey(id)
	if err != nil {
		return err
	}

	m.ungrabShortcut(s)
	dbus.Emit(m, "Deleted", id, ty)
	return nil
}

// Disable cancel the special id accels
func (m *Manager) Disable(id string, ty int32) error {
	logger.Debugf("Disable '%s' type '%v'", id, ty)
	s := m.listAll().GetById(id, ty)
	if s == nil {
		return fmt.Errorf("Invalid shortcut id '%s'", id)
	}

	m.ungrabAccels(s.Accels)
	s.Disable()

	// system/media update by gsetting changed
	if ty == shortcuts.KeyTypeCustom {
		m.grabedList = m.grabedList.Add(id, ty)
	}
	dbus.Emit(m, "Changed", id, ty)

	return nil
}

// CheckAvaliable check the accel whether conflict
func (m *Manager) CheckAvaliable(accel string) (bool, string) {
	logger.Debug("Check accel:", accel)
	s := m.listAll().GetByAccel(accel)
	if s == nil {
		return true, ""
	}

	return false, s.Id
}

// ModifiedName modify the special id name, only for custom shortcut
func (m *Manager) ModifiedName(id string, ty int32, name string) error {
	logger.Debugf("Modify name '%s' type '%v' value '%s'", id, ty, name)
	if ty != shortcuts.KeyTypeCustom {
		return fmt.Errorf("Invalid shortcut type '%v'", ty)
	}

	s := m.grabedList.GetById(id, ty)
	if s == nil {
		return fmt.Errorf("Invalid shortcut id '%s'", id)
	}

	s.SetName(name)
	m.grabedList = m.grabedList.Add(s.Id, s.Type)
	dbus.Emit(m, "Changed", id, ty)
	return nil
}

// ModifiedAction modify the special id action, only for custom shortcut
func (m *Manager) ModifiedAction(id string, ty int32, action string) error {
	logger.Debugf("Modify action '%s' type '%v' value '%s'", id, ty, action)
	if ty != shortcuts.KeyTypeCustom {
		return fmt.Errorf("Invalid shortcut type '%v'", ty)
	}

	s := m.grabedList.GetById(id, ty)
	if s == nil {
		return fmt.Errorf("Invalid shortcut id '%s'", id)
	}

	s.SetAction(action)
	m.grabedList = m.grabedList.Add(s.Id, s.Type)
	dbus.Emit(m, "Changed", id, ty)
	return nil
}

// ModifiedAccel modify the special id action
//
// id: the special id
// ty: the special type
// accel: new accel
// grabed: if true, add accel for the special id; else delete it
func (m *Manager) ModifiedAccel(id string, ty int32, accel string, grabed bool) (bool, string, error) {
	logger.Debugf("Modify accel '%s' type '%v' value '%s' grabed: %v", id, ty, accel, grabed)
	if !grabed {
		return false, "", m.deleteAccel(id, ty, accel)
	}

	avaliable, conflict := m.CheckAvaliable(accel)
	if !avaliable {
		return true, conflict, nil
	}

	return false, "", m.addAccel(id, ty, accel)
}

// Query query shortcut detail info by id and type
func (m *Manager) Query(id string, ty int32) (string, error) {
	logger.Debug("Query by:", id, ty)
	var v interface{}
	switch ty {
	case shortcuts.KeyTypeSystem:
		v = shortcuts.ListSystemShortcut().GetById(id, ty)
	case shortcuts.KeyTypeWM:
		v = shortcuts.ListWMShortcut().GetById(id, ty)
	case shortcuts.KeyTypeMedia:
		v = shortcuts.ListMediaShortcut().GetById(id, ty)
	case shortcuts.KeyTypeMetacity:
		v = shortcuts.ListMetacityShortcut().GetById(id, ty)
	case shortcuts.KeyTypeCustom:
		v = shortcuts.ListCustomKey().Get(id)
	default:
		return "", fmt.Errorf("Invalid shortcut type:", ty)
	}

	if v == nil {
		return "", fmt.Errorf("Not found the id:", id, ty)
	}

	ret, _ := doMarshal(v)
	return ret, nil
}

// GrabScreen grab screen for getting the key pressed
func (m *Manager) GrabScreen() error {
	return m.doGrabScreen()
}

func (m *Manager) addAccel(id string, ty int32, accel string) error {
	s := m.listAll().GetById(id, ty)
	if s == nil {
		return fmt.Errorf("Invalid id '%s' or type '%v'", id, ty)
	}

	err := m.grabAccels([]string{accel}, m.handleKeyEvent)
	if err != nil {
		return err
	}

	s.AddAccel(accel)

	// system/media update by gsetting changed
	if ty == shortcuts.KeyTypeCustom {
		m.grabedList = m.grabedList.Add(id, ty)
	}
	dbus.Emit(m, "Changed", id, ty)

	return nil
}

func (m *Manager) deleteAccel(id string, ty int32, accel string) error {
	s := m.listAll().GetById(id, ty)
	if s == nil {
		return fmt.Errorf("Invalid id '%s' or type '%v'", id, ty)
	}

	m.ungrabAccels([]string{accel})
	s.DeleteAccel(accel)

	// system/media update by gsetting changed
	if ty == shortcuts.KeyTypeCustom {
		m.grabedList = m.grabedList.Add(id, ty)
	}
	dbus.Emit(m, "Changed", id, ty)

	return nil
}
