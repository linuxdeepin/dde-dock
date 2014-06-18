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

package themes

import (
	"dlib/gio-2.0"
	"github.com/howeyc/fsnotify"
	"os/exec"
	"path"
	"strconv"
)

type Theme struct {
	Name        string
	GtkTheme    string
	IconTheme   string
	SoundTheme  string
	CursorTheme string
	Background  string
	FontSize    int32
	Type        int32
	Preview     []string
	filePath    string
	objectPath  string

	watcher  *fsnotify.Watcher
	quitFlag chan bool
}

const (
	QT_CONFIG_FILE    = ".config/Trolltech.conf"
	DEFAULT_FONT_SIZE = " 11"

	QT_KEY_GROUP   = "Qt"
	QT_KEY_STYLE   = "stype"
	QT_STYLE_VALUE = "GTK+"
	QT_KEY_FONT    = "font"
	QT_FONT_ARGS   = ",-1,5,50,0,0,0,0,0"

	DEFAULT_FONT      = "WenQuanYi Micro Hei"
	DEFAULT_FONT_MONO = "WenQuanYi Micro Hei Mono"
)

var (
	wmPreSettings = gio.NewSettings("org.gnome.desktop.wm.preferences")
)

func (obj *Theme) setAllThemes() {
	obj.setGtkTheme()
	obj.setIconTheme()
	obj.setCursorTheme()
	obj.setFontName()

	bg := themeSettings.GetString(GS_KEY_CURRENT_BG)
	bg = decodeURI(bg)
	if obj.Background != bg {
		themeSettings.SetString(GS_KEY_CURRENT_BG,
			encodeURI(obj.Background))
	}
}

func (obj *Theme) setGtkTheme() {
	objXS.SetString("Net/ThemeName", obj.GtkTheme)
	wmPreSettings.SetString("theme", obj.GtkTheme)
	homeDir, _ := objUtil.GetHomeDir()
	if ok := objUtil.WriteKeyToKeyFile(path.Join(homeDir, QT_CONFIG_FILE),
		QT_KEY_GROUP, QT_KEY_STYLE, QT_STYLE_VALUE); !ok {
		Logger.Error("Set QT Style Failed")
		return
	}
}

func (obj *Theme) setIconTheme() {
	objXS.SetString("Net/IconThemeName", obj.IconTheme)
}

func (obj *Theme) setCursorTheme() {
	objXS.SetString("Gtk/CursorThemeName", obj.CursorTheme)
}

func (obj *Theme) setFontName() {
	size := ""
	if obj.FontSize < 1 {
		size = DEFAULT_FONT_SIZE
	} else {
		size = strconv.FormatInt(int64(obj.FontSize), 10)
	}
	objXS.SetString("Gtk/FontName", DEFAULT_FONT+" "+size)
	wmPreSettings.SetString("titlebar-font", DEFAULT_FONT+" Bold "+size)
	homeDir, _ := objUtil.GetHomeDir()
	if ok := objUtil.WriteKeyToKeyFile(path.Join(homeDir, QT_CONFIG_FILE),
		QT_KEY_GROUP, QT_KEY_FONT,
		"\""+DEFAULT_FONT+","+size+QT_FONT_ARGS+"\""); !ok {
		Logger.Error("Set QT Font Failed")
		return
	}
	setMonoFont(DEFAULT_FONT_MONO, size)
}

func setMonoFont(name, size string) {
	if len(name) <= 0 {
		return
	}

	if len(size) <= 0 {
		size = "10"
	}

	args := []string{}
	args = append(args, "-t")
	args = append(args, "string")
	args = append(args, "-s")
	args = append(args, "/desktop/gnome/interface/monospace_font_name")
	args = append(args, name+" "+size)

	exec.Command("/usr/bin/gconftool", args...).Run()
}

func newTheme(info ThemeInfo) *Theme {
	t := &Theme{}

	t.Type = info.T
	t.filePath = info.Path

	t.setAllProps()
	t.objectPath = THEME_PATH + t.Name
	//t.setAllThemes()

	if t.Type == THEME_TYPE_LOCAL {
		var err error
		t.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			Logger.Errorf("New Watcher Failed: %v", err)
			panic(err)
		}

		t.quitFlag = make(chan bool)

		t.startWatch()
		go t.handleEvent()
	}

	return t
}
