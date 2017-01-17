/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

import (
	"dbus/com/deepin/daemon/helper/backlight"
	"gir/gio-2.0"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"path/filepath"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/dde/daemon/keybinding/xrecord"
	"pkg.deepin.io/lib/dbus"
	dutils "pkg.deepin.io/lib/utils"
	"pkg.deepin.io/lib/xdg/basedir"
	"time"
)

const (
	// shortcut signals:
	shortcutSignalChanged = "Changed"
	shortcutSignalAdded   = "Added"
	shortcutSignalDeleted = "Deleted"

	systemSchema   = "com.deepin.dde.keybinding.system"
	mediakeySchema = "com.deepin.dde.keybinding.mediakey"
	wmSchema       = "com.deepin.wrap.gnome.desktop.wm.keybindings"
	metacitySchema = "com.deepin.wrap.gnome.metacity.keybindings"
	galaSchema     = "com.deepin.wrap.pantheon.desktop.gala.keybindings"

	customConfigFile = "deepin/dde-daemon/keybinding/custom.ini"
)

type Manager struct {
	Added   func(string, int32)
	Deleted func(string, int32)
	Changed func(string, int32)

	// (pressed, accel)
	KeyEvent func(bool, string)

	xu *xgbutil.XUtil

	sysSetting            *gio.Settings
	mediaSetting          *gio.Settings
	wmSetting             *gio.Settings
	metacitySetting       *gio.Settings
	enableListenGSettings bool

	customShortcutManager *shortcuts.CustomShortcutManager

	blDaemon *backlight.Backlight
	// controllers
	audioController       *AudioController
	mediaPlayerController *MediaPlayerController
	displayController     *DisplayController
	kbdLightController    *KbdLightController
	touchpadController    *TouchpadController

	shortcuts *shortcuts.Shortcuts
	// shortcut action handlers
	handlers               []shortcuts.KeyEventFunc
	lastKeyEventTime       time.Time
	grabScreenPressedAccel *shortcuts.ParsedAccel
}

func NewManager() (*Manager, error) {
	var m = Manager{
		enableListenGSettings: true,
	}

	xu, err := xgbutil.NewConn()
	if err != nil {
		return nil, err
	}
	m.xu = xu
	keybind.Initialize(xu)

	m.sysSetting = gio.NewSettings(systemSchema)
	m.mediaSetting = gio.NewSettings(mediakeySchema)
	m.wmSetting = gio.NewSettings(wmSchema)

	m.metacitySetting, _ = dutils.CheckAndNewGSettings(galaSchema)
	if m.metacitySetting == nil {
		// try metacitySchema
		m.metacitySetting, _ = dutils.CheckAndNewGSettings(metacitySchema)
	}

	customConfigFilePath := filepath.Join(basedir.GetUserConfigDir(), customConfigFile)
	m.customShortcutManager = shortcuts.NewCustomShortcutManager(customConfigFilePath)

	//m.media = &Mediakey{}
	m.shortcuts = shortcuts.NewShortcuts(xu, m.handleKeyEvent)
	m.shortcuts.AddSystem(m.sysSetting)
	m.shortcuts.AddMedia(m.mediaSetting)
	m.shortcuts.AddCustom(m.customShortcutManager)
	m.shortcuts.AddWM(m.wmSetting)

	if m.metacitySetting != nil {
		m.shortcuts.AddMetacity(m.metacitySetting)
	} else {
		// TODO
		logger.Warning("Manager.metacitySetting is nil")
	}

	m.audioController, err = NewAudioController()
	if err != nil {
		logger.Warning("NewAudioController failed:", err)
	}

	m.mediaPlayerController, err = NewMediaPlayerController()
	if err != nil {
		logger.Warning("NewMediaPlayerController failed:", err)
	}

	m.blDaemon, err = backlight.NewBacklight("com.deepin.daemon.helper.Backlight",
		"/com/deepin/daemon/helper/Backlight")
	if err != nil {
		logger.Warning("NewBacklight failed:", err)
	}

	m.displayController, err = NewDisplayController(m.blDaemon)
	if err != nil {
		logger.Warning("NewDisplayController failed:", err)
	}

	m.kbdLightController = NewKbdLightController(m.blDaemon)

	m.touchpadController, err = NewTouchpadController()
	if err != nil {
		logger.Warning("NewTouchpadController failed:", err)
	}

	m.initHandlers()
	m.shortcuts.ListenXEvents()

	// listen gsetting changed event
	m.listenGSettingsChanged(m.sysSetting, shortcuts.ShortcutTypeSystem)
	m.listenGSettingsChanged(m.mediaSetting, shortcuts.ShortcutTypeMedia)
	m.listenGSettingsChanged(m.wmSetting, shortcuts.ShortcutTypeWM)
	m.listenGSettingsChanged(m.metacitySetting, shortcuts.ShortcutTypeMetacity)

	go xevent.Main(m.xu)

	// init package xrecord
	xrecord.Initialize()
	xrecord.SetKeyReleaseCallback(m.shortcuts.HandleXRecordKeyRelease)

	return &m, nil
}

func (m *Manager) destroy() {
	// TODO ungrab all shortcuts
	xrecord.Finalize()
	xrecord.SetKeyReleaseCallback(nil)

	// destroy settings
	if m.sysSetting != nil {
		m.sysSetting.Unref()
		m.sysSetting = nil
	}

	if m.mediaSetting != nil {
		m.mediaSetting.Unref()
		m.mediaSetting = nil
	}

	if m.wmSetting != nil {
		m.wmSetting.Unref()
		m.wmSetting = nil
	}

	if m.metacitySetting != nil {
		m.metacitySetting.Unref()
		m.metacitySetting = nil
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

func (m *Manager) emitShortcutSignal(signalName string, shortcut shortcuts.Shortcut) {
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

		shortcut := m.shortcuts.GetByIdType(key, type_)
		if shortcut == nil {
			return
		}

		accelStrv := gsettings.GetStrv(key)
		m.shortcuts.ModifyShortcutAccels(shortcut, shortcuts.ParseStandardAccels(accelStrv))
		m.emitShortcutSignal(shortcutSignalChanged, shortcut)
	})
}
