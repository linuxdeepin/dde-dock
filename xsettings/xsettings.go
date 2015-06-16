/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package xsettings

import (
	"dbus/com/deepin/sessionmanager"
	"sync"
)

const (
	//string
	NetThemeName         = "Net/ThemeName"
	NetIconTheme         = "Net/IconThemeName"
	NetFallbackIconTheme = "Net/FallbackIconTheme"
	NetSoundTheme        = "Net/SoundThemeName"

	GtkThemeName       = "Gtk/GtkThemeName"
	GtkCursorTheme     = "Gtk/CursorThemeName"
	GtkFontName        = "Gtk/FontName"
	GtkKeyTheme        = "Gtk/KeyThemeName"
	GtkColorPalette    = "Gtk/ColorPalette"
	GtkToolbarStyle    = "Gtk/ToolbarStyle"
	GtkToolbarIconSize = "Gtk/ToolbarIconSize"
	GtkColorScheme     = "Gtk/ColorScheme"
	GtkIMPreeditStyle  = "Gtk/IMPreeditStyle"
	GtkIMStatusStyle   = "Gtk/IMStatusStyle"
	GtkIMModule        = "Gtk/IMModule"
	GtkMenuBarAccel    = "Gtk/MenuBarAccel"

	XftHintStyle = "xft/HintStyle"
	XftRgba      = "xft/RGBA"

	//integer
	NetCursorBlinkTime    = "Net/CursorBlinkTime"
	NetCursorBlinkTimeout = "Net/CursorBlinkTimeout"
	NetDoubleClick        = "Net/DoubleClickTime"
	NetDragThreshold      = "Net/DndDragThreshold"

	GtkCursorThemeSize   = "Gtk/CursorThemeSize"
	GtkTimeoutInitial    = "Gtk/TimeoutInitial"
	GtkTimeoutRepeat     = "Gtk/TimeoutRepeat"
	GtkRecentFilesMaxAge = "Gtk/RecentFilesMaxAge"

	//bool
	NetCursorBlink               = "Net/CursorBlink"
	NetEnableEventSounds         = "Net/EnableEventSounds"
	NetEnableInputFeedbackSounds = "Net/EnableInputFeedbackSounds"

	GtkCanChangeAccels     = "Gtk/CanChangeAccels"
	GtkMenuImages          = "Gtk/MenuImages"
	GtkButtonImages        = "Gtk/ButtonImages"
	GtkEnableAnimations    = "Gtk/EnableAnimations"
	GtkShowInputMethodMenu = "Gtk/ShowInputMethodMenu"
	GtkShowUnicodeMenu     = "Gtk/ShowUnicodeMenu"
	GtkAutoMnemonics       = "Gtk/AutoMnemonics"
	GtkEnableRecentFiles   = "Gtk/RecentFilesEnabled"

	XftAntialias = "xft/Antialias"
	XftHinting   = "xft/HintStyle"
)

type XSProxy struct {
	xs     *sessionmanager.XSettings
	cnt    uint32
	locker sync.Mutex
}

var _proxy *XSProxy

func NewXSProxy() (*XSProxy, error) {
	if _proxy != nil {
		_proxy.refer()
		return _proxy, nil
	}

	xs, err := sessionmanager.NewXSettings(
		"com.deepin.SessionManager",
		"/com/deepin/XSettings",
	)
	if err != nil {
		return nil, err
	}

	_proxy = &XSProxy{xs: xs}
	_proxy.refer()

	return _proxy, nil
}

func (proxy *XSProxy) Free() {
	proxy.unref()
}

func (proxy *XSProxy) SetString(name, value string) error {
	return proxy.xs.SetString(name, value)
}

func (proxy *XSProxy) SetInteger(name string, value uint32) error {
	return proxy.xs.SetInteger(name, value)
}

func (proxy *XSProxy) SetColor(name string, value []byte) error {
	return proxy.xs.SetColor(name, value)
}

func (proxy *XSProxy) GetString(name string) (string, error) {
	v, _, err := proxy.xs.GetString(name)
	return v, err
}

func (proxy *XSProxy) GetInteger(name string) (uint32, error) {
	v, _, err := proxy.xs.GetInteger(name)
	return v, err
}

func (proxy *XSProxy) GetColor(name string) ([]byte, error) {
	v, _, err := proxy.xs.GetColor(name)
	return v, err
}

func (proxy *XSProxy) refer() {
	proxy.locker.Lock()
	defer proxy.locker.Unlock()

	proxy.cnt++
}

func (proxy *XSProxy) unref() {
	proxy.locker.Lock()
	defer proxy.locker.Unlock()

	proxy.cnt--
	if proxy.cnt == 0 {
		sessionmanager.DestroyXSettings(proxy.xs)
		proxy = nil
	}
}
