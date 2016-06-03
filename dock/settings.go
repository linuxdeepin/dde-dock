/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"fmt"
	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
	"sync"
)

const (
	HideModeKey    string = "hide-mode"
	DisplayModeKey string = "display-mode"
)

// Setting存储dock相关的设置。
type Setting struct {
	dockManager     *DockManager
	core            *gio.Settings
	hideModeLock    sync.RWMutex
	displayModeLock sync.RWMutex
	// HideModeChanged在dock的隐藏模式改变后触发，返回改变后的模式。
	HideModeChanged func(mode int32)
	// DisplayModeChanged在dock的显示模式改变后触发，返回改变后的模式。
	DisplayModeChanged func(mode int32)
}

func NewSetting(dockManager *DockManager) *Setting {
	s := &Setting{
		dockManager: dockManager,
	}

	s.core = gio.NewSettings(dockSchema)
	if s.core == nil {
		return nil
	}
	return s
}

func (s *Setting) listenSettingChange(key string, handler func(*gio.Settings, string)) {
	signalDetial := fmt.Sprintf("changed::%s", key)
	logger.Debugf("connect to %s signal", signalDetial)
	s.core.Connect(signalDetial, handler)
}

func (s *Setting) listenSettingsChanged() {
	// listen hide mode change
	s.listenSettingChange(HideModeKey, func(g *gio.Settings, key string) {
		mode := HideModeType(g.GetEnum(key))
		logger.Debug(key, "changed to", mode)
		s.setManagerHideMode(mode)
	})

	// listen display mode change
	s.listenSettingChange(DisplayModeKey, func(g *gio.Settings, key string) {
		mode := DisplayModeType(g.GetEnum(key))
		logger.Debug(key, "changed to", mode)
		s.setManagerDisplayMode(mode)
	})
}

// func (s *Setting) initManagerDisplayHideMode() {
// 	// TODO: test it
// 	// at least one read operation must be called after signal connected, otherwise,
// 	// the signal connection won't work from glib 2.43.
// 	// NB: https://github.com/GNOME/glib/commit/8ff5668a458344da22d30491e3ce726d861b3619
// 	// s.displayMode = DisplayModeType(s.core.GetEnum(DisplayModeKey))
// 	// s.hideMode = HideModeType(s.core.GetEnum(HideModeKey))
// 	// if s.hideMode == HideModeAutoHide {
// 	// 	s.hideMode = HideModeSmartHide
// 	// 	s.core.SetEnum(HideModeKey, int32(HideModeSmartHide))
// 	// }
// 	logger.Debug("initManagerDisplayHideMode")
//
// 	s.setManagerDisplayMode(DisplayModeType(s.core.GetEnum(DisplayModeKey)))
// 	s.setManagerHideMode(HideModeType(s.core.GetEnum(HideModeKey)))
// }

func (s *Setting) setManagerDisplayMode(mode DisplayModeType) bool {
	s.displayModeLock.Lock()
	defer s.displayModeLock.Unlock()

	logger.Debug("[Setting.SetDisplayMode]:", mode)
	modeChanged := s.dockManager.setDisplayMode(mode)
	if modeChanged {
		dbus.Emit(s, "DisplayModeChanged", int32(mode))
	}
	return modeChanged
}

func (s *Setting) setDisplayMode(mode DisplayModeType) bool {
	if s.setManagerDisplayMode(mode) {
		// mode changed, save setting
		return s.core.SetEnum(DisplayModeKey, int32(mode))
	}
	return false
}

func (s *Setting) setManagerHideMode(mode HideModeType) bool {
	s.hideModeLock.Lock()
	defer s.hideModeLock.Unlock()

	logger.Debug("[Setting.SetHideMode]:", mode)
	modeChanged := s.dockManager.setHideMode(mode)
	if modeChanged {
		dbus.Emit(s, "HideModeChanged", int32(mode))
	}
	return modeChanged
}

func (s *Setting) setHideMode(mode HideModeType) bool {
	if s.setManagerHideMode(mode) {
		// mode changed, save setting
		return s.core.SetEnum(HideModeKey, int32(mode))
	}
	return false
}

// GetHideMode返回当前的隐藏模式。
func (s *Setting) GetHideMode() int32 {
	return int32(s.dockManager.hideMode)
}

// SetHideMode设置dock的隐藏模式。
func (s *Setting) SetHideMode(mode int32) bool {
	if validHideModeNum(mode) {
		return s.setHideMode(HideModeType(mode))
	}
	logger.Warning("Invalid hide mode", mode)
	return false
}

// GetDisplayMode获取dock当前的显示模式。
func (s *Setting) GetDisplayMode() int32 {
	return int32(s.dockManager.displayMode)
}

// SetDisplayMode设置dock的显示模式。
func (s *Setting) SetDisplayMode(mode int32) bool {
	if validDisplayModeNum(mode) {
		return s.setDisplayMode(DisplayModeType(mode))
	}
	logger.Warning("Invalid display mode", mode)
	return false
}

func (s *Setting) destroy() {
	if s.core != nil {
		s.core.Unref()
		s.core = nil
	}
	dbus.UnInstallObject(s)
}
