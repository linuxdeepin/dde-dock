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

package xsettings_wrapper

import (
	"dbus/com/deepin/sessionmanager"
	"fmt"
)

const (
	NetStringThemeName               = "Net/ThemeName"
	NetStringIconTheme               = "Net/IconThemeName"
	NetStringFallbackIconTheme       = "Net/FallbackIconTheme"
	NetBoolCursorBlink               = "Net/CursorBlink"
	NetIntCursorBlinkTime            = "Net/CursorBlinkTime"
	NetIntCursorBlinkTimeout         = "Net/CursorBlinkTimeout"
	NetIntDoubleClick                = "Net/DoubleClickTime"
	NetIntDragThreshold              = "Net/DndDragThreshold"
	NetStringSoundTheme              = "Net/SoundThemeName"
	NetBoolEnableEventSounds         = "Net/EnableEventSounds"
	NetBoolEnableInputFeedbackSounds = "Net/EnableInputFeedbackSounds"

	GtkStringThemeName         = "Gtk/GtkThemeName"
	GtkStringCursorTheme       = "Gtk/CursorThemeName"
	GtkIntCursorThemeSize      = "Gtk/CursorThemeSize"
	GtkStringFontName          = "Gtk/FontName"
	GtkStringKeyTheme          = "Gtk/KeyThemeName"
	GtkStringToolbarStyle      = "Gtk/ToolbarStyle"
	GtkStringToolbarIconSize   = "Gtk/ToolbarIconSize"
	GtkBoolCanChangeAccels     = "Gtk/CanChangeAccels"
	GtkStringColorPalette      = "Gtk/ColorPalette"
	GtkIntTimeoutInitial       = "Gtk/TimeoutInitial"
	GtkIntTimeoutRepeat        = "Gtk/TimeoutRepeat"
	GtkStringColorScheme       = "Gtk/ColorScheme"
	GtkStringIMPreeditStyle    = "Gtk/IMPreeditStyle"
	GtkStringIMStatusStyle     = "Gtk/IMStatusStyle"
	GtkStringIMModule          = "Gtk/IMModule"
	GtkBoolMenuImages          = "Gtk/MenuImages"
	GtkBoolButtonImages        = "Gtk/ButtonImages"
	GtkStringMenuBarAccel      = "Gtk/MenuBarAccel"
	GtkBoolEnableAnimations    = "Gtk/EnableAnimations"
	GtkBoolShowInputMethodMenu = "Gtk/ShowInputMethodMenu"
	GtkBoolShowUnicodeMenu     = "Gtk/ShowUnicodeMenu"
	GtkBoolAutoMnemonics       = "Gtk/AutoMnemonics"
	GtkIntRecentFilesMaxAge    = "Gtk/RecentFilesMaxAge"
	GtkBoolEnableRecentFiles   = "Gtk/RecentFilesEnabled"

	XftBoolAntialias   = "xft/Antialias"
	XftBoolHinting     = "xft/HintStyle"
	XftStringHintStyle = "xft/HintStyle"
	XftStringRgba      = "xft/RGBA"
)

var (
	errUninitialized = fmt.Errorf("XSettings uninitialized")
)

var _xsettings *sessionmanager.XSettings

// Must be called before using other methods
func InitXSettings() error {
	if _xsettings != nil {
		return nil
	}

	xsettings, err := sessionmanager.NewXSettings(
		"com.deepin.SessionManager",
		"/com/deepin/XSettings",
	)
	if err != nil {
		return err
	}
	_xsettings = xsettings

	return nil
}

func SetString(name string, value string) error {
	if _xsettings == nil {
		return errUninitialized
	}

	return _xsettings.SetString(name, value)
}

func SetInteger(name string, value uint32) error {
	if _xsettings == nil {
		return errUninitialized
	}

	return _xsettings.SetInteger(name, value)
}

func SetColor(name string, value []byte) error {
	if _xsettings == nil {
		return errUninitialized
	}

	return _xsettings.SetColor(name, value)
}
