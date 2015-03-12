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

package fonts

import (
	"fmt"
	"os/exec"
	. "pkg.linuxdeepin.com/dde-daemon/appearance/utils"
	"pkg.linuxdeepin.com/dde-daemon/xsettings"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	FontTypeStandard   = "font-standard"
	FontTypeMonospaced = "font-mono"

	wmGSettingsSchema = "org.gnome.desktop.wm.preferences"
)

func (font *FontManager) IsStandardFontValid(name string) bool {
	for _, info := range font.standardList {
		if name == info.Id {
			return true
		}
	}

	return false
}

func (font *FontManager) IsMonospacedFontValid(name string) bool {
	for _, info := range font.monospaceList {
		if name == info.Id {
			return true
		}
	}

	return false
}

// fontsize: 9 ~ 26
func (font *FontManager) IsFontSizeValid(size int32) bool {
	if size < 9 || size > 26 {
		return false
	}

	return true
}

func (font *FontManager) Set(fontType, name string, size int32) {
	value := fmt.Sprintf("%s %v", name, size)
	switch fontType {
	case FontTypeStandard:
		if !font.IsStandardFontValid(name) ||
			!font.IsFontSizeValid(size) {
			return
		}

		if font.xsProxy == nil {
			return
		}
		font.xsProxy.SetString(xsettings.GtkFontName, value)

		settings := CheckAndNewGSettings(wmGSettingsSchema)
		if settings != nil {
			settings.SetString("titlebar-font", value)
			//Unref(settings)
		}
		setQt4Font(GetUserQt4Config(), name, size)
	case FontTypeMonospaced:
		if !font.IsMonospacedFontValid(name) ||
			!font.IsFontSizeValid(size) {
			return
		}

		setMonoFont(name, size)
	}
}

/**
 * xft-antialias default 1(true)
 * xft-hinting default 1(true)
 * xft-hintstyle default "hintfull"
 * xft-rgba default "rgb"
 * lcdfilter default "lcddefault"
 */
func (font *FontManager) SetXft(anti, hinting uint32, hintstyle, rgba string) {
	if font.xsProxy == nil {
		return
	}
	font.xsProxy.SetInteger(xsettings.XftAntialias, anti)
	font.xsProxy.SetInteger(xsettings.XftHinting, hinting)
	font.xsProxy.SetString(xsettings.XftHintStyle, hintstyle)
	font.xsProxy.SetString(xsettings.XftRgba, rgba)
}

func (font *FontManager) GetStyleListByName(name string) []string {
	infos := font.standardList
	infos = append(infos, font.monospaceList...)

	return getStyleList(name, infos)
}

func (font *FontManager) GetNameList(fontType string) []string {
	switch fontType {
	case FontTypeStandard:
		font.standardList = getStandardFonts()
		return getNameStrList(font.standardList)
	case FontTypeMonospaced:
		font.monospaceList = getMonospaceFonts()
		return getNameStrList(font.monospaceList)
	}

	return nil
}

func (font *FontManager) Destroy() {
	if font.xsProxy != nil {
		font.xsProxy.Free()
		font.xsProxy = nil
	}
}

func setQt4Font(config, name string, size int32) {
	value := fmt.Sprintf("\"%s, %v, -1, 5, 50, 0, 0, 0, 0, 0\"",
		name, size)
	dutils.WriteKeyToKeyFile(config, "Qt", "font", value)
}

/**
 * gconf
 * /desktop/gnome/interface/monospace_font_name
 */
func setMonoFont(name string, size int32) {
	value := fmt.Sprintf("%s %v", name, size)
	cmdline := "gconftool -t string -s /desktop/gnome/interface/monospace_font_name \"" + value + "\""
	exec.Command("/bin/sh", "-c", cmdline).Run()
}
