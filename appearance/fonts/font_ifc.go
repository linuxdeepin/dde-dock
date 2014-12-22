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
	xsettings "pkg.linuxdeepin.com/dde-daemon/xsettings_wrapper"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	FontTypeStandard   = "font-standard"
	FontTypeMonospaced = "font-mono"
)

func (font *FontManager) IsStandardFontValid(name string) bool {
	for _, info := range getStandardFonts() {
		if name == info.Id {
			return true
		}
	}

	return false
}

func (font *FontManager) IsMonospacedFontValid(name string) bool {
	for _, info := range getMonospaceFonts() {
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

		xsettings.SetString(xsettings.GtkStringFontName, value)
		WMSetString(WMTitlebarFont, value)
		setQt4Font(GetUserQt4Config(), name, size)
		WriteUserGtk3Config(GetUserGtk3Config(),
			"gtk-font-name", value)
		WriteUserGtk2Config(GetUserGtk2Config(),
			"gtk-font-name", value)
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
	xsettings.SetInteger(xsettings.XftBoolAntialias, anti)
	xsettings.SetInteger(xsettings.XftBoolHinting, hinting)
	xsettings.SetString(xsettings.XftStringHintStyle, hintstyle)
	xsettings.SetString(xsettings.XftStringRgba, rgba)

	WriteUserGtk2Config(GetUserGtk2Config(),
		"gtk-xft-antialias", fmt.Sprintf("%v", anti))
	WriteUserGtk2Config(GetUserGtk2Config(),
		"gtk-xft-hinting", fmt.Sprintf("%v", hinting))
	WriteUserGtk2Config(GetUserGtk2Config(),
		"gtk-xft-hintstyle", hintstyle)
	WriteUserGtk2Config(GetUserGtk2Config(),
		"gtk-xft-rgba", rgba)

	WriteUserGtk3Config(GetUserGtk3Config(),
		"gtk-xft-antialias", fmt.Sprintf("%v", anti))
	WriteUserGtk3Config(GetUserGtk3Config(),
		"gtk-xft-hinting", fmt.Sprintf("%v", hinting))
	WriteUserGtk3Config(GetUserGtk3Config(),
		"gtk-xft-hintstyle", hintstyle)
	WriteUserGtk3Config(GetUserGtk3Config(),
		"gtk-xft-rgba", rgba)
}

func (font *FontManager) GetStyleListByName(name string) []string {
	infos := getStandardFonts()
	infos = append(infos, getMonospaceFonts()...)

	return getStyleList(name, infos)
}

func (font *FontManager) GetNameList(fontType string) []string {
	switch fontType {
	case FontTypeStandard:
		return getNameStrList(getStandardFonts())
	case FontTypeMonospaced:
		return getNameStrList(getMonospaceFonts())
	}

	return nil
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
