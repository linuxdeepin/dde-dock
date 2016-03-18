/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package appearance

import (
	"dbus/com/deepin/daemon/accounts"
	"encoding/json"
	"fmt"
	"gir/gio-2.0"
	"gir/glib-2.0"
	"github.com/howeyc/fsnotify"
	"os/user"
	"path"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/appearance/fonts"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/lib/dbus/property"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	TypeGtkTheme          = "gtk"
	TypeIconTheme         = "icon"
	TypeCursorTheme       = "cursor"
	TypeBackground        = "background"
	TypeGreeterBackground = "greeterbackground"
	TypeStandardFont      = "standardfont"
	TypeMonospaceFont     = "monospacefont"
	TypeFontSize          = "fontsize"
)

const (
	wrapBgSchema    = "com.deepin.wrap.gnome.desktop.background"
	gnomeBgSchema   = "org.gnome.desktop.background"
	gsKeyBackground = "picture-uri"

	appearanceSchema   = "com.deepin.dde.appearance"
	gsKeyGtkTheme      = "gtk-theme"
	gsKeyIconTheme     = "icon-theme"
	gsKeyCursorTheme   = "cursor-theme"
	gsKeyFontStandard  = "font-standard"
	gsKeyFontMonospace = "font-monospace"
	gsKeyFontSize      = "font-size"
)

const (
	defaultStandardFont  = "sans-serif"
	defaultMonospaceFont = "monospace"
)

type Manager struct {
	GtkTheme      *property.GSettingsStringProperty `access:"readwrite"`
	IconTheme     *property.GSettingsStringProperty `access:"readwrite"`
	CursorTheme   *property.GSettingsStringProperty `access:"readwrite"`
	Background    *property.GSettingsStringProperty `access:"readwrite"`
	StandardFont  *property.GSettingsStringProperty `access:"readwrite"`
	MonospaceFont *property.GSettingsStringProperty `access:"readwrite"`

	FontSize *property.GSettingsIntProperty `access:"readwrite"`

	// Theme changed signal
	// ty, name
	Changed func(string, string)

	userObj *accounts.User

	setting        *gio.Settings
	wrapBgSetting  *gio.Settings
	gnomeBgSetting *gio.Settings

	watcher    *fsnotify.Watcher
	endWatcher chan struct{}
}

func NewManager() *Manager {
	var m = new(Manager)
	m.setting = gio.NewSettings(appearanceSchema)
	m.wrapBgSetting = gio.NewSettings(wrapBgSchema)

	m.GtkTheme = property.NewGSettingsStringProperty(
		m, "GtkTheme",
		m.setting, gsKeyGtkTheme)
	m.IconTheme = property.NewGSettingsStringProperty(
		m, "IconTheme",
		m.setting, gsKeyIconTheme)
	m.CursorTheme = property.NewGSettingsStringProperty(
		m, "CursorTheme",
		m.setting, gsKeyCursorTheme)
	m.StandardFont = property.NewGSettingsStringProperty(
		m, "StandardFont",
		m.setting, gsKeyFontStandard)
	m.MonospaceFont = property.NewGSettingsStringProperty(
		m, "MonospaceFont",
		m.setting, gsKeyFontMonospace)
	m.Background = property.NewGSettingsStringProperty(
		m, "Background",
		m.wrapBgSetting, gsKeyBackground)

	m.FontSize = property.NewGSettingsIntProperty(
		m, "FontSize",
		m.setting, gsKeyFontSize)

	m.gnomeBgSetting, _ = dutils.CheckAndNewGSettings(gnomeBgSchema)

	var err error
	m.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		logger.Warning("New file watcher failed:", err)
	} else {
		m.endWatcher = make(chan struct{})
	}

	cur, err := user.Current()
	if err != nil {
		logger.Warning("Get current user info failed:", err)
	} else {
		m.userObj, err = ddbus.NewUserByUid(cur.Uid)
		if err != nil {
			logger.Warning("New user object failed:", cur.Name, err)
			m.userObj = nil
		}
	}

	m.init()

	return m
}

func (m *Manager) destroy() {
	if m.setting != nil {
		m.setting.Unref()
		m.setting = nil
	}

	if m.wrapBgSetting != nil {
		m.wrapBgSetting.Unref()
		m.wrapBgSetting = nil
	}

	if m.gnomeBgSetting != nil {
		m.gnomeBgSetting.Unref()
		m.gnomeBgSetting = nil
	}

	if m.userObj != nil {
		ddbus.DestroyUser(m.userObj)
		m.userObj = nil
	}

	if m.watcher != nil {
		close(m.endWatcher)
		m.watcher.Close()
		m.watcher = nil
	}
}

func (m *Manager) init() {
	m.correctFontName()

	var file = path.Join(glib.GetUserConfigDir(), "fontconfig", "fonts.conf")
	if dutils.IsFileExist(file) {
		return
	}

	err := fonts.SetFamily(m.StandardFont.Get(), m.MonospaceFont.Get(),
		m.FontSize.Get())
	if err != nil {
		logger.Debug("[init]----------- font failed:", err)
		return
	}
}

func (m *Manager) correctFontName() {
	families := fonts.ListAllFamily()
	stand := families.Get(m.StandardFont.Get())
	if stand != nil {
		if stand.Id != m.StandardFont.Get() {
			m.StandardFont.Set(stand.Id)
		}
	} else {
		m.StandardFont.Set(defaultStandardFont)
	}

	mono := families.Get(m.MonospaceFont.Get())
	if mono != nil {
		if mono.Id != m.MonospaceFont.Get() {
			m.MonospaceFont.Set(mono.Id)
		}
	} else {
		m.MonospaceFont.Set(defaultMonospaceFont)
	}
}

func (m *Manager) doSetGtkTheme(value string) error {
	if !subthemes.IsGtkTheme(value) {
		return fmt.Errorf("Invalid gtk theme '%v'", value)
	}

	return subthemes.SetGtkTheme(value)
}

func (m *Manager) doSetIconTheme(value string) error {
	if !subthemes.IsIconTheme(value) {
		return fmt.Errorf("Invalid icon theme '%v'", value)
	}

	return subthemes.SetIconTheme(value)
}

func (m *Manager) doSetCursorTheme(value string) error {
	if !subthemes.IsCursorTheme(value) {
		return fmt.Errorf("Invalid cursor theme '%v'", value)
	}

	return subthemes.SetCursorTheme(value)
}

func (m *Manager) doSetBackground(value string) (string, error) {
	if !background.IsBackgroundFile(value) {
		return "", fmt.Errorf("Invalid background file '%v'", value)
	}

	uri, err := background.ListBackground().EnsureExists(value)
	if err != nil {
		logger.Debugf("[doSetBackground] set '%s' failed: %v", value, uri, err)
		return "", err
	}

	if m.userObj != nil {
		m.userObj.SetBackgroundFile(uri)
	}
	return uri, nil
}

func (m *Manager) doSetGreeterBackground(value string) error {
	if m.userObj == nil {
		return fmt.Errorf("Create user object failed")
	}

	_, err := m.userObj.SetGreeterBackground(value)
	return err
}

func (m *Manager) doSetStandardFont(value string) error {
	if !fonts.IsFontFamily(value) {
		return fmt.Errorf("Invalid font family '%v'", value)
	}

	return fonts.SetFamily(value, m.MonospaceFont.Get(), m.FontSize.Get())
}

func (m *Manager) doSetMonnospaceFont(value string) error {
	if !fonts.IsFontFamily(value) {
		return fmt.Errorf("Invalid font family '%v'", value)
	}

	return fonts.SetFamily(m.StandardFont.Get(), value, m.FontSize.Get())
}

func (m *Manager) doSetFontSize(size int32) error {
	if !fonts.IsFontSizeValid(size) {
		logger.Debug("[doSetFontSize] invalid size:", size)
		return fmt.Errorf("Invalid font size '%v'", size)
	}

	return fonts.SetFamily(m.StandardFont.Get(), m.MonospaceFont.Get(), size)
}

func (*Manager) doShow(ifc interface{}) (string, error) {
	if ifc == nil {
		return "", fmt.Errorf("Not found target")
	}
	content, err := json.Marshal(ifc)
	return string(content), err
}
