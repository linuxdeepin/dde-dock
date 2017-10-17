/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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
	"path/filepath"
	"time"

	"dbus/com/deepin/daemon/helper/backlight"

	"gir/gio-2.0"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
	"pkg.deepin.io/lib/xdg/basedir"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/keysyms"
)

const (
	// shortcut signals:
	shortcutSignalChanged = "Changed"
	shortcutSignalAdded   = "Added"
	shortcutSignalDeleted = "Deleted"

	gsSchemaKeyboard          = "com.deepin.dde.keyboard"
	gsKeyNumLockState         = "numlock-state"
	gsKeySaveNumLockState     = "save-numlock-state"
	gsKeyShortcutSwitchLayout = "shortcut-switch-layout"
	gsKeyShowCapsLockOSD      = "capslock-toggle"

	gsSchemaSystem   = "com.deepin.dde.keybinding.system"
	gsSchemaMediaKey = "com.deepin.dde.keybinding.mediakey"
	gsSchemaGnomeWM  = "com.deepin.wrap.gnome.desktop.wm.keybindings"

	customConfigFile = "deepin/dde-daemon/keybinding/custom.ini"
)

type Manager struct {
	// properties
	NumLockState         *property.GSettingsEnumProperty
	ShortcutSwitchLayout *property.GSettingsUintProperty `access:"readwrite"`

	// Signals
	Added   func(string, int32)
	Deleted func(string, int32)
	Changed func(string, int32)

	// (pressed, keystroke)
	KeyEvent func(bool, string)

	conn       *x.Conn
	keySymbols *keysyms.KeySymbols

	gsKeyboard *gio.Settings
	gsSystem   *gio.Settings
	gsMediaKey *gio.Settings
	gsGnomeWM  *gio.Settings

	enableListenGSettings bool

	customShortcutManager *shortcuts.CustomShortcutManager

	backlightHelper *backlight.Backlight
	// controllers
	audioController       *AudioController
	mediaPlayerController *MediaPlayerController
	displayController     *DisplayController
	kbdLightController    *KbdLightController
	touchpadController    *TouchpadController

	shortcutManager *shortcuts.ShortcutManager
	// shortcut action handlers
	handlers            []shortcuts.KeyEventFunc
	lastKeyEventTime    time.Time
	grabScreenKeystroke *shortcuts.Keystroke

	// for switch kbd layout
	switchKbdLayoutState SKLState
	sklWaitQuit          chan int
}

// SKLState Switch keyboard Layout state
type SKLState uint

const (
	SKLStateNone SKLState = iota
	SKLStateWait
	SKLStateOSDShown
)

func NewManager() (*Manager, error) {
	conn, err := x.NewConn()
	if err != nil {
		return nil, err
	}

	var m = Manager{
		enableListenGSettings: true,
		conn:       conn,
		keySymbols: keysyms.NewKeySymbols(conn),
	}
	return &m, nil
}

func (m *Manager) init() {
	m.gsKeyboard = gio.NewSettings(gsSchemaKeyboard)
	// init numlock state
	m.NumLockState = property.NewGSettingsEnumProperty(m, "NumLockState", m.gsKeyboard, gsKeyNumLockState)
	if m.gsKeyboard.GetBoolean(gsKeySaveNumLockState) {
		nlState := NumLockState(m.NumLockState.Get())
		if nlState == NumLockUnknown {
			state, err := queryNumLockState(m.conn)
			if err != nil {
				logger.Warning("queryNumLockState failed:", err)
			} else {
				m.NumLockState.Set(int32(state))
			}
		} else {
			err := setNumLockState(m.conn, m.keySymbols, nlState)
			if err != nil {
				logger.Warning("setNumLockState failed:", err)
			}
		}
	}
	m.ShortcutSwitchLayout = property.NewGSettingsUintProperty(m, "ShortcutSwitchLayout",
		m.gsKeyboard, gsKeyShortcutSwitchLayout)

	// init settings
	m.gsSystem = gio.NewSettings(gsSchemaSystem)
	m.gsMediaKey = gio.NewSettings(gsSchemaMediaKey)
	m.gsGnomeWM = gio.NewSettings(gsSchemaGnomeWM)

	m.shortcutManager = shortcuts.NewShortcutManager(m.conn, m.keySymbols, m.handleKeyEvent)
	m.shortcutManager.AddSpecial()
	m.shortcutManager.AddSystem(m.gsSystem)
	m.shortcutManager.AddMedia(m.gsMediaKey)
	m.shortcutManager.AddWM(m.gsGnomeWM)

	customConfigFilePath := filepath.Join(basedir.GetUserConfigDir(), customConfigFile)
	m.customShortcutManager = shortcuts.NewCustomShortcutManager(customConfigFilePath)
	m.shortcutManager.AddCustom(m.customShortcutManager)

	var err error
	m.audioController, err = NewAudioController()
	if err != nil {
		logger.Warning("NewAudioController failed:", err)
	}

	m.mediaPlayerController, err = NewMediaPlayerController()
	if err != nil {
		logger.Warning("NewMediaPlayerController failed:", err)
	}

	m.backlightHelper, err = backlight.NewBacklight("com.deepin.daemon.helper.Backlight",
		"/com/deepin/daemon/helper/Backlight")
	if err != nil {
		logger.Warning("NewBacklight failed:", err)
	}

	m.displayController, err = NewDisplayController(m.backlightHelper)
	if err != nil {
		logger.Warning("NewDisplayController failed:", err)
	}

	m.kbdLightController = NewKbdLightController(m.backlightHelper)

	m.touchpadController, err = NewTouchpadController()
	if err != nil {
		logger.Warning("NewTouchpadController failed:", err)
	}
}

func (m *Manager) destroy() {
	m.shortcutManager.Destroy()

	// destroy settings
	if m.gsSystem != nil {
		m.gsSystem.Unref()
		m.gsSystem = nil
	}

	if m.gsMediaKey != nil {
		m.gsMediaKey.Unref()
		m.gsMediaKey = nil
	}

	if m.gsGnomeWM != nil {
		m.gsGnomeWM.Unref()
		m.gsGnomeWM = nil
	}

	if m.audioController != nil {
		m.audioController.Destroy()
		m.audioController = nil
	}

	if m.mediaPlayerController != nil {
		m.mediaPlayerController.Destroy()
		m.mediaPlayerController = nil
	}

	if m.displayController != nil {
		m.displayController.Destroy()
		m.displayController = nil
	}

	if m.touchpadController != nil {
		m.touchpadController.Destroy()
		m.touchpadController = nil
	}
}

func (m *Manager) handleKeyEvent(ev *shortcuts.KeyEvent) {
	now := time.Now()
	duration := now.Sub(m.lastKeyEventTime)
	logger.Debug("duration:", duration)
	if 0 < duration && duration < 200*time.Millisecond {
		logger.Debug("handleKeyEvent ignore key event")
		return
	}
	m.lastKeyEventTime = now

	logger.Debugf("handleKeyEvent ev: %#v", ev)
	action := ev.Shortcut.GetAction()
	if action == nil {
		logger.Warning("action is nil")
		return
	}
	logger.Debugf("shortcut action: %#v", action)
	if handler := m.handlers[int(action.Type)]; handler != nil {
		handler(ev)
	} else {
		logger.Warning("handler is nil")
	}
}

func shouldEmitSignalChanged(shortcut shortcuts.Shortcut) bool {
	return shortcut.GetType() == shortcuts.ShortcutTypeCustom
}

func (m *Manager) emitShortcutSignal(signalName string, shortcut shortcuts.Shortcut) {
	logger.Debug("emit DBus signal", signalName, shortcut.GetId(), shortcut.GetType())
	dbus.Emit(m, signalName, shortcut.GetId(), shortcut.GetType())
}

func (m *Manager) enableListenGSettingsChanged(val bool) {
	m.enableListenGSettings = val
}

func (m *Manager) listenGSettingsChanged(gsettings *gio.Settings, type_ int32) {
	gsettings.Connect("changed", func(s *gio.Settings, key string) {
		if !m.enableListenGSettings {
			return
		}

		shortcut := m.shortcutManager.GetByIdType(key, type_)
		if shortcut == nil {
			return
		}

		keystrokes := gsettings.GetStrv(key)
		m.shortcutManager.ModifyShortcutKeystrokes(shortcut, shortcuts.ParseKeystrokes(keystrokes))
		m.emitShortcutSignal(shortcutSignalChanged, shortcut)
	})
}
