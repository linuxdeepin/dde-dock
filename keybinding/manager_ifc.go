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
	"errors"
	"fmt"
	"strings"

	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest     = "com.deepin.daemon.Keybinding"
	bindDBusPath = "/com/deepin/daemon/Keybinding"
	bindDBusIFC  = "com.deepin.daemon.Keybinding"
)

type ErrInvalidShortcutType struct {
	Type int32
}

func (err ErrInvalidShortcutType) Error() string {
	return fmt.Sprintf("shortcut type %v is invalid", err.Type)
}

type ErrShortcutNotFound struct {
	Id   string
	Type int32
}

func (err ErrShortcutNotFound) Error() string {
	return fmt.Sprintf("shortcut id %q type %v is not found", err.Id, err.Type)
}

var errTypeAssertionFail = errors.New("type assertion failed")
var errShortcutKeystrokesUnmodifiable = errors.New("keystrokes of this shortcut is unmodifiable")
var errKeystrokeUsed = errors.New("keystroke have been used")

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: bindDBusPath,
		Interface:  bindDBusIFC,
	}
}

// Reset reset all shortcut
func (m *Manager) Reset() {
	m.shortcutManager.UngrabAll()

	m.enableListenGSettingsChanged(false)
	// reset all gsettings
	resetGSettings(m.gsSystem)
	resetGSettings(m.gsMediaKey)
	resetGSettings(m.gsGnomeWM)

	// disable all custom shortcuts
	m.customShortcutManager.DisableAll()

	changes := m.shortcutManager.ReloadAllShortcutsKeystrokes()
	m.enableListenGSettingsChanged(true)
	m.shortcutManager.GrabAll()
	for _, shortcut := range changes {
		m.emitShortcutSignal(shortcutSignalChanged, shortcut)
	}
}

func (m *Manager) ListAllShortcuts() (string, error) {
	list := m.shortcutManager.List()
	ret, err := doMarshal(list)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (m *Manager) ListShortcutsByType(type0 int32) (string, error) {
	list := m.shortcutManager.ListByType(type0)
	ret, err := doMarshal(list)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (m *Manager) AddCustomShortcut(name, action, keystroke string) (id string, type0 int32, err error) {
	logger.Debugf("Add custom key: %q %q %q", name, action, keystroke)
	ks, err := shortcuts.ParseKeystroke(keystroke)
	if err != nil {
		return
	}

	conflictKeystroke, err := m.shortcutManager.FindConflictingKeystroke(ks)
	if err != nil {
		return
	}
	if conflictKeystroke != nil {
		err = errKeystrokeUsed
		return
	}

	shortcut, err := m.customShortcutManager.Add(name, action, []*shortcuts.Keystroke{ks})
	if err != nil {
		return
	}
	m.shortcutManager.Add(shortcut)
	m.emitShortcutSignal(shortcutSignalAdded, shortcut)
	id = shortcut.GetId()
	type0 = shortcut.GetType()
	return
}

func (m *Manager) DeleteCustomShortcut(id string) error {
	shortcut := m.shortcutManager.GetByIdType(id, shortcuts.ShortcutTypeCustom)
	if err := m.customShortcutManager.Delete(shortcut.GetId()); err != nil {
		return err
	}
	m.shortcutManager.Delete(shortcut)
	m.emitShortcutSignal(shortcutSignalDeleted, shortcut)
	return nil
}

func (m *Manager) ClearShortcutKeystrokes(id string, type0 int32) error {
	logger.Debug("ClearShortcutKeystrokes", id, type0)
	shortcut := m.shortcutManager.GetByIdType(id, type0)
	if shortcut == nil {
		return ErrShortcutNotFound{id, type0}
	}
	m.shortcutManager.ModifyShortcutKeystrokes(shortcut, nil)
	err := shortcut.SaveKeystrokes()
	if err != nil {
		return err
	}
	if shouldEmitSignalChanged(shortcut) {
		m.emitShortcutSignal(shortcutSignalChanged, shortcut)
	}
	return nil
}

func (m *Manager) LookupConflictingShortcut(keystroke string) (string, error) {
	ks, err := shortcuts.ParseKeystroke(keystroke)
	if err != nil {
		// parse keystroke error
		return "", err
	}

	conflictKeystroke, err := m.shortcutManager.FindConflictingKeystroke(ks)
	if err != nil {
		return "", err
	}
	if conflictKeystroke != nil {
		detail, err := doMarshal(conflictKeystroke.Shortcut)
		if err != nil {
			return "", err
		}
		return detail, nil
	}
	return "", nil
}

// ModifyCustomShortcut modify custom shortcut
//
// id: shortcut id
// name: new name
// cmd: new commandline
// keystroke: new keystroke
func (m *Manager) ModifyCustomShortcut(id, name, cmd, keystroke string) error {
	logger.Debugf("ModifyCustomShortcut id: %q, name: %q, cmd: %q, keystroke: %q", id, name, cmd, keystroke)
	const ty = shortcuts.ShortcutTypeCustom
	// get the shortcut
	shortcut := m.shortcutManager.GetByIdType(id, ty)
	if shortcut == nil {
		return ErrShortcutNotFound{id, ty}
	}
	customShortcut, ok := shortcut.(*shortcuts.CustomShortcut)
	if !ok {
		return errTypeAssertionFail
	}

	var keystrokes []*shortcuts.Keystroke
	if keystroke != "" {
		ks, err := shortcuts.ParseKeystroke(keystroke)
		if err != nil {
			return err
		}
		// check conflicting
		conflictKeystroke, err := m.shortcutManager.FindConflictingKeystroke(ks)
		if err != nil {
			return err
		}
		if conflictKeystroke != nil && conflictKeystroke.Shortcut != shortcut {
			return errKeystrokeUsed
		}
		keystrokes = []*shortcuts.Keystroke{ks}
	}

	// modify then save
	customShortcut.Name = name
	customShortcut.Cmd = cmd
	m.shortcutManager.ModifyShortcutKeystrokes(shortcut, keystrokes)
	err := customShortcut.Save()
	if err != nil {
		return err
	}
	m.emitShortcutSignal(shortcutSignalChanged, shortcut)
	return nil
}

func (m *Manager) AddShortcutKeystroke(id string, type0 int32, keystroke string) error {
	logger.Debug("AddShortcutKeystroke", id, type0, keystroke)
	shortcut := m.shortcutManager.GetByIdType(id, type0)
	if shortcut == nil {
		return ErrShortcutNotFound{id, type0}
	}
	if !shortcut.GetKeystrokesModifiable() {
		return errShortcutKeystrokesUnmodifiable
	}

	ks, err := shortcuts.ParseKeystroke(keystroke)
	if err != nil {
		// parse keystroke error
		return err
	}
	logger.Debug("keystroke:", ks.DebugString())

	if type0 == shortcuts.ShortcutTypeWM && ks.Mods == 0 {
		keyLower := strings.ToLower(ks.Keystr)
		if keyLower == "super_l" || keyLower == "super_r" {
			return errors.New("keystroke of shortcut which type is wm can not be set to the Super key")
		}
	}

	conflictKeystroke, err := m.shortcutManager.FindConflictingKeystroke(ks)
	if err != nil {
		return err
	}
	if conflictKeystroke == nil {
		m.shortcutManager.AddShortcutKeystroke(shortcut, ks)
		err := shortcut.SaveKeystrokes()
		if err != nil {
			return err
		}
		if shouldEmitSignalChanged(shortcut) {
			m.emitShortcutSignal(shortcutSignalChanged, shortcut)
		}
	} else if conflictKeystroke.Shortcut != shortcut {
		return errKeystrokeUsed
	}
	return nil
}

func (m *Manager) DeleteShortcutKeystroke(id string, type0 int32, keystroke string) error {
	logger.Debug("DeleteShortcutKeystroke", id, type0, keystroke)
	shortcut := m.shortcutManager.GetByIdType(id, type0)
	if shortcut == nil {
		return ErrShortcutNotFound{id, type0}
	}
	if !shortcut.GetKeystrokesModifiable() {
		return errShortcutKeystrokesUnmodifiable
	}

	ks, err := shortcuts.ParseKeystroke(keystroke)
	if err != nil {
		// parse keystroke error
		return err
	}
	logger.Debug("keystroke:", ks.DebugString())

	m.shortcutManager.DeleteShortcutKeystroke(shortcut, ks)
	err = shortcut.SaveKeystrokes()
	if err != nil {
		return err
	}
	if shouldEmitSignalChanged(shortcut) {
		m.emitShortcutSignal(shortcutSignalChanged, shortcut)
	}
	return nil
}

func (m *Manager) GetShortcut(id string, type0 int32) (string, error) {
	shortcut := m.shortcutManager.GetByIdType(id, type0)
	if shortcut == nil {
		return "", ErrShortcutNotFound{id, type0}
	}
	return doMarshal(shortcut)
}

func (m *Manager) SelectKeystroke() error {
	logger.Debug("SelectKeystroke")
	return m.selectKeystroke()
}

func (m *Manager) SetNumLockState(state int32) error {
	logger.Debug("SetNumLockState", state)
	return setNumLockState(m.conn, m.keySymbols, NumLockState(state))
}
